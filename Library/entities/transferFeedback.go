package entities

import (
	entities2 "github.com/EricBastos/ProjetoTG/Library/pkg/entities"
	"time"
)

type TransferFeedback struct {
	ID                *entities2.ID
	LogId             string
	TransferId        string
	TransferStatus    string
	TransferUpdatedAt *time.Time
	LogType           string
	LogCreatedAt      *time.Time
	WorkspaceId       string
	TaxId             string
	Fee               int
}

func NewTransferFeedback(
	logId string,
	transferId string,
	transferStatus string,
	transferUpdatedAt *time.Time,
	logType string,
	logCreatedAt *time.Time,
	workspaceId string,
	taxId string,
	fee int,
) *TransferFeedback {
	id := entities2.NewID()
	return &TransferFeedback{
		ID:                &id,
		LogId:             logId,
		TransferId:        transferId,
		TransferStatus:    transferStatus,
		TransferUpdatedAt: transferUpdatedAt,
		LogType:           logType,
		LogCreatedAt:      logCreatedAt,
		WorkspaceId:       workspaceId,
		TaxId:             taxId,
		Fee:               fee,
	}
}
