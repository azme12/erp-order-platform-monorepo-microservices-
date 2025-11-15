package main

import (
	"context"
	"microservice-challenge/gateway/router"
	"microservice-challenge/package/config"
	"microservice-challenge/package/log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go.uber.org/zap"
)

func main() {
	ctx := context.Background()

	// Initialize logger
	logger := log.Init("gateway")
	logger.Info(ctx, "starting API Gateway")

	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		logger.Fatal(ctx, "failed to load config", zap.Error(err))
	}

	// Initialize router
	r := router.NewRouter(cfg, logger)

	// Get port from environment or use default
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Create HTTP server
	srv := &http.Server{
		Addr:    ":" + port,
		Handler: r,
	}

	// Start server in goroutine
	go func() {
		logger.Info(ctx, "starting HTTP server", zap.String("port", port))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal(ctx, "failed to start server", zap.Error(err))
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info(ctx, "shutting down server")

	// Graceful shutdown
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		logger.Error(ctx, "server forced to shutdown", zap.Error(err))
	}

	logger.Info(ctx, "server exited")
}
