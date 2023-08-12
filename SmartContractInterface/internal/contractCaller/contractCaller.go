package contractCaller

import (
	"errors"
	"github.com/EricBastos/ProjetoTG/Library/entities"
	"github.com/EricBastos/ProjetoTG/SmartContractInterface/internal/contracts"
)

var CC *ContractCaller

type ContractCaller struct {
	Contracts map[string]contracts.BRLAContract
}

func InitializeContractCaller(
	ethConfig *contracts.EthContractConfig,
	polyConfig *contracts.PolygonContractConfig) {
	c := ContractCaller{
		Contracts: map[string]contracts.BRLAContract{
			"Ethereum": contracts.NewBRLAEthereumContract(ethConfig),
			"Polygon":  contracts.NewBRLAPolygonContract(polyConfig),
		},
	}
	CC = &c
}

func (c *ContractCaller) Mint(txData *contracts.OnGoingTransactionData, skipTest bool, chain string, address string, amount int) (*contracts.TxResult, error) {
	contract, ok := c.Contracts[chain]
	if !ok {
		return nil, errors.New("chain not implemented")
	}
	return contract.Mint(txData, skipTest, address, amount)
}

func (c *ContractCaller) Burn(txData *contracts.OnGoingTransactionData, skipTest bool, chain string, address string, amount int, permitData *entities.PermitData) (*contracts.TxResult, error) {
	contract, ok := c.Contracts[chain]
	if !ok {
		return nil, errors.New("chain not implemented")
	}
	return contract.BurnWithPermit(txData, skipTest, address, amount, permitData)
}

func (c *ContractCaller) GetWalletBalance(chain string, address string, tokenAddress string) (string, error) {
	contract, ok := c.Contracts[chain]
	if !ok {
		return "", errors.New("chain not implemented")
	}
	return contract.GetWalletBalance(tokenAddress, address)
}

func (c *ContractCaller) IsWaitingPermit(chain string, walletAddress string) bool {
	contract, ok := c.Contracts[chain]
	if !ok {
		return true
	}
	return contract.WaitingPermit(walletAddress)
}
