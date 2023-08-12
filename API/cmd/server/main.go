package main

import (
	"context"
	"fmt"
	"github.com/EricBastos/ProjetoTG/API/configs"
	"github.com/EricBastos/ProjetoTG/API/internal/infra/webserver/handlers"
	"github.com/EricBastos/ProjetoTG/API/internal/infra/webserver/middlewares"
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
	)
	if err != nil {
		log.Fatal(err.Error())
	}

	userDb := database.NewUserDB(db)
	staticDepositDb := database.NewStaticDepositDB(db)

	userHandler := handlers.NewUserHandler(userDb, staticDepositDb)

	userGeneralAuthenticator := middlewares.NewAuthenticator("USER", userDb)

	r := chi.NewRouter()

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

					//r.Route("/mint", func(r chi.Router) {
					//
					//	r.Route("/static-pix", func(r chi.Router) {
					//		r.Group(func(r chi.Router) {
					//			r.Use(userGeneralAuthenticator.OnlyKycVerified())
					//			r.Use(middlewares.NewPromHttpMiddleware("/user_deposit_static_pix"))
					//			r.Post("/", depositHandler.CreateUserStaticPixDeposit)
					//		})
					//		r.Group(func(r chi.Router) {
					//			r.Use(middlewares.NewPromHttpMiddleware("/user_deposit_history"))
					//			r.Get("/history", userHandler.GetStaticDepositLogs)
					//		})
					//	})
					//
					//	r.Route("/pix-to-usd", func(r chi.Router) {
					//		r.Group(func(r chi.Router) {
					//			r.Use(userGeneralAuthenticator.OnlyKycVerified())
					//			r.Use(middlewares.NewPromHttpMiddleware("/user_deposit_pix_to_usd"))
					//			r.Post("/", pixToUsdHandler.GetUserWebsocketToken)
					//		})
					//		r.Group(func(r chi.Router) {
					//			r.Use(middlewares.NewPromHttpMiddleware("/user_deposit_pix_to_usd_history"))
					//			r.Get("/history", userHandler.GetPixToUsdDepositLogs)
					//		})
					//	})
					//
					//	//r.Post("/dynamic-pix", depositHandler.CreateDynamicPixDeposit)
					//	//r.Post("/boleto", depositHandler.CreateBoletoDeposit)
					//})
					//
					//r.Route("/burn", func(r chi.Router) {
					//
					//	r.Group(func(r chi.Router) {
					//		r.Use(middlewares.NewPromHttpMiddleware("/user_transfers_history"))
					//		r.Get("/history", userHandler.GetTransfersLogs)
					//	})
					//
					//	r.Group(func(r chi.Router) {
					//		r.Use(userGeneralAuthenticator.OnlyKycVerified())
					//		r.Use(middlewares.NewPromHttpMiddleware("/user_withdraw"))
					//		r.Post("/", withdrawHandler.CreateUserWithdraw)
					//	})
					//})

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
