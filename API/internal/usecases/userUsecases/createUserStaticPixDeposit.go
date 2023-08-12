package userUsecases

import (
	"errors"
	"github.com/EricBastos/ProjetoTG/API/internal/dtos"
	"github.com/EricBastos/ProjetoTG/API/internal/utils"
	"github.com/EricBastos/ProjetoTG/Library/database"
	"github.com/EricBastos/ProjetoTG/Library/entities"
	entities3 "github.com/EricBastos/ProjetoTG/Library/pkg/entities"
	"net/http"
	"time"
)

type CreateUserStaticPixDepositUsecase struct {
	userInfo        *utils.UserInformation
	staticDepositDb database.StaticDepositInterface
}

func NewCreateUserStaticPixDepositUsecase(
	userInfo *utils.UserInformation,
	staticDepositDb database.StaticDepositInterface) *CreateUserStaticPixDepositUsecase {
	return &CreateUserStaticPixDepositUsecase{
		userInfo:        userInfo,
		staticDepositDb: staticDepositDb,
	}
}

func (u *CreateUserStaticPixDepositUsecase) CreateDeposit(input *dtos.CreateUserStaticPixDepositInput) (string, error, int) {

	id, creationError, errorCode := u.createAsUser(input)

	if creationError != nil {
		return "", creationError, errorCode
	}

	return id, nil, 0
}

func (u *CreateUserStaticPixDepositUsecase) createAsUser(input *dtos.CreateUserStaticPixDepositInput) (string, error, int) {

	validateChainFunc, ok := utils.ValidChains[input.Chain]
	if !ok {
		return "", errors.New(utils.InvalidChain), http.StatusBadRequest
	}
	if err := validateChainFunc(&input.WalletAddress); err != nil {
		return "", err, http.StatusBadRequest
	}

	createdAt := time.Now()
	due := createdAt.Add(5 * time.Hour)

	uId, err := entities3.ParseID(u.userInfo.UserId)
	if err != nil {
		return "", errors.New(utils.InternalError), http.StatusInternalServerError
	}

	staticDeposit := entities.NewStaticDeposit(
		&uId,
		input.WalletAddress,
		input.Amount,
		u.userInfo.TaxId,
		&due,
		&createdAt,
		input.Chain,
	)

	err = u.staticDepositDb.Create(staticDeposit)
	if err != nil {
		return "", errors.New(utils.DepositCreationError), http.StatusInternalServerError
	}
	return staticDeposit.Id.String(), nil, 0

}
