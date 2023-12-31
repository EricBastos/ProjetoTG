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
	GetLogs(taxId, responsibleUser string, page, pageSize int) ([]entities.StaticDepositAPI, error)
}

type BurnOpInterface interface {
	Create(op *entities.BurnOp) error
	CreateEmit(op *entities.BurnOp, f func() error) error
	Get(id string) (*entities.BurnOp, error)
	GetLogs(docData, responsibleUser string, page, pageSize int) ([]entities.BurnOpAPI, error)
}

type BridgeOpInterface interface {
	Create(op *entities.BridgeOp) error
	CreateEmit(op *entities.BridgeOp, f func() error) error
	Get(id string) (*entities.BridgeOp, error)
	GetLogs(responsibleUser string, page, pageSize int) ([]entities.BridgeOpAPI, error)
}

type TransferInterface interface {
	Create(transfer *entities.Transfer) error
	FindById(id string) (*entities.Transfer, error)
}

type TransferFeedbackInterface interface {
	Create(feedback *entities.TransferFeedback) error
	FindById(id string) (*entities.TransferFeedback, error)
}
