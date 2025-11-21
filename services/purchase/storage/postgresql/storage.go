package postgresql

import (
	"context"
	"database/sql"
	"microservice-challenge/package/errors"
	"microservice-challenge/services/purchase/model"
	"microservice-challenge/services/purchase/storage/postgresql/db"
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

// convertDBOrderToModel converts sqlc generated db.PurchaseOrder to model.PurchaseOrder
func convertDBOrderToModel(dbOrder db.PurchaseOrder) model.PurchaseOrder {
	order := model.PurchaseOrder{
		ID:        dbOrder.ID,
		VendorID:  dbOrder.VendorID,
		Status:    model.PurchaseOrderStatus(dbOrder.Status),
		CreatedAt: dbOrder.CreatedAt,
		UpdatedAt: dbOrder.UpdatedAt,
	}

	if totalAmount, err := strconv.ParseFloat(dbOrder.TotalAmount, 64); err == nil {
		order.TotalAmount = totalAmount
	}

	return order
}

// convertModelOrderToCreateParams converts model.PurchaseOrder to sqlc CreateOrderParams
func convertModelOrderToCreateParams(order model.PurchaseOrder) db.CreateOrderParams {
	return db.CreateOrderParams{
		ID:          order.ID,
		VendorID:    order.VendorID,
		Status:      string(order.Status),
		TotalAmount: strconv.FormatFloat(order.TotalAmount, 'f', 2, 64),
		CreatedAt:   order.CreatedAt,
		UpdatedAt:   order.UpdatedAt,
	}
}

// convertModelOrderToUpdateParams converts model.PurchaseOrder to sqlc UpdateOrderParams
func convertModelOrderToUpdateParams(order model.PurchaseOrder) db.UpdateOrderParams {
	return db.UpdateOrderParams{
		ID:          order.ID,
		VendorID:    order.VendorID,
		Status:      string(order.Status),
		TotalAmount: strconv.FormatFloat(order.TotalAmount, 'f', 2, 64),
		UpdatedAt:   order.UpdatedAt,
	}
}

// convertDBOrderItemToModel converts sqlc generated db.PurchaseOrderItem to model.PurchaseOrderItem
func convertDBOrderItemToModel(dbItem db.PurchaseOrderItem) model.PurchaseOrderItem {
	item := model.PurchaseOrderItem{
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

// convertModelOrderItemToCreateParams converts model.PurchaseOrderItem to sqlc CreateOrderItemParams
func convertModelOrderItemToCreateParams(item model.PurchaseOrderItem) db.CreateOrderItemParams {
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

func (s *Storage) CreateOrder(ctx context.Context, order model.PurchaseOrder) error {
	params := convertModelOrderToCreateParams(order)
	if err := s.queries.CreateOrder(ctx, params); err != nil {
		return errors.ErrInternalServerError
	}
	return nil
}

func (s *Storage) GetOrderByID(ctx context.Context, id string) (model.PurchaseOrder, error) {
	orderID, err := uuid.Parse(id)
	if err != nil {
		return model.PurchaseOrder{}, errors.ErrBadRequest
	}

	dbOrder, err := s.queries.GetOrderByID(ctx, orderID)
	if err == sql.ErrNoRows {
		return model.PurchaseOrder{}, errors.ErrNotFound
	}
	if err != nil {
		return model.PurchaseOrder{}, errors.ErrInternalServerError
	}

	return convertDBOrderToModel(dbOrder), nil
}

func (s *Storage) ListOrders(ctx context.Context, limit, offset int) ([]model.PurchaseOrder, error) {
	params := db.ListOrdersParams{
		Limit:  int32(limit),
		Offset: int32(offset),
	}

	dbOrders, err := s.queries.ListOrders(ctx, params)
	if err != nil {
		return nil, errors.ErrInternalServerError
	}

	orders := make([]model.PurchaseOrder, 0, len(dbOrders))
	for _, dbOrder := range dbOrders {
		orders = append(orders, convertDBOrderToModel(dbOrder))
	}

	return orders, nil
}

func (s *Storage) UpdateOrder(ctx context.Context, order model.PurchaseOrder) error {
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

func (s *Storage) UpdateOrderStatus(ctx context.Context, id string, status model.PurchaseOrderStatus) error {
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

func (s *Storage) CreateOrderItem(ctx context.Context, item model.PurchaseOrderItem) error {
	params := convertModelOrderItemToCreateParams(item)
	if err := s.queries.CreateOrderItem(ctx, params); err != nil {
		return errors.ErrInternalServerError
	}
	return nil
}

func (s *Storage) CreateOrderItems(ctx context.Context, items []model.PurchaseOrderItem) error {
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

func (s *Storage) GetOrderItemsByOrderID(ctx context.Context, orderID string) ([]model.PurchaseOrderItem, error) {
	orderUUID, err := uuid.Parse(orderID)
	if err != nil {
		return nil, errors.ErrBadRequest
	}

	dbItems, err := s.queries.GetOrderItemsByOrderID(ctx, orderUUID)
	if err != nil {
		return nil, errors.ErrInternalServerError
	}

	items := make([]model.PurchaseOrderItem, 0, len(dbItems))
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
