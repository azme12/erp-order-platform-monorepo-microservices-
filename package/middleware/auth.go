package middleware

import (
	"context"
	"microservice-challenge/package/errors"
	"microservice-challenge/package/jwt"
	"microservice-challenge/package/log"
	"microservice-challenge/package/response"
	"net/http"
	"strings"

	"go.uber.org/zap"
)

type contextKey string

const (
	userIDKey    contextKey = "user_id"
	emailKey     contextKey = "email"
	roleKey      contextKey = "role"
	tokenKey     contextKey = "token"
	tokenTypeKey contextKey = "token_type"
)

type AuthMiddleware struct {
	jwtSecret string
	logger    log.Logger
}

func NewAuthMiddleware(jwtSecret string, logger log.Logger) *AuthMiddleware {
	return &AuthMiddleware{
		jwtSecret: jwtSecret,
		logger:    logger,
	}
}

func (m *AuthMiddleware) ValidateToken(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			response.SendErrorResponse(w, errors.ErrUnauthorized)
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			response.SendErrorResponse(w, errors.ErrUnauthorized)
			return
		}

		tokenString := parts[1]

		claims, err := jwt.ValidateToken(tokenString, m.jwtSecret)
		if err != nil {
			m.logger.Error(ctx, "token validation failed", zap.Error(err))
			response.SendErrorResponse(w, errors.ErrUnauthorized)
			return
		}

		ctx = context.WithValue(ctx, userIDKey, claims.UserID)
		ctx = context.WithValue(ctx, emailKey, claims.Email)
		ctx = context.WithValue(ctx, roleKey, claims.Role)
		ctx = context.WithValue(ctx, tokenKey, tokenString)
		ctx = context.WithValue(ctx, tokenTypeKey, claims.Type)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (m *AuthMiddleware) RequireRole(allowedRoles ...string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()

			// Allow service tokens to bypass role checks for inter-service communication
			tokenType, _ := ctx.Value(tokenTypeKey).(string)
			if tokenType == "service" {
				next.ServeHTTP(w, r)
				return
			}

			role, ok := ctx.Value(roleKey).(string)
			if !ok || role == "" {
				response.SendErrorResponse(w, errors.ErrForbidden)
				return
			}

			for _, allowedRole := range allowedRoles {
				if role == allowedRole {
					next.ServeHTTP(w, r)
					return
				}
			}

			response.SendErrorResponse(w, errors.ErrForbidden)
		})
	}
}

func GetRoleFromContext(ctx context.Context) string {
	role, _ := ctx.Value(roleKey).(string)
	return role
}

func GetUserIDFromContext(ctx context.Context) string {
	userID, _ := ctx.Value(userIDKey).(string)
	return userID
}

func GetTokenKey() contextKey {
	return tokenKey
}
