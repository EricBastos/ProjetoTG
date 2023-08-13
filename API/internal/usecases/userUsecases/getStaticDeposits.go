package userUsecases

import (
	"github.com/EricBastos/ProjetoTG/API/internal/dtos"
	"github.com/EricBastos/ProjetoTG/Library/database"
	"github.com/EricBastos/ProjetoTG/Library/entities"
)

type GetStaticDepositsUsecase struct {
	taxId           string
	userId          string
	staticDepositDb database.StaticDepositInterface
}

func NewGetStaticDepositsUsecase(
	taxId string,
	userId string,
	staticDepositDb database.StaticDepositInterface) *GetStaticDepositsUsecase {
	return &GetStaticDepositsUsecase{
		taxId:           taxId,
		userId:          userId,
		staticDepositDb: staticDepositDb,
	}
}

func (u *GetStaticDepositsUsecase) GetStaticDepositsLogs(input *dtos.GetDepositsLogsInput) (*dtos.GetDepositsLogsOutput, error, int) {
	depositLogs, err := u.staticDepositDb.GetLogs(u.taxId, u.userId, input.Page, input.PageSize)
	if err != nil {
		return &dtos.GetDepositsLogsOutput{DepositsLogs: []entities.StaticDepositAPI{}}, nil, 0
	}
	return &dtos.GetDepositsLogsOutput{DepositsLogs: depositLogs}, nil, 0
}
