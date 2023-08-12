package database

import (
	"github.com/EricBastos/ProjetoTG/Library/entities"
)

type UserInterface interface {
	Create(user *entities.User) error
	FindByEmail(email string) (*entities.User, error)
	FindById(id string) (*entities.User, error)
}

type SmartcontractOperationInterface interface {
	Create(op *entities.SmartcontractOperation) error
}

type FeedbackInterface interface {
	Create(feedback *entities.Feedback) error
	FindById(id string) (*entities.Feedback, error)
}

type MintOpInterface interface {
	Create(op *entities.MintOp) error
	CreateEmit(op *entities.MintOp, f func() error) error
	Get(id string) (*entities.MintOp, error)
}

type StaticDepositFeedbackInterface interface {
	Create(feedback *entities.StaticDepositFeedback) error
	FindById(id string) (*entities.StaticDepositFeedback, error)
}

type StaticDepositInterface interface {
	Create(invoice *entities.StaticDeposit) error
	Update(updatedDeposit *entities.StaticDeposit) error
	FindById(id string) (*entities.StaticDeposit, error)
	FindUnpaidByTaxIdAndAmount(taxId string, amount int) (*entities.StaticDeposit, error)
	FindUnpaidByTaxId(taxId string) (*entities.StaticDeposit, error)
}

type BurnOpInterface interface {
	Create(op *entities.BurnOp) error
	CreateEmit(op *entities.BurnOp, f func() error) error
	GetLogs(docData, responsibleUser string, page, pageSize int) ([]entities.BurnOpAPI, error)
	GetBusinessLogs(responsibleUser string, page, pageSize int) ([]entities.BurnOpAPI, error)
	Get(id string) (*entities.BurnOp, error)
}

type TransferInterface interface {
	Create(transfer *entities.Transfer) error
	FindById(id string) (*entities.Transfer, error)
}

type TransferFeedbackInterface interface {
	Create(feedback *entities.TransferFeedback) error
	FindById(id string) (*entities.TransferFeedback, error)
}

type BurnFeedbackInterface interface {
	Create(feedback *entities.BurnFeedback) error
	FindById(id string) (*entities.BurnFeedback, error)
}
