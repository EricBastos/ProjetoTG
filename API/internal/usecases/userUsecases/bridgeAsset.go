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
	"math/big"
	"net/http"
)

type CreateBridgeAssetUsecase struct {
	userInfo           *utils.UserInformation
	bridgeOperationsDb database.BridgeOpInterface
	rabbitClient       *rabbitmqClient.RabbitMQClient
}

func NewCreateBridgeAssetUsecase(
	userInfo *utils.UserInformation,
	bridgeOperationsDb database.BridgeOpInterface,
	rabbitClient *rabbitmqClient.RabbitMQClient,
) *CreateBridgeAssetUsecase {
	return &CreateBridgeAssetUsecase{
		userInfo:           userInfo,
		bridgeOperationsDb: bridgeOperationsDb,
		rabbitClient:       rabbitClient,
	}
}

func (u *CreateBridgeAssetUsecase) Bridge(input *dtos.BridgeAssetInput) (string, error, int) {

	id, creationError, creationCode := u.createAsUser(input)

	if creationError != nil {
		return "", creationError, creationCode
	}

	return id, nil, 0
}

func (u *CreateBridgeAssetUsecase) createAsUser(input *dtos.BridgeAssetInput) (string, error, int) {

	if !utils.VerifyBurnPermit(input.WalletAddress, input.InputChain, input.Amount, input.Permit, &utils.Contracts{
		EthereumWalletAddress: configs.Cfg.EthereumWalletAddress,
		EthereumTokenAddress:  configs.Cfg.EthereumTokenContract,
		PolygonWalletAddress:  configs.Cfg.PolygonWalletAddress,
		PolygonTokenAddress:   configs.Cfg.PolygonTokenContract,
	}) {
		return "", errors.New(utils.SignatureError), http.StatusUnauthorized
	}

	var err error

	waiting := false

	switch input.InputChain {
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

	var op *entities.BridgeOp

	op = entities.NewBridge(
		input.WalletAddress,
		input.Amount,
		input.InputChain,
		input.OutputChain,
		&uId,
		input.Permit,
	)

	err = u.bridgeOperationsDb.Create(op)
	if err != nil {
		return "", errors.New(utils.OperationCreationError), http.StatusInternalServerError
	}

	burnOp := entities.NewBurnWithPermit(
		input.WalletAddress,
		input.Amount,
		"",
		"",
		"",
		"",
		"",
		input.InputChain,
		&uId,
		input.Permit,
	)

	burnOp.Id = op.Id

	err = u.rabbitClient.CallSmartcontract(burnOp, entities.BRIDGE)
	if err != nil {
		return "", errors.New(utils.OperationCreationError), http.StatusInternalServerError
	}
	return op.Id.String(), nil, 0
}
