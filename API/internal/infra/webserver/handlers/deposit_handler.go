package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/EricBastos/ProjetoTG/API/configs"
	"github.com/EricBastos/ProjetoTG/API/internal/dtos"
	"github.com/EricBastos/ProjetoTG/API/internal/usecases/userUsecases"
	"github.com/EricBastos/ProjetoTG/API/internal/utils"
	"github.com/EricBastos/ProjetoTG/Library/database"
	"log"
	"net/http"
	"time"
)

type DepositHandler struct {
	staticDepositDb database.StaticDepositInterface
	httpClient      *http.Client
}

func NewDepositHandler(
	staticDepositDb database.StaticDepositInterface,
) *DepositHandler {
	return &DepositHandler{
		staticDepositDb: staticDepositDb,
		httpClient:      &http.Client{},
	}
}

func (h *DepositHandler) CreatePixDeposit(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	userInfo := &utils.UserInformation{
		UserId: r.Context().Value("subject").(string),
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

	// Activate Webhook after 10 seconds
	go func() {
		time.Sleep(10 * time.Second)
		err := postStaticPixDepositWebhook(input.Amount, userInfo.TaxId, h.httpClient)
		if err != nil {
			log.Println("(SANDBOX) Error posting pseudo webhook for deposit: " + err.Error())
		}
	}()

	resp := dtos.CreateStaticPixDepositOutput{Id: depositId}
	_ = json.NewEncoder(w).Encode(resp)
	w.WriteHeader(http.StatusCreated)
}

func postStaticPixDepositWebhook(amount int, taxId string, httpClient *http.Client) error {

	body, _ := json.Marshal(map[string]interface{}{
		"subscription": "deposit",
		"amount":       amount,
		"taxId":        taxId,
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
	resp, err := httpClient.Do(req)
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
	if input.Amount != 0 {
		return errors.New(utils.InvalidMintAmount)
	}
	return nil
}
