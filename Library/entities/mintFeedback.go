package entities

import (
	"github.com/EricBastos/ProjetoTG/Library/pkg/entities"
	"time"
)

type MintFeedback struct {
	ResponsibleUser     *entities.ID `json:"responsibleUser"`
	ID                  *entities.ID `json:"id"`
	Chain               string       `json:"chain"`
	WalletAddress       string       `json:"walletAddress"`
	Amount              int          `json:"amount"`
	Success             bool         `json:"success"`
	OperationId         string       `json:"operationId"`
	SmartcontractCallId string       `json:"smartcontract_call_id"`
	ErrorMsg            string       `json:"errorMsg"`
	IsRetry             bool         `json:"isRetry"`
	CreatedAt           time.Time    `json:"createdAt"`
}

func NewMintFeedback(
	responsibleUser *entities.ID,
	walletAddress string,
	amount int,
	success bool,
	operationId string,
	smartcontractCallId string,
	errorMsg string,
	isRetry bool,
	chain string) *MintFeedback {
	feedbackId := entities.NewID()
	return &MintFeedback{
		ResponsibleUser:     responsibleUser,
		ID:                  &feedbackId,
		WalletAddress:       walletAddress,
		Amount:              amount,
		Success:             success,
		OperationId:         operationId,
		SmartcontractCallId: smartcontractCallId,
		ErrorMsg:            errorMsg,
		IsRetry:             isRetry,
		Chain:               chain,
	}
}
