package main

import (
	"context"
	"microservice-challenge/package/config"
	"microservice-challenge/package/database"
	"microservice-challenge/package/log"
	"microservice-challenge/package/nats"
	"microservice-challenge/services/contact/httphandler"
	"microservice-challenge/services/contact/router"
	"microservice-challenge/services/contact/storage/postgresql"
	"microservice-challenge/services/contact/usecase/contact"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go.uber.org/zap"
)

// @title           Contact Service API
// @version         1.0
// @description     Contact Service for managing Customers and Vendors
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    http://www.example.com/support
// @contact.email  support@example.com

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8000
// @BasePath  /

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token. Example: "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	logger := log.Init("contact-service")
	logger.Info(ctx, "starting contact service")

	cfg, err := config.LoadConfig()
	if err != nil {
		logger.Fatal(ctx, "failed to load config", zap.Error(err))
	}

	dbHost := os.Getenv("DB_HOST")
	if dbHost == "" {
		dbHost = "localhost"
	}
	dbPort := os.Getenv("DB_PORT")
	if dbPort == "" {
		dbPort = "5432"
	}
	dbUser := os.Getenv("DB_USER")
	if dbUser == "" {
		dbUser = "microservice"
	}
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	if dbName == "" {
		dbName = "microservice"
	}

	if dbPassword == "" {
		logger.Fatal(ctx, "DB_PASSWORD environment variable is required")
	}

	dbConnStr := database.BuildConnectionString(
		dbHost,
		dbPort,
		dbUser,
		dbPassword,
		dbName,
	)

	db, err := database.NewPostgresConnection(dbConnStr)
	if err != nil {
		logger.Fatal(ctx, "failed to connect to database", zap.Error(err))
	}
	defer db.Close()

	logger.Info(ctx, "connected to database")

	natsClient, err := nats.NewClient(cfg.NATS.URL)
	if err != nil {
		logger.Fatal(ctx, "failed to connect to NATS", zap.Error(err))
	}
	defer natsClient.Close()

	logger.Info(ctx, "connected to NATS")

	storage := postgresql.NewStorage(db)

	usecase := contact.NewUsecase(storage, natsClient, logger)

	handler := httphandler.NewHandler(usecase, logger)
	r := router.NewRouter(handler, logger, cfg.JWT.Secret, db)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}

	srv := &http.Server{
		Addr:         ":" + port,
		Handler:      r,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	serverErrors := make(chan error, 1)

	go func() {
		logger.Info(ctx, "starting HTTP server", zap.String("port", port))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			serverErrors <- err
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	select {
	case err := <-serverErrors:
		logger.Fatal(ctx, "server error", zap.Error(err))
	case sig := <-quit:
		logger.Info(ctx, "shutdown signal received", zap.String("signal", sig.String()))
	}

	logger.Info(ctx, "shutting down server gracefully")

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer shutdownCancel()

	cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		logger.Error(ctx, "server forced to shutdown", zap.Error(err))
		if closeErr := srv.Close(); closeErr != nil {
			logger.Error(ctx, "error closing server", zap.Error(closeErr))
		}
	} else {
		logger.Info(ctx, "server shutdown gracefully")
	}

	if err := db.Close(); err != nil {
		logger.Error(ctx, "error closing database", zap.Error(err))
	}

	natsClient.Close()

	logger.Info(ctx, "server exited")
}
