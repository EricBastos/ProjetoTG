package entities

import (
	entities2 "github.com/EricBastos/ProjetoTG/Library/pkg/entities"
	"github.com/EricBastos/ProjetoTG/Library/utils"
	SBTransfer "github.com/starkbank/sdk-go/starkbank/transfer"
	"time"
)

type Transfer struct {
	WorkspaceId      string             `json:"workspaceId"`
	Chain            string             `json:"chain"`
	WalletAddress    string             `json:"walletAddress"`
	Amount           int                `json:"amount"`
	Name             string             `json:"name"`
	TaxId            string             `json:"taxId"`
	BankCode         string             `json:"bankCode"`
	BranchCode       string             `json:"branchCode"`
	AccountNumber    string             `json:"accountNumber"`
	AccountType      string             `json:"accountType"`
	ExternalId       string             `json:"externalId"`
	Id               string             `json:"id"`
	CreatedAt        *time.Time         `json:"createdAt"`
	Scheduled        *time.Time         `json:"scheduled"`
	ResponsibleUser  *entities2.ID      `json:"responsibleUser"`
	AssociatedBurnId string             `json:"associatedBurnId"`
	Feedbacks        []TransferFeedback `json:"feedbacks" gorm:"foreignKey:TransferId"`
	NotifyEmail      bool               `json:"notifyEmail"`
	Fee              int                `json:"fee"`
}

func NewTransfer(walletAddress string, transfer SBTransfer.Transfer, responsibleUser *entities2.ID, chain, workspaceId, associatedBurnId string) *Transfer {
	return &Transfer{
		Chain:            chain,
		Id:               transfer.Id,
		CreatedAt:        transfer.Created,
		WalletAddress:    walletAddress,
		Amount:           transfer.Amount,
		Name:             transfer.Name,
		TaxId:            utils.TrimCpfCnpj(transfer.TaxId),
		BankCode:         transfer.BankCode,
		BranchCode:       transfer.BranchCode,
		AccountNumber:    transfer.AccountNumber,
		AccountType:      transfer.AccountType,
		ExternalId:       transfer.ExternalId,
		Scheduled:        transfer.Scheduled,
		ResponsibleUser:  responsibleUser,
		WorkspaceId:      workspaceId,
		AssociatedBurnId: associatedBurnId,
	}
}
