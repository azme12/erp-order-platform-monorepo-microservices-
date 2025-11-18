package postgresql

import (
	"context"
	"database/sql"
	"fmt"
	"microservice-challenge/package/errors"
	"microservice-challenge/services/sales/model"
	"strings"

	"github.com/google/uuid"
)

type Storage struct {
	db *sql.DB
}

func NewStorage(db *sql.DB) *Storage {
	return &Storage{db: db}
}

func (s *Storage) DB() *sql.DB {
	return s.db
}

func (s *Storage) CreateOrder(ctx context.Context, order model.SalesOrder) error {
	query := `
		INSERT INTO sales_orders (id, customer_id, status, total_amount, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)
	`

	_, err := s.db.ExecContext(ctx, query,
		order.ID,
		order.CustomerID,
		order.Status,
		order.TotalAmount,
		order.CreatedAt,
		order.UpdatedAt,
	)

	if err != nil {
		return errors.ErrInternalServerError
	}

	return nil
}

func (s *Storage) GetOrderByID(ctx context.Context, id string) (model.SalesOrder, error) {
	var order model.SalesOrder

	orderID, err := uuid.Parse(id)
	if err != nil {
		return order, errors.ErrBadRequest
	}

	query := `
		SELECT id, customer_id, status, total_amount, created_at, updated_at
		FROM sales_orders
		WHERE id = $1
	`

	err = s.db.QueryRowContext(ctx, query, orderID).Scan(
		&order.ID,
		&order.CustomerID,
		&order.Status,
		&order.TotalAmount,
		&order.CreatedAt,
		&order.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return order, errors.ErrNotFound
	}

	if err != nil {
		return order, errors.ErrInternalServerError
	}

	return order, nil
}

func (s *Storage) ListOrders(ctx context.Context, limit, offset int) ([]model.SalesOrder, error) {
	query := `
		SELECT id, customer_id, status, total_amount, created_at, updated_at
		FROM sales_orders
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := s.db.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, errors.ErrInternalServerError
	}
	defer rows.Close()

	orders := make([]model.SalesOrder, 0, limit)
	for rows.Next() {
		var order model.SalesOrder
		if err := rows.Scan(
			&order.ID,
			&order.CustomerID,
			&order.Status,
			&order.TotalAmount,
			&order.CreatedAt,
			&order.UpdatedAt,
		); err != nil {
			return nil, errors.ErrInternalServerError
		}
		orders = append(orders, order)
	}

	if err := rows.Err(); err != nil {
		return nil, errors.ErrInternalServerError
	}

	return orders, nil
}

func (s *Storage) UpdateOrder(ctx context.Context, order model.SalesOrder) error {
	query := `
		UPDATE sales_orders
		SET customer_id = $2, status = $3, total_amount = $4, updated_at = $5
		WHERE id = $1
	`

	result, err := s.db.ExecContext(ctx, query,
		order.ID,
		order.CustomerID,
		order.Status,
		order.TotalAmount,
		order.UpdatedAt,
	)

	if err != nil {
		return errors.ErrInternalServerError
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return errors.ErrInternalServerError
	}

	if rowsAffected == 0 {
		return errors.ErrNotFound
	}

	return nil
}

func (s *Storage) UpdateOrderStatus(ctx context.Context, id string, status model.OrderStatus) error {
	orderID, err := uuid.Parse(id)
	if err != nil {
		return errors.ErrBadRequest
	}

	query := `
		UPDATE sales_orders
		SET status = $2, updated_at = CURRENT_TIMESTAMP
		WHERE id = $1
	`

	result, err := s.db.ExecContext(ctx, query, orderID, status)
	if err != nil {
		return errors.ErrInternalServerError
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return errors.ErrInternalServerError
	}

	if rowsAffected == 0 {
		return errors.ErrNotFound
	}

	return nil
}

func (s *Storage) CreateOrderItem(ctx context.Context, item model.OrderItem) error {
	query := `
		INSERT INTO order_items (id, order_id, item_id, quantity, unit_price, subtotal, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`

	_, err := s.db.ExecContext(ctx, query,
		item.ID,
		item.OrderID,
		item.ItemID,
		item.Quantity,
		item.UnitPrice,
		item.Subtotal,
		item.CreatedAt,
		item.UpdatedAt,
	)

	if err != nil {
		return errors.ErrInternalServerError
	}

	return nil
}

func (s *Storage) CreateOrderItems(ctx context.Context, items []model.OrderItem) error {
	if len(items) == 0 {
		return nil
	}

	query := `
		INSERT INTO order_items (id, order_id, item_id, quantity, unit_price, subtotal, created_at, updated_at)
		VALUES `

	values := make([]interface{}, 0, len(items)*8)
	placeholders := make([]string, 0, len(items))

	for i, item := range items {
		offset := i * 8
		placeholders = append(placeholders, fmt.Sprintf("($%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d)",
			offset+1, offset+2, offset+3, offset+4, offset+5, offset+6, offset+7, offset+8))
		values = append(values,
			item.ID,
			item.OrderID,
			item.ItemID,
			item.Quantity,
			item.UnitPrice,
			item.Subtotal,
			item.CreatedAt,
			item.UpdatedAt,
		)
	}

	query += strings.Join(placeholders, ", ")

	_, err := s.db.ExecContext(ctx, query, values...)
	if err != nil {
		return errors.ErrInternalServerError
	}

	return nil
}

func (s *Storage) GetOrderItemsByOrderID(ctx context.Context, orderID string) ([]model.OrderItem, error) {
	orderUUID, err := uuid.Parse(orderID)
	if err != nil {
		return nil, errors.ErrBadRequest
	}

	query := `
		SELECT id, order_id, item_id, quantity, unit_price, subtotal, created_at, updated_at
		FROM order_items
		WHERE order_id = $1
		ORDER BY created_at ASC
	`

	rows, err := s.db.QueryContext(ctx, query, orderUUID)
	if err != nil {
		return nil, errors.ErrInternalServerError
	}
	defer rows.Close()

	items := make([]model.OrderItem, 0)
	for rows.Next() {
		var item model.OrderItem
		if err := rows.Scan(
			&item.ID,
			&item.OrderID,
			&item.ItemID,
			&item.Quantity,
			&item.UnitPrice,
			&item.Subtotal,
			&item.CreatedAt,
			&item.UpdatedAt,
		); err != nil {
			return nil, errors.ErrInternalServerError
		}
		items = append(items, item)
	}

	if err := rows.Err(); err != nil {
		return nil, errors.ErrInternalServerError
	}

	return items, nil
}

func (s *Storage) DeleteOrderItemsByOrderID(ctx context.Context, orderID string) error {
	orderUUID, err := uuid.Parse(orderID)
	if err != nil {
		return errors.ErrBadRequest
	}

	query := `DELETE FROM order_items WHERE order_id = $1`

	_, err = s.db.ExecContext(ctx, query, orderUUID)
	if err != nil {
		return errors.ErrInternalServerError
	}

	return nil
}
