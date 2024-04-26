package main

import (
	"context"
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
)

func main() {
	cfg := config.New()

	runGatewayServer(cfg)
}

func runGatewayServer(
	cfg *config.Config,
) {
	server := api.NewServer(cfg, nil)

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

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Println("Failed to shutdown server gracefully:", err)
	}

	<-ctx.Done()
	log.Println("Server shutdown gracefully")
}
