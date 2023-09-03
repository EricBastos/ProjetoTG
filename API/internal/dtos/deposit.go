package dtos

import (
	"github.com/EricBastos/ProjetoTG/Library/entities"
)

type CreateUserStaticPixDepositInput struct {
	WalletAddress string `json:"walletAddress"`
	Chain         string `json:"chain"`
	Amount        int    `json:"amount"`
}

type CreateStaticPixDepositOutput struct {
	Id string `json:"id"`
}

type GetDepositsLogsInput struct {
	Page     int `json:"page"`
	PageSize int `json:"pageSize"`
}

type GetDepositsLogsOutput struct {
	DepositsLogs []entities.StaticDepositAPI `json:"depositsLogs"`
}
