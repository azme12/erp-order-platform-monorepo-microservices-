package storage

import (
	"context"
	"microservice-challenge/services/sales/model"
)

type Storage interface {
	CreateOrder(ctx context.Context, order model.SalesOrder) error
	GetOrderByID(ctx context.Context, id string) (model.SalesOrder, error)
	ListOrders(ctx context.Context, limit, offset int) ([]model.SalesOrder, error)
	UpdateOrder(ctx context.Context, order model.SalesOrder) error
	UpdateOrderStatus(ctx context.Context, id string, status model.OrderStatus) error

	CreateOrderItem(ctx context.Context, item model.OrderItem) error
	CreateOrderItems(ctx context.Context, items []model.OrderItem) error
	GetOrderItemsByOrderID(ctx context.Context, orderID string) ([]model.OrderItem, error)
	DeleteOrderItemsByOrderID(ctx context.Context, orderID string) error
}
