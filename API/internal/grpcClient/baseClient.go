package grpcClient

import (
	"fmt"
	"github.com/EricBastos/ProjetoTG/API/configs"
	"github.com/EricBastos/ProjetoTG/API/internal/grpcClient/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type client struct {
	connectionSmartcontract *grpc.ClientConn
	EthereumClient          pb.EthereumServiceClient
	PolygonClient           pb.PolygonServiceClient
}

func InitializeServices() error {
	connSmartcontractString := fmt.Sprintf("dns:///%s:%s", configs.Cfg.GRPCSmartcontractHost, configs.Cfg.GRPCSmartcontractPort)
	connSmartcontract, err := grpc.Dial(connSmartcontractString, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return err
	}

	c := &client{
		connectionSmartcontract: connSmartcontract,
		EthereumClient:          pb.NewEthereumServiceClient(connSmartcontract),
		PolygonClient:           pb.NewPolygonServiceClient(connSmartcontract),
	}
	EthereumService = &EthereumGRPCService{pbClient: c.EthereumClient}
	PolygonService = &PolygonGRPCService{pbClient: c.PolygonClient}
	return nil
}
