package polygon

import (
	"crypto/ecdsa"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"log"
)

var (
	privateKey   *ecdsa.PrivateKey
	OwnerAddress common.Address
)

func SetupAccounts(ownerPrivateKey string) {
	// Set private key and infer public key and address from it
	var err error
	privateKey, err = crypto.HexToECDSA(ownerPrivateKey)
	if err != nil {
		log.Fatal(err.Error())
	}
	publicKey := privateKey.Public()
	publicKeyOk, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Fatal("Cannot assert type: publicKey is not of type *ecdsa.PublicKey")
	}
	OwnerAddress = crypto.PubkeyToAddress(*publicKeyOk)
}
