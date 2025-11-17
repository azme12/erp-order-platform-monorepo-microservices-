package purchase

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
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

type Usecase struct {
	storage         storage.Storage
	natsClient      *natsclient.Client
	contactClient   *client.ContactClient
	inventoryClient *client.InventoryClient
	authClient      *client.AuthClient
	logger          log.Logger
}

func NewUsecase(storage storage.Storage, natsClient *natsclient.Client, contactClient *client.ContactClient, inventoryClient *client.InventoryClient, authClient *client.AuthClient, logger log.Logger) *Usecase {
	return &Usecase{
		storage:         storage,
		natsClient:      natsClient,
		contactClient:   contactClient,
		inventoryClient: inventoryClient,
		authClient:      authClient,
		logger:          logger,
	}
}

func (u *Usecase) getTokenFromContext(ctx context.Context) (string, error) {

	token := ctx.Value(middleware.GetTokenKey())
	if tokenStr, ok := token.(string); ok && tokenStr != "" {
		return tokenStr, nil
	}

	serviceToken, err := u.authClient.GetServiceToken(ctx)
	if err != nil {
		u.logger.Error(ctx, "failed to get service token from auth service", zap.Error(err))
		return "", errors.ErrInternalServerError
	}

	return serviceToken, nil
}

func (u *Usecase) CreateOrder(ctx context.Context, req model.CreatePurchaseOrderRequest) (model.PurchaseOrderWithItems, error) {
	token, err := u.getTokenFromContext(ctx)
	if err != nil {
		u.logger.Error(ctx, "failed to get token from context", zap.Error(err))
		return model.PurchaseOrderWithItems{}, errors.ErrInternalServerError
	}

	_, err = u.contactClient.GetVendorByID(ctx, req.VendorID.String(), token)
	if err != nil {
		u.logger.Error(ctx, "failed to validate vendor", zap.String("vendor_id", req.VendorID.String()), zap.Error(err), zap.String("error_type", fmt.Sprintf("%T", err)), zap.String("error_msg", err.Error()))
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

	var totalAmount float64
	items := make([]model.PurchaseOrderItem, 0, len(req.Items))

	for _, itemReq := range req.Items {
		inventoryItem, err := u.inventoryClient.GetItemByID(ctx, itemReq.ItemID.String(), token)
		if err != nil {
			u.logger.Error(ctx, "failed to validate item", zap.String("item_id", itemReq.ItemID.String()), zap.Error(err), zap.String("error_type", fmt.Sprintf("%T", err)), zap.String("error_msg", err.Error()))
			if err == errors.ErrNotFound {
				return model.PurchaseOrderWithItems{}, errors.ErrBadRequest
			}
			return model.PurchaseOrderWithItems{}, errors.ErrInternalServerError
		}

		subtotal := inventoryItem.UnitPrice * float64(itemReq.Quantity)
		item := model.PurchaseOrderItem{
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

	if err := u.storage.CreateOrder(ctx, order); err != nil {
		u.logger.Error(ctx, "failed to store order", zap.Error(err))
		return model.PurchaseOrderWithItems{}, errors.ErrInternalServerError
	}

	for _, item := range items {
		if err := u.storage.CreateOrderItem(ctx, item); err != nil {
			u.logger.Error(ctx, "failed to store order item", zap.Error(err))
			return model.PurchaseOrderWithItems{}, errors.ErrInternalServerError
		}
	}

	return model.PurchaseOrderWithItems{
		PurchaseOrder: order,
		Items:         items,
	}, nil
}

func (u *Usecase) GetOrderByID(ctx context.Context, id string) (model.PurchaseOrderWithItems, error) {
	order, err := u.storage.GetOrderByID(ctx, id)
	if err != nil {
		return model.PurchaseOrderWithItems{}, err
	}

	items, err := u.storage.GetOrderItemsByOrderID(ctx, id)
	if err != nil {
		return model.PurchaseOrderWithItems{}, err
	}

	return model.PurchaseOrderWithItems{
		PurchaseOrder: order,
		Items:         items,
	}, nil
}

func (u *Usecase) ListOrders(ctx context.Context, limit, offset int) ([]model.PurchaseOrder, error) {
	return u.storage.ListOrders(ctx, limit, offset)
}

func (u *Usecase) UpdateOrder(ctx context.Context, id string, req model.UpdatePurchaseOrderRequest) (model.PurchaseOrderWithItems, error) {
	order, err := u.storage.GetOrderByID(ctx, id)
	if err != nil {
		return model.PurchaseOrderWithItems{}, err
	}

	if order.Status != model.PurchaseOrderStatusDraft {
		return model.PurchaseOrderWithItems{}, errors.ErrBadRequest
	}

	if err := u.storage.DeleteOrderItemsByOrderID(ctx, id); err != nil {
		return model.PurchaseOrderWithItems{}, err
	}

	token, err := u.getTokenFromContext(ctx)
	if err != nil {
		return model.PurchaseOrderWithItems{}, errors.ErrInternalServerError
	}

	var totalAmount float64
	items := make([]model.PurchaseOrderItem, 0, len(req.Items))

	for _, itemReq := range req.Items {
		inventoryItem, err := u.inventoryClient.GetItemByID(ctx, itemReq.ItemID.String(), token)
		if err != nil {
			u.logger.Error(ctx, "failed to validate item for update", zap.String("item_id", itemReq.ItemID.String()), zap.Error(err), zap.String("error_type", fmt.Sprintf("%T", err)), zap.String("error_msg", err.Error()))
			if err == errors.ErrNotFound {
				return model.PurchaseOrderWithItems{}, errors.ErrBadRequest
			}
			return model.PurchaseOrderWithItems{}, errors.ErrInternalServerError
		}

		subtotal := inventoryItem.UnitPrice * float64(itemReq.Quantity)
		item := model.PurchaseOrderItem{
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

	if err := u.storage.UpdateOrder(ctx, order); err != nil {
		u.logger.Error(ctx, "failed to update order in storage", zap.Error(err))
		return model.PurchaseOrderWithItems{}, errors.ErrInternalServerError
	}

	for _, item := range items {
		if err := u.storage.CreateOrderItem(ctx, item); err != nil {
			u.logger.Error(ctx, "failed to create order item for update", zap.Error(err))
			return model.PurchaseOrderWithItems{}, errors.ErrInternalServerError
		}
	}

	return model.PurchaseOrderWithItems{
		PurchaseOrder: order,
		Items:         items,
	}, nil
}

func (u *Usecase) ReceiveOrder(ctx context.Context, id string) (model.PurchaseOrderWithItems, error) {
	order, err := u.storage.GetOrderByID(ctx, id)
	if err != nil {
		return model.PurchaseOrderWithItems{}, err
	}

	if order.Status != model.PurchaseOrderStatusDraft {
		return model.PurchaseOrderWithItems{}, errors.ErrBadRequest
	}

	items, err := u.storage.GetOrderItemsByOrderID(ctx, id)
	if err != nil {
		return model.PurchaseOrderWithItems{}, err
	}

	if err := u.storage.UpdateOrderStatus(ctx, id, model.PurchaseOrderStatusReceived); err != nil {
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

	if err := u.natsClient.Publish("purchase.order.received", event); err != nil {
		u.logger.Error(ctx, "failed to publish purchase.order.received event", zap.Error(err))
	} else {
		u.logger.Info(ctx, "published purchase.order.received event",
			zap.String("order_id", order.ID.String()),
			zap.String("vendor_id", order.VendorID.String()),
		)
	}

	return model.PurchaseOrderWithItems{
		PurchaseOrder: order,
		Items:         items,
	}, nil
}

func (u *Usecase) PayOrder(ctx context.Context, id string) (model.PurchaseOrderWithItems, error) {
	order, err := u.storage.GetOrderByID(ctx, id)
	if err != nil {
		return model.PurchaseOrderWithItems{}, err
	}

	if order.Status != model.PurchaseOrderStatusReceived {
		return model.PurchaseOrderWithItems{}, errors.ErrBadRequest
	}

	if err := u.storage.UpdateOrderStatus(ctx, id, model.PurchaseOrderStatusPaid); err != nil {
		return model.PurchaseOrderWithItems{}, err
	}

	order.Status = model.PurchaseOrderStatusPaid
	order.UpdatedAt = time.Now()

	items, err := u.storage.GetOrderItemsByOrderID(ctx, id)
	if err != nil {
		return model.PurchaseOrderWithItems{}, err
	}

	return model.PurchaseOrderWithItems{
		PurchaseOrder: order,
		Items:         items,
	}, nil
}
