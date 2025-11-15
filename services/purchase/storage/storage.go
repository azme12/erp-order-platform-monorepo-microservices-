package storage

import (
	"context"
	"microservice-challenge/services/purchase/model"
)

type Storage interface {
	CreateOrder(ctx context.Context, order model.PurchaseOrder) error
	GetOrderByID(ctx context.Context, id string) (model.PurchaseOrder, error)
	ListOrders(ctx context.Context, limit, offset int) ([]model.PurchaseOrder, error)
	UpdateOrder(ctx context.Context, order model.PurchaseOrder) error
	UpdateOrderStatus(ctx context.Context, id string, status model.PurchaseOrderStatus) error

	CreateOrderItem(ctx context.Context, item model.PurchaseOrderItem) error
	GetOrderItemsByOrderID(ctx context.Context, orderID string) ([]model.PurchaseOrderItem, error)
	DeleteOrderItemsByOrderID(ctx context.Context, orderID string) error
}
