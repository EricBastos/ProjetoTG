package entities

import (
	entities2 "github.com/EricBastos/ProjetoTG/Library/pkg/entities"
	"github.com/EricBastos/ProjetoTG/Library/utils"
	"time"
)

type Transfer struct {
	Chain            string             `json:"chain"`
	WalletAddress    string             `json:"walletAddress"`
	Amount           int                `json:"amount"`
	Name             string             `json:"name"`
	TaxId            string             `json:"taxId"`
	BankCode         string             `json:"bankCode"`
	BranchCode       string             `json:"branchCode"`
	AccountNumber    string             `json:"accountNumber"`
	Id               string             `json:"id"`
	CreatedAt        time.Time          `json:"createdAt"`
	ResponsibleUser  *entities2.ID      `json:"responsibleUser"`
	AssociatedBurnId string             `json:"associatedBurnId"`
	Feedbacks        []TransferFeedback `json:"feedbacks" gorm:"foreignKey:TransferId"`
}

func NewTransfer(
	walletAddress string,
	amount int,
	name string,
	taxId string,
	bankCode string,
	branchCode string,
	accountNumber string,
	responsibleUser *entities2.ID, chain, associatedBurnId string) *Transfer {
	newId := entities2.NewID()
	return &Transfer{
		Chain:            chain,
		Id:               newId.String(),
		WalletAddress:    walletAddress,
		Amount:           amount,
		Name:             name,
		TaxId:            utils.TrimCpfCnpj(taxId),
		BankCode:         bankCode,
		BranchCode:       branchCode,
		AccountNumber:    accountNumber,
		ResponsibleUser:  responsibleUser,
		AssociatedBurnId: associatedBurnId,
	}
}
