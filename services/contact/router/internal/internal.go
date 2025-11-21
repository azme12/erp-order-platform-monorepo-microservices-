package internal

import (
	"microservice-challenge/package/middleware"
	routerpkg "microservice-challenge/package/router"
	"microservice-challenge/services/contact/httphandler"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func InitContactRoutes(router chi.Router, handler *httphandler.Handler, authMiddleware *middleware.AuthMiddleware) {
	routes := []routerpkg.Route{
		{
			Method:      http.MethodGet,
			Path:        "/customers",
			Handler:     handler.ListCustomers,
			Middlewares: []func(next http.Handler) http.Handler{authMiddleware.ValidateToken, authMiddleware.RequireRole("inventory_manager", "finance_manager")},
		},
		{
			Method:      http.MethodGet,
			Path:        "/customers/{id}",
			Handler:     handler.GetCustomer,
			Middlewares: []func(next http.Handler) http.Handler{authMiddleware.ValidateToken},
		},
		{
			Method:      http.MethodPost,
			Path:        "/customers",
			Handler:     handler.CreateCustomer,
			Middlewares: []func(next http.Handler) http.Handler{authMiddleware.ValidateToken, authMiddleware.RequireRole("inventory_manager", "finance_manager")},
		},
		{
			Method:      http.MethodPut,
			Path:        "/customers/{id}",
			Handler:     handler.UpdateCustomer,
			Middlewares: []func(next http.Handler) http.Handler{authMiddleware.ValidateToken, authMiddleware.RequireRole("inventory_manager", "finance_manager")},
		},
		{
			Method:      http.MethodDelete,
			Path:        "/customers/{id}",
			Handler:     handler.DeleteCustomer,
			Middlewares: []func(next http.Handler) http.Handler{authMiddleware.ValidateToken, authMiddleware.RequireRole("finance_manager")},
		},
		{
			Method:      http.MethodGet,
			Path:        "/vendors",
			Handler:     handler.ListVendors,
			Middlewares: []func(next http.Handler) http.Handler{authMiddleware.ValidateToken, authMiddleware.RequireRole("inventory_manager", "finance_manager")},
		},
		{
			Method:      http.MethodGet,
			Path:        "/vendors/{id}",
			Handler:     handler.GetVendor,
			Middlewares: []func(next http.Handler) http.Handler{authMiddleware.ValidateToken},
		},
		{
			Method:      http.MethodPost,
			Path:        "/vendors",
			Handler:     handler.CreateVendor,
			Middlewares: []func(next http.Handler) http.Handler{authMiddleware.ValidateToken, authMiddleware.RequireRole("inventory_manager", "finance_manager")},
		},
		{
			Method:      http.MethodPut,
			Path:        "/vendors/{id}",
			Handler:     handler.UpdateVendor,
			Middlewares: []func(next http.Handler) http.Handler{authMiddleware.ValidateToken, authMiddleware.RequireRole("inventory_manager", "finance_manager")},
		},
		{
			Method:      http.MethodDelete,
			Path:        "/vendors/{id}",
			Handler:     handler.DeleteVendor,
			Middlewares: []func(next http.Handler) http.Handler{authMiddleware.ValidateToken, authMiddleware.RequireRole("finance_manager")},
		},
	}

	routerpkg.RegisterRoutes(router, routes)
}
