package contractInterface

import (
	"github.com/EricBastos/ProjetoTG/Library/entities"
)

type BRLAContract interface {
	Mint(txData *OnGoingTransactionData, skipTest bool, address string, amount int) (*TxResult, error)
	BurnWithPermit(txData *OnGoingTransactionData, skipTest bool, address string, amount int, permitData *entities.PermitData) (*TxResult, error)
	GetWalletBalance(tokenAddress, walletAddress string) (string, error)
	WaitingPermit(walletAddress string) bool
}

type TxResult struct {
	Hash string
}

type OnGoingTransactionData struct {
	SmartcontractCallId string
	OperationOriginType string
	TxHash              string
	Success             bool
	ErrReason           string

	Operation     string
	WalletAddress string
	PostedTime    int64
	ID            string
}

type EthContractConfig struct {
	RPCHost           string
	ContractAddress   string
	OwnerPrivateKey   string
	TransactionNotify chan *OnGoingTransactionData
}

type PolygonContractConfig struct {
	RPCHost           string
	ContractAddress   string
	OwnerPrivateKey   string
	TransactionNotify chan *OnGoingTransactionData
}
