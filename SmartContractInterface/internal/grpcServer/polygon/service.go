package polygon

import (
	"context"
	"github.com/EricBastos/ProjetoTG/SmartContractInterface/internal/contractCaller"
	"github.com/EricBastos/ProjetoTG/SmartContractInterface/internal/contracts/polygon"
	"github.com/EricBastos/ProjetoTG/SmartContractInterface/internal/grpcServer/pb"
)

type PolygonService struct {
	pb.UnimplementedPolygonServiceServer
}

func NewPolygonService() *PolygonService {
	return &PolygonService{}
}

func (c *PolygonService) IsWaitingPermit(ctx context.Context, in *pb.WalletAddress) (*pb.IsWaitingPermitOutput, error) {
	return &pb.IsWaitingPermitOutput{Waiting: contractCaller.CC.IsWaitingPermit("Polygon", in.Wallet)}, nil
}

func (c *PolygonService) GetBalance(ctx context.Context, in *pb.GetBalanceInput) (*pb.GetBalanceOutput, error) {
	balance, err := polygon.GetWalletBalance(in.TokenAddress, in.WalletAddress)
	if err != nil {
		return nil, err
	}
	return &pb.GetBalanceOutput{Balance: balance}, nil
}

func (c *PolygonService) GetAllowance(ctx context.Context, in *pb.GetAllowanceInput) (*pb.GetAllowanceOutput, error) {
	allowance, err := polygon.GetAllowance(in.OwnerAddress, in.SpenderAddress, in.TokenAddress)
	if err != nil {
		return nil, err
	}
	return &pb.GetAllowanceOutput{Allowance: allowance}, nil
}
