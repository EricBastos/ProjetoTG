package handlers

import (
	"encoding/json"
	"errors"
	"github.com/EricBastos/ProjetoTG/API/configs"
	"github.com/EricBastos/ProjetoTG/API/internal/dtos"
	"github.com/EricBastos/ProjetoTG/API/internal/usecases/userUsecases"
	"github.com/EricBastos/ProjetoTG/Library/database"
	"github.com/EricBastos/ProjetoTG/Library/utils"
	"github.com/klassmann/cpfcnpj"
	"net/http"
	"time"
)

type UserHandler struct {
	userDb          database.UserInterface
	staticDepositDb database.StaticDepositInterface
}

func NewUserHandler(
	userDb database.UserInterface,
	staticDepositDb database.StaticDepositInterface) *UserHandler {
	return &UserHandler{
		userDb:          userDb,
		staticDepositDb: staticDepositDb,
	}
}

func (h *UserHandler) Create(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var input dtos.CreateUserInput
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

	if err := validateUserCreationInput(&input); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(struct {
			Error string `json:"error"`
		}{
			Error: err.Error(),
		})
		return
	}

	usecase := userUsecases.NewCreateUserUsecase(h.userDb)
	err, code := usecase.CreateUser(&input)
	if err != nil {
		w.WriteHeader(code)
		_ = json.NewEncoder(w).Encode(struct {
			Error string `json:"error"`
		}{
			Error: err.Error(),
		})
		return
	}
	w.WriteHeader(http.StatusCreated)
}

func validateUserCreationInput(input *dtos.CreateUserInput) error {
	input.TaxId = utils.TrimCpfCnpj(input.TaxId)
	if input.Password != input.ConfirmPassword {
		return errors.New(utils.PasswordsDontMatch)
	}
	if input.Email == "" {
		return errors.New(utils.MissingEmail)
	}
	if !cpfcnpj.ValidateCPF(input.TaxId) {
		return errors.New(utils.InvalidTaxId)
	}
	return nil
}

func (h *UserHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	subject := r.Context().Value("subject").(string)
	usecase := userUsecases.NewGetUserUsecase(subject, h.userDb)
	output, err, code := usecase.GetUser()
	if err != nil {
		w.WriteHeader(code)
		_ = json.NewEncoder(w).Encode(struct {
			Error string `json:"error"`
		}{
			Error: err.Error(),
		})
		return
	}

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(output)
}

func (h *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var input dtos.GetJwtInput
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

	if err := validateLoginInput(&input); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(struct {
			Error string `json:"error"`
		}{
			Error: err.Error(),
		})
		return
	}

	usecase := userUsecases.NewUserLoginUsecase(h.userDb)
	output, err, code := usecase.GetToken(&input)
	if err != nil {
		w.WriteHeader(code)
		_ = json.NewEncoder(w).Encode(struct {
			Error string `json:"error"`
		}{
			Error: err.Error(),
		})
		return
	}
	http.SetCookie(w, &http.Cookie{
		Name:     "jwt",
		Value:    output.AccessToken,
		Expires:  time.Now().Add(time.Second * time.Duration(configs.Cfg.JwtExpiration)),
		MaxAge:   configs.Cfg.JwtExpiration,
		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteNoneMode,
		Domain:   "brla.digital",
		Path:     "/",
	})
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(output)
}

func validateLoginInput(input *dtos.GetJwtInput) error {
	if input.Email == "" {
		return errors.New(utils.MissingEmail)
	}
	if input.Password == "" {
		return errors.New(utils.MissingPassword)
	}
	return nil
}
