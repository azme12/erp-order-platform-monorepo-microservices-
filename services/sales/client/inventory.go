package client

import (
	"context"
	"fmt"
	"microservice-challenge/package/client"
	"microservice-challenge/services/inventory/model"
)

type InventoryClient struct {
	*client.BaseClient
}

func NewInventoryClient(baseURL string) *InventoryClient {
	return &InventoryClient{
		BaseClient: client.NewBaseClient(baseURL),
	}
}

func (c *InventoryClient) GetItemByID(ctx context.Context, itemID string, token string) (model.Item, error) {
	var item model.Item
	path := fmt.Sprintf("/items/%s", itemID)
	if err := c.Get(ctx, path, token, &item); err != nil {
		return model.Item{}, err
	}
	return item, nil
}
