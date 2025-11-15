package client

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"microservice-challenge/package/errors"
	"microservice-challenge/services/inventory/model"
	"net/http"
)

type InventoryClient struct {
	baseURL string
	client  *http.Client
}

func NewInventoryClient(baseURL string) *InventoryClient {
	return &InventoryClient{
		baseURL: baseURL,
		client:  &http.Client{},
	}
}

func (c *InventoryClient) GetItemByID(ctx context.Context, itemID string, token string) (model.Item, error) {
	url := fmt.Sprintf("%s/items/%s", c.baseURL, itemID)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return model.Item{}, errors.ErrInternalServerError
	}

	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.client.Do(req)
	if err != nil {
		return model.Item{}, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	bodyBytes, _ := io.ReadAll(resp.Body)
	bodyStr := string(bodyBytes)

	if resp.StatusCode == http.StatusNotFound {
		return model.Item{}, errors.ErrNotFound
	}

	if resp.StatusCode == http.StatusUnauthorized {
		return model.Item{}, errors.ErrUnauthorized
	}

	if resp.StatusCode != http.StatusOK {
		return model.Item{}, fmt.Errorf("unexpected status code %d: %s", resp.StatusCode, bodyStr)
	}

	var response struct {
		Status  int        `json:"status"`
		Message string     `json:"message"`
		Data    model.Item `json:"data"`
	}

	if err := json.Unmarshal(bodyBytes, &response); err != nil {
		return model.Item{}, fmt.Errorf("failed to decode response: %w, body: %s", err, bodyStr)
	}

	return response.Data, nil
}
