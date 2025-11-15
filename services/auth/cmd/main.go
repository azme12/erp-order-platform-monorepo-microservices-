package main

import (
	"context"
	"microservice-challenge/package/config"
	"microservice-challenge/package/database"
	"microservice-challenge/package/log"
	"microservice-challenge/services/auth/httphandler"
	"microservice-challenge/services/auth/router"
	"microservice-challenge/services/auth/storage/postgresql"
	"microservice-challenge/services/auth/usecase/auth"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go.uber.org/zap"
)

// @title           Auth Service API
// @version         1.0
// @description     Authentication and Authorization Service for Microservices Challenge
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

	logger := log.Init("auth-service")
	logger.Info(ctx, "starting auth service")

	cfg, err := config.LoadConfig()
	if err != nil {
		logger.Fatal(ctx, "failed to load config", zap.Error(err))
	}

	dbConnStr := database.BuildConnectionString(
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.Name,
	)

	db, err := database.NewPostgresConnection(dbConnStr)
	if err != nil {
		logger.Fatal(ctx, "failed to connect to database", zap.Error(err))
	}
	defer db.Close()

	logger.Info(ctx, "connected to database")

	storage := postgresql.NewStorage(db)

	usecase := auth.NewUsecase(
		storage,
		cfg.JWT.Secret,
		cfg.JWT.UserExpirationHours,
	)

	handler := httphandler.NewHandler(usecase, logger)
	r := router.NewRouter(handler, logger)

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

	logger.Info(ctx, "server exited")
}
