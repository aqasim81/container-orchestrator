// Package main is the entry point for the container orchestrator API server.
package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/rs/zerolog"

	"github.com/github-builder/container-orchestrator/internal/api"
	"github.com/github-builder/container-orchestrator/internal/config"
	"github.com/github-builder/container-orchestrator/internal/store"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

func run() error {
	// Load configuration.
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("loading config: %w", err)
	}

	// Initialize logger.
	level, err := zerolog.ParseLevel(cfg.LogLevel)
	if err != nil {
		return fmt.Errorf("parsing log level: %w", err)
	}
	logger := zerolog.New(os.Stdout).
		With().
		Timestamp().
		Logger().
		Level(level)

	// Open store.
	s, err := store.NewBoltStore(cfg.DataDir)
	if err != nil {
		return fmt.Errorf("opening store: %w", err)
	}
	defer func() {
		if closeErr := s.Close(); closeErr != nil {
			logger.Error().Err(closeErr).Msg("failed to close store")
		}
	}()

	// Create router.
	router := api.NewRouter(&api.RouterConfig{
		Store:        s,
		Logger:       logger,
		DashboardURL: cfg.DashboardURL,
		APIKey:       cfg.APIKey,
	})

	// Start HTTP server.
	srv := &http.Server{
		Addr:              fmt.Sprintf(":%d", cfg.Port),
		Handler:           router,
		ReadHeaderTimeout: 10 * time.Second,
	}

	// Graceful shutdown on SIGINT/SIGTERM.
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	go func() {
		logger.Info().
			Int("port", cfg.Port).
			Str("data_dir", cfg.DataDir).
			Msg("starting API server")

		if listenErr := srv.ListenAndServe(); listenErr != nil && !errors.Is(listenErr, http.ErrServerClosed) {
			logger.Fatal().Err(listenErr).Msg("server failed")
		}
	}()

	// Wait for shutdown signal.
	<-ctx.Done()
	logger.Info().Msg("shutting down server")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		return fmt.Errorf("server shutdown: %w", err)
	}

	logger.Info().Msg("server stopped")
	return nil
}
