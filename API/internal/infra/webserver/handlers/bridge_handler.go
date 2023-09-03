package handlers

import (
	"encoding/json"
	"errors"
	"github.com/EricBastos/ProjetoTG/API/internal/dtos"
	"github.com/EricBastos/ProjetoTG/API/internal/infra/rabbitmqClient"
	"github.com/EricBastos/ProjetoTG/API/internal/usecases/userUsecases"
	"github.com/EricBastos/ProjetoTG/Library/database"
	"github.com/EricBastos/ProjetoTG/Library/utils"
	"net/http"
)

type BridgeHandler struct {
	bridgeOperationsDb database.BridgeOpInterface
	rabbitClient       *rabbitmqClient.RabbitMQClient
}

func NewBridgeHandler(
	bridgeOperationsDb database.BridgeOpInterface,
	rabbitClient *rabbitmqClient.RabbitMQClient,
) *BridgeHandler {
	return &BridgeHandler{
		bridgeOperationsDb: bridgeOperationsDb,
		rabbitClient:       rabbitClient,
	}
}

func (h *BridgeHandler) BridgeAsset(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	userInfo := &utils.UserInformation{
		UserId: r.Context().Value("subject").(string),
		Name:   r.Context().Value("name").(string),
		TaxId:  r.Context().Value("taxId").(string),
		Email:  r.Context().Value("email").(string),
	}
	var input dtos.BridgeAssetInput
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(struct {
			Error string `json:"error"`
		}{
			Error: utils.BadRequest,
		})
		return
	}
	if err := validateBridgeAssetInput(&input); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(struct {
			Error string `json:"error"`
		}{
			Error: err.Error(),
		})
		return
	}
	usecase := userUsecases.NewCreateBridgeAssetUsecase(userInfo, h.bridgeOperationsDb, h.rabbitClient)
	var errCode int
	var depositId string
	depositId, err, errCode = usecase.Bridge(&input)
	if err != nil {
		w.WriteHeader(errCode)
		_ = json.NewEncoder(w).Encode(struct {
			Error string `json:"error"`
		}{
			Error: err.Error(),
		})
		return
	}

	resp := dtos.BridgeAssetOutput{Id: depositId}
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(resp)
}

func validateBridgeAssetInput(input *dtos.BridgeAssetInput) error {
	if input.Amount < 0 {
		return errors.New(utils.InvalidAmount)
	}
	if input.WalletAddress == "" {
		return errors.New(utils.MissingWalletAddress)
	}
	if input.InputChain == "" {
		return errors.New(utils.MissingInputChain)
	}
	if input.OutputChain == "" {
		return errors.New(utils.MissingOutputChain)
	}
	if input.Permit == nil {
		return errors.New(utils.MissingPermit)
	}

	if input.InputChain == input.OutputChain {
		return errors.New(utils.ChainsAreEqual)
	}

	validateChainFunc, ok := utils.ValidChains[input.InputChain]
	if !ok {
		return errors.New(utils.InvalidChain)
	}
	if err := validateChainFunc(&input.WalletAddress); err != nil {
		return err
	}
	validateChainFunc, ok = utils.ValidChains[input.OutputChain]
	if !ok {
		return errors.New(utils.InvalidChain)
	}
	if err := validateChainFunc(&input.WalletAddress); err != nil {
		return err
	}
	return nil
}
