package service

import (
	"context"
	"encoding/json"
	"microservice-challenge/package/errors"
	"microservice-challenge/package/log"
	natsclient "microservice-challenge/package/nats"
	"microservice-challenge/services/inventory/model"
	"microservice-challenge/services/inventory/storage"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/nats-io/nats.go"
	"go.uber.org/zap"
)

type Service struct {
	storage    storage.Storage
	natsClient *natsclient.Client
	logger     log.Logger
}

func NewService(storage storage.Storage, natsClient *natsclient.Client, logger log.Logger) *Service {
	return &Service{
		storage:    storage,
		natsClient: natsClient,
		logger:     logger,
	}
}

func (s *Service) CreateItem(ctx context.Context, req model.CreateItemRequest) (model.Item, error) {
	sku := strings.ToUpper(strings.TrimSpace(req.SKU))

	_, err := s.storage.GetItemBySKU(ctx, sku)
	if err == nil {
		return model.Item{}, errors.ErrConflict
	}
	if err != errors.ErrNotFound {
		return model.Item{}, err
	}

	item := model.Item{
		ID:          uuid.New(),
		Name:        strings.TrimSpace(req.Name),
		Description: strings.TrimSpace(req.Description),
		SKU:         sku,
		UnitPrice:   req.UnitPrice,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if err := s.storage.CreateItem(ctx, item); err != nil {
		return model.Item{}, err
	}

	stock := model.Stock{
		ID:        uuid.New(),
		ItemID:    item.ID,
		Quantity:  0,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := s.storage.CreateStock(ctx, stock); err != nil {
		s.logger.Error(ctx, "failed to create initial stock", zap.Error(err))
	}

	return item, nil
}

func (s *Service) GetItemByID(ctx context.Context, id string) (model.Item, error) {
	return s.storage.GetItemByID(ctx, id)
}

func (s *Service) ListItems(ctx context.Context, limit, offset int) ([]model.Item, error) {
	return s.storage.ListItems(ctx, limit, offset)
}

func (s *Service) UpdateItem(ctx context.Context, id string, req model.UpdateItemRequest) (model.Item, error) {
	item, err := s.storage.GetItemByID(ctx, id)
	if err != nil {
		return model.Item{}, err
	}

	sku := strings.ToUpper(strings.TrimSpace(req.SKU))
	if sku != item.SKU {
		_, err := s.storage.GetItemBySKU(ctx, sku)
		if err == nil {
			return model.Item{}, errors.ErrConflict
		}
		if err != errors.ErrNotFound {
			return model.Item{}, err
		}
	}

	item.Name = strings.TrimSpace(req.Name)
	item.Description = strings.TrimSpace(req.Description)
	item.SKU = sku
	item.UnitPrice = req.UnitPrice
	item.UpdatedAt = time.Now()

	if err := s.storage.UpdateItem(ctx, item); err != nil {
		return model.Item{}, err
	}

	return item, nil
}

func (s *Service) DeleteItem(ctx context.Context, id string) error {
	return s.storage.DeleteItem(ctx, id)
}

func (s *Service) GetStockByItemID(ctx context.Context, itemID string) (model.Stock, error) {
	return s.storage.GetStockByItemID(ctx, itemID)
}

func (s *Service) AdjustStock(ctx context.Context, itemID string, quantity int) (model.Stock, error) {
	_, err := s.storage.GetItemByID(ctx, itemID)
	if err != nil {
		return model.Stock{}, err
	}

	if err := s.storage.AdjustStock(ctx, itemID, quantity); err != nil {
		return model.Stock{}, err
	}

	stock, err := s.storage.GetStockByItemID(ctx, itemID)
	if err != nil {
		return model.Stock{}, err
	}

	return stock, nil
}

func (s *Service) StartEventSubscriptions(ctx context.Context) error {
	salesSub, err := s.natsClient.Subscribe("sales.order.confirmed", func(msg *nats.Msg) {
		s.handleSalesOrderConfirmed(ctx, msg)
	})
	if err != nil {
		return err
	}
	s.logger.Info(ctx, "subscribed to sales.order.confirmed", zap.String("subscription", salesSub.Subject))

	purchaseSub, err := s.natsClient.Subscribe("purchase.order.received", func(msg *nats.Msg) {
		s.handlePurchaseOrderReceived(ctx, msg)
	})
	if err != nil {
		return err
	}
	s.logger.Info(ctx, "subscribed to purchase.order.received", zap.String("subscription", purchaseSub.Subject))

	return nil
}

func (s *Service) handleSalesOrderConfirmed(ctx context.Context, msg *nats.Msg) {
	var event map[string]interface{}
	if err := json.Unmarshal(msg.Data, &event); err != nil {
		s.logger.Error(ctx, "failed to unmarshal sales.order.confirmed event", zap.Error(err))
		return
	}

	items, ok := event["items"].([]interface{})
	if !ok {
		s.logger.Error(ctx, "invalid items format in sales.order.confirmed event")
		return
	}

	for _, itemData := range items {
		itemMap, ok := itemData.(map[string]interface{})
		if !ok {
			continue
		}

		itemID, ok := itemMap["item_id"].(string)
		if !ok {
			continue
		}

		quantity, ok := itemMap["quantity"].(float64)
		if !ok {
			continue
		}

		if err := s.storage.AdjustStock(ctx, itemID, -int(quantity)); err != nil {
			s.logger.Error(ctx, "failed to decrease stock for sales order",
				zap.String("item_id", itemID),
				zap.Int("quantity", int(quantity)),
				zap.Error(err),
			)
		} else {
			s.logger.Info(ctx, "decreased stock for sales order",
				zap.String("item_id", itemID),
				zap.Int("quantity", int(quantity)),
			)
		}
	}
}

func (s *Service) handlePurchaseOrderReceived(ctx context.Context, msg *nats.Msg) {
	var event map[string]interface{}
	if err := json.Unmarshal(msg.Data, &event); err != nil {
		s.logger.Error(ctx, "failed to unmarshal purchase.order.received event", zap.Error(err))
		return
	}

	items, ok := event["items"].([]interface{})
	if !ok {
		s.logger.Error(ctx, "invalid items format in purchase.order.received event")
		return
	}

	for _, itemData := range items {
		itemMap, ok := itemData.(map[string]interface{})
		if !ok {
			continue
		}

		itemID, ok := itemMap["item_id"].(string)
		if !ok {
			continue
		}

		quantity, ok := itemMap["quantity"].(float64)
		if !ok {
			continue
		}

		if err := s.storage.AdjustStock(ctx, itemID, int(quantity)); err != nil {
			s.logger.Error(ctx, "failed to increase stock for purchase order",
				zap.String("item_id", itemID),
				zap.Int("quantity", int(quantity)),
				zap.Error(err),
			)
		} else {
			s.logger.Info(ctx, "increased stock for purchase order",
				zap.String("item_id", itemID),
				zap.Int("quantity", int(quantity)),
			)
		}
	}
}
