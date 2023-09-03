package dtos

import (
	"github.com/EricBastos/ProjetoTG/Library/entities"
)

type BridgeAssetInput struct {
	WalletAddress string `json:"walletAddress"`
	InputChain    string `json:"inputChain"`
	OutputChain   string `json:"outputChain"`
	Amount        int    `json:"amount"`

	// Permit Info (EVM-compatible)
	Permit *entities.PermitData `json:"permit"`
}

type BridgeAssetOutput struct {
	Id string `json:"id"`
}

type GetBridgeLogsInput struct {
	Page     int `json:"page"`
	PageSize int `json:"pageSize"`
}

type GetBridgeLogsOutput struct {
	BridgeLogs []entities.BridgeOpAPI `json:"bridgeLogs"`
}
