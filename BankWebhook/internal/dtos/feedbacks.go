package dtos

type TransferFeedbackInput struct {
	TransferId string `json:"transferId"`
}

type DepositFeedbackInput struct {
	DepositId string `json:"depositId"`
	TaxId     string `json:"taxId"`
	Amount    int    `json:"amount"`
}
