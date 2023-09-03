package userUsecases

import (
	"errors"
	"github.com/EricBastos/ProjetoTG/API/internal/dtos"
	"github.com/EricBastos/ProjetoTG/Library/database"
	"github.com/EricBastos/ProjetoTG/Library/entities"
	"github.com/EricBastos/ProjetoTG/Library/utils"
	"net/http"
)

type CreateUserUsecase struct {
	userDb database.UserInterface
}

func NewCreateUserUsecase(userDB database.UserInterface) *CreateUserUsecase {
	return &CreateUserUsecase{
		userDb: userDB,
	}
}

func (u *CreateUserUsecase) CreateUser(input *dtos.CreateUserInput) (error, int) {

	var err error
	var user *entities.User

	user, err = entities.NewUser(
		input.Name,
		input.Email,
		input.Password,
		input.TaxId)

	if err != nil {
		return errors.New(utils.InternalError), http.StatusInternalServerError
	}

	err = u.userDb.Create(user)
	if err != nil {
		if err.Error() == utils.UserAlreadyExists {
			return err, http.StatusConflict
		}
		return errors.New(utils.InternalError), http.StatusInternalServerError
	}

	return nil, 0

}
