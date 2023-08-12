package entities

import (
	"github.com/EricBastos/ProjetoTG/Library/pkg/entities"
	"time"
)

type StaticDepositFeedback struct {
	ID        entities.ID
	TaxId     string
	Amount    int
	DepositId string
	CreatedAt time.Time
}

func NewStaticDepositFeedback(
	depositId string,
	taxId string,
	amount int,
) *StaticDepositFeedback {
	newId := entities.NewID()
	return &StaticDepositFeedback{
		ID:        newId,
		TaxId:     taxId,
		Amount:    amount,
		DepositId: depositId,
	}
}
