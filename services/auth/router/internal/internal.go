package internal

import (
	"microservice-challenge/services/auth/httphandler"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func InitAuthRoutes(router chi.Router, handler *httphandler.Handler) {
	routes := []Route{
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
	}

	RegisterRoutes(router, routes)
}
