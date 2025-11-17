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

func (c *ContactClient) GetVendorByID(ctx context.Context, vendorID string, token string) (model.Vendor, error) {
	var vendor model.Vendor
	path := fmt.Sprintf("/vendors/%s", vendorID)
	if err := c.Get(ctx, path, token, &vendor); err != nil {
		return model.Vendor{}, err
	}
	return vendor, nil
}
