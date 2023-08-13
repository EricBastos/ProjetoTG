package dtos

import (
	"github.com/EricBastos/ProjetoTG/Library/entities"
)

type CreateUserWithdrawInput struct {
	PixKey        string `json:"pixKey"`
	WalletAddress string `json:"walletAddress"`
	Chain         string `json:"chain"`
	Amount        int    `json:"amount"`

	// Permit Info (EVM-compatible)
	Permit *entities.PermitData `json:"permit"`
}

type CreateWithdrawOutput struct {
	Id string `json:"id"`
}
