package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/danyouknowme/assessment-tax/api"
	"github.com/danyouknowme/assessment-tax/config"
	"github.com/danyouknowme/assessment-tax/db"

	_ "github.com/lib/pq"
)

func main() {
	cfg := config.New()

	conn, err := sql.Open("postgres", cfg.DatabaseUrl)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer conn.Close()

	if err := db.PrepareDatabase(conn); err != nil {
		log.Fatalf("Failed to prepare database: %v", err)
	}

	store := db.NewStore(conn)

	runGatewayServer(cfg, store)
}

func runGatewayServer(cfg *config.Config, store db.Store) {
	server := api.NewServer(cfg, store)

	log.Println("Start listening for HTTP requests...")
	go func() {
		log.Printf("Server listening on port %s\n", cfg.Port)
		if err := server.Start(fmt.Sprintf(":%s", cfg.Port)); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	<-shutdown
	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Println("Failed to shutdown server gracefully:", err)
	}

	<-ctx.Done()
	log.Println("Server shutdown gracefully")
}
