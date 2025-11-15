package client

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"microservice-challenge/package/errors"
	"microservice-challenge/services/contact/model"
	"net/http"
)

type ContactClient struct {
	baseURL string
	client  *http.Client
}

func NewContactClient(baseURL string) *ContactClient {
	return &ContactClient{
		baseURL: baseURL,
		client:  &http.Client{},
	}
}

func (c *ContactClient) GetVendorByID(ctx context.Context, vendorID string, token string) (model.Vendor, error) {
	url := fmt.Sprintf("%s/vendors/%s", c.baseURL, vendorID)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return model.Vendor{}, errors.ErrInternalServerError
	}

	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.client.Do(req)
	if err != nil {
		return model.Vendor{}, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	bodyBytes, _ := io.ReadAll(resp.Body)
	bodyStr := string(bodyBytes)

	if resp.StatusCode == http.StatusNotFound {
		return model.Vendor{}, errors.ErrNotFound
	}

	if resp.StatusCode == http.StatusUnauthorized {
		return model.Vendor{}, errors.ErrUnauthorized
	}

	if resp.StatusCode != http.StatusOK {
		return model.Vendor{}, fmt.Errorf("unexpected status code %d: %s", resp.StatusCode, bodyStr)
	}

	var response struct {
		Status  int          `json:"status"`
		Message string       `json:"message"`
		Data    model.Vendor `json:"data"`
	}

	if err := json.Unmarshal(bodyBytes, &response); err != nil {
		return model.Vendor{}, fmt.Errorf("failed to decode response: %w, body: %s", err, bodyStr)
	}

	return response.Data, nil
}
