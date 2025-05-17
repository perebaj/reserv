// Package main ...
package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/clerk/clerk-sdk-go/v2"
	"github.com/cloudflare/cloudflare-go"
	"github.com/perebaj/reserv"
	"github.com/perebaj/reserv/handler"
	"github.com/perebaj/reserv/postgres"
)

// Config gathers all the configuration for the application.
type Config struct {
	// PostgresURL Represents the whole connection string for the database.
	PostgresURL string
	// LogLevel is the level of the logs. Avaiable values are: info, debug, error.
	LogLevel string
	// LogFormat is the format of the log. Available values are: json, logfmt, gcp
	LogFormat string
	// CloudFlareAPIKey is the API key for the Cloudflare API.
	CloudFlareAPIKey string
	// ClerkAPIKey is the private key for the Clerk API.
	ClerkAPIKey string
}

func main() {
	cfg := Config{
		// POSTGRES_URL Represents the whole connection string for the database.
		// TODO(@perebaj): Remove this default value when we have a proper configuration.
		PostgresURL:      getEnvWithDefault("POSTGRES_URL", "postgres://postgres:postgres@localhost:5432/postgres"),
		LogLevel:         getEnvWithDefault("LOG_LEVEL", "debug"),
		LogFormat:        getEnvWithDefault("LOG_FORMAT", "json"),
		CloudFlareAPIKey: getEnvWithDefault("CLOUDFLARE_API_KEY", ""),
		// ClerkAPIKey is the private key for the Clerk API.
		ClerkAPIKey: getEnvWithDefault("CLERK_API_KEY", ""),
	}

	if cfg.PostgresURL == "" || cfg.CloudFlareAPIKey == "" || cfg.ClerkAPIKey == "" {
		slog.Error("POSTGRES_URL or CLOUDFLARE_API_KEY or CLERK_API_KEY is not set")
		os.Exit(1)
	}
	clerk.SetKey(cfg.ClerkAPIKey)

	logger, err := reserv.NewLoggerSlog(
		reserv.ConfigSlog{
			Level:  cfg.LogLevel,
			Format: cfg.LogFormat,
		},
	)
	if err != nil {
		slog.Error("failed to create logger", "error", err)
		os.Exit(1)
	}

	slog.SetDefault(logger)

	db, err := postgres.OpenDB(postgres.Config{
		URL: cfg.PostgresURL,
	})
	if err != nil {
		slog.Error("failed to open db", "error", err)
		os.Exit(1)
	}

	err = postgres.Migrate(db.DB)
	if err != nil {
		slog.Error("failed to migrate db", "error", err)
		os.Exit(1)
	}

	repo := postgres.NewRepository(db)

	cloudFlareClient, err := cloudflare.NewWithAPIToken(cfg.CloudFlareAPIKey)
	if err != nil {
		slog.Error("failed to create cloudflare client", "error", err)
		os.Exit(1)
	}

	// TODO(@perebaj): Duplicating the repo object to turn easy on testing. But this is not a ideal solution.
	handler := handler.NewHandler(repo, cloudFlareClient, repo)

	mux := http.NewServeMux()
	handler.RegisterRoutes(mux)

	port := getEnvWithDefault("PORT", "")
	if port == "" {
		slog.Error("PORT is not set")
		os.Exit(1)
	}

	intPort, err := strconv.Atoi(port)
	if err != nil {
		slog.Error("failed to convert port to int", "error", err)
		os.Exit(1)
	}

	// cors is a middleware that adds the necessary headers to the response.
	// TODO(@perebaj): Remove this function when we have a proper CORS implementation.
	cors := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "X-CSRF-Token, X-Requested-With, Accept, Accept-Version, Content-Length, Content-MD5, Content-Type, Date, X-Api-Version, Authorization")

			if r.Method == "OPTIONS" {
				w.WriteHeader(http.StatusOK)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", intPort),
		Handler: cors(mux),
	}

	slog.Info("starting server", "address", srv.Addr)
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

func getEnvWithDefault(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
