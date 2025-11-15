package middleware

import (
	"context"
	"encoding/json"
	"microservice-challenge/package/log"
	"net/http"
	"strings"

	"go.uber.org/zap"
)

type AuthMiddleware struct {
	authServiceURL string
	logger         log.Logger
}

// NewAuthMiddleware creates a new JWT authentication middleware
func NewAuthMiddleware(authServiceURL string, logger log.Logger) *AuthMiddleware {
	return &AuthMiddleware{
		authServiceURL: authServiceURL,
		logger:         logger,
	}
}

// ValidateToken validates a JWT token with the Auth Service
func (m *AuthMiddleware) ValidateToken(ctx context.Context, token string) (bool, map[string]interface{}, error) {
	// Call Auth Service /validate endpoint
	reqBody := map[string]string{"token": token}
	jsonBody, _ := json.Marshal(reqBody)

	req, err := http.NewRequestWithContext(ctx, "POST", m.authServiceURL+"/validate", strings.NewReader(string(jsonBody)))
	if err != nil {
		return false, nil, err
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		m.logger.Error(ctx, "failed to call auth service", zap.Error(err))
		return false, nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return false, nil, nil
	}

	var validateResp struct {
		Data struct {
			Valid   bool   `json:"valid"`
			UserID  string `json:"user_id,omitempty"`
			Email   string `json:"email,omitempty"`
			Role    string `json:"role,omitempty"`
			Service string `json:"service,omitempty"`
			Type    string `json:"type"`
		} `json:"data"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&validateResp); err != nil {
		return false, nil, err
	}

	if !validateResp.Data.Valid {
		return false, nil, nil
	}

	// Return token claims
	claims := map[string]interface{}{
		"user_id": validateResp.Data.UserID,
		"email":   validateResp.Data.Email,
		"role":    validateResp.Data.Role,
		"service": validateResp.Data.Service,
		"type":    validateResp.Data.Type,
	}

	return true, claims, nil
}

// Middleware returns the authentication middleware handler
func (m *AuthMiddleware) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		// Skip auth for /api/auth/register, /api/auth/login, and /api/auth/validate
		path := r.URL.Path
		if path == "/api/auth/register" || path == "/api/auth/login" || path == "/api/auth/validate" {
			next.ServeHTTP(w, r)
			return
		}

		// Extract token from Authorization header
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// Extract token (Bearer <token>)
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		token := parts[1]

		// Validate token
		valid, claims, err := m.ValidateToken(ctx, token)
		if err != nil {
			m.logger.Error(ctx, "token validation error", zap.Error(err))
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		if !valid {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// Add claims to context
		for key, value := range claims {
			ctx = context.WithValue(ctx, key, value)
		}

		// Continue to next handler
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
