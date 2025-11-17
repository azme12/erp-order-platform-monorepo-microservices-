package postgresql

import (
	"context"
	"database/sql"
	"microservice-challenge/package/errors"
	"microservice-challenge/services/inventory/model"
	"strings"

	"github.com/google/uuid"
	"github.com/lib/pq"
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

func (s *Storage) CreateItem(ctx context.Context, item model.Item) error {
	query := `
		INSERT INTO items (id, name, description, sku, unit_price, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`

	_, err := s.db.ExecContext(ctx, query,
		item.ID,
		strings.TrimSpace(item.Name),
		strings.TrimSpace(item.Description),
		strings.ToUpper(strings.TrimSpace(item.SKU)),
		item.UnitPrice,
		item.CreatedAt,
		item.UpdatedAt,
	)

	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			if pqErr.Code == "23505" {
				return errors.ErrConflict
			}
		}
		return errors.ErrInternalServerError
	}

	return nil
}

func (s *Storage) GetItemByID(ctx context.Context, id string) (model.Item, error) {
	var item model.Item

	itemID, err := uuid.Parse(id)
	if err != nil {
		return item, errors.ErrBadRequest
	}

	query := `
		SELECT id, name, description, sku, unit_price, created_at, updated_at
		FROM items
		WHERE id = $1
	`

	err = s.db.QueryRowContext(ctx, query, itemID).Scan(
		&item.ID,
		&item.Name,
		&item.Description,
		&item.SKU,
		&item.UnitPrice,
		&item.CreatedAt,
		&item.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return item, errors.ErrNotFound
	}

	if err != nil {
		return item, errors.ErrInternalServerError
	}

	return item, nil
}

func (s *Storage) GetItemBySKU(ctx context.Context, sku string) (model.Item, error) {
	var item model.Item

	query := `
		SELECT id, name, description, sku, unit_price, created_at, updated_at
		FROM items
		WHERE sku = $1
	`

	err := s.db.QueryRowContext(ctx, query, strings.ToUpper(strings.TrimSpace(sku))).Scan(
		&item.ID,
		&item.Name,
		&item.Description,
		&item.SKU,
		&item.UnitPrice,
		&item.CreatedAt,
		&item.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return item, errors.ErrNotFound
	}

	if err != nil {
		return item, errors.ErrInternalServerError
	}

	return item, nil
}

func (s *Storage) ListItems(ctx context.Context, limit, offset int) ([]model.Item, error) {
	query := `
		SELECT id, name, description, sku, unit_price, created_at, updated_at
		FROM items
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := s.db.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, errors.ErrInternalServerError
	}
	defer rows.Close()

	items := make([]model.Item, 0, limit)
	for rows.Next() {
		var item model.Item
		if err := rows.Scan(
			&item.ID,
			&item.Name,
			&item.Description,
			&item.SKU,
			&item.UnitPrice,
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

func (s *Storage) UpdateItem(ctx context.Context, item model.Item) error {
	query := `
		UPDATE items
		SET name = $2, description = $3, sku = $4, unit_price = $5, updated_at = $6
		WHERE id = $1
	`

	result, err := s.db.ExecContext(ctx, query,
		item.ID,
		strings.TrimSpace(item.Name),
		strings.TrimSpace(item.Description),
		strings.ToUpper(strings.TrimSpace(item.SKU)),
		item.UnitPrice,
		item.UpdatedAt,
	)

	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			if pqErr.Code == "23505" {
				return errors.ErrConflict
			}
		}
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

func (s *Storage) DeleteItem(ctx context.Context, id string) error {
	itemID, err := uuid.Parse(id)
	if err != nil {
		return errors.ErrBadRequest
	}

	query := `DELETE FROM items WHERE id = $1`

	result, err := s.db.ExecContext(ctx, query, itemID)
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

func (s *Storage) GetStockByItemID(ctx context.Context, itemID string) (model.Stock, error) {
	var stock model.Stock

	itemUUID, err := uuid.Parse(itemID)
	if err != nil {
		return stock, errors.ErrBadRequest
	}

	query := `
		SELECT id, item_id, quantity, created_at, updated_at
		FROM stock
		WHERE item_id = $1
	`

	err = s.db.QueryRowContext(ctx, query, itemUUID).Scan(
		&stock.ID,
		&stock.ItemID,
		&stock.Quantity,
		&stock.CreatedAt,
		&stock.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return stock, errors.ErrNotFound
	}

	if err != nil {
		return stock, errors.ErrInternalServerError
	}

	return stock, nil
}

func (s *Storage) CreateStock(ctx context.Context, stock model.Stock) error {
	query := `
		INSERT INTO stock (id, item_id, quantity, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5)
	`

	_, err := s.db.ExecContext(ctx, query,
		stock.ID,
		stock.ItemID,
		stock.Quantity,
		stock.CreatedAt,
		stock.UpdatedAt,
	)

	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			if pqErr.Code == "23505" {
				return errors.ErrConflict
			}
		}
		return errors.ErrInternalServerError
	}

	return nil
}

func (s *Storage) UpdateStock(ctx context.Context, stock model.Stock) error {
	query := `
		UPDATE stock
		SET quantity = $2, updated_at = $3
		WHERE item_id = $1
	`

	result, err := s.db.ExecContext(ctx, query,
		stock.ItemID,
		stock.Quantity,
		stock.UpdatedAt,
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

func (s *Storage) AdjustStock(ctx context.Context, itemID string, quantityDelta int) error {
	itemUUID, err := uuid.Parse(itemID)
	if err != nil {
		return errors.ErrBadRequest
	}

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return errors.ErrInternalServerError
	}
	defer tx.Rollback()

	var currentQuantity int
	query := `SELECT quantity FROM stock WHERE item_id = $1 FOR UPDATE`
	err = tx.QueryRowContext(ctx, query, itemUUID).Scan(&currentQuantity)
	if err == sql.ErrNoRows {
		return errors.ErrNotFound
	}
	if err != nil {
		return errors.ErrInternalServerError
	}

	newQuantity := currentQuantity + quantityDelta
	if newQuantity < 0 {
		return errors.ErrBadRequest
	}

	updateQuery := `
		UPDATE stock
		SET quantity = $2, updated_at = CURRENT_TIMESTAMP
		WHERE item_id = $1
	`

	_, err = tx.ExecContext(ctx, updateQuery, itemUUID, newQuantity)
	if err != nil {
		return errors.ErrInternalServerError
	}

	if err := tx.Commit(); err != nil {
		return errors.ErrInternalServerError
	}

	return nil
}
