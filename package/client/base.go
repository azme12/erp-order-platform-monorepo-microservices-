package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"microservice-challenge/package/errors"
	"net/http"
	"time"
)

type BaseClient struct {
	baseURL string
	client  *http.Client
}

func NewBaseClient(baseURL string) *BaseClient {
	return &BaseClient{
		baseURL: baseURL,
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

func (c *BaseClient) DoRequest(ctx context.Context, method, path string, body interface{}, token string) (*http.Response, error) {
	var bodyReader io.Reader
	if body != nil {
		jsonData, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal request body: %w", err)
		}
		bodyReader = bytes.NewReader(jsonData)
	}

	url := fmt.Sprintf("%s%s", c.baseURL, path)
	req, err := http.NewRequestWithContext(ctx, method, url, bodyReader)
	if err != nil {
		return nil, errors.ErrInternalServerError
	}

	req.Header.Set("Content-Type", "application/json")
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}

	return resp, nil
}

func (c *BaseClient) ParseResponse(resp *http.Response, target interface{}) error {
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %w", err)
	}

	bodyStr := string(bodyBytes)

	switch resp.StatusCode {
	case http.StatusNotFound:
		return errors.ErrNotFound
	case http.StatusUnauthorized:
		return errors.ErrUnauthorized
	case http.StatusForbidden:
		return errors.ErrForbidden
	case http.StatusConflict:
		return errors.ErrConflict
	case http.StatusBadRequest:
		return errors.ErrBadRequest
	case http.StatusOK, http.StatusCreated:

		var responseWrapper struct {
			Status  int             `json:"status"`
			Message string          `json:"message"`
			Data    json.RawMessage `json:"data"`
		}

		if err := json.Unmarshal(bodyBytes, &responseWrapper); err != nil {
			return fmt.Errorf("failed to decode response: %w, body: %s", err, bodyStr)
		}

		if err := json.Unmarshal(responseWrapper.Data, target); err != nil {
			return fmt.Errorf("failed to unmarshal into target: %w, body: %s", err, bodyStr)
		}

		return nil
	default:
		return fmt.Errorf("unexpected status code %d: %s", resp.StatusCode, bodyStr)
	}
}

func (c *BaseClient) Get(ctx context.Context, path string, token string, target interface{}) error {
	resp, err := c.DoRequest(ctx, http.MethodGet, path, nil, token)
	if err != nil {
		return err
	}
	return c.ParseResponse(resp, target)
}

func (c *BaseClient) Post(ctx context.Context, path string, body interface{}, token string, target interface{}) error {
	resp, err := c.DoRequest(ctx, http.MethodPost, path, body, token)
	if err != nil {
		return err
	}
	return c.ParseResponse(resp, target)
}

func (c *BaseClient) Put(ctx context.Context, path string, body interface{}, token string, target interface{}) error {
	resp, err := c.DoRequest(ctx, http.MethodPut, path, body, token)
	if err != nil {
		return err
	}
	return c.ParseResponse(resp, target)
}

func (c *BaseClient) Delete(ctx context.Context, path string, token string) error {
	resp, err := c.DoRequest(ctx, http.MethodDelete, path, nil, token)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		bodyBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("unexpected status code %d: failed to read response body: %w", resp.StatusCode, err)
		}
		return fmt.Errorf("unexpected status code %d: %s", resp.StatusCode, string(bodyBytes))
	}

	return nil
}
