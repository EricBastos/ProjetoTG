package entities

import (
	entities2 "github.com/EricBastos/ProjetoTG/Library/pkg/entities"
	"time"
)

type TransferFeedback struct {
	ID         *entities2.ID
	TransferId string
	CreatedAt  time.Time
}

func NewTransferFeedback(
	transferId string,
) *TransferFeedback {
	id := entities2.NewID()
	return &TransferFeedback{
		ID:         &id,
		TransferId: transferId,
	}
}
