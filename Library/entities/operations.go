package entities

import (
	"encoding/json"
	"github.com/EricBastos/ProjetoTG/Library/pkg/entities"
	"time"
)

type OperationOriginType string

const (
	MINT OperationOriginType = "MINT"
	BURN OperationOriginType = "BURN"
)

type SmartContractOp interface {
	GetDataInJson() string
	GetOperationType() string
	GetResponsibleUser() *entities.ID
	GetID() *entities.ID
	GetChain() string
}

type SmartcontractOperation struct {
	ID            *entities.ID
	OperationName string
	OperationID   *entities.ID
	OperationType string
	Executed      bool
	Tx            string
	Reason        string
	CreatedAt     time.Time `json:"createdAt"`
	IsRetry       bool      `json:"isRetry"`

	Feedback *Feedback `json:"feedback"`
}

type PermitData struct {
	Deadline int64  `json:"deadline"`
	Nonce    int64  `json:"nonce"`
	R        string `json:"r"` //[32]byte
	S        string `json:"s"` //[32]byte
	V        uint8  `json:"v"`
}

type BurnOp struct {
	Id              *entities.ID `json:"id"`
	ResponsibleUser *entities.ID `json:"userId"`
	Chain           string       `json:"chain"`
	WalletAddress   string       `json:"address"`

	Amount        int    `json:"amount"`
	UserName      string `json:"userName"`
	UserTaxId     string `json:"userTaxId"`
	AccBankCode   string `json:"accBankCode"`
	AccBranchCode string `json:"accBranchCode"`
	AccNumber     string `json:"accNumber"`

	Permit *PermitData `json:"permit" gorm:"embedded"`

	SmartContractOps []SmartcontractOperation `json:"smartContractOps" gorm:"polymorphic:Operation"`

	Transfers []Transfer `json:"transfers" gorm:"foreignKey:AssociatedBurnId"` // MUST FETCH THE MOST RECENT ONE

	CreatedAt time.Time `json:"createdAt"`
}

type MintOp struct {
	Id                            *entities.ID `json:"id"`
	ResponsibleUser               *entities.ID `json:"userId"`
	Chain                         string       `json:"chain"`
	WalletAddress                 string       `json:"address"`
	Amount                        int          `json:"amount"`
	Reason                        string       `json:"reason"`
	CreatedAt                     time.Time    `json:"createdAt"`
	AssociatedBankTransactionID   string       `json:"associatedBankTransactionId"`
	AssociatedBankTransactionType string

	SmartContractOps []SmartcontractOperation `json:"smartContractOps" gorm:"polymorphic:Operation"`
}

func NewSmartcontractOperation(op, opOrigin, opId string, executed bool, tx, reason string, isRetry bool) *SmartcontractOperation {
	parsedOpId, _ := entities.ParseID(opId)
	id := entities.NewID()
	opType := ""
	switch opOrigin {
	case "MINT":
		opType = "mint_ops"
	case "BURN":
		opType = "burn_ops"
	}
	return &SmartcontractOperation{
		ID:            &id,
		OperationName: op,
		OperationID:   &parsedOpId,
		OperationType: opType,
		Executed:      executed,
		Tx:            tx,
		Reason:        reason,
		IsRetry:       isRetry,
	}
}

func NewBurnWithPermit(walletAddress string, amount int, userName, userTaxId, accBankCode, accBranchCode, accNumber, chain string, responsible *entities.ID, permit *PermitData) *BurnOp {
	opId := entities.NewID()
	return &BurnOp{
		Id:              &opId,
		ResponsibleUser: responsible,
		Chain:           chain,
		WalletAddress:   walletAddress,
		Amount:          amount,
		UserName:        userName,
		UserTaxId:       userTaxId,
		AccBankCode:     accBankCode,
		AccBranchCode:   accBranchCode,
		AccNumber:       accNumber,
		Permit:          permit,
	}
}

func NewBurnFromJson(jsonData []byte) (*BurnOp, error) {
	var burnOp BurnOp
	err := json.Unmarshal(jsonData, &burnOp)
	if err != nil {
		return nil, err
	}
	return &burnOp, nil
}

func (op *BurnOp) GetDataInJson() string {
	data, _ := json.Marshal(op)
	return string(data)
}

func (op *BurnOp) GetOperationType() string {
	return "BURN"
}

func (op *BurnOp) GetResponsibleUser() *entities.ID {
	return op.ResponsibleUser
}

func (op *BurnOp) GetID() *entities.ID {
	return op.Id
}

func (op *BurnOp) GetChain() string {
	return op.Chain
}

func NewMint(walletAddress string, amount int, chain, reason string, responsible *entities.ID, associatedBankTransactionId string) *MintOp {
	opId := entities.NewID()
	return &MintOp{
		Id:                          &opId,
		Chain:                       chain,
		WalletAddress:               walletAddress,
		Amount:                      amount,
		Reason:                      reason,
		AssociatedBankTransactionID: associatedBankTransactionId,
		ResponsibleUser:             responsible,
	}
}
func NewMintFromJson(jsonData []byte) (*MintOp, error) {
	var mintOp MintOp
	err := json.Unmarshal(jsonData, &mintOp)
	if err != nil {
		return nil, err
	}
	return &mintOp, nil
}

func (op *MintOp) GetDataInJson() string {
	data, _ := json.Marshal(op)
	return string(data)
}

func (op *MintOp) GetOperationType() string {
	return "MINT"
}

func (op *MintOp) GetResponsibleUser() *entities.ID {
	return op.ResponsibleUser
}

func (op *MintOp) GetID() *entities.ID {
	return op.Id
}

func (op *MintOp) GetChain() string {
	return op.Chain
}

type BurnOpAPI struct {
	Id               *entities.ID                `json:"id"`
	Chain            string                      `json:"chain"`
	WalletAddress    string                      `json:"walletAddress"`
	Amount           int                         `json:"amount"`
	Permit           *PermitData                 `json:"permit" gorm:"embedded"`
	SmartContractOps []SmartcontractOperationAPI `json:"smartContractOps" gorm:"polymorphic:Operation"`
	Transfers        []TransferAPI               `json:"transfers" gorm:"foreignKey:AssociatedBurnId"`
	CreatedAt        time.Time                   `json:"createdAt"`
	Fee              int                         `json:"fee"`
}

type SmartcontractOperationAPI struct {
	ID            *entities.ID `json:"id"`
	OperationName string       `json:"operationName"`
	Executed      bool         `json:"posted"`
	Tx            string       `json:"tx"`
	Reason        string       `json:"notPostedReason"`
	CreatedAt     time.Time    `json:"createdAt"`
	IsRetry       bool         `json:"isRetry"`

	Feedback *FeedbackAPI `json:"feedback" gorm:"foreignKey:SmartcontractOperationId"`
}

type FeedbackAPI struct {
	ID                       *entities.ID `json:"id"`
	SmartcontractOperationId string       `json:"-"`
	Success                  bool         `json:"success"`
	ErrorMsg                 string       `json:"errorMsg"`
	CreatedAt                time.Time    `json:"createdAt"`
}

type TransferAPI struct {
	Amount        int                   `json:"amount"`
	Name          string                `json:"name"`
	TaxId         string                `json:"taxId"`
	BankCode      string                `json:"bankCode"`
	BranchCode    string                `json:"branchCode"`
	AccountNumber string                `json:"accountNumber"`
	Id            string                `json:"id"`
	CreatedAt     *time.Time            `json:"createdAt"`
	Feedbacks     []TransferFeedbackAPI `json:"feedbacks" gorm:"foreignKey:TransferId"`
}

type TransferFeedbackAPI struct {
	ID                *entities.ID `json:"id"`
	TransferStatus    string       `json:"transferStatus"`
	TransferUpdatedAt *time.Time   `json:"updatedAt"`
	LogType           string       `json:"logType"`
}

type StaticDepositAPI struct {
	Chain         string       `json:"chain"`
	WalletAddress string       `json:"walletAddress"`
	Amount        int          `json:"amount"`
	TaxId         string       `json:"taxId"`
	Due           *time.Time   `json:"due"`
	Id            *entities.ID `json:"id"`
	CreatedAt     *time.Time   `json:"createdAt"`
	Status        string       `json:"status"`
	PayerName     string       `json:"payerName"`
	UpdatedAt     time.Time    `json:"updatedAt"`

	MintOps []MintOpAPI `json:"mintOps"`
}

type MintOpAPI struct {
	Id               *entities.ID                `json:"id"`
	Amount           int                         `json:"amount"`
	Reason           string                      `json:"createdReason"`
	CreatedAt        time.Time                   `json:"createdAt"`
	Fee              int                         `json:"fee"`
	SmartContractOps []SmartcontractOperationAPI `json:"smartContractOps"`
}
