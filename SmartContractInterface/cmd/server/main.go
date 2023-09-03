package main

import (
	"fmt"
	"github.com/EricBastos/ProjetoTG/Library/database"
	"github.com/EricBastos/ProjetoTG/Library/entities"
	"github.com/EricBastos/ProjetoTG/SmartContractInterface/configs"
	"github.com/EricBastos/ProjetoTG/SmartContractInterface/internal/grpcServer/ethereum"
	"github.com/EricBastos/ProjetoTG/SmartContractInterface/internal/grpcServer/pb"
	"github.com/EricBastos/ProjetoTG/SmartContractInterface/internal/grpcServer/polygon"
	"github.com/EricBastos/ProjetoTG/SmartContractInterface/internal/operator"
	"github.com/EricBastos/ProjetoTG/SmartContractInterface/internal/rabbitmqClient"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"net"
	"os"
	"os/signal"
	"time"
)

func main() {

	config, err := configs.LoadConfig(".")
	if err != nil {
		log.Fatal(err.Error())
	}

	rabbit, err := rabbitmqClient.NewRabbitMQClient(
		config.RABBITUser,
		config.RABBITPassword,
		config.RABBITHost,
		config.RABBITPort,
		config.RABBITCallExchange,
		map[string]string{
			"Ethereum": config.RABBITCallQueueEthereum,
			"Polygon":  config.RABBITCallQueuePolygon,
		},
	)
	if err != nil {
		log.Println("Waiting rabbitmq:", err.Error())
		time.Sleep(3 * time.Second)
		os.Exit(1)
	}

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		config.DBHost,
		config.DBPort,
		config.DBUser,
		config.DBPassword,
		config.DBName)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(err.Error())
	}
	a, err := db.DB()
	if err != nil {
		log.Fatal(err.Error())
	}
	a.SetMaxOpenConns(30)
	a.SetConnMaxLifetime(2 * time.Minute)

	err = db.AutoMigrate(
		&entities.SmartcontractOperation{},
		&entities.Feedback{},
		&entities.BurnOp{},
		&entities.BridgeOp{},
		&entities.Transfer{},
	)
	if err != nil {
		log.Fatal(err.Error())
	}

	notifyOffline := make(chan bool)

	smartcontractOpDb := database.NewSmartcontractOperationDB(db)
	feedbackDb := database.NewFeedbackDB(db)
	burnOpDb := database.NewBurnOperationsDB(db)
	bridgeOpDb := database.NewBridgeOperationsDB(db)
	transferDb := database.NewTransferDB(db)

	op := operator.NewOperator(smartcontractOpDb, feedbackDb, burnOpDb, bridgeOpDb, transferDb, rabbit, notifyOffline)

	go op.Start()

	ethereumService := ethereum.NewEthereumService()
	polygonService := polygon.NewPolygonService()

	grpcServer := grpc.NewServer()
	reflection.Register(grpcServer)
	pb.RegisterEthereumServiceServer(grpcServer, ethereumService)
	pb.RegisterPolygonServiceServer(grpcServer, polygonService)

	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatal(err.Error())
	}

	go func() {
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatal(err.Error())
		}
	}()

	log.Println("System up")

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

ForLoop:
	for {
		select {
		case <-interrupt:
			break ForLoop
		case <-notifyOffline:
			break ForLoop
		}
	}

	grpcServer.GracefulStop()

}
