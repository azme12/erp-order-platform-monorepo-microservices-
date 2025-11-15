package client

import (
	"bytes"
	"context"
	"io"
	"microservice-challenge/package/log"
	"net/http"
	"time"

	"go.uber.org/zap"
)

type Client struct {
	httpClient *http.Client
	logger     log.Logger
}

// NewClient creates a new HTTP client
func NewClient(logger log.Logger) *Client {
	return &Client{
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		logger: logger,
	}
}

// ForwardRequest forwards a request to a target service
func (c *Client) ForwardRequest(ctx context.Context, targetURL string, r *http.Request) (*http.Response, error) {
	// Create new request
	req, err := http.NewRequestWithContext(ctx, r.Method, targetURL, r.Body)
	if err != nil {
		return nil, err
	}

	// Copy headers (except Host)
	for key, values := range r.Header {
		if key != "Host" {
			for _, value := range values {
				req.Header.Add(key, value)
			}
		}
	}

	// Forward the request
	resp, err := c.httpClient.Do(req)
	if err != nil {
		c.logger.Error(ctx, "failed to forward request", zap.String("target", targetURL), zap.Error(err))
		return nil, err
	}

	return resp, nil
}

// ForwardRequestWithBody forwards a request with a new body
func (c *Client) ForwardRequestWithBody(ctx context.Context, targetURL string, method string, body []byte, headers map[string]string) (*http.Response, error) {
	var bodyReader io.Reader
	if body != nil {
		bodyReader = bytes.NewReader(body)
	}

	req, err := http.NewRequestWithContext(ctx, method, targetURL, bodyReader)
	if err != nil {
		return nil, err
	}

	// Set headers
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	// Forward the request
	resp, err := c.httpClient.Do(req)
	if err != nil {
		c.logger.Error(ctx, "failed to forward request", zap.String("target", targetURL), zap.Error(err))
		return nil, err
	}

	return resp, nil
}
