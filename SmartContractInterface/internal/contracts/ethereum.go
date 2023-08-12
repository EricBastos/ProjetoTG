package contracts

import (
	"context"
	"errors"
	"github.com/EricBastos/ProjetoTG/Library/entities"
	ethereum2 "github.com/EricBastos/ProjetoTG/SmartContractInterface/internal/contracts/ethereum"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"log"
	"os"
	"os/signal"
	"sync"
	"time"
)

type BRLAEthereumContract struct {
	mu                 sync.Mutex
	onGoingTransactons map[*types.Transaction]*OnGoingTransactionData
	waitingForPermit   map[string]bool
	transactionNotify  chan *OnGoingTransactionData
}

func NewBRLAEthereumContract(ethConfig *EthContractConfig) *BRLAEthereumContract {
	ethereum2.SetupClients(ethConfig.RPCHost, ethConfig.ContractAddress)
	ethereum2.SetupAccounts(ethConfig.OwnerPrivateKey)
	ethereum2.InitializeBRLA()
	return (&BRLAEthereumContract{}).initialize(
		ethConfig.TransactionNotify)
}

func (c *BRLAEthereumContract) Mint(txData *OnGoingTransactionData, skipTest bool, address string, amount int) (*TxResult, error) {
	tx, err := ethereum2.Mint(skipTest, address, amount)
	if err != nil {
		return nil, err
	}
	c.mu.Lock()
	txData.TxHash = tx.Hash().String()
	c.onGoingTransactons[tx] = txData
	c.mu.Unlock()
	return &TxResult{Hash: tx.Hash().String()}, nil
}

func (c *BRLAEthereumContract) BurnWithPermit(
	txData *OnGoingTransactionData,
	skipTest bool,
	address string,
	amount int,
	permit *entities.PermitData) (*TxResult, error) {

	c.mu.Lock()
	_, burnExists := c.waitingForPermit[address]
	if burnExists {
		c.mu.Unlock()
		return nil, errors.New("must wait for past permit transaction to finish")
	}
	c.waitingForPermit[address] = true
	c.mu.Unlock()

	r := common.HexToHash(permit.R)
	s := common.HexToHash(permit.S)

	tx, err := ethereum2.Burn(skipTest, address, amount, permit.Deadline, permit.V, r, s)
	if err != nil {
		c.mu.Lock()
		delete(c.waitingForPermit, address)
		c.mu.Unlock()
		return nil, err
	}
	c.mu.Lock()
	txData.TxHash = tx.Hash().String()
	c.onGoingTransactons[tx] = txData
	c.mu.Unlock()
	return &TxResult{Hash: tx.Hash().String()}, nil
}

func (c *BRLAEthereumContract) GetWalletBalance(tokenAddress, walletAddress string) (string, error) {
	return ethereum2.GetWalletBalance(tokenAddress, walletAddress)
}

func (c *BRLAEthereumContract) initialize(tNC chan *OnGoingTransactionData) *BRLAEthereumContract {
	c.transactionNotify = tNC
	c.onGoingTransactons = make(map[*types.Transaction]*OnGoingTransactionData)
	c.waitingForPermit = make(map[string]bool)
	go c.trackTransactions()
	return c
}

func (c *BRLAEthereumContract) WaitingPermit(walletAddress string) bool {
	c.mu.Lock()
	res := c.waitingForPermit[walletAddress]
	c.mu.Unlock()
	return res
}

func (c *BRLAEthereumContract) trackTransactions() {
	pollTicker := time.NewTicker(15 * time.Second)
	itt := make(chan os.Signal, 1)
	signal.Notify(itt, os.Interrupt)
	defer pollTicker.Stop()
	for {
		select {
		case <-itt:
			return
		case <-pollTicker.C:
			currTime := time.Now().UnixMilli()
			c.mu.Lock()
			for tx, info := range c.onGoingTransactons {
				if currTime-info.PostedTime < 15*1000 {
					continue
				}
				receipt, err := ethereum2.Eth.TransactionReceipt(context.Background(), tx.Hash())
				// Block with transaction mined
				if err == nil {
					log.Println("(Ethereum) Receipt received for transaction:", tx.Hash().Hex())

					if receipt.Status == 0 {
						errReason, err := ethereum2.ErrorReason(context.Background(), ethereum2.Eth, tx, receipt.BlockNumber)
						if err != nil {
							errReason = err.Error()
						}
						info.Success = false
						info.ErrReason = errReason
					} else {
						info.Success = true
					}

					c.transactionNotify <- info
					delete(c.onGoingTransactons, tx)
					delete(c.waitingForPermit, info.WalletAddress)
					continue
				}

				// 5 minutes threshold to send error (maybe need to readjust gas, for example)
				if currTime-info.PostedTime >= 1*60*1000 {
					info.Success = false
					info.ErrReason = "more than 5 minutes passed since transaction was posted"
					log.Println("More than 5 minutes passed since transaction was posted:", tx.Hash().Hex())

					c.transactionNotify <- info
					delete(c.onGoingTransactons, tx)
					delete(c.waitingForPermit, info.WalletAddress)
					// Todo: maybe post cancel transaction?
				}
			}
			c.mu.Unlock()
		}
	}
}
