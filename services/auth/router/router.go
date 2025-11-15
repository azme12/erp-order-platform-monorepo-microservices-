package router

import (
	"encoding/json"
	"microservice-challenge/package/log"
	"microservice-challenge/services/auth/docs"
	"microservice-challenge/services/auth/httphandler"
	"microservice-challenge/services/auth/router/internal"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	httpSwagger "github.com/swaggo/http-swagger"
)

func NewRouter(handler *httphandler.Handler, logger log.Logger) http.Handler {
	r := chi.NewRouter()

	docs.SwaggerInfo.Schemes = []string{"http", "https"}
	docs.SwaggerInfo.Host = "localhost:8000"

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.AllowContentType("application/json"))

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{
			"status":  "healthy",
			"service": "auth-service",
		})
	})

	r.Route("/", func(r chi.Router) {
		r.Handle("/swagger/*", httpSwagger.WrapHandler)
		internal.InitAuthRoutes(r, handler)
	})

	return r
}
