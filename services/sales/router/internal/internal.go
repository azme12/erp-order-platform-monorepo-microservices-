package internal

import (
	"microservice-challenge/package/middleware"
	"microservice-challenge/services/sales/httphandler"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func InitSalesRoutes(router chi.Router, handler *httphandler.Handler, authMiddleware *middleware.AuthMiddleware) {
	routes := []Route{
		{
			Method:      http.MethodGet,
			Path:        "/orders",
			Handler:     handler.ListOrders,
			Middlewares: []func(next http.Handler) http.Handler{authMiddleware.ValidateToken, authMiddleware.RequireRole("inventory_manager", "finance_manager")},
		},
		{
			Method:      http.MethodGet,
			Path:        "/orders/{id}",
			Handler:     handler.GetOrder,
			Middlewares: []func(next http.Handler) http.Handler{authMiddleware.ValidateToken, authMiddleware.RequireRole("inventory_manager", "finance_manager")},
		},
		{
			Method:      http.MethodPost,
			Path:        "/orders",
			Handler:     handler.CreateOrder,
			Middlewares: []func(next http.Handler) http.Handler{authMiddleware.ValidateToken, authMiddleware.RequireRole("inventory_manager", "finance_manager")},
		},
		{
			Method:      http.MethodPut,
			Path:        "/orders/{id}",
			Handler:     handler.UpdateOrder,
			Middlewares: []func(next http.Handler) http.Handler{authMiddleware.ValidateToken, authMiddleware.RequireRole("inventory_manager", "finance_manager")},
		},
		{
			Method:      http.MethodPost,
			Path:        "/orders/{id}/confirm",
			Handler:     handler.ConfirmOrder,
			Middlewares: []func(next http.Handler) http.Handler{authMiddleware.ValidateToken, authMiddleware.RequireRole("finance_manager")},
		},
		{
			Method:      http.MethodPost,
			Path:        "/orders/{id}/pay",
			Handler:     handler.PayOrder,
			Middlewares: []func(next http.Handler) http.Handler{authMiddleware.ValidateToken, authMiddleware.RequireRole("finance_manager")},
		},
	}

	RegisterRoutes(router, routes)
}
