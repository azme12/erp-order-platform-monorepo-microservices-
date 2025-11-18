package router

import (
	"io"
	"microservice-challenge/gateway/client"
	"microservice-challenge/package/config"
	"microservice-challenge/package/log"
	"microservice-challenge/package/middleware"
	"net/http"
	"strings"

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

func NewRouter(cfg *config.Config, logger log.Logger) http.Handler {
	r := chi.NewRouter()

	r.Use(chiMiddleware.RequestID)
	r.Use(chiMiddleware.RealIP)
	r.Use(chiMiddleware.Logger)
	r.Use(chiMiddleware.Recoverer)
	r.Use(chiMiddleware.AllowContentType("application/json"))

	httpClient := client.NewClient(logger)
	authMiddleware := middleware.NewAuthMiddleware(cfg.JWT.Secret, logger)

	router := &Router{
		client:         httpClient,
		authMiddleware: authMiddleware,
		config:         cfg,
		logger:         logger,
	}

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	r.Route("/api", func(r chi.Router) {
		r.Route("/auth", func(r chi.Router) {
			r.Post("/register", router.forwardToService("auth", "/register"))
			r.Post("/login", router.forwardToService("auth", "/login"))
			r.Post("/forgot-password", router.forwardToService("auth", "/forgot-password"))
			r.Post("/reset-password", router.forwardToService("auth", "/reset-password"))
			r.Post("/service-token", router.forwardToService("auth", "/service-token"))
		})

		r.Group(func(r chi.Router) {
			r.Use(authMiddleware.ValidateToken)

			r.Route("/customers", func(r chi.Router) {
				r.Get("/", router.forwardToService("contact", "/customers"))
				r.Get("/{id}", router.forwardToService("contact", "/customers/{id}"))
				r.Post("/", router.forwardToService("contact", "/customers"))
				r.Put("/{id}", router.forwardToService("contact", "/customers/{id}"))
				r.Delete("/{id}", router.forwardToService("contact", "/customers/{id}"))
			})

			r.Route("/vendors", func(r chi.Router) {
				r.Get("/", router.forwardToService("contact", "/vendors"))
				r.Get("/{id}", router.forwardToService("contact", "/vendors/{id}"))
				r.Post("/", router.forwardToService("contact", "/vendors"))
				r.Put("/{id}", router.forwardToService("contact", "/vendors/{id}"))
				r.Delete("/{id}", router.forwardToService("contact", "/vendors/{id}"))
			})

			r.Route("/items", func(r chi.Router) {
				r.Get("/", router.forwardToService("inventory", "/items"))
				r.Get("/{id}", router.forwardToService("inventory", "/items/{id}"))
				r.Post("/", router.forwardToService("inventory", "/items"))
				r.Put("/{id}", router.forwardToService("inventory", "/items/{id}"))
				r.Delete("/{id}", router.forwardToService("inventory", "/items/{id}"))
				r.Get("/{item_id}/stock", router.forwardToService("inventory", "/items/{item_id}/stock"))
				r.Put("/{item_id}/stock", router.forwardToService("inventory", "/items/{item_id}/stock"))
			})

			r.Route("/sales/orders", func(r chi.Router) {
				r.Get("/", router.forwardToService("sales", "/orders"))
				r.Get("/{id}", router.forwardToService("sales", "/orders/{id}"))
				r.Post("/", router.forwardToService("sales", "/orders"))
				r.Put("/{id}", router.forwardToService("sales", "/orders/{id}"))
				r.Post("/{id}/confirm", router.forwardToService("sales", "/orders/{id}/confirm"))
				r.Post("/{id}/pay", router.forwardToService("sales", "/orders/{id}/pay"))
			})

			r.Route("/purchase/orders", func(r chi.Router) {
				r.Get("/", router.forwardToService("purchase", "/orders"))
				r.Get("/{id}", router.forwardToService("purchase", "/orders/{id}"))
				r.Post("/", router.forwardToService("purchase", "/orders"))
				r.Put("/{id}", router.forwardToService("purchase", "/orders/{id}"))
				r.Post("/{id}/receive", router.forwardToService("purchase", "/orders/{id}/receive"))
				r.Post("/{id}/pay", router.forwardToService("purchase", "/orders/{id}/pay"))
			})
		})
	})

	return r
}

func (rt *Router) forwardToService(serviceName string, path string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var baseURL string
		switch serviceName {
		case "auth":
			baseURL = rt.config.Services.Auth.URL
		case "contact":
			baseURL = rt.config.Services.Contact.URL
		case "inventory":
			baseURL = rt.config.Services.Inventory.URL
		case "sales":
			baseURL = rt.config.Services.Sales.URL
		case "purchase":
			baseURL = rt.config.Services.Purchase.URL
		default:
			http.Error(w, "Unknown service", http.StatusBadGateway)
			return
		}

		targetURL := baseURL + path

		// Replace all path parameters dynamically
		// Support multiple parameter names: id, item_id, order_id, etc.
		ctx := r.Context()
		if strings.Contains(path, "{id}") {
			id := chi.URLParamFromCtx(ctx, "id")
			targetURL = strings.ReplaceAll(targetURL, "{id}", id)
		}
		if strings.Contains(path, "{item_id}") {
			itemID := chi.URLParamFromCtx(ctx, "item_id")
			targetURL = strings.ReplaceAll(targetURL, "{item_id}", itemID)
		}
		if strings.Contains(path, "{order_id}") {
			orderID := chi.URLParamFromCtx(ctx, "order_id")
			targetURL = strings.ReplaceAll(targetURL, "{order_id}", orderID)
		}
		// Generic replacement for any remaining {param} patterns
		for {
			start := strings.Index(targetURL, "{")
			if start == -1 {
				break
			}
			end := strings.Index(targetURL[start:], "}")
			if end == -1 {
				break
			}
			paramName := targetURL[start+1 : start+end]
			paramValue := chi.URLParamFromCtx(ctx, paramName)
			targetURL = targetURL[:start] + paramValue + targetURL[start+end+1:]
		}

		resp, err := rt.client.ForwardRequest(r.Context(), targetURL, r)
		if err != nil {
			rt.logger.Error(r.Context(), "failed to forward request", zap.String("service", serviceName), zap.String("url", targetURL), zap.Error(err))
			http.Error(w, "Service unavailable", http.StatusBadGateway)
			return
		}
		defer resp.Body.Close()

		for key, values := range resp.Header {
			for _, value := range values {
				w.Header().Add(key, value)
			}
		}
		w.WriteHeader(resp.StatusCode)
		io.Copy(w, resp.Body)
	}
}
