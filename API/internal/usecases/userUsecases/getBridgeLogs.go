package userUsecases

import (
	"github.com/EricBastos/ProjetoTG/API/internal/dtos"
	"github.com/EricBastos/ProjetoTG/Library/database"
	"github.com/EricBastos/ProjetoTG/Library/entities"
)

type GetBridgeLogsUsecase struct {
	userId      string
	bridgeOpsDb database.BridgeOpInterface
}

func NewGetBridgeLogsUsecase(
	userId string,
	bridgeOpsDb database.BridgeOpInterface) *GetBridgeLogsUsecase {
	return &GetBridgeLogsUsecase{
		userId:      userId,
		bridgeOpsDb: bridgeOpsDb,
	}
}

func (u *GetBridgeLogsUsecase) GetTransfersLogs(input *dtos.GetBridgeLogsInput) (*dtos.GetBridgeLogsOutput, error, int) {

	TransfersLogs, err := u.bridgeOpsDb.GetLogs(u.userId, input.Page, input.PageSize)
	if err != nil {
		return &dtos.GetBridgeLogsOutput{BridgeLogs: []entities.BridgeOpAPI{}}, nil, 0
	}
	return &dtos.GetBridgeLogsOutput{BridgeLogs: TransfersLogs}, nil, 0

}
