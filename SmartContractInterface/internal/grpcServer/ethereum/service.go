package ethereum

import (
	"context"
	"github.com/EricBastos/ProjetoTG/SmartContractInterface/internal/contractCaller"
	"github.com/EricBastos/ProjetoTG/SmartContractInterface/internal/contractInterface/ethereum"
	"github.com/EricBastos/ProjetoTG/SmartContractInterface/internal/grpcServer/pb"
)

type EthereumService struct {
	pb.UnimplementedEthereumServiceServer
}

func NewEthereumService() *EthereumService {
	return &EthereumService{}
}

func (c *EthereumService) IsWaitingPermit(ctx context.Context, in *pb.WalletAddress) (*pb.IsWaitingPermitOutput, error) {
	return &pb.IsWaitingPermitOutput{Waiting: contractCaller.CC.IsWaitingPermit("Ethereum", in.Wallet)}, nil
}

func (c *EthereumService) GetBalance(ctx context.Context, in *pb.GetBalanceInput) (*pb.GetBalanceOutput, error) {
	balance, err := ethereum.GetWalletBalance(in.TokenAddress, in.WalletAddress)
	if err != nil {
		return nil, err
	}
	return &pb.GetBalanceOutput{Balance: balance}, nil
}

func (c *EthereumService) GetAllowance(ctx context.Context, in *pb.GetAllowanceInput) (*pb.GetAllowanceOutput, error) {
	allowance, err := ethereum.GetAllowance(in.OwnerAddress, in.SpenderAddress, in.TokenAddress)
	if err != nil {
		return nil, err
	}
	return &pb.GetAllowanceOutput{Allowance: allowance}, nil
}
