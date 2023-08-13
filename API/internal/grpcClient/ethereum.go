package grpcClient

import (
	"context"
	"github.com/EricBastos/ProjetoTG/API/internal/grpcClient/pb"
	"math/big"
)

var EthereumService *EthereumGRPCService

type EthereumGRPCService struct {
	pbClient pb.EthereumServiceClient
}

func (c *EthereumGRPCService) IsWaitingPermit(wallet string) (bool, error) {
	resp, err := c.pbClient.IsWaitingPermit(context.Background(), &pb.WalletAddress{Wallet: wallet})
	if err != nil {
		return true, err
	}
	return resp.Waiting, nil
}

func (c *EthereumGRPCService) GetBalance(tokenAddress, walletAddress string) (*big.Int, error) {
	resp, err := c.pbClient.GetBalance(context.Background(), &pb.GetBalanceInput{TokenAddress: tokenAddress, WalletAddress: walletAddress})
	if err != nil {
		return nil, err
	}
	res := new(big.Int)
	res.SetString(resp.Balance, 10)
	return res, nil
}

func (c *EthereumGRPCService) GetAllowance(ownerAddress, spenderAddress, tokenAddress string) (*big.Int, error) {
	resp, err := c.pbClient.GetAllowance(context.Background(), &pb.GetAllowanceInput{
		OwnerAddress:   ownerAddress,
		SpenderAddress: spenderAddress,
		TokenAddress:   tokenAddress,
	})
	if err != nil {
		return nil, err
	}
	res := new(big.Int)
	res.SetString(resp.Allowance, 10)
	return res, nil
}
