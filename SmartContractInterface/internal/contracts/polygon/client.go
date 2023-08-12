package polygon

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"log"
)

var (
	Polygon         *ethclient.Client
	Contract        *brlatoken.BRLAContract
	ContractAddress common.Address
)

func SetupClients(
	RPCHost string,
	contAddress string) {
	ethClient, err := ethclient.Dial(RPCHost)
	if err != nil {
		log.Fatal(err.Error())
	}
	Polygon = ethClient

	contractAddress := common.HexToAddress(contAddress)
	ContractAddress = contractAddress
	ethContract, err := brlatoken.NewStableContract(contractAddress, Polygon)
	if err != nil {
		log.Fatal(err.Error())
	}

	Contract = ethContract
}
