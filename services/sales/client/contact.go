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

func (c *ContactClient) GetCustomerByID(ctx context.Context, customerID string, token string) (model.Customer, error) {
	url := fmt.Sprintf("%s/customers/%s", c.baseURL, customerID)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return model.Customer{}, errors.ErrInternalServerError
	}

	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.client.Do(req)
	if err != nil {
		return model.Customer{}, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	bodyBytes, _ := io.ReadAll(resp.Body)
	bodyStr := string(bodyBytes)

	if resp.StatusCode == http.StatusNotFound {
		return model.Customer{}, errors.ErrNotFound
	}

	if resp.StatusCode == http.StatusUnauthorized {
		return model.Customer{}, errors.ErrUnauthorized
	}

	if resp.StatusCode != http.StatusOK {
		return model.Customer{}, fmt.Errorf("unexpected status code %d: %s", resp.StatusCode, bodyStr)
	}

	var response struct {
		Status  int            `json:"status"`
		Message string         `json:"message"`
		Data    model.Customer `json:"data"`
	}

	if err := json.Unmarshal(bodyBytes, &response); err != nil {
		return model.Customer{}, fmt.Errorf("failed to decode response: %w, body: %s", err, bodyStr)
	}

	return response.Data, nil
}
