package storage

import (
	"context"
	"microservice-challenge/services/inventory/model"
)

type Storage interface {
	CreateItem(ctx context.Context, item model.Item) error
	GetItemByID(ctx context.Context, id string) (model.Item, error)
	GetItemBySKU(ctx context.Context, sku string) (model.Item, error)
	ListItems(ctx context.Context, limit, offset int) ([]model.Item, error)
	UpdateItem(ctx context.Context, item model.Item) error
	DeleteItem(ctx context.Context, id string) error

	GetStockByItemID(ctx context.Context, itemID string) (model.Stock, error)
	CreateStock(ctx context.Context, stock model.Stock) error
	UpdateStock(ctx context.Context, stock model.Stock) error
	AdjustStock(ctx context.Context, itemID string, quantityDelta int) error
}
