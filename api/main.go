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
		log.Fatal(err)
	}
}

func run() error {
	upstreamBaseURL := os.Getenv("BASE_URL")
	if upstreamBaseURL == "" {
		return fmt.Errorf("BASE_URL is not set")
	}
	getHandler, err := handlers.NewGetHandler(upstreamBaseURL)
	if err != nil {
		return fmt.Errorf("configure GET handlers: %w", err)
	}

	if err := repository.InitDB(context.Background()); err != nil {
		return fmt.Errorf("initialize database: %w", err)
	}
	defer repository.DB.Close()

	server := &http.Server{
		Addr:              ":8081",
		Handler:           getHandler,
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
			return fmt.Errorf("HTTP server stopped: %w", err)
		}
		return nil
	case <-shutdownContext.Done():
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		_ = server.Close()
		return fmt.Errorf("shut down HTTP server: %w", err)
	}

	if err := <-serverErrors; err != nil && !errors.Is(err, http.ErrServerClosed) {
		return fmt.Errorf("HTTP server stopped: %w", err)
	}
	return nil
}
