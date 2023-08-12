package contractCaller

import (
	"errors"
	"github.com/EricBastos/ProjetoTG/Library/entities"
	"github.com/EricBastos/ProjetoTG/SmartContractInterface/internal/contractInterface"
)

var CC *ContractCaller

type ContractCaller struct {
	Contracts map[string]contractInterface.BRLAContract
}

func InitializeContractCaller(
	ethConfig *contractInterface.EthContractConfig,
	polyConfig *contractInterface.PolygonContractConfig) {
	c := ContractCaller{
		Contracts: map[string]contractInterface.BRLAContract{
			"Ethereum": contractInterface.NewBRLAEthereumContract(ethConfig),
			"Polygon":  contractInterface.NewBRLAPolygonContract(polyConfig),
		},
	}
	CC = &c
}

func (c *ContractCaller) Mint(txData *contractInterface.OnGoingTransactionData, skipTest bool, chain string, address string, amount int) (*contractInterface.TxResult, error) {
	contract, ok := c.Contracts[chain]
	if !ok {
		return nil, errors.New("chain not implemented")
	}
	return contract.Mint(txData, skipTest, address, amount)
}

func (c *ContractCaller) Burn(txData *contractInterface.OnGoingTransactionData, skipTest bool, chain string, address string, amount int, permitData *entities.PermitData) (*contractInterface.TxResult, error) {
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
