package client

import (
	"context"
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
	req, err := http.NewRequestWithContext(ctx, r.Method, targetURL, r.Body)
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
	if err != nil {
		c.logger.Error(ctx, "failed to forward request", zap.String("target", targetURL), zap.Error(err))
		return nil, err
	}

	return resp, nil
}
