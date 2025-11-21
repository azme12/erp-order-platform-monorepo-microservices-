package internal

import (
	routerpkg "microservice-challenge/package/router"
	"microservice-challenge/services/auth/httphandler"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func InitAuthRoutes(router chi.Router, handler *httphandler.Handler) {
	routes := []routerpkg.Route{
		{
			Method:  http.MethodPost,
			Path:    "/register",
			Handler: handler.Register,
		},
		{
			Method:  http.MethodPost,
			Path:    "/login",
			Handler: handler.Login,
		},
		{
			Method:  http.MethodPost,
			Path:    "/forgot-password",
			Handler: handler.ForgotPassword,
		},
		{
			Method:  http.MethodPost,
			Path:    "/reset-password",
			Handler: handler.ResetPassword,
		},
		{
			Method:  http.MethodPost,
			Path:    "/service-token",
			Handler: handler.GenerateServiceToken,
		},
	}

	routerpkg.RegisterRoutes(router, routes)
}
