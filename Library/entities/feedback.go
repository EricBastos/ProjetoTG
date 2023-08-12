package entities

import (
	"github.com/EricBastos/ProjetoTG/Library/pkg/entities"
	"time"
)

type Feedback struct {
	ID                       *entities.ID `json:"id"`
	FeedbackType             string       `json:"feedbackType"`
	OperationId              string       `json:"operationId"`
	SmartcontractOperationId string       `json:"smartcontractOperationId"`
	Success                  bool         `json:"success"`
	ErrorMsg                 string       `json:"errorMsg"`
	CreatedAt                time.Time    `json:"createdAt"`
}

func NewFeedback(
	operationId string,
	feedbackType string,
	smartContractCallId string,
	success bool,
	errorMsg string) *Feedback {
	feedbackId := entities.NewID()
	return &Feedback{
		OperationId:              operationId,
		FeedbackType:             feedbackType,
		ID:                       &feedbackId,
		Success:                  success,
		SmartcontractOperationId: smartContractCallId,
		ErrorMsg:                 errorMsg,
	}
}
