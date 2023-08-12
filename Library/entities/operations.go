package entities

import (
	"encoding/json"
	"github.com/EricBastos/ProjetoTG/Library/pkg/entities"
	"time"
)

type OperationOriginType string

const (
	MINT     OperationOriginType = "MINT"
	BURN     OperationOriginType = "BURN"
	SWAP     OperationOriginType = "SWAP"
	PIXTOUSD OperationOriginType = "PIX-TO-USD"
	TRANSFER OperationOriginType = "TRANSFER-WITH-PERMIT"
)

type SmartContractOp interface {
	GetDataInJson() string
	GetOperationType() string
	GetResponsibleUser() *entities.ID
	GetID() *entities.ID
	GetChain() string
}

type OperationToRetry struct {
	Id *entities.ID `json:"id"`
	// must retry
	// retrying
	// success
	// failed
	Status              string       `json:"status"`
	Operation           string       `json:"operation"`
	OperationOriginType string       `json:"operationOriginType"`
	OperationId         *entities.ID `json:"operationId"`
	//AccountId   *entities.ID `json:"accountId"`
	DataJson string `json:"data"`
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

type UsdcTransferWithAuthData struct {
	FromWallet  string `json:"from"`
	ToWallet    string `json:"to"`
	Value       int    `json:"value"` // Must include decimals already
	ValidAfter  int    `json:"validAfter"`
	ValidBefore int    `json:"validBefore"`
	Nonce       string `json:"nonce"`
	R           string `json:"r"` //[32]byte
	S           string `json:"s"` //[32]byte
	V           uint8  `json:"v"`
}

type UsdtTransferWithAuthData struct {
	FromWallet        string `json:"from"`
	ToWallet          string `json:"to"`
	Value             int    `json:"value"` // Must include decimals already
	Deadline          int64  `json:"deadline"`
	FunctionSignature string `json:"functionSignature"`
	Nonce             int64  `json:"nonce"`
	R                 string `json:"r"` //[32]byte
	S                 string `json:"s"` //[32]byte
	V                 uint8  `json:"v"`
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
	//WorkspaceId     string       `json:"-"`

	Amount        int    `json:"amount"`
	UserName      string `json:"userName"`
	UserTaxId     string `json:"userTaxId"`
	AccBankCode   string `json:"accBankCode"`
	AccBranchCode string `json:"accBranchCode"`
	AccNumber     string `json:"accNumber"`

	Permit       *PermitData `json:"permit" gorm:"embedded"`
	IsWaasWallet bool        `json:"isWaasWallet"`

	SmartContractOps []SmartcontractOperation `json:"smartContractOps" gorm:"polymorphic:Operation"`

	Transfers []Transfer `json:"transfers" gorm:"foreignKey:AssociatedBurnId"` // MUST FETCH THE MOST RECENT ONE

	CreatedAt   time.Time `json:"createdAt"`
	NotifyEmail bool      `json:"notifyEmail"`

	ElbowWallet string `json:"elbowWallet"`
	Fee         int    `json:"fee"`
	MarkupFee   int    `json:"markupFee"`
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
	//WorkspaceId      string                   `json:"-"`
	NotifyEmail bool   `json:"notifyEmail"`
	ElbowWallet string `json:"elbowWallet"`
	Fee         int    `json:"fee"`
}

type SwapOp struct {
	Id              *entities.ID `json:"id"`
	ResponsibleUser *entities.ID `json:"userId"`
	Chain           string       `json:"chain"`
	WalletAddress   string       `json:"address"`
	ReceiverAddress string       `json:"receiverAddress"`

	BrlaAmount    int    `json:"brlaAmount"`
	UsdAmount     int    `json:"usdAmount"`
	UsdToBrla     bool   `json:"usdToBrla"`
	Coin          string `json:"coin"`
	IsWholesale   bool   `json:"isWholesale"`
	WholesaleCode string `json:"wholesaleCode"`
	UserDocument  string `json:"userDocument"`
	BasePrice     string `json:"basePrice"`
	BaseFee       string `json:"baseFee"`
	MarkupFee     string `json:"markupFee"`

	SmartContractOps []SmartcontractOperation `json:"smartContractOps" gorm:"polymorphic:Operation"`

	// After calling hugoX
	OrderId string `json:"orderId"`

	Permit       *PermitData `json:"permit" gorm:"embedded"`
	IsWaasWallet bool        `json:"isWaasWallet"`

	CreatedAt time.Time `json:"createdAt"`

	//WorkspaceId string `json:"-"`
	NotifyEmail bool `json:"notifyEmail"`
}

type PixToUsdOp struct {
	Id              *entities.ID `json:"id"`
	ResponsibleUser *entities.ID `json:"userId"`
	Chain           string       `json:"chain"`
	WalletAddress   string       `json:"address"`
	ReceiverAddress string       `json:"receiverAddress"`

	BrlaAmount                  int                      `json:"brlaAmount"`
	UsdAmount                   int                      `json:"usdAmount"`
	Coin                        string                   `json:"coin"`
	CreatedAt                   time.Time                `json:"createdAt"`
	AssociatedBankTransactionId string                   `json:"associatedBankTransactionId"`
	SmartContractOps            []SmartcontractOperation `json:"smartContractOps" gorm:"polymorphic:Operation"`
	//WorkspaceId                 string                   `json:"-"`
	NotifyEmail bool `json:"notifyEmail"`

	Permit *PermitData `json:"permit" gorm:"-:all"`
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
	case "SWAP":
		opType = "swap_ops"
	case "PIX-TO-USD":
		opType = "pix_to_usd_ops"
	case "PARTIAL-USD-SWAP":
		opType = "swap_ops"
	case "CONVERT-USD":
		opType = "swap_ops"
	case "TRANSFER-WITH-PERMIT":
		opType = "blockchain_transfer_ops"
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

func NewRetryOp(op SmartContractOp, opOriginType OperationOriginType) *OperationToRetry {
	retryId := entities.NewID()
	return &OperationToRetry{
		Id:                  &retryId,
		Status:              "must retry",
		OperationOriginType: string(opOriginType),
		Operation:           op.GetOperationType(),
		OperationId:         op.GetID(),
		DataJson:            op.GetDataInJson(),
	}
}

func NewBurn(walletAddress string, amount int, userName, userTaxId, accBankCode, accBranchCode, accNumber, chain string, responsible *entities.ID) *BurnOp {
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

func NewSwap(walletAddress, receiverAddress string, brlaAmount, usdAmount int, chain string, usdToBrla, isWholesale bool, wholesaleCode, orderId string, responsible *entities.ID, permit *PermitData, userDocument string, basePrice, baseFee, markupFee, coin string) *SwapOp {
	opId := entities.NewID()
	return &SwapOp{
		Id:              &opId,
		ResponsibleUser: responsible,
		Chain:           chain,
		ReceiverAddress: receiverAddress,
		WalletAddress:   walletAddress,
		BrlaAmount:      brlaAmount,
		UsdAmount:       usdAmount,
		UsdToBrla:       usdToBrla,
		IsWholesale:     isWholesale,
		WholesaleCode:   wholesaleCode,
		Permit:          permit,
		OrderId:         orderId,
		UserDocument:    userDocument,
		BasePrice:       basePrice,
		BaseFee:         baseFee,
		MarkupFee:       markupFee,
		Coin:            coin,
	}
}

func NewSwapFromJson(jsonData []byte) (*SwapOp, error) {
	var swapOp SwapOp
	err := json.Unmarshal(jsonData, &swapOp)
	if err != nil {
		return nil, err
	}
	return &swapOp, nil
}

func (op *SwapOp) GetDataInJson() string {
	data, _ := json.Marshal(op)
	return string(data)
}

func (op *SwapOp) GetOperationType() string {
	return "SWAP"
}

func (op *SwapOp) GetResponsibleUser() *entities.ID {
	return op.ResponsibleUser
}

func (op *SwapOp) GetID() *entities.ID {
	return op.Id
}

func (op *SwapOp) GetChain() string {
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

func NewPixToUsdOp(responsible *entities.ID, chain, walletAddress, receiverAddress string, brlaAmount, usdAmount int, coin, associatedBankTransactionId string) *PixToUsdOp {
	opId := entities.NewID()
	return &PixToUsdOp{
		Id:                          &opId,
		ResponsibleUser:             responsible,
		Chain:                       chain,
		WalletAddress:               walletAddress,
		BrlaAmount:                  brlaAmount,
		UsdAmount:                   usdAmount,
		Coin:                        coin,
		ReceiverAddress:             receiverAddress,
		AssociatedBankTransactionId: associatedBankTransactionId,
	}
}
func NewNewPixToUsdOpFromJson(jsonData []byte) (*PixToUsdOp, error) {
	var pixToUsdOp PixToUsdOp
	err := json.Unmarshal(jsonData, &pixToUsdOp)
	if err != nil {
		return nil, err
	}
	return &pixToUsdOp, nil
}

func (op *PixToUsdOp) GetDataInJson() string {
	data, _ := json.Marshal(op)
	return string(data)
}

func (op *PixToUsdOp) GetOperationType() string {
	return "PIX-TO-USD"
}

func (op *PixToUsdOp) GetResponsibleUser() *entities.ID {
	return op.ResponsibleUser
}

func (op *PixToUsdOp) GetID() *entities.ID {
	return op.Id
}

func (op *PixToUsdOp) GetChain() string {
	return op.Chain
}

type PartialUsdSwapOp struct {
	Id              *entities.ID `json:"id"`
	Chain           string       `json:"chain"`
	WalletAddress   string       `json:"walletAddress"`
	UsdAmount       int          `json:"usdAmount"`
	Coin            string       `json:"coin"`
	WorkspaceId     string
	ResponsibleUser *entities.ID
}

func (op *PartialUsdSwapOp) GetDataInJson() string {
	data, _ := json.Marshal(op)
	return string(data)
}

func (op *PartialUsdSwapOp) GetOperationType() string {
	return "PARTIAL-USD-SWAP"
}

func (op *PartialUsdSwapOp) GetResponsibleUser() *entities.ID {
	return op.ResponsibleUser
}

func (op *PartialUsdSwapOp) GetID() *entities.ID {
	return op.Id
}

func (op *PartialUsdSwapOp) GetChain() string {
	return op.Chain
}

func (op *PartialUsdSwapOp) GetWorkspaceId() string {
	return op.WorkspaceId
}

type BlockchainTransferOp struct {
	Id              *entities.ID `json:"id"`
	ResponsibleUser *entities.ID `json:"userId"`
	Chain           string       `json:"chain"`
	FromWallet      string       `json:"from"`
	ToWallet        string       `json:"to"`
	Value           int          `json:"value"` // Must include decimals already
	Coin            string       `json:"coin"`  // Purely informational
	CreatedAt       time.Time    `json:"createdAt"`

	UsdcPermit *UsdcTransferWithAuthData `json:"usdcPermit" gorm:"-:all"`
	UsdtPermit *UsdtTransferWithAuthData `json:"usdtPermit" gorm:"-:all"`

	SmartContractOps []SmartcontractOperation `json:"smartContractOps" gorm:"polymorphic:Operation"`
	//WorkspaceId      string                   `json:"-"`
	NotifyEmail bool `json:"notifyEmail"`
}

func NewBlockchainTransferOp(responsible *entities.ID, chain, fromWallet, toWallet string, value int) *BlockchainTransferOp {
	opId := entities.NewID()
	return &BlockchainTransferOp{
		Id:              &opId,
		ResponsibleUser: responsible,
		Chain:           chain,
		FromWallet:      fromWallet,
		ToWallet:        toWallet,
		Value:           value,
	}
}

func (op *BlockchainTransferOp) GetDataInJson() string {
	data, _ := json.Marshal(op)
	return string(data)
}

func (op *BlockchainTransferOp) GetOperationType() string {
	return "TRANSFER-WITH-PERMIT"
}

func (op *BlockchainTransferOp) GetResponsibleUser() *entities.ID {
	return op.ResponsibleUser
}

func (op *BlockchainTransferOp) GetID() *entities.ID {
	return op.Id
}

func (op *BlockchainTransferOp) GetChain() string {
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

type PixToUsdDepositAPI struct {
	Id              *entities.ID `json:"id"`
	Chain           string       `json:"chain"`
	WalletAddress   string       `json:"walletAddress"`
	ReceiverAddress string       `json:"receiverAddress"`
	Coin            string       `json:"coin"`
	AmountBrl       int          `json:"amountBrl"`
	AmountUsd       int          `json:"amountUsd"`
	TaxId           string       `json:"taxId"`
	Due             *time.Time   `json:"due"`
	CreatedAt       *time.Time   `json:"createdAt"`
	Status          string       `json:"status"`
	Permit          *PermitData  `json:"permit"`
	UpdatedAt       time.Time    `json:"updatedAt"`

	PixToUsdOps []PixToUsdOpAPI `json:"pixToUsdOps"`
}

type MintOpAPI struct {
	Id               *entities.ID                `json:"id"`
	Amount           int                         `json:"amount"`
	Reason           string                      `json:"createdReason"`
	CreatedAt        time.Time                   `json:"createdAt"`
	Fee              int                         `json:"fee"`
	SmartContractOps []SmartcontractOperationAPI `json:"smartContractOps"`
}

type PixToUsdOpAPI struct {
	Id               *entities.ID                `json:"id"`
	BrlaAmount       int                         `json:"brlaAmount"`
	UsdAmount        int                         `json:"usdAmount"`
	Coin             string                      `json:"coin"`
	CreatedAt        time.Time                   `json:"createdAt"`
	SmartContractOps []SmartcontractOperationAPI `json:"smartContractOps"`
}

type SwapOpAPI struct {
	Id              *entities.ID `json:"id"`
	Chain           string       `json:"chain"`
	WalletAddress   string       `json:"walletAddress"`
	ReceiverAddress string       `json:"receiverAddress"`

	BrlaAmount   int    `json:"brlaAmount"`
	UsdAmount    int    `json:"usdAmount"`
	UsdToBrla    bool   `json:"usdToBrla"`
	Coin         string `json:"coin"`
	UserDocument string `json:"userDocument"`
	BasePrice    string `json:"basePrice"`
	BaseFee      string `json:"baseFee"`

	SmartContractOps []SmartcontractOperationAPI `json:"smartContractOps"`

	Permit *PermitData `json:"permit"`

	CreatedAt time.Time `json:"createdAt"`
}
