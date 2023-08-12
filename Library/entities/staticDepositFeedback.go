package entities

import (
	"time"
)

type StaticDepositFeedback struct {
	ID               string
	DepositId        string
	DepositStatus    string
	DepositUpdatedAt *time.Time
	LogType          string
	LogCreatedAt     *time.Time
	Amount           int
	TaxId            string
	StaticDepositId  string
	WorkspaceId      string
	Fee              int
}

func NewStaticDepositFeedback(
	id string,
	depositId string,
	depositStatus string,
	depositUpdatedAt *time.Time,
	logType string,
	logCreatedAt *time.Time,
	amount int,
	taxId string,
	staticDepositId string,
	workspaceId string,
	fee int,
) *StaticDepositFeedback {
	return &StaticDepositFeedback{
		ID:               id,
		DepositId:        depositId,
		DepositStatus:    depositStatus,
		DepositUpdatedAt: depositUpdatedAt,
		LogType:          logType,
		LogCreatedAt:     logCreatedAt,
		Amount:           amount,
		TaxId:            taxId,
		StaticDepositId:  staticDepositId,
		WorkspaceId:      workspaceId,
		Fee:              fee,
	}
}
