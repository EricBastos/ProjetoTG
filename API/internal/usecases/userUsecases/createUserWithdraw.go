package userUsecases

import (
	"errors"
	"github.com/EricBastos/ProjetoTG/API/configs"
	"github.com/EricBastos/ProjetoTG/API/internal/dtos"
	"github.com/EricBastos/ProjetoTG/API/internal/grpcClient"
	"github.com/EricBastos/ProjetoTG/API/internal/infra/rabbitmqClient"
	"github.com/EricBastos/ProjetoTG/Library/database"
	"github.com/EricBastos/ProjetoTG/Library/entities"
	entities2 "github.com/EricBastos/ProjetoTG/Library/pkg/entities"
	"github.com/EricBastos/ProjetoTG/Library/utils"
	"log"
	"math/big"
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

	if !utils.VerifyBurnPermit(input.WalletAddress, input.Chain, input.Amount, input.Permit, &utils.Contracts{
		EthereumWalletAddress: configs.Cfg.EthereumWalletAddress,
		EthereumTokenAddress:  configs.Cfg.EthereumTokenContract,
		PolygonWalletAddress:  configs.Cfg.PolygonWalletAddress,
		PolygonTokenAddress:   configs.Cfg.PolygonTokenContract,
	}) {
		return "", errors.New(utils.SignatureError), http.StatusUnauthorized
	}

	var err error

	waiting := false

	switch input.Chain {
	case "Ethereum":
		balance, balErr := grpcClient.EthereumService.GetBalance(configs.Cfg.EthereumTokenContract, input.WalletAddress)
		if balErr != nil {
			return "", errors.New(utils.InternalError), http.StatusInternalServerError
		}
		bigInputAmount := new(big.Int).Mul(big.NewInt(int64(input.Amount)), big.NewInt(10^16))
		if balance.Cmp(bigInputAmount) < 0 {
			return "", errors.New("wallet balance must be greater than burn amount"), http.StatusBadRequest
		}
		waiting, err = grpcClient.EthereumService.IsWaitingPermit(input.WalletAddress)
		if err != nil {
			waiting = true
		}
	case "Polygon":
		balance, balErr := grpcClient.PolygonService.GetBalance(configs.Cfg.PolygonTokenContract, input.WalletAddress)
		if balErr != nil {
			return "", errors.New(utils.InternalError), http.StatusInternalServerError
		}
		log.Println(balance)
		bigInputAmount := new(big.Int).Mul(big.NewInt(int64(input.Amount)), big.NewInt(10000000000000000))
		if balance.Cmp(bigInputAmount) < 0 {
			return "", errors.New("wallet balance must be greater than burn amount"), http.StatusBadRequest
		}
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
