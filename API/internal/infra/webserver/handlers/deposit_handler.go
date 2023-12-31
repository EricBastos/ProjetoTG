package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/EricBastos/ProjetoTG/API/configs"
	"github.com/EricBastos/ProjetoTG/API/internal/dtos"
	"github.com/EricBastos/ProjetoTG/API/internal/usecases/userUsecases"
	"github.com/EricBastos/ProjetoTG/Library/database"
	"github.com/EricBastos/ProjetoTG/Library/utils"
	"log"
	"net/http"
	"time"
)

type DepositHandler struct {
	staticDepositDb database.StaticDepositInterface
}

func NewDepositHandler(
	staticDepositDb database.StaticDepositInterface,
) *DepositHandler {
	return &DepositHandler{
		staticDepositDb: staticDepositDb,
	}
}

func (h *DepositHandler) CreatePixDeposit(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	userInfo := &utils.UserInformation{
		UserId: r.Context().Value("subject").(string),
		Name:   r.Context().Value("name").(string),
		TaxId:  r.Context().Value("taxId").(string),
		Email:  r.Context().Value("email").(string),
	}
	var input dtos.CreateUserStaticPixDepositInput
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
	if err := firstValidationCreateUserStaticPixDepositInput(&input); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(struct {
			Error string `json:"error"`
		}{
			Error: err.Error(),
		})
		return
	}
	usecase := userUsecases.NewCreateUserStaticPixDepositUsecase(userInfo, h.staticDepositDb)
	var errCode int
	var depositId string
	depositId, err, errCode = usecase.CreateDeposit(&input)
	if err != nil {
		w.WriteHeader(errCode)
		_ = json.NewEncoder(w).Encode(struct {
			Error string `json:"error"`
		}{
			Error: err.Error(),
		})
		return
	}

	// Activate Webhook after 10 seconds mocking a PIX deposit
	go func() {
		time.Sleep(10 * time.Second)
		err := postStaticPixDepositWebhook(input.Amount, userInfo.TaxId, depositId)
		if err != nil {
			log.Println("(SANDBOX) Error posting pseudo webhook for deposit: " + err.Error())
		}
	}()

	resp := dtos.CreateStaticPixDepositOutput{Id: depositId}
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(resp)
}

func postStaticPixDepositWebhook(amount int, taxId string, depositId string) error {
	body, _ := json.Marshal(map[string]interface{}{
		"subscription": "deposit",
		"data": map[string]interface{}{
			"amount":    amount,
			"taxId":     taxId,
			"depositId": depositId,
		},
	})

	webhookUrl := "http://" + configs.Cfg.BankWebhookHost + ":" + configs.Cfg.BankWebhookPort
	req, err := http.NewRequest("POST", webhookUrl, bytes.NewReader(body))
	req.Close = true
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	if err != nil {
		return err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	_ = resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return errors.New(http.StatusText(http.StatusInternalServerError))
	}
	return nil
}

func firstValidationCreateUserStaticPixDepositInput(input *dtos.CreateUserStaticPixDepositInput) error {
	if input.Amount <= 0 {
		return errors.New(utils.InvalidAmount)
	}
	validateChainFunc, ok := utils.ValidChains[input.Chain]
	if !ok {
		return errors.New(utils.InvalidChain)
	}
	if err := validateChainFunc(&input.WalletAddress); err != nil {
		return err
	}
	return nil
}
