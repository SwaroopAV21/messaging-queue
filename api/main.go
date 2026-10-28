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

	"github.com/SwaroopAV21/messaging-queue/handlers"
	"github.com/SwaroopAV21/messaging-queue/repository"
)

func main() {
	if err := run(); err != nil {
		fmt.Println("application error:", err)
		log.Fatal(err)
	}
}

func run() error {
	upstreamBaseURL := os.Getenv("BASE_URL")
	if upstreamBaseURL == "" {
		fmt.Println("error: BASE_URL is not set")
		return fmt.Errorf("BASE_URL is not set")
	}
	getHandler, err := handlers.NewGetHandler(upstreamBaseURL)
	if err != nil {
		fmt.Println("configure GET handlers:", err)
		return fmt.Errorf("configure GET handlers: %w", err)
	}

	if err := repository.InitDB(context.Background()); err != nil {
		fmt.Println("initialize database:", err)
		return fmt.Errorf("initialize database: %w", err)
	}
	defer repository.DB.Close()

	corsOrigin := os.Getenv("CORS_ORIGIN")
	if corsOrigin == "" {
		corsOrigin = "http://localhost:5173"
	}

	server := &http.Server{
		Addr:              ":8081",
		Handler:           corsMiddleware(getHandler, corsOrigin),
		ReadHeaderTimeout: 5 * time.Second,
	}

	serverErrors := make(chan error, 1)
	go func() {
		serverErrors <- server.ListenAndServe()
	}()

	fmt.Println("Starting message brokers http middlerware server")

	shutdownContext, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	select {
	case err := <-serverErrors:
		if !errors.Is(err, http.ErrServerClosed) {
			fmt.Println("HTTP server stopped:", err)
			return fmt.Errorf("HTTP server stopped: %w", err)
		}
		return nil
	case <-shutdownContext.Done():
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		fmt.Println("shut down HTTP server:", err)
		_ = server.Close()
		return fmt.Errorf("shut down HTTP server: %w", err)
	}

	if err := <-serverErrors; err != nil && !errors.Is(err, http.ErrServerClosed) {
		fmt.Println("HTTP server stopped:", err)
		return fmt.Errorf("HTTP server stopped: %w", err)
	}
	return nil
}
