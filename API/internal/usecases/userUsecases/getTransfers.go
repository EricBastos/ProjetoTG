package userUsecases

import (
	"github.com/EricBastos/ProjetoTG/API/internal/dtos"
	"github.com/EricBastos/ProjetoTG/Library/database"
	"github.com/EricBastos/ProjetoTG/Library/entities"
)

type GetTransfersLogsUsecase struct {
	taxId     string
	userId    string
	burnOpsDb database.BurnOpInterface
}

func NewGetTransfersLogsUsecase(
	taxId string,
	userId string,
	burnOpsDb database.BurnOpInterface) *GetTransfersLogsUsecase {
	return &GetTransfersLogsUsecase{
		taxId:     taxId,
		userId:    userId,
		burnOpsDb: burnOpsDb,
	}
}

func (u *GetTransfersLogsUsecase) GetTransfersLogs(input *dtos.GetTransfersLogsInput) (*dtos.GetTransfersLogsOutput, error, int) {

	TransfersLogs, err := u.burnOpsDb.GetLogs(u.taxId, u.userId, input.Page, input.PageSize)
	if err != nil {
		return &dtos.GetTransfersLogsOutput{TransfersLogs: []entities.BurnOpAPI{}}, nil, 0
	}
	return &dtos.GetTransfersLogsOutput{TransfersLogs: TransfersLogs}, nil, 0

}
