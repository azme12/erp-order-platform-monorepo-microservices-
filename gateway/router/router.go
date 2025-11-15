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
		})

		r.Group(func(r chi.Router) {
			r.Use(authMiddleware.ValidateToken)

			r.Route("/customers", func(r chi.Router) {
				r.Get("/*", router.forwardToService("contact", ""))
				r.Post("/*", router.forwardToService("contact", ""))
				r.Put("/*", router.forwardToService("contact", ""))
				r.Delete("/*", router.forwardToService("contact", ""))
			})
			r.Route("/vendors", func(r chi.Router) {
				r.Get("/*", router.forwardToService("contact", ""))
				r.Post("/*", router.forwardToService("contact", ""))
				r.Put("/*", router.forwardToService("contact", ""))
				r.Delete("/*", router.forwardToService("contact", ""))
			})

			r.Route("/inventory", func(r chi.Router) {
				r.Get("/*", router.forwardToService("inventory", ""))
				r.Post("/*", router.forwardToService("inventory", ""))
				r.Put("/*", router.forwardToService("inventory", ""))
				r.Delete("/*", router.forwardToService("inventory", ""))
			})

			r.Route("/sales", func(r chi.Router) {
				r.Get("/*", router.forwardToService("sales", ""))
				r.Post("/*", router.forwardToService("sales", ""))
				r.Put("/*", router.forwardToService("sales", ""))
				r.Delete("/*", router.forwardToService("sales", ""))
			})

			r.Route("/purchase", func(r chi.Router) {
				r.Get("/*", router.forwardToService("purchase", ""))
				r.Post("/*", router.forwardToService("purchase", ""))
				r.Put("/*", router.forwardToService("purchase", ""))
				r.Delete("/*", router.forwardToService("purchase", ""))
			})
		})
	})

	return r
}

func (rt *Router) forwardToService(serviceName, path string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

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

		targetPath := path
		if targetPath == "" {
			requestPath := r.URL.Path
			if strings.HasPrefix(requestPath, "/api/") {
				pathAfterAPI := requestPath[5:]
				servicePrefix := serviceName + "/"
				if strings.HasPrefix(pathAfterAPI, servicePrefix) {
					targetPath = "/" + pathAfterAPI[len(servicePrefix):]
				} else {
					targetPath = "/" + pathAfterAPI
				}
			} else {
				targetPath = requestPath
			}
		}
		if !strings.HasPrefix(targetPath, "/") {
			targetPath = "/" + targetPath
		}
		targetURL := serviceURL + targetPath

		if r.URL.RawQuery != "" {
			targetURL += "?" + r.URL.RawQuery
		}

		rt.logger.Info(ctx, "forwarding request", zap.String("service", serviceName), zap.String("target", targetURL))

		resp, err := rt.client.ForwardRequest(ctx, targetURL, r)
		if err != nil {
			rt.logger.Error(ctx, "failed to forward request", zap.Error(err))
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		defer resp.Body.Close()

		for key, values := range resp.Header {
			for _, value := range values {
				w.Header().Add(key, value)
			}
		}

		w.WriteHeader(resp.StatusCode)

		_, err = io.Copy(w, resp.Body)
		if err != nil {
			rt.logger.Error(ctx, "failed to copy response body", zap.Error(err))
		}
	}
}
