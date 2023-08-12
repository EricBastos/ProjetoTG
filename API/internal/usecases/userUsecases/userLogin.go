package userUsecases

import (
	"errors"
	"github.com/EricBastos/ProjetoTG/API/configs"
	"github.com/EricBastos/ProjetoTG/API/internal/dtos"
	"github.com/EricBastos/ProjetoTG/Library/database"
	"github.com/EricBastos/ProjetoTG/Library/utils"
	"net/http"
	"time"
)

type UserLoginUsecase struct {
	userDb database.UserInterface
}

func NewUserLoginUsecase(userDB database.UserInterface) *UserLoginUsecase {
	return &UserLoginUsecase{
		userDb: userDB,
	}
}

func (u *UserLoginUsecase) GetToken(input *dtos.GetJwtInput) (*dtos.GetJwtOutput, error, int) {

	user, err := u.userDb.FindByEmail(input.Email)
	if err != nil {
		return nil, errors.New(utils.InvalidCredentials), http.StatusUnauthorized
	}
	if !user.ValidatePassword(input.Password) {
		return nil, errors.New(utils.InvalidCredentials), http.StatusUnauthorized
	}

	_, tokenString, err := configs.Cfg.TokenAuthUser.Encode(map[string]interface{}{
		"sub":   user.ID.String(),
		"taxId": user.TaxId,
		"email": user.Email,
		"iat":   time.Now().Unix(),
		"exp":   time.Now().Add(time.Second * time.Duration(configs.Cfg.JwtExpiration)).Unix(),
	})
	if err != nil {
		return nil, err, http.StatusInternalServerError
	}

	return &dtos.GetJwtOutput{AccessToken: tokenString}, nil, 0

}
