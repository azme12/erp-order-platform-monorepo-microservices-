package postgresql

import (
	"context"
	"database/sql"
	"microservice-challenge/package/errors"
	"microservice-challenge/services/inventory/model"
	"microservice-challenge/services/inventory/storage/postgresql/db"
	"strconv"
	"strings"

	"github.com/google/uuid"
	"github.com/lib/pq"
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

// convertDBItemToModel converts sqlc generated db.Item to model.Item
func convertDBItemToModel(dbItem db.Item) model.Item {
	item := model.Item{
		ID:        dbItem.ID,
		Name:      dbItem.Name,
		SKU:       dbItem.Sku,
		CreatedAt: dbItem.CreatedAt,
		UpdatedAt: dbItem.UpdatedAt,
	}

	if dbItem.Description.Valid {
		item.Description = dbItem.Description.String
	}

	if unitPrice, err := strconv.ParseFloat(dbItem.UnitPrice, 64); err == nil {
		item.UnitPrice = unitPrice
	}

	return item
}

// convertModelItemToCreateParams converts model.Item to sqlc CreateItemParams
func convertModelItemToCreateParams(item model.Item) db.CreateItemParams {
	params := db.CreateItemParams{
		ID:        item.ID,
		Name:      item.Name,
		Sku:       item.SKU,
		UnitPrice: strconv.FormatFloat(item.UnitPrice, 'f', 2, 64),
		CreatedAt: item.CreatedAt,
		UpdatedAt: item.UpdatedAt,
	}

	if item.Description != "" {
		params.Description = sql.NullString{
			String: item.Description,
			Valid:  true,
		}
	}

	return params
}

// convertModelItemToUpdateParams converts model.Item to sqlc UpdateItemParams
func convertModelItemToUpdateParams(item model.Item) db.UpdateItemParams {
	params := db.UpdateItemParams{
		ID:        item.ID,
		Name:      item.Name,
		Sku:       item.SKU,
		UnitPrice: strconv.FormatFloat(item.UnitPrice, 'f', 2, 64),
		UpdatedAt: item.UpdatedAt,
	}

	if item.Description != "" {
		params.Description = sql.NullString{
			String: item.Description,
			Valid:  true,
		}
	}

	return params
}

// convertDBStockToModel converts sqlc generated db.Stock to model.Stock
func convertDBStockToModel(dbStock db.Stock) model.Stock {
	return model.Stock{
		ID:        dbStock.ID,
		ItemID:    dbStock.ItemID,
		Quantity:  int(dbStock.Quantity),
		CreatedAt: dbStock.CreatedAt,
		UpdatedAt: dbStock.UpdatedAt,
	}
}

// convertModelStockToCreateParams converts model.Stock to sqlc CreateStockParams
func convertModelStockToCreateParams(stock model.Stock) db.CreateStockParams {
	return db.CreateStockParams{
		ID:        stock.ID,
		ItemID:    stock.ItemID,
		Quantity:  int32(stock.Quantity),
		CreatedAt: stock.CreatedAt,
		UpdatedAt: stock.UpdatedAt,
	}
}

// convertModelStockToUpdateParams converts model.Stock to sqlc UpdateStockParams
func convertModelStockToUpdateParams(stock model.Stock) db.UpdateStockParams {
	return db.UpdateStockParams{
		ItemID:    stock.ItemID,
		Quantity:  int32(stock.Quantity),
		UpdatedAt: stock.UpdatedAt,
	}
}

func (s *Storage) CreateItem(ctx context.Context, item model.Item) error {
	item.SKU = strings.ToUpper(strings.TrimSpace(item.SKU))
	item.Name = strings.TrimSpace(item.Name)
	item.Description = strings.TrimSpace(item.Description)

	params := convertModelItemToCreateParams(item)
	err := s.queries.CreateItem(ctx, params)

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
	itemID, err := uuid.Parse(id)
	if err != nil {
		return model.Item{}, errors.ErrBadRequest
	}

	dbItem, err := s.queries.GetItemByID(ctx, itemID)
	if err == sql.ErrNoRows {
		return model.Item{}, errors.ErrNotFound
	}
	if err != nil {
		return model.Item{}, errors.ErrInternalServerError
	}

	return convertDBItemToModel(dbItem), nil
}

func (s *Storage) GetItemBySKU(ctx context.Context, sku string) (model.Item, error) {
	normalizedSKU := strings.ToUpper(strings.TrimSpace(sku))
	dbItem, err := s.queries.GetItemBySKU(ctx, normalizedSKU)
	if err == sql.ErrNoRows {
		return model.Item{}, errors.ErrNotFound
	}
	if err != nil {
		return model.Item{}, errors.ErrInternalServerError
	}

	return convertDBItemToModel(dbItem), nil
}

func (s *Storage) ListItems(ctx context.Context, limit, offset int) ([]model.Item, error) {
	params := db.ListItemsParams{
		Limit:  int32(limit),
		Offset: int32(offset),
	}

	dbItems, err := s.queries.ListItems(ctx, params)
	if err != nil {
		return nil, errors.ErrInternalServerError
	}

	items := make([]model.Item, 0, len(dbItems))
	for _, dbItem := range dbItems {
		items = append(items, convertDBItemToModel(dbItem))
	}

	return items, nil
}

func (s *Storage) UpdateItem(ctx context.Context, item model.Item) error {
	_, err := s.queries.GetItemByID(ctx, item.ID)
	if err == sql.ErrNoRows {
		return errors.ErrNotFound
	}
	if err != nil {
		return errors.ErrInternalServerError
	}

	item.SKU = strings.ToUpper(strings.TrimSpace(item.SKU))
	item.Name = strings.TrimSpace(item.Name)
	item.Description = strings.TrimSpace(item.Description)

	params := convertModelItemToUpdateParams(item)
	err = s.queries.UpdateItem(ctx, params)

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

func (s *Storage) DeleteItem(ctx context.Context, id string) error {
	itemID, err := uuid.Parse(id)
	if err != nil {
		return errors.ErrBadRequest
	}

	_, err = s.queries.GetItemByID(ctx, itemID)
	if err == sql.ErrNoRows {
		return errors.ErrNotFound
	}
	if err != nil {
		return errors.ErrInternalServerError
	}

	err = s.queries.DeleteItem(ctx, itemID)
	if err != nil {
		return errors.ErrInternalServerError
	}

	return nil
}

func (s *Storage) GetStockByItemID(ctx context.Context, itemID string) (model.Stock, error) {
	itemUUID, err := uuid.Parse(itemID)
	if err != nil {
		return model.Stock{}, errors.ErrBadRequest
	}

	dbStock, err := s.queries.GetStockByItemID(ctx, itemUUID)
	if err == sql.ErrNoRows {
		return model.Stock{}, errors.ErrNotFound
	}
	if err != nil {
		return model.Stock{}, errors.ErrInternalServerError
	}

	return convertDBStockToModel(dbStock), nil
}

func (s *Storage) CreateStock(ctx context.Context, stock model.Stock) error {
	params := convertModelStockToCreateParams(stock)
	err := s.queries.CreateStock(ctx, params)

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
	params := convertModelStockToUpdateParams(stock)
	err := s.queries.UpdateStock(ctx, params)

	if err != nil {
		return errors.ErrInternalServerError
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

	qtx := s.queries.WithTx(tx)

	currentQuantity, err := qtx.GetStockQuantityForUpdate(ctx, itemUUID)
	if err == sql.ErrNoRows {
		return errors.ErrNotFound
	}
	if err != nil {
		return errors.ErrInternalServerError
	}

	newQuantity := int(currentQuantity) + quantityDelta
	if newQuantity < 0 {
		return errors.ErrBadRequest
	}

	adjustParams := db.AdjustStockParams{
		ItemID:   itemUUID,
		Quantity: int32(quantityDelta),
	}
	if err := qtx.AdjustStock(ctx, adjustParams); err != nil {
		return errors.ErrInternalServerError
	}

	if err := tx.Commit(); err != nil {
		return errors.ErrInternalServerError
	}

	return nil
}
