package router

import (
	"database/sql"
	"microservice-challenge/package/health"
	"microservice-challenge/package/log"
	"microservice-challenge/package/middleware"
	"microservice-challenge/services/contact/docs"
	"microservice-challenge/services/contact/httphandler"
	"microservice-challenge/services/contact/router/internal"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	chiMiddleware "github.com/go-chi/chi/v5/middleware"
	httpSwagger "github.com/swaggo/http-swagger"
)

func NewRouter(handler *httphandler.Handler, logger log.Logger, jwtSecret string, db *sql.DB) http.Handler {
	r := chi.NewRouter()

	docs.SwaggerInfo.Schemes = []string{"http", "https"}
	docs.SwaggerInfo.Host = "localhost:8001"

	r.Use(chiMiddleware.RequestID)
	r.Use(chiMiddleware.RealIP)
	r.Use(middleware.TraceMiddleware(logger))
	r.Use(chiMiddleware.Logger)
	r.Use(chiMiddleware.Recoverer)
	r.Use(chiMiddleware.AllowContentType("application/json"))
	r.Use(middleware.TimeoutMiddleware(30 * time.Second))

	healthCheck := health.New("contact-service", "1.0.0")
	if db != nil {
		healthCheck.Register("database", health.NewDBChecker(db))
	}

	r.Get("/health", healthCheck.Handler())
	r.Get("/health/ready", healthCheck.ReadinessHandler())
	r.Get("/health/live", healthCheck.LivenessHandler())

	authMiddleware := middleware.NewAuthMiddleware(jwtSecret, logger)

	r.Route("/", func(r chi.Router) {
		r.Handle("/swagger/*", httpSwagger.Handler(
			httpSwagger.URL("http://localhost:8001/swagger/doc.json"),
		))

		r.Group(func(r chi.Router) {
			r.Use(authMiddleware.ValidateToken)
			internal.InitContactRoutes(r, handler, authMiddleware)
		})
	})

	return r
}
