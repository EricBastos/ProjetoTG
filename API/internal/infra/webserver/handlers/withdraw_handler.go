package handlers

import (
	"encoding/json"
	"errors"
	"github.com/EricBastos/ProjetoTG/API/internal/dtos"
	"github.com/EricBastos/ProjetoTG/API/internal/infra/rabbitmqClient"
	"github.com/EricBastos/ProjetoTG/API/internal/usecases/userUsecases"
	"github.com/EricBastos/ProjetoTG/API/internal/utils"
	"github.com/EricBastos/ProjetoTG/Library/database"
	"net/http"
)

type WithdrawHandler struct {
	burnOperationsDb database.BurnOpInterface
	rabbitClient     *rabbitmqClient.RabbitMQClient
}

func NewWithdrawHandler(
	burnOperationsDb database.BurnOpInterface,
	rabbitClient *rabbitmqClient.RabbitMQClient,
) *WithdrawHandler {
	return &WithdrawHandler{
		burnOperationsDb: burnOperationsDb,
		rabbitClient:     rabbitClient,
	}
}

func (h *WithdrawHandler) CreateUserWithdraw(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	userInfo := &utils.UserInformation{
		UserId: r.Context().Value("subject").(string),
		Name:   r.Context().Value("name").(string),
		TaxId:  r.Context().Value("taxId").(string),
		Email:  r.Context().Value("email").(string),
	}
	var input dtos.CreateUserWithdrawInput
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
	if err := firstValidationCreateUserWithdrawInput(&input); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(struct {
			Error string `json:"error"`
		}{
			Error: err.Error(),
		})
		return
	}
	usecase := userUsecases.NewCreateUserWithdrawUsecase(
		userInfo,
		h.burnOperationsDb,
		h.rabbitClient,
	)

	id, err, code := usecase.Create(&input)
	if err != nil {
		w.WriteHeader(code)
		_ = json.NewEncoder(w).Encode(struct {
			Error string `json:"error"`
		}{
			Error: err.Error(),
		})
		return
	}

	res := dtos.CreateWithdrawOutput{Id: id}
	_ = json.NewEncoder(w).Encode(res)
	w.WriteHeader(http.StatusCreated)
}

func firstValidationCreateUserWithdrawInput(input *dtos.CreateUserWithdrawInput) error {
	if input.Amount < 0 {
		return errors.New(utils.InvalidAmount)
	}
	if input.PixKey == "" {
		return errors.New(utils.MissingPixKey)
	}
	if input.WalletAddress == "" {
		return errors.New(utils.MissingWalletAddress)
	}
	if input.Chain == "" {
		return errors.New(utils.MissingChain)
	}
	if input.Permit == nil {
		return errors.New(utils.MissingPermit)
	}
	return nil
}
