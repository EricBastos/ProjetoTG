package userUsecases

import (
	"errors"
	"github.com/EricBastos/ProjetoTG/API/internal/dtos"
	"github.com/EricBastos/ProjetoTG/Library/database"
	"github.com/EricBastos/ProjetoTG/Library/utils"
	"net/http"
)

type GetUserUsecase struct {
	userId string
	userDb database.UserInterface
}

func NewGetUserUsecase(userId string, userDb database.UserInterface) *GetUserUsecase {
	return &GetUserUsecase{
		userId: userId,
		userDb: userDb,
	}
}

func (u *GetUserUsecase) GetUser() (*dtos.GetUserOutput, error, int) {

	user, err := u.userDb.FindById(u.userId)
	if err != nil {
		return nil, errors.New(utils.UserNotFoundError), http.StatusInternalServerError
	}

	return &dtos.GetUserOutput{
		Name:      user.Name,
		Email:     user.Email,
		TaxId:     user.TaxId,
		CreatedAt: user.CreatedAt,
	}, nil, 0

}
