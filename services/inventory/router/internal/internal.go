package internal

import (
	"microservice-challenge/package/middleware"
	"microservice-challenge/services/inventory/httphandler"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func InitInventoryRoutes(router chi.Router, handler *httphandler.Handler, authMiddleware *middleware.AuthMiddleware) {
	routes := []Route{
		{
			Method:      http.MethodGet,
			Path:        "/items",
			Handler:     handler.ListItems,
			Middlewares: []func(next http.Handler) http.Handler{authMiddleware.ValidateToken, authMiddleware.RequireRole("inventory_manager", "finance_manager")},
		},
		{
			Method:      http.MethodGet,
			Path:        "/items/{id}",
			Handler:     handler.GetItem,
			Middlewares: []func(next http.Handler) http.Handler{authMiddleware.ValidateToken},
		},
		{
			Method:      http.MethodPost,
			Path:        "/items",
			Handler:     handler.CreateItem,
			Middlewares: []func(next http.Handler) http.Handler{authMiddleware.ValidateToken, authMiddleware.RequireRole("inventory_manager", "finance_manager")},
		},
		{
			Method:      http.MethodPut,
			Path:        "/items/{id}",
			Handler:     handler.UpdateItem,
			Middlewares: []func(next http.Handler) http.Handler{authMiddleware.ValidateToken, authMiddleware.RequireRole("inventory_manager", "finance_manager")},
		},
		{
			Method:      http.MethodDelete,
			Path:        "/items/{id}",
			Handler:     handler.DeleteItem,
			Middlewares: []func(next http.Handler) http.Handler{authMiddleware.ValidateToken, authMiddleware.RequireRole("finance_manager")},
		},
		{
			Method:      http.MethodGet,
			Path:        "/items/{item_id}/stock",
			Handler:     handler.GetStock,
			Middlewares: []func(next http.Handler) http.Handler{authMiddleware.ValidateToken},
		},
		{
			Method:      http.MethodPut,
			Path:        "/items/{item_id}/stock",
			Handler:     handler.AdjustStock,
			Middlewares: []func(next http.Handler) http.Handler{authMiddleware.ValidateToken, authMiddleware.RequireRole("inventory_manager", "finance_manager")},
		},
	}

	RegisterRoutes(router, routes)
}
