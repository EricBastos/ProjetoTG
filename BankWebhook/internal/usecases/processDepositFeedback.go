package usecases

import (
	"fmt"
	"github.com/EricBastos/ProjetoTG/BankWebhook/internal/dtos"
	"github.com/EricBastos/ProjetoTG/BankWebhook/internal/infra/rabbitmqClient"
	"github.com/EricBastos/ProjetoTG/Library/database"
	"github.com/EricBastos/ProjetoTG/Library/entities"
	"github.com/EricBastos/ProjetoTG/Library/utils"
	"log"
	"sync"
)

type ProcessDepositFeedbackUsecase struct {
	mintOperationsDb        database.MintOpInterface
	staticDepositDb         database.StaticDepositInterface
	staticDepositFeedbackDb database.StaticDepositFeedbackInterface
	rabbitClient            *rabbitmqClient.RabbitMQClient
}

func NewProcessDepositFeedbackUsecase(
	mintOperationsDb database.MintOpInterface,
	staticDepositDb database.StaticDepositInterface,
	staticDepositFeedbackDb database.StaticDepositFeedbackInterface,
	rabbitClient *rabbitmqClient.RabbitMQClient,
) *ProcessDepositFeedbackUsecase {
	return &ProcessDepositFeedbackUsecase{
		mintOperationsDb:        mintOperationsDb,
		staticDepositDb:         staticDepositDb,
		staticDepositFeedbackDb: staticDepositFeedbackDb,
		rabbitClient:            rabbitClient,
	}
}

func (u *ProcessDepositFeedbackUsecase) Process(deposit *dtos.DepositFeedbackInput) error {

	deposit.TaxId = utils.TrimCpfCnpj(deposit.TaxId)

	depFeedback := entities.NewStaticDepositFeedback(
		deposit.DepositId,
		deposit.TaxId,
		deposit.Amount,
	)
	err := u.staticDepositFeedbackDb.Create(depFeedback)
	if err != nil {
		depData := fmt.Sprintf("%+v", depFeedback)
		log.Println("(Deposit Feedback) Error creating deposit feedback. Err: " + err.Error() + ", deposit data: " + depData)
		return err
	}

	go func() {

		fetchedDeposit, err := u.staticDepositDb.FindUnpaidByTaxIdAndAmount(deposit.TaxId, deposit.Amount)
		if err != nil {
			log.Println("Error fetching data related to received deposit:", err.Error())
			return
		}

		var op entities.SmartContractOp

		var updateDepositErr error
		var createEmitErr error
		wg := sync.WaitGroup{}
		wg.Add(2)

		go func() {
			defer wg.Done()
			fetchedDeposit.Status = "PAID"
			updateDepositErr = u.staticDepositDb.Update(fetchedDeposit)
		}()

		go func() {
			defer wg.Done()

			mintOp := entities.NewMint(
				fetchedDeposit.WalletAddress,
				fetchedDeposit.Amount,
				fetchedDeposit.Chain,
				"deposit paid",
				fetchedDeposit.ResponsibleUser,
				fetchedDeposit.Id.String(),
			)
			mintOp.AssociatedBankTransactionType = "static_deposits"
			op = mintOp

			createEmitErr = u.mintOperationsDb.Create(mintOp)
			if createEmitErr != nil {
				return
			}
			createEmitErr = u.rabbitClient.CallSmartcontract(op, entities.MINT)
			if createEmitErr != nil {
				return
			}

		}()

		wg.Wait()

		if updateDepositErr != nil {
			dData := fmt.Sprintf("%+v", fetchedDeposit)
			log.Println("(Deposit Feedback [paid|credited]) Error updating deposit in DB, deposit data: " + dData + ", err: " + updateDepositErr.Error())
		}

		if createEmitErr != nil {
			// Critical error: we couldn't post the mint operation.
			// We can't mint the user's coins

			opData := fmt.Sprintf("%+v", op)
			log.Println("(Deposit Feedback [paid|credited]) Error sending op to rabbitmq, op data: " + opData + ", err: " + createEmitErr.Error())
		}

	}()

	return nil
}
