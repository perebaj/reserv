// Package main is the entry point for the API.
package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/clerk/clerk-sdk-go/v2"
	"github.com/perebaj/reserv"
	"github.com/perebaj/reserv/handler"
)

// config gather all important configuration for the API.
type config struct {
	// ClerkKey is the private key for the Clerk API.
	ClerkKey string
}

func main() {
	logger, err := reserv.NewLoggerSlog(reserv.ConfigSlog{
		Level:  reserv.LevelInfo,
		Format: reserv.FormatJSON,
	})

	slog.SetDefault(logger)

	if err != nil {
		slog.Error("failed to create logger", "err", err)
		os.Exit(1)
	}

	cfg := config{
		ClerkKey: getEnvVar("CLERK_KEY", ""),
	}

	if cfg.ClerkKey == "" {
		slog.Error("Clerk key is not set, unable to start the server")
		os.Exit(1)
	}

	clerk.SetKey(cfg.ClerkKey)
	mux := handler.Router()

	srv := &http.Server{
		Addr:    ":8000",
		Handler: mux,
	}
	slog.Info("starting server", "addr", srv.Addr)

	// serverErrors is a channel to receive errors from the server.
	// It is buffered to avoid blocking the goroutine that starts the server.
	// This allows us to handle server errors without blocking the main goroutine.
	serverErrors := make(chan error, 1)
	go func() {
		serverErrors <- srv.ListenAndServe()
	}()

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM, syscall.SIGINT, syscall.SIGHUP)

	select {
	case err := <-serverErrors:
		slog.Error("server error", "err", err)
	case <-shutdown:
		slog.Info("shutting down server")
		// Create a deadline to wait for a graceful shutdown.
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		// Shutdown is a built-in method for http.Server. It stops the server from accepting new connections using a deadline context.
		if err := srv.Shutdown(ctx); err != nil {
			slog.Error("Graceful shutdown failed", "err", err)

			if err := srv.Close(); err != nil {
				slog.Error("Forcing server shutdown", "err", err)
			}
		}
	}
}

// getEnvVar returns the value of an environment variable.
// If the variable is not set, it returns the default value.
func getEnvVar(key string, defaultValue string) string {
	envVar := os.Getenv(key)
	if envVar == "" {
		return defaultValue
	}
	return envVar
}
