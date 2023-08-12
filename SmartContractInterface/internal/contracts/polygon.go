package contracts

import (
	"context"
	"errors"
	"github.com/EricBastos/ProjetoTG/Library/entities"
	"github.com/EricBastos/ProjetoTG/SmartContractInterface/internal/contracts/polygon"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"log"
	"os"
	"os/signal"
	"sync"
	"time"
)

type BRLAPolygonContract struct {
	mu                 sync.Mutex
	onGoingTransactons map[*types.Transaction]*OnGoingTransactionData
	waitingForPermit   map[string]bool
	transactionNotify  chan *OnGoingTransactionData
}

func NewBRLAPolygonContract(polyConfig *PolygonContractConfig) *BRLAPolygonContract {
	polygon.SetupClients(polyConfig.RPCHost, polyConfig.ContractAddress)
	polygon.SetupAccounts(polyConfig.OwnerPrivateKey)
	polygon.InitializeBRLA()
	return (&BRLAPolygonContract{}).initialize(
		polyConfig.TransactionNotify)
}

func (c *BRLAPolygonContract) Mint(txData *OnGoingTransactionData, skipTest bool, address string, amount int) (*TxResult, error) {
	tx, err := polygon.Mint(skipTest, address, amount)
	if err != nil {
		return nil, err
	}
	c.mu.Lock()
	txData.TxHash = tx.Hash().String()
	c.onGoingTransactons[tx] = txData
	c.mu.Unlock()
	return &TxResult{Hash: tx.Hash().String()}, nil
}

func (c *BRLAPolygonContract) BurnWithPermit(txData *OnGoingTransactionData, skipTest bool, address string, amount int, permit *entities.PermitData) (*TxResult, error) {

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

	tx, err := polygon.Burn(skipTest, address, amount, permit.Deadline, permit.V, r, s)
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

func (c *BRLAPolygonContract) GetWalletBalance(tokenAddress, walletAddress string) (string, error) {
	return polygon.GetWalletBalance(tokenAddress, walletAddress)
}

func (c *BRLAPolygonContract) WaitingPermit(walletAddress string) bool {
	c.mu.Lock()
	res := c.waitingForPermit[walletAddress]
	c.mu.Unlock()
	return res
}

func (c *BRLAPolygonContract) initialize(tNC chan *OnGoingTransactionData) *BRLAPolygonContract {
	c.transactionNotify = tNC
	c.onGoingTransactons = make(map[*types.Transaction]*OnGoingTransactionData)
	c.waitingForPermit = make(map[string]bool)
	go c.trackTransactions()
	return c
}

func (c *BRLAPolygonContract) trackTransactions() {
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
				receipt, err := polygon.Polygon.TransactionReceipt(context.Background(), tx.Hash())
				// Block with transaction mined
				if err == nil {
					log.Println("(Polygon) Receipt received for transaction:", tx.Hash().Hex())

					if receipt.Status == 0 {
						errReason, err := polygon.ErrorReason(context.Background(), polygon.Polygon, tx, receipt.BlockNumber)
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
				if currTime-info.PostedTime >= 1*30*1000 {
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
