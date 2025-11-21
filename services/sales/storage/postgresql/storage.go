package postgresql

import (
	"context"
	"database/sql"
	"microservice-challenge/package/errors"
	"microservice-challenge/services/sales/model"
	"microservice-challenge/services/sales/storage/postgresql/db"
	"strconv"

	"github.com/google/uuid"
)

type Storage struct {
	queries *db.Queries
	db      *sql.DB
}

func NewStorage(database *sql.DB) *Storage {
	return &Storage{
		queries: db.New(database),
		db:      database,
	}
}

// convertDBOrderToModel converts sqlc generated db.SalesOrder to model.SalesOrder
func convertDBOrderToModel(dbOrder db.SalesOrder) model.SalesOrder {
	order := model.SalesOrder{
		ID:         dbOrder.ID,
		CustomerID: dbOrder.CustomerID,
		Status:     model.OrderStatus(dbOrder.Status),
		CreatedAt:  dbOrder.CreatedAt,
		UpdatedAt:  dbOrder.UpdatedAt,
	}

	if totalAmount, err := strconv.ParseFloat(dbOrder.TotalAmount, 64); err == nil {
		order.TotalAmount = totalAmount
	}

	return order
}

// convertModelOrderToCreateParams converts model.SalesOrder to sqlc CreateOrderParams
func convertModelOrderToCreateParams(order model.SalesOrder) db.CreateOrderParams {
	return db.CreateOrderParams{
		ID:          order.ID,
		CustomerID:  order.CustomerID,
		Status:      string(order.Status),
		TotalAmount: strconv.FormatFloat(order.TotalAmount, 'f', 2, 64),
		CreatedAt:   order.CreatedAt,
		UpdatedAt:   order.UpdatedAt,
	}
}

// convertModelOrderToUpdateParams converts model.SalesOrder to sqlc UpdateOrderParams
func convertModelOrderToUpdateParams(order model.SalesOrder) db.UpdateOrderParams {
	return db.UpdateOrderParams{
		ID:          order.ID,
		CustomerID:  order.CustomerID,
		Status:      string(order.Status),
		TotalAmount: strconv.FormatFloat(order.TotalAmount, 'f', 2, 64),
		UpdatedAt:   order.UpdatedAt,
	}
}

// convertDBOrderItemToModel converts sqlc generated db.OrderItem to model.OrderItem
func convertDBOrderItemToModel(dbItem db.OrderItem) model.OrderItem {
	item := model.OrderItem{
		ID:        dbItem.ID,
		OrderID:   dbItem.OrderID,
		ItemID:    dbItem.ItemID,
		Quantity:  int(dbItem.Quantity),
		CreatedAt: dbItem.CreatedAt,
		UpdatedAt: dbItem.UpdatedAt,
	}

	if unitPrice, err := strconv.ParseFloat(dbItem.UnitPrice, 64); err == nil {
		item.UnitPrice = unitPrice
	}
	if subtotal, err := strconv.ParseFloat(dbItem.Subtotal, 64); err == nil {
		item.Subtotal = subtotal
	}

	return item
}

// convertModelOrderItemToCreateParams converts model.OrderItem to sqlc CreateOrderItemParams
func convertModelOrderItemToCreateParams(item model.OrderItem) db.CreateOrderItemParams {
	return db.CreateOrderItemParams{
		ID:        item.ID,
		OrderID:   item.OrderID,
		ItemID:    item.ItemID,
		Quantity:  int32(item.Quantity),
		UnitPrice: strconv.FormatFloat(item.UnitPrice, 'f', 2, 64),
		Subtotal:  strconv.FormatFloat(item.Subtotal, 'f', 2, 64),
		CreatedAt: item.CreatedAt,
		UpdatedAt: item.UpdatedAt,
	}
}

func (s *Storage) CreateOrder(ctx context.Context, order model.SalesOrder) error {
	params := convertModelOrderToCreateParams(order)
	if err := s.queries.CreateOrder(ctx, params); err != nil {
		return errors.ErrInternalServerError
	}
	return nil
}

func (s *Storage) GetOrderByID(ctx context.Context, id string) (model.SalesOrder, error) {
	orderID, err := uuid.Parse(id)
	if err != nil {
		return model.SalesOrder{}, errors.ErrBadRequest
	}

	dbOrder, err := s.queries.GetOrderByID(ctx, orderID)
	if err == sql.ErrNoRows {
		return model.SalesOrder{}, errors.ErrNotFound
	}
	if err != nil {
		return model.SalesOrder{}, errors.ErrInternalServerError
	}

	return convertDBOrderToModel(dbOrder), nil
}

func (s *Storage) ListOrders(ctx context.Context, limit, offset int) ([]model.SalesOrder, error) {
	params := db.ListOrdersParams{
		Limit:  int32(limit),
		Offset: int32(offset),
	}

	dbOrders, err := s.queries.ListOrders(ctx, params)
	if err != nil {
		return nil, errors.ErrInternalServerError
	}

	orders := make([]model.SalesOrder, 0, len(dbOrders))
	for _, dbOrder := range dbOrders {
		orders = append(orders, convertDBOrderToModel(dbOrder))
	}

	return orders, nil
}

func (s *Storage) UpdateOrder(ctx context.Context, order model.SalesOrder) error {
	_, err := s.queries.GetOrderByID(ctx, order.ID)
	if err == sql.ErrNoRows {
		return errors.ErrNotFound
	}
	if err != nil {
		return errors.ErrInternalServerError
	}

	params := convertModelOrderToUpdateParams(order)
	if err := s.queries.UpdateOrder(ctx, params); err != nil {
		return errors.ErrInternalServerError
	}

	return nil
}

func (s *Storage) UpdateOrderStatus(ctx context.Context, id string, status model.OrderStatus) error {
	orderID, err := uuid.Parse(id)
	if err != nil {
		return errors.ErrBadRequest
	}

	params := db.UpdateOrderStatusParams{
		ID:     orderID,
		Status: string(status),
	}

	if err := s.queries.UpdateOrderStatus(ctx, params); err != nil {
		return errors.ErrInternalServerError
	}

	_, err = s.queries.GetOrderByID(ctx, orderID)
	if err == sql.ErrNoRows {
		return errors.ErrNotFound
	}
	if err != nil {
		return errors.ErrInternalServerError
	}

	return nil
}

func (s *Storage) CreateOrderItem(ctx context.Context, item model.OrderItem) error {
	params := convertModelOrderItemToCreateParams(item)
	if err := s.queries.CreateOrderItem(ctx, params); err != nil {
		return errors.ErrInternalServerError
	}
	return nil
}

func (s *Storage) CreateOrderItems(ctx context.Context, items []model.OrderItem) error {
	if len(items) == 0 {
		return nil
	}

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return errors.ErrInternalServerError
	}
	defer tx.Rollback()

	qtx := s.queries.WithTx(tx)
	for _, item := range items {
		params := convertModelOrderItemToCreateParams(item)
		if err := qtx.CreateOrderItem(ctx, params); err != nil {
			return errors.ErrInternalServerError
		}
	}

	if err := tx.Commit(); err != nil {
		return errors.ErrInternalServerError
	}

	return nil
}

func (s *Storage) GetOrderItemsByOrderID(ctx context.Context, orderID string) ([]model.OrderItem, error) {
	orderUUID, err := uuid.Parse(orderID)
	if err != nil {
		return nil, errors.ErrBadRequest
	}

	dbItems, err := s.queries.GetOrderItemsByOrderID(ctx, orderUUID)
	if err != nil {
		return nil, errors.ErrInternalServerError
	}

	items := make([]model.OrderItem, 0, len(dbItems))
	for _, dbItem := range dbItems {
		items = append(items, convertDBOrderItemToModel(dbItem))
	}

	return items, nil
}

func (s *Storage) DeleteOrderItemsByOrderID(ctx context.Context, orderID string) error {
	orderUUID, err := uuid.Parse(orderID)
	if err != nil {
		return errors.ErrBadRequest
	}

	if err := s.queries.DeleteOrderItemsByOrderID(ctx, orderUUID); err != nil {
		return errors.ErrInternalServerError
	}

	return nil
}
