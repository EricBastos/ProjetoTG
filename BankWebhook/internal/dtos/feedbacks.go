package dtos

//type TransferFeedbackInput struct {
//	CreatedAt *time.Time           `json:"created,omitempty"`
//	Errors    []string             `json:"errors,omitempty"`
//	LogId     string               `json:"id,omitempty"`
//	Transfer  *SBTransfer.Transfer `json:"transfer,omitempty"`
//	Type      string               `json:"type,omitempty"`
//}

type DepositFeedbackInput struct {
	DepositId string `json:"depositId"`
	TaxId     string `json:"taxId"`
	Amount    int    `json:"amount"`
}
