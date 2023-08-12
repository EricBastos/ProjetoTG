package main

import (
	"context"
	"fmt"
	"github.com/EricBastos/ProjetoTG/BankWebhook/configs"
	"github.com/EricBastos/ProjetoTG/BankWebhook/internal/infra/rabbitmqClient"
	"github.com/EricBastos/ProjetoTG/BankWebhook/internal/infra/webserver/handlers"
	"github.com/EricBastos/ProjetoTG/Library/database"
	"github.com/EricBastos/ProjetoTG/Library/entities"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

func main() {

	config, err := configs.LoadConfig(".")
	if err != nil {
		log.Fatal(err.Error())
	}

	// Rabbitmq setup
	rabbitClient, err := rabbitmqClient.NewRabbitMQClient(
		&rabbitmqClient.RabbitMQClientConfig{
			User:             config.RABBITUser,
			Pass:             config.RABBITPassword,
			Host:             config.RABBITHost,
			Port:             config.RABBITPort,
			ProducerExchange: config.RABBITCallExchange,
		},
	)

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
		&entities.MintOp{},
		&entities.StaticDeposit{},
		&entities.StaticDepositFeedback{},
		&entities.Transfer{},
		&entities.TransferFeedback{},
	)
	if err != nil {
		log.Fatal(err.Error())
	}

	mintOperationsDb := database.NewMintOperationsDB(db)
	transferDb := database.NewTransferDB(db)
	staticDepositDb := database.NewStaticDepositDB(db)
	staticDepositFeedbackDb := database.NewStaticDepositFeedbackDB(db)
	transferFeedbackDb := database.NewTransferFeedbackDB(db)

	bankWebhookHandler := handlers.NewWebhookHandler(
		mintOperationsDb,
		transferDb,
		transferFeedbackDb,
		staticDepositDb,
		staticDepositFeedbackDb,
		rabbitClient)

	r := chi.NewRouter()

	r.Group(func(r chi.Router) {
		r.Use(middleware.Logger)
		r.Use(middleware.Recoverer)
		r.Route("/", func(r chi.Router) {
			r.Post("/", bankWebhookHandler.Listen)
		})
	})

	server := &http.Server{Addr: ":8080", Handler: r}

	go func() {
		err = server.ListenAndServe()
		if err != nil {
			log.Fatal(err.Error())
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)

	log.Println("System up")

	<-stop

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		log.Println("(Shutdown) Errors shutting server down: " + err.Error())
	}

}
