package ethereum

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"golang.org/x/crypto/sha3"
	"log"
	"math/big"
	"strings"
	"sync"
	"time"
)

var (
	ChainID   *big.Int
	txLock    sync.Mutex
	currNonce int64
)

func InitializeBRLA() {
	var err error
	// Get chain ID
	ChainID, err = Eth.ChainID(context.Background())
	if err != nil {
		log.Fatal(err.Error())
	}
	if err := UpdateNonce(); err != nil {
		log.Fatal(err)
	}
}

func UpdateNonce() error {
	nonce, err := Eth.PendingNonceAt(context.Background(), OwnerAddress)
	if err != nil {
		return err
	}
	currNonce = int64(nonce)
	return nil
}

func Mint(skipTest bool, address string, amount int) (res *types.Transaction, err error) {
	txLock.Lock()
	defer txLock.Unlock()
	res, err = mint(skipTest, address, amount)
	if err != nil && (strings.Contains(err.Error(), "nonce too low") || strings.Contains(err.Error(), "transaction underpriced")) {
		for i := 0; UpdateNonce() != nil && i < 3; i++ {
			time.Sleep(1 * time.Second)
		}
	}
	return
}

func mint(skipTest bool, address string, amount int) (*types.Transaction, error) {
	to := common.HexToAddress(address)
	qty := big.NewInt(int64(amount))
	dec := new(big.Int)
	dec, _ = dec.SetString("10000000000000000", 10)
	qty = qty.Mul(qty, dec)
	auth, err := setupTransaction(0, 300000)
	if err != nil {
		return nil, err
	}
	auth.NoSend = true

	log.Println("Testing transaction with auth:", auth.GasLimit, auth.GasTipCap.String(), auth.GasFeeCap.String())
	txTest, err := Contract.Mint(auth, to, qty)
	if err != nil {
		return nil, err
	}

	if !skipTest {
		err = testTransaction(txTest)
		if err != nil {
			return nil, err
		}
	}

	auth.NoSend = false
	log.Println("Sending transaction with auth:", auth.GasLimit, auth.GasTipCap.String(), auth.GasFeeCap.String())
	tx, err := Contract.Mint(auth, to, qty)
	if err != nil {
		return nil, err
	}
	currNonce++
	return tx, nil
}

func Burn(skipTest bool, address string, amount int, deadline int64, v uint8, r, s [32]byte) (res *types.Transaction, err error) {
	txLock.Lock()
	defer txLock.Unlock()
	res, err = burn(skipTest, address, amount, deadline, v, r, s)
	if err != nil && (strings.Contains(err.Error(), "nonce too low") || strings.Contains(err.Error(), "transaction underpriced")) {
		for i := 0; UpdateNonce() != nil && i < 3; i++ {
			time.Sleep(1 * time.Second)
		}
	}
	return
}

func burn(skipTest bool, address string, amount int, deadline int64, v uint8, r, s [32]byte) (*types.Transaction, error) {
	target := common.HexToAddress(address)
	qty := big.NewInt(int64(amount))
	dec := new(big.Int)
	dec, _ = dec.SetString("10000000000000000", 10)
	qty = qty.Mul(qty, dec)
	deadlineBigInt := big.NewInt(deadline)
	auth, err := setupTransaction(0, 300000)
	if err != nil {
		return nil, err
	}
	auth.NoSend = true

	log.Println("Testing transaction with auth:", auth.GasLimit, auth.GasTipCap.String(), auth.GasFeeCap.String())
	txTest, err := Contract.BurnFromWithPermit(auth, target, OwnerAddress, qty, deadlineBigInt, v, r, s)
	if err != nil {
		return nil, err
	}

	if !skipTest {
		err = testTransaction(txTest)
		if err != nil {
			return nil, err
		}
	}

	auth.NoSend = false
	log.Println("Sending transaction with auth:", auth.GasLimit, auth.GasTipCap.String(), auth.GasFeeCap.String())
	tx, err := Contract.BurnFromWithPermit(auth, target, OwnerAddress, qty, deadlineBigInt, v, r, s)
	if err != nil {
		return nil, err
	}
	currNonce++
	return tx, nil
}

func EstimateGasPrice() (string, error) {

	_, gasFeeCap, err := getEIP1559Gas()
	if err != nil {
		return "", fmt.Errorf("[EstimateGasPrice] internal error: %w", err)
	}

	return gasFeeCap.String(), nil

}

func GetWalletBalance(tokenAddress, walletAddress string) (string, error) {

	toAddress := common.HexToAddress(walletAddress)
	transferFnSignature := []byte("balanceOf(address)")
	hash := sha3.NewLegacyKeccak256()
	hash.Write(transferFnSignature)
	methodID := hash.Sum(nil)[:4]
	paddedAddress := common.LeftPadBytes(toAddress.Bytes(), 32)

	var data []byte
	data = append(data, methodID...)
	data = append(data, paddedAddress...)

	tokAddress := common.HexToAddress(tokenAddress)

	tx := &types.DynamicFeeTx{
		To:    &tokAddress,
		Value: big.NewInt(0),
		Data:  data,
	}

	msg := ethereum.CallMsg{
		From:  toAddress,
		To:    tx.To,
		Gas:   tx.Gas,
		Value: tx.Value,
		Data:  tx.Data,
	}
	res, err := Eth.CallContract(context.Background(), msg, nil)
	if err != nil {
		return "", err
	}

	i := new(big.Int)
	i.SetString(common.Bytes2Hex(res), 16)

	return i.String(), nil
}

func GetAllowance(ownerAddress, spenderAddress, tokenAddress string) (string, error) {

	owner := common.HexToAddress(ownerAddress)
	spender := common.HexToAddress(spenderAddress)
	transferFnSignature := []byte("allowance(address,address)")
	hash := sha3.NewLegacyKeccak256()
	hash.Write(transferFnSignature)
	methodID := hash.Sum(nil)[:4]
	paddedOwnerAddress := common.LeftPadBytes(owner.Bytes(), 32)
	paddedSpenderAddress := common.LeftPadBytes(spender.Bytes(), 32)

	var data []byte
	data = append(data, methodID...)
	data = append(data, paddedOwnerAddress...)
	data = append(data, paddedSpenderAddress...)

	tokenAddressHex := common.HexToAddress(tokenAddress)

	tx := &types.DynamicFeeTx{
		To:    &tokenAddressHex,
		Value: big.NewInt(0),
		Data:  data,
	}

	msg := ethereum.CallMsg{
		From:  owner,
		To:    tx.To,
		Gas:   tx.Gas,
		Value: tx.Value,
		Data:  tx.Data,
	}
	res, err := Eth.CallContract(context.Background(), msg, nil)
	if err != nil {
		return "", err
	}

	i := new(big.Int)
	i.SetString(common.Bytes2Hex(res), 16)

	return i.String(), nil
}

func setupTransaction(value, gasLimit int64) (*bind.TransactOpts, error) {
	// Set transactor
	auth, err := bind.NewKeyedTransactorWithChainID(privateKey, ChainID)
	if err != nil {
		return nil, err
	}

	gasTipCap, gasFeeCap, err := getEIP1559Gas()
	if err != nil {
		return nil, fmt.Errorf("[Eth.setupTransaction] internal error: %w", err)
	}

	auth.Nonce = big.NewInt(currNonce)
	auth.Value = big.NewInt(value)
	auth.GasLimit = uint64(gasLimit)
	auth.GasTipCap = gasTipCap
	auth.GasFeeCap = gasFeeCap
	return auth, nil
}

func getGasTipCap() (*big.Int, error) {
	gasTip, err := Eth.SuggestGasTipCap(context.Background())
	if err != nil {
		return nil, fmt.Errorf("[Eth.getGasTipCap] internal error: %w", err)
	}
	if gasTip.Cmp(big.NewInt(30000000000)) < 0 {
		gasTip = big.NewInt(30000000000)
	}
	return gasTip, nil
}

func getEIP1559Gas() (*big.Int, *big.Int, error) {

	gasFeeCap := new(big.Int)

	gasTipCap, err := getGasTipCap()
	if err != nil {
		return nil, nil, fmt.Errorf("[Eth.getEIP1559Gas] internal error: %w", err)
	}

	head, err := Eth.HeaderByNumber(context.Background(), nil)
	if err != nil {
		return nil, nil, fmt.Errorf("[Eth.HeaderByNumber] internal error: %w", err)
	}
	gasFeeCap = new(big.Int).Add(
		gasTipCap,
		new(big.Int).Mul(head.BaseFee, big.NewInt(4)),
	)
	if gasFeeCap.Cmp(gasTipCap) < 0 {
		return nil, nil, fmt.Errorf("maxFeePerGas (%v) < maxPriorityFeePerGas (%v)", gasFeeCap, gasTipCap)
	}

	return gasTipCap, gasFeeCap, nil

}

func ErrorReason(ctx context.Context, b ethereum.ContractCaller, tx *types.Transaction, blockNum *big.Int) (string, error) {
	msg := ethereum.CallMsg{
		From:     OwnerAddress,
		To:       tx.To(),
		Gas:      tx.Gas(),
		GasPrice: tx.GasPrice(),
		Value:    tx.Value(),
		Data:     tx.Data(),
	}
	res, err := b.CallContract(ctx, msg, blockNum)
	if err != nil {
		return "", err
	}
	if len(res) < 4 {
		return "transaction error", nil
	}
	return unpackError(res)
}

var (
	errorSig     = []byte{0x08, 0xc3, 0x79, 0xa0} // Keccak256("Error(string)")[:4]
	abiString, _ = abi.NewType("string", "", nil)
)

func unpackError(result []byte) (string, error) {
	if !bytes.Equal(result[:4], errorSig) {
		return "<tx result not Error(string)>", errors.New("TX result not of type Error(string)")
	}
	vs, err := abi.Arguments{{Type: abiString}}.UnpackValues(result[4:])
	if err != nil {
		return "<invalid tx result>", err
	}
	return vs[0].(string), nil
}

func testTransaction(tx *types.Transaction) error {
	from, err := types.Sender(types.LatestSignerForChainID(tx.ChainId()), tx)
	if err != nil {
		return err
	}

	callMsg := ethereum.CallMsg{
		From:       from,
		To:         tx.To(),
		Gas:        tx.Gas(),
		GasPrice:   tx.GasPrice(),
		GasFeeCap:  tx.GasFeeCap(),
		GasTipCap:  tx.GasTipCap(),
		Value:      tx.Value(),
		Data:       tx.Data(),
		AccessList: tx.AccessList(),
	}

	_, err = Eth.EstimateGas(context.Background(), callMsg)

	if err != nil {
		return err
	}

	return nil
}
