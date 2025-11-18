package service

import (
	"context"
	"fmt"
	"microservice-challenge/package/errors"
	"microservice-challenge/package/log"
	"microservice-challenge/package/middleware"
	natsclient "microservice-challenge/package/nats"
	"microservice-challenge/services/purchase/client"
	"microservice-challenge/services/purchase/model"
	"microservice-challenge/services/purchase/storage"
	"sync"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

type Service struct {
	storage         storage.Storage
	natsClient      *natsclient.Client
	contactClient   *client.ContactClient
	inventoryClient *client.InventoryClient
	authClient      *client.AuthClient
	logger          log.Logger
}

func NewService(storage storage.Storage, natsClient *natsclient.Client, contactClient *client.ContactClient, inventoryClient *client.InventoryClient, authClient *client.AuthClient, logger log.Logger) *Service {
	return &Service{
		storage:         storage,
		natsClient:      natsClient,
		contactClient:   contactClient,
		inventoryClient: inventoryClient,
		authClient:      authClient,
		logger:          logger,
	}
}

func (s *Service) getTokenFromContext(ctx context.Context) (string, error) {

	token := ctx.Value(middleware.GetTokenKey())
	if tokenStr, ok := token.(string); ok && tokenStr != "" {
		return tokenStr, nil
	}

	serviceToken, err := s.authClient.GetServiceToken(ctx)
	if err != nil {
		s.logger.Error(ctx, "failed to get service token from auth service", zap.Error(err))
		return "", errors.ErrInternalServerError
	}

	return serviceToken, nil
}

func (s *Service) CreateOrder(ctx context.Context, req model.CreatePurchaseOrderRequest) (model.PurchaseOrderWithItems, error) {
	token, err := s.getTokenFromContext(ctx)
	if err != nil {
		s.logger.Error(ctx, "failed to get token from context", zap.Error(err))
		return model.PurchaseOrderWithItems{}, errors.ErrInternalServerError
	}

	_, err = s.contactClient.GetVendorByID(ctx, req.VendorID.String(), token)
	if err != nil {
		s.logger.Error(ctx, "failed to validate vendor", zap.String("vendor_id", req.VendorID.String()), zap.Error(err), zap.String("error_type", fmt.Sprintf("%T", err)), zap.String("error_msg", err.Error()))
		if err == errors.ErrNotFound {
			return model.PurchaseOrderWithItems{}, errors.ErrBadRequest
		}
		return model.PurchaseOrderWithItems{}, errors.ErrInternalServerError
	}

	order := model.PurchaseOrder{
		ID:          uuid.New(),
		VendorID:    req.VendorID,
		Status:      model.PurchaseOrderStatusDraft,
		TotalAmount: 0,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	// Validate all items in parallel for better performance
	type itemResult struct {
		item     model.PurchaseOrderItem
		subtotal float64
		err      error
		itemReq  model.CreatePurchaseOrderItemRequest
	}

	results := make(chan itemResult, len(req.Items))
	var wg sync.WaitGroup

	for _, itemReq := range req.Items {
		wg.Add(1)
		go func(ir model.CreatePurchaseOrderItemRequest) {
			defer wg.Done()
			inventoryItem, err := s.inventoryClient.GetItemByID(ctx, ir.ItemID.String(), token)
			if err != nil {
				results <- itemResult{err: err, itemReq: ir}
				return
			}

			subtotal := inventoryItem.UnitPrice * float64(ir.Quantity)
			item := model.PurchaseOrderItem{
				ID:        uuid.New(),
				OrderID:   order.ID,
				ItemID:    ir.ItemID,
				Quantity:  ir.Quantity,
				UnitPrice: inventoryItem.UnitPrice,
				Subtotal:  subtotal,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			}
			results <- itemResult{item: item, subtotal: subtotal}
		}(itemReq)
	}

	wg.Wait()
	close(results)

	var totalAmount float64
	items := make([]model.PurchaseOrderItem, 0, len(req.Items))

	for result := range results {
		if result.err != nil {
			s.logger.Error(ctx, "failed to validate item", zap.String("item_id", result.itemReq.ItemID.String()), zap.Error(result.err))
			if result.err == errors.ErrNotFound {
				return model.PurchaseOrderWithItems{}, errors.ErrBadRequest
			}
			return model.PurchaseOrderWithItems{}, errors.ErrInternalServerError
		}
		items = append(items, result.item)
		totalAmount += result.subtotal
	}

	order.TotalAmount = totalAmount

	if err := s.storage.CreateOrder(ctx, order); err != nil {
		s.logger.Error(ctx, "failed to store order", zap.Error(err))
		return model.PurchaseOrderWithItems{}, errors.ErrInternalServerError
	}

	if err := s.storage.CreateOrderItems(ctx, items); err != nil {
		s.logger.Error(ctx, "failed to store order items", zap.Error(err))
		return model.PurchaseOrderWithItems{}, errors.ErrInternalServerError
	}

	return model.PurchaseOrderWithItems{
		PurchaseOrder: order,
		Items:         items,
	}, nil
}

func (s *Service) GetOrderByID(ctx context.Context, id string) (model.PurchaseOrderWithItems, error) {
	order, err := s.storage.GetOrderByID(ctx, id)
	if err != nil {
		return model.PurchaseOrderWithItems{}, err
	}

	items, err := s.storage.GetOrderItemsByOrderID(ctx, id)
	if err != nil {
		return model.PurchaseOrderWithItems{}, err
	}

	return model.PurchaseOrderWithItems{
		PurchaseOrder: order,
		Items:         items,
	}, nil
}

func (s *Service) ListOrders(ctx context.Context, limit, offset int) ([]model.PurchaseOrder, error) {
	return s.storage.ListOrders(ctx, limit, offset)
}

func (s *Service) UpdateOrder(ctx context.Context, id string, req model.UpdatePurchaseOrderRequest) (model.PurchaseOrderWithItems, error) {
	order, err := s.storage.GetOrderByID(ctx, id)
	if err != nil {
		return model.PurchaseOrderWithItems{}, err
	}

	if order.Status != model.PurchaseOrderStatusDraft {
		return model.PurchaseOrderWithItems{}, errors.ErrBadRequest
	}

	if err := s.storage.DeleteOrderItemsByOrderID(ctx, id); err != nil {
		return model.PurchaseOrderWithItems{}, err
	}

	token, err := s.getTokenFromContext(ctx)
	if err != nil {
		return model.PurchaseOrderWithItems{}, errors.ErrInternalServerError
	}

	// Validate all items in parallel for better performance
	type itemResult struct {
		item     model.PurchaseOrderItem
		subtotal float64
		err      error
		itemReq  model.CreatePurchaseOrderItemRequest
	}

	results := make(chan itemResult, len(req.Items))
	var wg sync.WaitGroup

	for _, itemReq := range req.Items {
		wg.Add(1)
		go func(ir model.CreatePurchaseOrderItemRequest) {
			defer wg.Done()
			inventoryItem, err := s.inventoryClient.GetItemByID(ctx, ir.ItemID.String(), token)
			if err != nil {
				results <- itemResult{err: err, itemReq: ir}
				return
			}

			subtotal := inventoryItem.UnitPrice * float64(ir.Quantity)
			item := model.PurchaseOrderItem{
				ID:        uuid.New(),
				OrderID:   order.ID,
				ItemID:    ir.ItemID,
				Quantity:  ir.Quantity,
				UnitPrice: inventoryItem.UnitPrice,
				Subtotal:  subtotal,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			}
			results <- itemResult{item: item, subtotal: subtotal}
		}(itemReq)
	}

	wg.Wait()
	close(results)

	var totalAmount float64
	items := make([]model.PurchaseOrderItem, 0, len(req.Items))

	for result := range results {
		if result.err != nil {
			s.logger.Error(ctx, "failed to validate item", zap.String("item_id", result.itemReq.ItemID.String()), zap.Error(result.err))
			if result.err == errors.ErrNotFound {
				return model.PurchaseOrderWithItems{}, errors.ErrBadRequest
			}
			return model.PurchaseOrderWithItems{}, errors.ErrInternalServerError
		}
		items = append(items, result.item)
		totalAmount += result.subtotal
	}

	order.TotalAmount = totalAmount
	order.UpdatedAt = time.Now()

	if err := s.storage.UpdateOrder(ctx, order); err != nil {
		s.logger.Error(ctx, "failed to update order in storage", zap.Error(err))
		return model.PurchaseOrderWithItems{}, errors.ErrInternalServerError
	}

	if err := s.storage.DeleteOrderItemsByOrderID(ctx, id); err != nil {
		s.logger.Error(ctx, "failed to delete existing order items", zap.Error(err))
		return model.PurchaseOrderWithItems{}, errors.ErrInternalServerError
	}

	if err := s.storage.CreateOrderItems(ctx, items); err != nil {
		s.logger.Error(ctx, "failed to create order items", zap.Error(err))
		return model.PurchaseOrderWithItems{}, errors.ErrInternalServerError
	}

	return model.PurchaseOrderWithItems{
		PurchaseOrder: order,
		Items:         items,
	}, nil
}

func (s *Service) ReceiveOrder(ctx context.Context, id string) (model.PurchaseOrderWithItems, error) {
	order, err := s.storage.GetOrderByID(ctx, id)
	if err != nil {
		return model.PurchaseOrderWithItems{}, err
	}

	if order.Status != model.PurchaseOrderStatusDraft {
		return model.PurchaseOrderWithItems{}, errors.ErrBadRequest
	}

	items, err := s.storage.GetOrderItemsByOrderID(ctx, id)
	if err != nil {
		return model.PurchaseOrderWithItems{}, err
	}

	if err := s.storage.UpdateOrderStatus(ctx, id, model.PurchaseOrderStatusReceived); err != nil {
		return model.PurchaseOrderWithItems{}, err
	}

	order.Status = model.PurchaseOrderStatusReceived
	order.UpdatedAt = time.Now()

	eventItems := make([]map[string]interface{}, 0, len(items))
	for _, item := range items {
		eventItems = append(eventItems, map[string]interface{}{
			"item_id":    item.ItemID.String(),
			"quantity":   item.Quantity,
			"unit_price": item.UnitPrice,
			"subtotal":   item.Subtotal,
		})
	}

	event := map[string]interface{}{
		"event_type":   "purchase.order.received",
		"order_id":     order.ID.String(),
		"vendor_id":    order.VendorID.String(),
		"items":        eventItems,
		"total_amount": order.TotalAmount,
		"timestamp":    time.Now().Format(time.RFC3339),
	}

	if err := s.natsClient.Publish("purchase.order.received", event); err != nil {
		s.logger.Error(ctx, "failed to publish purchase.order.received event", zap.Error(err))
	} else {
		s.logger.Info(ctx, "published purchase.order.received event",
			zap.String("order_id", order.ID.String()),
			zap.String("vendor_id", order.VendorID.String()),
		)
	}

	return model.PurchaseOrderWithItems{
		PurchaseOrder: order,
		Items:         items,
	}, nil
}

func (s *Service) PayOrder(ctx context.Context, id string) (model.PurchaseOrderWithItems, error) {
	order, err := s.storage.GetOrderByID(ctx, id)
	if err != nil {
		return model.PurchaseOrderWithItems{}, err
	}

	if order.Status != model.PurchaseOrderStatusReceived {
		return model.PurchaseOrderWithItems{}, errors.ErrBadRequest
	}

	if err := s.storage.UpdateOrderStatus(ctx, id, model.PurchaseOrderStatusPaid); err != nil {
		return model.PurchaseOrderWithItems{}, err
	}

	order.Status = model.PurchaseOrderStatusPaid
	order.UpdatedAt = time.Now()

	items, err := s.storage.GetOrderItemsByOrderID(ctx, id)
	if err != nil {
		return model.PurchaseOrderWithItems{}, err
	}

	return model.PurchaseOrderWithItems{
		PurchaseOrder: order,
		Items:         items,
	}, nil
}
