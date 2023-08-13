package dtos

import (
	"github.com/EricBastos/ProjetoTG/Library/entities"
)

type GetTransfersLogsInput struct {
	Page     int `json:"page"`
	PageSize int `json:"pageSize"`
}

type GetTransfersLogsOutput struct {
	TransfersLogs []entities.BurnOpAPI `json:"transfersLogs"`
}

type GetDepositsLogsInput struct {
	Page     int `json:"page"`
	PageSize int `json:"pageSize"`
}

type GetDepositsLogsOutput struct {
	DepositsLogs []entities.StaticDepositAPI `json:"depositsLogs"`
}
