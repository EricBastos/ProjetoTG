package userUsecases

import (
	"errors"
	"github.com/EricBastos/ProjetoTG/API/internal/dtos"
	"github.com/EricBastos/ProjetoTG/API/internal/grpcClient"
	"github.com/EricBastos/ProjetoTG/API/internal/infra/rabbitmqClient"
	"github.com/EricBastos/ProjetoTG/API/internal/utils"
	"github.com/EricBastos/ProjetoTG/Library/database"
	"github.com/EricBastos/ProjetoTG/Library/entities"
	entities2 "github.com/EricBastos/ProjetoTG/Library/pkg/entities"
	"net/http"
)

type CreateUserWithdrawUsecase struct {
	userInfo         *utils.UserInformation
	burnOperationsDB database.BurnOpInterface
	rabbitClient     *rabbitmqClient.RabbitMQClient
}

func NewCreateUserWithdrawUsecase(
	userInfo *utils.UserInformation,
	burnOperationsDB database.BurnOpInterface,
	rabbitClient *rabbitmqClient.RabbitMQClient,
) *CreateUserWithdrawUsecase {
	return &CreateUserWithdrawUsecase{
		userInfo:         userInfo,
		burnOperationsDB: burnOperationsDB,
		rabbitClient:     rabbitClient,
	}
}

func (u *CreateUserWithdrawUsecase) Create(input *dtos.CreateUserWithdrawInput) (string, error, int) {

	id, creationError, creationCode := u.createAsUser(input)

	if creationError != nil {
		return "", creationError, creationCode
	}

	return id, nil, 0
}

func (u *CreateUserWithdrawUsecase) createAsUser(input *dtos.CreateUserWithdrawInput) (string, error, int) {

	if !utils.VerifyBurnPermit(input.WalletAddress, input.Chain, input.Amount, input.Permit) {
		return "", errors.New(utils.SignatureError), http.StatusUnauthorized
	}

	var err error

	waiting := false

	switch input.Chain {
	case "Ethereum":
		waiting, err = grpcClient.EthereumService.IsWaitingPermit(input.WalletAddress)
		if err != nil {
			waiting = true
		}
	case "Polygon":
		waiting, err = grpcClient.PolygonService.IsWaitingPermit(input.WalletAddress)
		if err != nil {
			waiting = true
		}
	}

	if waiting {
		return "", errors.New(utils.UserWaitingPermit), http.StatusBadRequest
	}

	uId, err := entities2.ParseID(u.userInfo.UserId)
	if err != nil {
		return "", errors.New(utils.InternalError), http.StatusInternalServerError
	}

	var op *entities.BurnOp

	op = entities.NewBurnWithPermit(
		input.WalletAddress,
		input.Amount,
		u.userInfo.Name,
		u.userInfo.TaxId,
		"Mocked bank code",
		"Mocked branch code",
		"Mocked account number",
		input.Chain,
		&uId,
		input.Permit,
	)

	err = u.burnOperationsDB.Create(op)
	if err != nil {
		return "", errors.New(utils.OperationCreationError), http.StatusInternalServerError
	}
	err = u.rabbitClient.CallSmartcontract(op, entities.BURN)
	if err != nil {
		return "", errors.New(utils.OperationCreationError), http.StatusInternalServerError
	}
	return op.Id.String(), nil, 0
}
