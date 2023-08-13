package main

import (
	"context"
	"fmt"
	"github.com/EricBastos/ProjetoTG/API/configs"
	"github.com/EricBastos/ProjetoTG/API/internal/grpcClient"
	"github.com/EricBastos/ProjetoTG/API/internal/infra/rabbitmqClient"
	"github.com/EricBastos/ProjetoTG/API/internal/infra/webserver/handlers"
	"github.com/EricBastos/ProjetoTG/API/internal/infra/webserver/middlewares"
	"github.com/EricBastos/ProjetoTG/Library/database"
	"github.com/EricBastos/ProjetoTG/Library/entities"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
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
	if err != nil {
		log.Println("Waiting rabbitmq:", err.Error())
		time.Sleep(3 * time.Second)
		os.Exit(1)
	}

	err = grpcClient.InitializeServices()
	if err != nil {
		log.Fatal(err.Error())
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
		&entities.User{},
		&entities.StaticDeposit{},
		&entities.BurnOp{},
	)
	if err != nil {
		log.Fatal(err.Error())
	}

	userDb := database.NewUserDB(db)
	staticDepositDb := database.NewStaticDepositDB(db)
	burnOpsDb := database.NewBurnOperationsDB(db)

	userHandler := handlers.NewUserHandler(userDb, staticDepositDb, burnOpsDb)
	depositHandler := handlers.NewDepositHandler(staticDepositDb)
	withdrawHandler := handlers.NewWithdrawHandler(burnOpsDb, rabbitClient)

	userGeneralAuthenticator := middlewares.NewAuthenticator("USER", userDb)

	r := chi.NewRouter()

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: true,
	}))

	r.Route("/api", func(r chi.Router) {

		r.Use(middleware.Logger)
		r.Use(middleware.Recoverer)

		r.Route("/user", func(r chi.Router) {

			// Unauthenticated
			r.Group(func(r chi.Router) {
				r.Group(func(r chi.Router) {
					r.Post("/create", userHandler.Create)
				})
				r.Group(func(r chi.Router) {
					r.Post("/login", userHandler.Login)
				})
			})

			// Authenticated
			r.Group(func(r chi.Router) {

				r.Group(func(r chi.Router) {

					r.Use(userGeneralAuthenticator.Authenticate("jwt"))

					r.Group(func(r chi.Router) {
						r.Get("/info", userHandler.GetUser)
					})

					r.Route("/mint", func(r chi.Router) {
						r.Route("/static-pix", func(r chi.Router) {
							r.Group(func(r chi.Router) {
								r.Post("/", depositHandler.CreatePixDeposit)
							})
							r.Group(func(r chi.Router) {
								r.Get("/history", userHandler.GetStaticDepositLogs)
							})
						})
					})

					r.Route("/burn", func(r chi.Router) {
						r.Group(func(r chi.Router) {
							r.Post("/", withdrawHandler.CreateUserWithdraw)
						})
						r.Group(func(r chi.Router) {
							r.Get("/history", userHandler.GetTransfersLogs)
						})
					})

				})

			})

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
