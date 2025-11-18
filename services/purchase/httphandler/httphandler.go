package httphandler

import (
	"encoding/json"
	"microservice-challenge/package/errors"
	"microservice-challenge/package/log"
	"microservice-challenge/package/pagination"
	"microservice-challenge/package/response"
	"microservice-challenge/services/purchase/model"
	purchaseservice "microservice-challenge/services/purchase/service/purchase"
	"net/http"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

type Handler struct {
	service *purchaseservice.Service
	logger  log.Logger
}

func NewHandler(service *purchaseservice.Service, logger log.Logger) *Handler {
	return &Handler{
		service: service,
		logger:  logger,
	}
}

// ListOrders godoc
// @Summary      List purchase orders
// @Description  Get a paginated list of purchase orders
// @Tags         purchase
// @Accept       json
// @Produce      json
// @Param        limit query int false "Limit" default(10)
// @Param        offset query int false "Offset" default(0)
// @Success      200 {object} response.Response{data=[]model.PurchaseOrder}
// @Failure      401 {object} response.Response
// @Failure      403 {object} response.Response
// @Failure      500 {object} response.Response
// @Router       /orders [get]
// @Security     BearerAuth
func (h *Handler) ListOrders(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	limit, offset := pagination.GetLimitOffset(r)

	orders, err := h.service.ListOrders(ctx, limit, offset)
	if err != nil {
		h.logger.Error(ctx, "failed to list orders", zap.Error(err))
		response.SendErrorResponse(w, err)
		return
	}

	response.SendSuccessResponse(w, http.StatusOK, "Purchase orders retrieved successfully", orders, nil)
}

// GetOrder godoc
// @Summary      Get purchase order by ID
// @Description  Get a single purchase order with its items by order ID
// @Tags         purchase
// @Accept       json
// @Produce      json
// @Param        id path string true "Order ID"
// @Success      200 {object} response.Response{data=model.PurchaseOrderWithItems}
// @Failure      401 {object} response.Response
// @Failure      403 {object} response.Response
// @Failure      404 {object} response.Response
// @Failure      500 {object} response.Response
// @Router       /orders/{id} [get]
// @Security     BearerAuth
func (h *Handler) GetOrder(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := chi.URLParam(r, "id")

	order, err := h.service.GetOrderByID(ctx, id)
	if err != nil {
		h.logger.Error(ctx, "failed to get order", zap.Error(err))
		response.SendErrorResponse(w, err)
		return
	}

	response.SendSuccessResponse(w, http.StatusOK, "Purchase order retrieved successfully", order, nil)
}

// CreateOrder godoc
// @Summary      Create a new purchase order
// @Description  Create a new purchase order linked to a vendor with order items
// @Tags         purchase
// @Accept       json
// @Produce      json
// @Param        request body model.CreatePurchaseOrderRequest true "Order creation request"
// @Success      201 {object} response.Response{data=model.PurchaseOrderWithItems}
// @Failure      400 {object} response.Response
// @Failure      401 {object} response.Response
// @Failure      403 {object} response.Response
// @Failure      500 {object} response.Response
// @Router       /orders [post]
// @Security     BearerAuth
func (h *Handler) CreateOrder(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	r.Body = http.MaxBytesReader(w, r.Body, 1<<20)

	var req model.CreatePurchaseOrderRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.SendErrorResponse(w, errors.ErrBadRequest)
		return
	}

	if err := req.Validate(); err != nil {
		response.SendErrorResponse(w, err)
		return
	}

	order, err := h.service.CreateOrder(ctx, req)
	if err != nil {
		h.logger.Error(ctx, "failed to create order", zap.Error(err))
		response.SendErrorResponse(w, err)
		return
	}

	response.SendSuccessResponse(w, http.StatusCreated, "Purchase order created successfully", order, nil)
}

// UpdateOrder godoc
// @Summary      Update a purchase order
// @Description  Update a draft purchase order's items and details
// @Tags         purchase
// @Accept       json
// @Produce      json
// @Param        id path string true "Order ID"
// @Param        request body model.UpdatePurchaseOrderRequest true "Order update request"
// @Success      200 {object} response.Response{data=model.PurchaseOrderWithItems}
// @Failure      400 {object} response.Response
// @Failure      401 {object} response.Response
// @Failure      403 {object} response.Response
// @Failure      404 {object} response.Response
// @Failure      500 {object} response.Response
// @Router       /orders/{id} [put]
// @Security     BearerAuth
func (h *Handler) UpdateOrder(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := chi.URLParam(r, "id")
	r.Body = http.MaxBytesReader(w, r.Body, 1<<20)

	var req model.UpdatePurchaseOrderRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.SendErrorResponse(w, errors.ErrBadRequest)
		return
	}

	if err := req.Validate(); err != nil {
		response.SendErrorResponse(w, err)
		return
	}

	order, err := h.service.UpdateOrder(ctx, id, req)
	if err != nil {
		h.logger.Error(ctx, "failed to update order", zap.Error(err))
		response.SendErrorResponse(w, err)
		return
	}

	response.SendSuccessResponse(w, http.StatusOK, "Purchase order updated successfully", order, nil)
}

// ReceiveOrder godoc
// @Summary      Receive a purchase order
// @Description  Receive a draft purchase order and publish purchase.order.received event to increase inventory stock
// @Tags         purchase
// @Accept       json
// @Produce      json
// @Param        id path string true "Order ID"
// @Success      200 {object} response.Response{data=model.PurchaseOrderWithItems}
// @Failure      400 {object} response.Response
// @Failure      401 {object} response.Response
// @Failure      403 {object} response.Response
// @Failure      404 {object} response.Response
// @Failure      500 {object} response.Response
// @Router       /orders/{id}/receive [post]
// @Security     BearerAuth
func (h *Handler) ReceiveOrder(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := chi.URLParam(r, "id")

	order, err := h.service.ReceiveOrder(ctx, id)
	if err != nil {
		h.logger.Error(ctx, "failed to receive order", zap.Error(err))
		response.SendErrorResponse(w, err)
		return
	}

	response.SendSuccessResponse(w, http.StatusOK, "Purchase order received successfully", order, nil)
}

// PayOrder godoc
// @Summary      Pay a purchase order
// @Description  Mark a received purchase order as paid
// @Tags         purchase
// @Accept       json
// @Produce      json
// @Param        id path string true "Order ID"
// @Success      200 {object} response.Response{data=model.PurchaseOrderWithItems}
// @Failure      400 {object} response.Response
// @Failure      401 {object} response.Response
// @Failure      403 {object} response.Response
// @Failure      404 {object} response.Response
// @Failure      500 {object} response.Response
// @Router       /orders/{id}/pay [post]
// @Security     BearerAuth
func (h *Handler) PayOrder(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := chi.URLParam(r, "id")

	order, err := h.service.PayOrder(ctx, id)
	if err != nil {
		h.logger.Error(ctx, "failed to pay order", zap.Error(err))
		response.SendErrorResponse(w, err)
		return
	}

	response.SendSuccessResponse(w, http.StatusOK, "Purchase order paid successfully", order, nil)
}
