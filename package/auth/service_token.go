package auth

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"microservice-challenge/package/errors"
	"microservice-challenge/package/log"
	"net/http"
	"sync"
	"time"

	"go.uber.org/zap"
)

type ServiceTokenClient struct {
	authServiceURL string
	serviceName    string
	serviceSecret  string
	httpClient     *http.Client
	logger         log.Logger
	tokenCache     *tokenCache
}

type tokenCache struct {
	token     string
	expiresAt time.Time
	mu        sync.RWMutex
}

func NewServiceTokenClient(authServiceURL, serviceName, serviceSecret string, logger log.Logger) *ServiceTokenClient {
	return &ServiceTokenClient{
		authServiceURL: authServiceURL,
		serviceName:    serviceName,
		serviceSecret:  serviceSecret,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
		logger:     logger,
		tokenCache: &tokenCache{},
	}
}

func (c *ServiceTokenClient) GetToken(ctx context.Context) (string, error) {
	c.tokenCache.mu.RLock()
	if c.tokenCache.token != "" && time.Now().Before(c.tokenCache.expiresAt.Add(-5*time.Minute)) {
		token := c.tokenCache.token
		c.tokenCache.mu.RUnlock()
		return token, nil
	}
	c.tokenCache.mu.RUnlock()

	return c.fetchToken(ctx)
}

func (c *ServiceTokenClient) fetchToken(ctx context.Context) (string, error) {
	c.tokenCache.mu.Lock()
	defer c.tokenCache.mu.Unlock()

	if c.tokenCache.token != "" && time.Now().Before(c.tokenCache.expiresAt.Add(-5*time.Minute)) {
		return c.tokenCache.token, nil
	}

	reqBody := map[string]string{
		"service_name":   c.serviceName,
		"service_secret": c.serviceSecret,
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request: %w", err)
	}

	url := fmt.Sprintf("%s/service-token", c.authServiceURL)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		c.logger.Error(ctx, "failed to fetch service token", zap.Error(err))
		return "", fmt.Errorf("failed to fetch service token: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		c.logger.Error(ctx, "failed to get service token", zap.Int("status", resp.StatusCode), zap.String("body", string(bodyBytes)))
		return "", errors.ErrUnauthorized
	}

	var response struct {
		Status  int    `json:"status"`
		Message string `json:"message"`
		Data    struct {
			Token     string `json:"token"`
			ExpiresIn int    `json:"expires_in"`
		} `json:"data"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return "", fmt.Errorf("failed to decode response: %w", err)
	}

	c.tokenCache.token = response.Data.Token
	c.tokenCache.expiresAt = time.Now().Add(time.Duration(response.Data.ExpiresIn) * time.Second)

	return response.Data.Token, nil
}
