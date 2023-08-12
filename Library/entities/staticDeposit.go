package entities

import (
	entities2 "github.com/EricBastos/ProjetoTG/Library/pkg/entities"
	"time"
)

type StaticDeposit struct {
	Chain           string        `json:"chain"`
	WalletAddress   string        `json:"walletAddress"`
	Amount          int           `json:"amount"`
	ResponsibleUser *entities2.ID `json:"responsibleUser"`
	TaxId           string        `json:"taxId"`
	Due             *time.Time    `json:"due"`
	Id              *entities2.ID `json:"id"`
	CreatedAt       *time.Time    `json:"createdAt"`
	Status          string        `json:"status"`
	UpdatedAt       time.Time     `json:"updatedAt"`

	MintOps []MintOp `json:"mintOps" gorm:"polymorphic:AssociatedBankTransaction"`
}

func NewStaticDeposit(
	responsibleUser *entities2.ID,
	walletAddress string,
	amount int,
	taxId string,
	due *time.Time,
	createdAt *time.Time,
	chain string) *StaticDeposit {
	id := entities2.NewID()
	return &StaticDeposit{
		Chain:           chain,
		ResponsibleUser: responsibleUser,
		Id:              &id,
		WalletAddress:   walletAddress,
		Amount:          amount,
		TaxId:           taxId,
		Due:             due,
		CreatedAt:       createdAt,
		Status:          "UNPAID",
	}
}
