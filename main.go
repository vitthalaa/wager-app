package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	env "github.com/joho/godotenv"

	"github.com/vitthalaa/wager-app/internal/config"
	"github.com/vitthalaa/wager-app/internal/db"
	"github.com/vitthalaa/wager-app/internal/handlers"
	"github.com/vitthalaa/wager-app/internal/repo"
	"github.com/vitthalaa/wager-app/internal/services"
)

const envFile = ".env"

// Changing to overload to override from .env file(Should be load from env in prod)
var loadEnv = env.Overload

func run() (s *http.Server) {
	err := loadEnv(envFile)
	if err != nil {
		log.Fatal(err)
	}

	// Load config
	conf := config.GetAppConfig()
	if conf.Port == 0 {
		log.Fatal("no port specified")
	}

	// Connect DB
	conn, err := db.OpenConnection(&conf.DataBaseConfig)
	if err != nil {
		log.Fatalf("DB connection error: %s\n", err)
	}

	// Init Repos
	wagerRepo := repo.NewWagerRepo(conn)
	purchaseRepo := repo.NewPurchaseRepo(conn)

	// Init Services
	wagerService := services.NewWagerService(wagerRepo)
	purchaseService := services.NewPurchaseService(purchaseRepo, wagerRepo)

	// Init handlers
	wagerHandler := handlers.NewWagersHandler(wagerService)
	purchasehandler := handlers.NewPurchasesHandler(purchaseService)

	mux := http.NewServeMux()
	mux.HandleFunc("/wagers", wagerHandler.Handle)
	mux.HandleFunc("/buy/", purchasehandler.Handle)

	address := fmt.Sprintf(":%d", conf.Port)
	s = &http.Server{
		Addr:         address,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		Handler:      mux,
	}

	go func() {
		log.Printf("Starting HTTP listener on: %s", address)
		if err := s.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("error listening on port: %s\n", err)
		}
	}()

	return
}

func main() {
	s := run()
	quit := make(chan os.Signal)
	// kill (no param) default send syscall.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall.SIGKILL but can't be catch, so don't need add it
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("shutting down server...")

	// The context is used to inform the server it has 5 seconds to finish
	// the request it is currently handling
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer func() {
		cancel()
	}()
	if err := s.Shutdown(ctx); err != nil {
		log.Fatal("server forced to shut down")
	}
	log.Println("server exiting")
}
