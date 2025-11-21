package service

import (
	"context"
	"fmt"
	"microservice-challenge/package/errors"
	"microservice-challenge/package/log"
	"microservice-challenge/package/middleware"
	natsclient "microservice-challenge/package/nats"
	"microservice-challenge/services/sales/client"
	"microservice-challenge/services/sales/model"
	"microservice-challenge/services/sales/storage"
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

func (s *Service) CreateOrder(ctx context.Context, req model.CreateOrderRequest) (model.SalesOrderWithItems, error) {
	token, err := s.getTokenFromContext(ctx)
	if err != nil {
		s.logger.Error(ctx, "failed to get token from context", zap.Error(err))
		return model.SalesOrderWithItems{}, errors.ErrInternalServerError
	}

	_, err = s.contactClient.GetCustomerByID(ctx, req.CustomerID.String(), token)
	if err != nil {
		s.logger.Error(ctx, "failed to validate customer", zap.String("customer_id", req.CustomerID.String()), zap.Error(err), zap.String("error_type", fmt.Sprintf("%T", err)), zap.String("error_msg", err.Error()))
		if err == errors.ErrNotFound {
			return model.SalesOrderWithItems{}, errors.ErrBadRequest
		}
		return model.SalesOrderWithItems{}, errors.ErrInternalServerError
	}

	order := model.SalesOrder{
		ID:          uuid.New(),
		CustomerID:  req.CustomerID,
		Status:      model.OrderStatusDraft,
		TotalAmount: 0,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	type itemResult struct {
		item     model.OrderItem
		subtotal float64
		err      error
		itemReq  model.CreateOrderItemRequest
	}

	results := make(chan itemResult, len(req.Items))
	var wg sync.WaitGroup

	for _, itemReq := range req.Items {
		wg.Add(1)
		go func(ir model.CreateOrderItemRequest) {
			defer wg.Done()
			inventoryItem, err := s.inventoryClient.GetItemByID(ctx, ir.ItemID.String(), token)
			if err != nil {
				results <- itemResult{err: err, itemReq: ir}
				return
			}

			subtotal := inventoryItem.UnitPrice * float64(ir.Quantity)
			item := model.OrderItem{
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
	items := make([]model.OrderItem, 0, len(req.Items))

	for result := range results {
		if result.err != nil {
			s.logger.Error(ctx, "failed to validate item", zap.String("item_id", result.itemReq.ItemID.String()), zap.Error(result.err))
			if result.err == errors.ErrNotFound {
				return model.SalesOrderWithItems{}, errors.ErrBadRequest
			}
			return model.SalesOrderWithItems{}, result.err
		}
		items = append(items, result.item)
		totalAmount += result.subtotal
	}

	order.TotalAmount = totalAmount

	if err := s.storage.CreateOrder(ctx, order); err != nil {
		return model.SalesOrderWithItems{}, err
	}

	if err := s.storage.CreateOrderItems(ctx, items); err != nil {
		return model.SalesOrderWithItems{}, err
	}

	return model.SalesOrderWithItems{
		SalesOrder: order,
		Items:      items,
	}, nil
}

func (s *Service) GetOrderByID(ctx context.Context, id string) (model.SalesOrderWithItems, error) {
	order, err := s.storage.GetOrderByID(ctx, id)
	if err != nil {
		return model.SalesOrderWithItems{}, err
	}

	items, err := s.storage.GetOrderItemsByOrderID(ctx, id)
	if err != nil {
		return model.SalesOrderWithItems{}, err
	}

	return model.SalesOrderWithItems{
		SalesOrder: order,
		Items:      items,
	}, nil
}

func (s *Service) ListOrders(ctx context.Context, limit, offset int) ([]model.SalesOrder, error) {
	return s.storage.ListOrders(ctx, limit, offset)
}

func (s *Service) UpdateOrder(ctx context.Context, id string, req model.UpdateOrderRequest) (model.SalesOrderWithItems, error) {
	order, err := s.storage.GetOrderByID(ctx, id)
	if err != nil {
		return model.SalesOrderWithItems{}, err
	}

	if order.Status != model.OrderStatusDraft {
		return model.SalesOrderWithItems{}, errors.ErrBadRequest
	}

	if err := s.storage.DeleteOrderItemsByOrderID(ctx, id); err != nil {
		return model.SalesOrderWithItems{}, err
	}

	token, err := s.getTokenFromContext(ctx)
	if err != nil {
		return model.SalesOrderWithItems{}, errors.ErrInternalServerError
	}

	var totalAmount float64
	items := make([]model.OrderItem, 0, len(req.Items))

	for _, itemReq := range req.Items {
		inventoryItem, err := s.inventoryClient.GetItemByID(ctx, itemReq.ItemID.String(), token)
		if err != nil {
			s.logger.Error(ctx, "failed to validate item", zap.String("item_id", itemReq.ItemID.String()), zap.Error(err))
			if err == errors.ErrNotFound {
				return model.SalesOrderWithItems{}, errors.ErrBadRequest
			}
			return model.SalesOrderWithItems{}, err
		}

		subtotal := inventoryItem.UnitPrice * float64(itemReq.Quantity)
		item := model.OrderItem{
			ID:        uuid.New(),
			OrderID:   order.ID,
			ItemID:    itemReq.ItemID,
			Quantity:  itemReq.Quantity,
			UnitPrice: inventoryItem.UnitPrice,
			Subtotal:  subtotal,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		items = append(items, item)
		totalAmount += subtotal
	}

	order.TotalAmount = totalAmount
	order.UpdatedAt = time.Now()

	if err := s.storage.UpdateOrder(ctx, order); err != nil {
		return model.SalesOrderWithItems{}, err
	}

	if err := s.storage.DeleteOrderItemsByOrderID(ctx, id); err != nil {
		return model.SalesOrderWithItems{}, err
	}

	if err := s.storage.CreateOrderItems(ctx, items); err != nil {
		return model.SalesOrderWithItems{}, err
	}

	return model.SalesOrderWithItems{
		SalesOrder: order,
		Items:      items,
	}, nil
}

func (s *Service) ConfirmOrder(ctx context.Context, id string) (model.SalesOrderWithItems, error) {
	order, err := s.storage.GetOrderByID(ctx, id)
	if err != nil {
		return model.SalesOrderWithItems{}, err
	}

	if order.Status != model.OrderStatusDraft {
		return model.SalesOrderWithItems{}, errors.ErrBadRequest
	}

	items, err := s.storage.GetOrderItemsByOrderID(ctx, id)
	if err != nil {
		return model.SalesOrderWithItems{}, err
	}

	if err := s.storage.UpdateOrderStatus(ctx, id, model.OrderStatusConfirmed); err != nil {
		return model.SalesOrderWithItems{}, err
	}

	order.Status = model.OrderStatusConfirmed
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
		"event_type":   "sales.order.confirmed",
		"order_id":     order.ID.String(),
		"customer_id":  order.CustomerID.String(),
		"items":        eventItems,
		"total_amount": order.TotalAmount,
		"timestamp":    time.Now().Format(time.RFC3339),
	}

	if err := s.natsClient.Publish("sales.order.confirmed", event); err != nil {
		s.logger.Error(ctx, "failed to publish sales.order.confirmed event", zap.Error(err))
	} else {
		s.logger.Info(ctx, "published sales.order.confirmed event",
			zap.String("order_id", order.ID.String()),
			zap.String("customer_id", order.CustomerID.String()),
		)
	}

	return model.SalesOrderWithItems{
		SalesOrder: order,
		Items:      items,
	}, nil
}

func (s *Service) PayOrder(ctx context.Context, id string) (model.SalesOrderWithItems, error) {
	order, err := s.storage.GetOrderByID(ctx, id)
	if err != nil {
		return model.SalesOrderWithItems{}, err
	}

	if order.Status != model.OrderStatusConfirmed {
		return model.SalesOrderWithItems{}, errors.ErrBadRequest
	}

	if err := s.storage.UpdateOrderStatus(ctx, id, model.OrderStatusPaid); err != nil {
		return model.SalesOrderWithItems{}, err
	}

	order.Status = model.OrderStatusPaid
	order.UpdatedAt = time.Now()

	items, err := s.storage.GetOrderItemsByOrderID(ctx, id)
	if err != nil {
		return model.SalesOrderWithItems{}, err
	}

	return model.SalesOrderWithItems{
		SalesOrder: order,
		Items:      items,
	}, nil
}
