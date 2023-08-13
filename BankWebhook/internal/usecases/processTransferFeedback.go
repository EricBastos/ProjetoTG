package usecases

import (
	"fmt"
	"github.com/EricBastos/ProjetoTG/BankWebhook/internal/dtos"
	"github.com/EricBastos/ProjetoTG/Library/database"
	"github.com/EricBastos/ProjetoTG/Library/entities"
	"log"
)

type ProcessTransferFeedbackUsecase struct {
	transferFeedbackDb database.TransferFeedbackInterface
}

func NewProcessTransferFeedbackUsecase(
	transferFeedbackDb database.TransferFeedbackInterface,
) *ProcessTransferFeedbackUsecase {
	return &ProcessTransferFeedbackUsecase{
		transferFeedbackDb: transferFeedbackDb,
	}
}

func (u *ProcessTransferFeedbackUsecase) Process(transfer *dtos.TransferFeedbackInput) error {

	transfFeedback := entities.NewTransferFeedback(
		transfer.TransferId,
	)

	err := u.transferFeedbackDb.Create(transfFeedback)
	if err != nil {
		transfData := fmt.Sprintf("%+v", transfFeedback)
		log.Println("(Transfer Feedback) Error creating transfer feedback, op data: " + transfData + ", err: " + err.Error())
		return err
	}
	return nil

}
