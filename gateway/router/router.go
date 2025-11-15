package router

import (
	"io"
	"microservice-challenge/gateway/client"
	"microservice-challenge/gateway/middleware"
	"microservice-challenge/package/config"
	"microservice-challenge/package/log"
	"net/http"

	"github.com/go-chi/chi/v5"
	chiMiddleware "github.com/go-chi/chi/v5/middleware"
	"go.uber.org/zap"
)

type Router struct {
	client         *client.Client
	authMiddleware *middleware.AuthMiddleware
	config         *config.Config
	logger         log.Logger
}

// NewRouter creates a new API Gateway router
func NewRouter(cfg *config.Config, logger log.Logger) http.Handler {
	r := chi.NewRouter()

	// Middleware
	r.Use(chiMiddleware.RequestID)
	r.Use(chiMiddleware.RealIP)
	r.Use(chiMiddleware.Logger)
	r.Use(chiMiddleware.Recoverer)
	r.Use(chiMiddleware.AllowContentType("application/json"))

	// Initialize client and middleware
	httpClient := client.NewClient(logger)
	authMiddleware := middleware.NewAuthMiddleware(cfg.Services.Auth.URL+"/validate", logger)

	router := &Router{
		client:         httpClient,
		authMiddleware: authMiddleware,
		config:         cfg,
		logger:         logger,
	}

	// Health check
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	// API routes with authentication middleware
	r.Route("/api", func(r chi.Router) {
		// Apply auth middleware to all routes except register/login
		r.Use(authMiddleware.Middleware)

		// Auth Service routes
		r.Route("/auth", func(r chi.Router) {
			r.Post("/register", router.forwardToService("auth", "/register"))
			r.Post("/login", router.forwardToService("auth", "/login"))
			r.Post("/validate", router.forwardToService("auth", "/validate"))
			r.Post("/service-token", router.forwardToService("auth", "/service-token"))
		})

		// Contact Service routes (will be implemented in Day 3)
		r.Route("/contacts", func(r chi.Router) {
			r.Get("/*", router.forwardToService("contact", ""))
			r.Post("/*", router.forwardToService("contact", ""))
			r.Put("/*", router.forwardToService("contact", ""))
			r.Delete("/*", router.forwardToService("contact", ""))
		})

		// Inventory Service routes (will be implemented in Day 4)
		r.Route("/inventory", func(r chi.Router) {
			r.Get("/*", router.forwardToService("inventory", ""))
			r.Post("/*", router.forwardToService("inventory", ""))
			r.Put("/*", router.forwardToService("inventory", ""))
			r.Delete("/*", router.forwardToService("inventory", ""))
		})

		// Sales Service routes (will be implemented in Day 5)
		r.Route("/sales", func(r chi.Router) {
			r.Get("/*", router.forwardToService("sales", ""))
			r.Post("/*", router.forwardToService("sales", ""))
			r.Put("/*", router.forwardToService("sales", ""))
			r.Delete("/*", router.forwardToService("sales", ""))
		})

		// Purchase Service routes (will be implemented in Day 5)
		r.Route("/purchase", func(r chi.Router) {
			r.Get("/*", router.forwardToService("purchase", ""))
			r.Post("/*", router.forwardToService("purchase", ""))
			r.Put("/*", router.forwardToService("purchase", ""))
			r.Delete("/*", router.forwardToService("purchase", ""))
		})
	})

	return r
}

// forwardToService forwards a request to the specified service
func (rt *Router) forwardToService(serviceName, path string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		// Get service URL from config
		var serviceURL string
		switch serviceName {
		case "auth":
			serviceURL = rt.config.Services.Auth.URL
		case "contact":
			serviceURL = rt.config.Services.Contact.URL
		case "inventory":
			serviceURL = rt.config.Services.Inventory.URL
		case "sales":
			serviceURL = rt.config.Services.Sales.URL
		case "purchase":
			serviceURL = rt.config.Services.Purchase.URL
		default:
			http.Error(w, "Service not found", http.StatusNotFound)
			return
		}

		// Build target URL
		targetPath := path
		if targetPath == "" {
			// Extract path after /api/{service}/
			requestPath := r.URL.Path
			// Remove /api/{service} prefix
			prefix := "/api/" + serviceName
			if len(requestPath) > len(prefix) {
				targetPath = requestPath[len(prefix):]
			}
		}
		targetURL := serviceURL + targetPath

		// Add query parameters
		if r.URL.RawQuery != "" {
			targetURL += "?" + r.URL.RawQuery
		}

		rt.logger.Info(ctx, "forwarding request", zap.String("service", serviceName), zap.String("target", targetURL))

		// Forward request
		resp, err := rt.client.ForwardRequest(ctx, targetURL, r)
		if err != nil {
			rt.logger.Error(ctx, "failed to forward request", zap.Error(err))
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		defer resp.Body.Close()

		// Copy response headers
		for key, values := range resp.Header {
			for _, value := range values {
				w.Header().Add(key, value)
			}
		}

		// Copy status code
		w.WriteHeader(resp.StatusCode)

		// Copy response body
		_, err = io.Copy(w, resp.Body)
		if err != nil {
			rt.logger.Error(ctx, "failed to copy response body", zap.Error(err))
		}
	}
}
