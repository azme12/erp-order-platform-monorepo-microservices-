package client

import (
	"context"
	"fmt"
	"microservice-challenge/package/client"
	"microservice-challenge/services/contact/model"
)

type ContactClient struct {
	*client.BaseClient
}

func NewContactClient(baseURL string) *ContactClient {
	return &ContactClient{
		BaseClient: client.NewBaseClient(baseURL),
	}
}

func (c *ContactClient) GetCustomerByID(ctx context.Context, customerID string, token string) (model.Customer, error) {
	var customer model.Customer
	path := fmt.Sprintf("/customers/%s", customerID)
	if err := c.Get(ctx, path, token, &customer); err != nil {
		return model.Customer{}, err
	}
	return customer, nil
}
