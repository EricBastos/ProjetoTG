package ethereum

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"log"
)

var (
	Eth             *ethclient.Client
	Contract        *brlatoken.StableContract
	ContractAddress common.Address
)

func SetupClients(
	RPCHost string,
	contAddress string) {
	ethClient, err := ethclient.Dial(RPCHost)
	if err != nil {
		log.Fatal(err.Error())
	}
	Eth = ethClient

	contractAddress := common.HexToAddress(contAddress)
	ContractAddress = contractAddress
	ethContract, err := brlatoken.NewStableContract(contractAddress, Eth)
	if err != nil {
		log.Fatal(err.Error())
	}

	Contract = ethContract
}
