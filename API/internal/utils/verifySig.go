package utils

import (
	"crypto/ecdsa"
	"errors"
	"fmt"
	"github.com/EricBastos/ProjetoTG/API/configs"
	"github.com/EricBastos/ProjetoTG/Library/entities"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/common/math"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/signer/core/apitypes"
	"log"
	"math/big"
	"strconv"
)

var ValidChains = map[string]func(walletAddress *string) error{
	"Ethereum": validateEth,
	"Polygon":  validatePolygon,
}

func VerifyBurnPermit(from, chain string, amountInt int, permit *entities.PermitData) bool {

	var result bool

	// EVM-compatible will be verified later when trying to execute smartcontract operation
	if chain == "Ethereum" {
		var amount = new(big.Int)
		amount = amount.Mul(big.NewInt(int64(amountInt)), big.NewInt(10000000000000000))
		result = ValidatePermit(
			"StableCoin",
			"1",
			false,
			common.HexToAddress(from),
			common.HexToAddress(configs.Cfg.EthereumWalletAddress),
			common.HexToAddress(configs.Cfg.EthereumTokenContract),
			11155111,
			amount,
			permit)
	} else if chain == "Polygon" {
		var amount = new(big.Int)
		amount = amount.Mul(big.NewInt(int64(amountInt)), big.NewInt(10000000000000000))
		result = ValidatePermit(
			"StableCoin",
			"1",
			false,
			common.HexToAddress(from),
			common.HexToAddress(configs.Cfg.PolygonWalletAddress),
			common.HexToAddress(configs.Cfg.PolygonTokenContract),
			80001,
			amount,
			permit)
	}
	return result
}

func validateEth(walletAddress *string) error {
	if !common.IsHexAddress(*walletAddress) {
		return errors.New(InvalidAddrError)
	}
	addr := common.HexToAddress(*walletAddress)
	*walletAddress = addr.String()
	//blacklisted, err := grpcClient.EthereumService.IsBlackListed(walletAddress)
	//if err != nil {
	//	return errors.New(InternalError)
	//}
	//if blacklisted {
	//	return errors.New(BlackListed)
	//}
	return nil
}

func validatePolygon(walletAddress *string) error {
	if !common.IsHexAddress(*walletAddress) {
		return errors.New(InvalidAddrError)
	}
	addr := common.HexToAddress(*walletAddress)
	*walletAddress = addr.String()
	//blacklisted, err := grpcClient.PolygonService.IsBlackListed(walletAddress)
	//if err != nil {
	//	return errors.New(InternalError)
	//}
	//if blacklisted {
	//	return errors.New(BlackListed)
	//}
	return nil
}

func GeneratePermitHash(
	contractName string,
	contractVersion string,
	legacyPermit bool,
	owner common.Address,
	spender common.Address,
	verifyingContract common.Address,
	chainId int64,
	value *big.Int,
	nonce int64,
	deadline int64,
) (common.Hash, error) {

	val := math.HexOrDecimal256(*value)

	var domain apitypes.TypedDataDomain
	var typesPermit apitypes.Types

	if !legacyPermit {
		domain = apitypes.TypedDataDomain{
			Name:              contractName,
			Version:           contractVersion,
			ChainId:           math.NewHexOrDecimal256(chainId),
			VerifyingContract: verifyingContract.String(),
		}
		typesPermit = apitypes.Types{
			"EIP712Domain": []apitypes.Type{
				{Name: "name", Type: "string"},
				{Name: "version", Type: "string"},
				{Name: "chainId", Type: "uint256"},
				{Name: "verifyingContract", Type: "address"},
			},
			"Permit": []apitypes.Type{
				{Name: "owner", Type: "address"},
				{Name: "spender", Type: "address"},
				{Name: "value", Type: "uint256"},
				{Name: "nonce", Type: "uint256"},
				{Name: "deadline", Type: "uint256"},
			},
		}
	} else {
		domain = apitypes.TypedDataDomain{
			Name:              contractName,
			Version:           contractVersion,
			Salt:              common.HexToHash(strconv.FormatInt(chainId, 16)).String(),
			VerifyingContract: verifyingContract.String(),
		}
		typesPermit = apitypes.Types{
			"EIP712Domain": []apitypes.Type{
				{Name: "name", Type: "string"},
				{Name: "version", Type: "string"},
				{Name: "verifyingContract", Type: "address"},
				{Name: "salt", Type: "bytes32"},
			},
			"Permit": []apitypes.Type{
				{Name: "owner", Type: "address"},
				{Name: "spender", Type: "address"},
				{Name: "value", Type: "uint256"},
				{Name: "nonce", Type: "uint256"},
				{Name: "deadline", Type: "uint256"},
			},
		}
	}

	signerData := apitypes.TypedData{
		Types:       typesPermit,
		PrimaryType: "Permit",
		Domain:      domain,
		Message: apitypes.TypedDataMessage{
			"owner":    owner.String(),
			"spender":  spender.String(),
			"value":    &val,
			"nonce":    math.NewHexOrDecimal256(nonce),
			"deadline": math.NewHexOrDecimal256(deadline),
		},
	}

	log.Println(signerData.Map())
	log.Println(nonce, deadline, val)

	domainSeparator, err := signerData.HashStruct("EIP712Domain", signerData.Domain.Map())

	if err != nil {
		return common.Hash{}, err
	}

	typedDataHash, err := signerData.HashStruct(signerData.PrimaryType, signerData.Message)
	if err != nil {
		return common.Hash{}, err
	}

	rawData := []byte(fmt.Sprintf("\x19\x01%s%s", string(domainSeparator), string(typedDataHash)))
	hash := common.BytesToHash(crypto.Keccak256(rawData))

	log.Println(hash.String())
	return hash, nil
}

func GenerateSignedPermit(
	contractName string,
	contractVersion string,
	legacyPermit bool,
	owner common.Address,
	spender common.Address,
	verifyingContract common.Address,
	chainId int64,
	value *big.Int,
	nonce int64,
	deadline int64,
	privateKey *ecdsa.PrivateKey,
) (r string, s string, v uint8, dl int64, hash common.Hash, err error) {

	hash, err = GeneratePermitHash(contractName, contractVersion, legacyPermit, owner, spender, verifyingContract, chainId, value, nonce, deadline)
	if err != nil {
		return
	}

	signatureBytes, err := crypto.Sign(hash.Bytes(), privateKey)
	if err != nil {
		log.Fatal(err.Error())
	}

	r = hexutil.Encode(signatureBytes[:32])
	s = hexutil.Encode(signatureBytes[32:64])
	v = uint8(int(signatureBytes[64])) + 27
	dl = deadline

	return

}

func ValidatePermit(
	contractName string,
	contractVersion string,
	legacyPermit bool,
	owner common.Address,
	spender common.Address,
	verifyingContract common.Address,
	chainId int64,
	value *big.Int,
	permit *entities.PermitData) bool {

	if permit == nil {
		return false
	}

	rByte, err := hexutil.Decode(permit.R)
	if err != nil {
		return false
	}
	sByte, err := hexutil.Decode(permit.S)
	if err != nil {
		return false
	}

	sig := append(rByte, sByte...)
	sig = append(sig, permit.V-27)

	log.Println(contractName, contractVersion, legacyPermit, owner.String(), spender.String(), verifyingContract.String(), chainId, value.String())

	hash, err := GeneratePermitHash(contractName, contractVersion, legacyPermit, owner, spender, verifyingContract, chainId, value, permit.Nonce, permit.Deadline)
	if err != nil {
		return false
	}

	log.Println("Hash validate:", hash.String())

	pKeyRes, err := crypto.SigToPub(hash.Bytes(), sig)
	if err != nil {
		return false
	}

	addr := crypto.PubkeyToAddress(*pKeyRes)

	log.Println("Addr validated:", addr.String())
	log.Println("Expected:", owner.String())

	return addr.String() == owner.String()
}
