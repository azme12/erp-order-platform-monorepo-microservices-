package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"microservice-challenge/package/errors"
	"net/http"
	"sync"
	"time"
)

type AuthClient struct {
	baseURL       string
	client        *http.Client
	serviceName   string
	serviceSecret string
	cachedToken   *cachedToken
	mu            sync.RWMutex
}

type cachedToken struct {
	token     string
	expiresAt time.Time
}

type ServiceTokenRequest struct {
	ServiceName   string `json:"service_name"`
	ServiceSecret string `json:"service_secret"`
}

type ServiceTokenResponse struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
	Data    struct {
		Token     string `json:"token"`
		ExpiresIn int    `json:"expires_in"`
	} `json:"data"`
}

func NewAuthClient(baseURL, serviceName, serviceSecret string) *AuthClient {
	return &AuthClient{
		baseURL:       baseURL,
		client:        &http.Client{Timeout: 10 * time.Second},
		serviceName:   serviceName,
		serviceSecret: serviceSecret,
	}
}

func (c *AuthClient) GetServiceToken(ctx context.Context) (string, error) {

	c.mu.RLock()
	if c.cachedToken != nil && time.Now().Before(c.cachedToken.expiresAt.Add(-5*time.Minute)) {
		token := c.cachedToken.token
		c.mu.RUnlock()
		return token, nil
	}
	c.mu.RUnlock()

	c.mu.Lock()
	defer c.mu.Unlock()

	if c.cachedToken != nil && time.Now().Before(c.cachedToken.expiresAt.Add(-5*time.Minute)) {
		return c.cachedToken.token, nil
	}

	reqBody := ServiceTokenRequest{
		ServiceName:   c.serviceName,
		ServiceSecret: c.serviceSecret,
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return "", fmt.Errorf("failed to marshal service token request: %w", err)
	}

	url := fmt.Sprintf("%s/service-token", c.baseURL)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(jsonData))
	if err != nil {
		return "", fmt.Errorf("failed to create service token request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := c.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to request service token from auth service: %w", err)
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read service token response: %w", err)
	}
	bodyStr := string(bodyBytes)

	if resp.StatusCode == http.StatusUnauthorized {
		return "", errors.ErrUnauthorized
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("auth service returned unexpected status code %d: %s", resp.StatusCode, bodyStr)
	}

	var response ServiceTokenResponse
	if err := json.Unmarshal(bodyBytes, &response); err != nil {
		return "", fmt.Errorf("failed to decode service token response: %w, body: %s", err, bodyStr)
	}

	if response.Data.Token == "" {
		return "", fmt.Errorf("auth service returned empty token in response")
	}

	c.cachedToken = &cachedToken{
		token:     response.Data.Token,
		expiresAt: time.Now().Add(time.Duration(response.Data.ExpiresIn) * time.Second),
	}

	return response.Data.Token, nil
}
