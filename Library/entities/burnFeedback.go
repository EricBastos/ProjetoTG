package entities

import (
	"github.com/EricBastos/ProjetoTG/Library/pkg/entities"
	"time"
)

type BurnFeedback struct {
	ID                  *entities.ID `json:"id"`
	ResponsibleUser     *entities.ID `json:"responsibleUser"`
	Chain               string       `json:"chain"`
	WalletAddress       string       `json:"walletAddress"`
	Amount              int          `json:"amount"`
	Success             bool         `json:"success"`
	OperationId         string       `json:"operationId"`
	SmartcontractCallId string       `json:"smartcontract_call_id"`
	ErrorMsg            string       `json:"errorMsg"`
	IsRetry             bool         `json:"isRetry"`
	UserName            string       `json:"userName"`
	UserTaxId           string       `json:"userTaxId"`
	AccBankCode         string       `json:"accBankCode"`
	AccBranchCode       string       `json:"accBranchCode"`
	AccNumber           string       `json:"accNumber"`
	CreatedAt           time.Time    `json:"createdAt"`
}

func NewBurnFeedback(
	responsibleUser *entities.ID,
	walletAddress string,
	amount int,
	success bool,
	operationId string,
	smartcontractCallId string,
	errorMsg string,
	isRetry bool,
	userName,
	userTaxId,
	accBankCode,
	accBranchCode,
	accNumber,
	chain string) *BurnFeedback {
	feedbackId := entities.NewID()
	return &BurnFeedback{
		ResponsibleUser:     responsibleUser,
		ID:                  &feedbackId,
		WalletAddress:       walletAddress,
		Amount:              amount,
		Success:             success,
		IsRetry:             isRetry,
		OperationId:         operationId,
		SmartcontractCallId: smartcontractCallId,
		ErrorMsg:            errorMsg,
		UserName:            userName,
		UserTaxId:           userTaxId,
		AccBankCode:         accBankCode,
		AccBranchCode:       accBranchCode,
		AccNumber:           accNumber,
		Chain:               chain,
	}
}
