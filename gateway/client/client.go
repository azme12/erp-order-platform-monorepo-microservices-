package client

import (
	"bytes"
	"context"
	"fmt"
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

func NewClient(logger log.Logger) *Client {
	return &Client{
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		logger: logger,
	}
}

func (c *Client) ForwardRequest(ctx context.Context, targetURL string, r *http.Request) (*http.Response, error) {
	return c.ForwardRequestWithRetry(ctx, targetURL, r, 3)
}

func (c *Client) ForwardRequestWithRetry(ctx context.Context, targetURL string, r *http.Request, maxRetries int) (*http.Response, error) {
	var bodyBytes []byte
	var err error

	if r.Body != nil {
		bodyBytes, err = io.ReadAll(r.Body)
		if err != nil {
			return nil, fmt.Errorf("failed to read request body: %w", err)
		}
		r.Body.Close()
	}

	var lastErr error
	for attempt := 0; attempt < maxRetries; attempt++ {
		if attempt > 0 {
			// Exponential backoff: 100ms, 200ms, 400ms
			backoff := time.Duration(attempt) * 100 * time.Millisecond
			select {
			case <-ctx.Done():
				return nil, ctx.Err()
			case <-time.After(backoff):
			}
		}

		var bodyReader io.Reader
		if bodyBytes != nil {
			bodyReader = bytes.NewReader(bodyBytes)
		}

		req, err := http.NewRequestWithContext(ctx, r.Method, targetURL, bodyReader)
		if err != nil {
			return nil, err
		}

		for key, values := range r.Header {
			if key != "Host" {
				for _, value := range values {
					req.Header.Add(key, value)
				}
			}
		}

		resp, err := c.httpClient.Do(req)
		if err == nil {
			// Success - but check if it's a retryable error status
			if resp.StatusCode >= 500 && resp.StatusCode < 600 && attempt < maxRetries-1 {
				resp.Body.Close()
				lastErr = fmt.Errorf("server error %d, retrying", resp.StatusCode)
				c.logger.Warn(ctx, "retrying request", zap.String("target", targetURL), zap.Int("attempt", attempt+1), zap.Int("status", resp.StatusCode))
				continue
			}
			return resp, nil
		}

		lastErr = err
		// Don't retry on context cancellation or timeout
		if err == context.Canceled || err == context.DeadlineExceeded {
			return nil, err
		}

		c.logger.Warn(ctx, "retrying request", zap.String("target", targetURL), zap.Int("attempt", attempt+1), zap.Error(err))
	}

	c.logger.Error(ctx, "failed to forward request after retries", zap.String("target", targetURL), zap.Int("attempts", maxRetries), zap.Error(lastErr))
	return nil, fmt.Errorf("failed after %d attempts: %w", maxRetries, lastErr)
}
