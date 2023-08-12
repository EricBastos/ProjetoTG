package polygon

import (
	"github.com/EricBastos/ProjetoTG/SmartContractInterface/internal/smartContract"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"log"
)

var (
	Polygon         *ethclient.Client
	Contract        *smartContract.StableCoin
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
	ethContract, err := smartContract.NewStableCoin(contractAddress, Polygon)
	if err != nil {
		log.Fatal(err.Error())
	}

	Contract = ethContract
}
