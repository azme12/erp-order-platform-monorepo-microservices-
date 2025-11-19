package httphandler

import (
	"encoding/json"
	"microservice-challenge/package/errors"
	"microservice-challenge/package/log"
	"microservice-challenge/package/pagination"
	"microservice-challenge/package/response"
	"microservice-challenge/services/sales/model"
	salesservice "microservice-challenge/services/sales/service/sales"
	"net/http"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

type Handler struct {
	service *salesservice.Service
	logger  log.Logger
}

func NewHandler(service *salesservice.Service, logger log.Logger) *Handler {
	return &Handler{
		service: service,
		logger:  logger,
	}
}

// ListOrders godoc
// @Summary      List sales orders
// @Description  Get a paginated list of sales orders
// @Tags         sales
// @Accept       json
// @Produce      json
// @Param        limit query int false "Limit" default(10)
// @Param        offset query int false "Offset" default(0)
// @Success      200 {object} response.SuccessResponse{data=[]model.SalesOrder}
// @Failure      401 {object} response.SimpleErrorResponse
// @Failure      403 {object} response.SimpleErrorResponse
// @Failure      500 {object} response.SimpleErrorResponse
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

	response.SendSuccessResponse(w, http.StatusOK, "Orders retrieved successfully", orders, nil)
}

// GetOrder godoc
// @Summary      Get sales order by ID
// @Description  Get a single sales order with its items by order ID
// @Tags         sales
// @Accept       json
// @Produce      json
// @Param        id path string true "Order ID"
// @Success      200 {object} response.SuccessResponse{data=model.SalesOrderWithItems}
// @Failure      401 {object} response.SimpleErrorResponse
// @Failure      403 {object} response.SimpleErrorResponse
// @Failure      404 {object} response.SimpleErrorResponse
// @Failure      500 {object} response.SimpleErrorResponse
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

	response.SendSuccessResponse(w, http.StatusOK, "Order retrieved successfully", order, nil)
}

// CreateOrder godoc
// @Summary      Create a new sales order
// @Description  Create a new sales order linked to a customer with order items
// @Tags         sales
// @Accept       json
// @Produce      json
// @Param        request body model.CreateOrderRequest true "Order creation request"
// @Success      201 {object} response.SuccessResponse{data=model.SalesOrderWithItems}
// @Failure      400 {object} response.ValidationErrorResponse
// @Failure      401 {object} response.SimpleErrorResponse
// @Failure      403 {object} response.SimpleErrorResponse
// @Failure      500 {object} response.SimpleErrorResponse
// @Router       /orders [post]
// @Security     BearerAuth
func (h *Handler) CreateOrder(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	r.Body = http.MaxBytesReader(w, r.Body, 1<<20)

	var req model.CreateOrderRequest
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

	response.SendSuccessResponse(w, http.StatusCreated, "Order created successfully", order, nil)
}

// UpdateOrder godoc
// @Summary      Update a sales order
// @Description  Update a draft sales order's items and details
// @Tags         sales
// @Accept       json
// @Produce      json
// @Param        id path string true "Order ID"
// @Param        request body model.UpdateOrderRequest true "Order update request"
// @Success      200 {object} response.SuccessResponse{data=model.SalesOrderWithItems}
// @Failure      400 {object} response.ValidationErrorResponse
// @Failure      401 {object} response.SimpleErrorResponse
// @Failure      403 {object} response.SimpleErrorResponse
// @Failure      404 {object} response.SimpleErrorResponse
// @Failure      500 {object} response.SimpleErrorResponse
// @Router       /orders/{id} [put]
// @Security     BearerAuth
func (h *Handler) UpdateOrder(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := chi.URLParam(r, "id")
	r.Body = http.MaxBytesReader(w, r.Body, 1<<20)

	var req model.UpdateOrderRequest
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

	response.SendSuccessResponse(w, http.StatusOK, "Order updated successfully", order, nil)
}

// ConfirmOrder godoc
// @Summary      Confirm a sales order
// @Description  Confirm a draft sales order and publish sales.order.confirmed event to decrease inventory stock
// @Tags         sales
// @Accept       json
// @Produce      json
// @Param        id path string true "Order ID"
// @Success      200 {object} response.SuccessResponse{data=model.SalesOrderWithItems}
// @Failure      400 {object} response.ValidationErrorResponse
// @Failure      401 {object} response.SimpleErrorResponse
// @Failure      403 {object} response.SimpleErrorResponse
// @Failure      404 {object} response.SimpleErrorResponse
// @Failure      500 {object} response.SimpleErrorResponse
// @Router       /orders/{id}/confirm [post]
// @Security     BearerAuth
func (h *Handler) ConfirmOrder(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := chi.URLParam(r, "id")

	order, err := h.service.ConfirmOrder(ctx, id)
	if err != nil {
		h.logger.Error(ctx, "failed to confirm order", zap.Error(err))
		response.SendErrorResponse(w, err)
		return
	}

	response.SendSuccessResponse(w, http.StatusOK, "Order confirmed successfully", order, nil)
}

// PayOrder godoc
// @Summary      Pay a sales order
// @Description  Mark a confirmed sales order as paid
// @Tags         sales
// @Accept       json
// @Produce      json
// @Param        id path string true "Order ID"
// @Success      200 {object} response.SuccessResponse{data=model.SalesOrderWithItems}
// @Failure      400 {object} response.ValidationErrorResponse
// @Failure      401 {object} response.SimpleErrorResponse
// @Failure      403 {object} response.SimpleErrorResponse
// @Failure      404 {object} response.SimpleErrorResponse
// @Failure      500 {object} response.SimpleErrorResponse
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

	response.SendSuccessResponse(w, http.StatusOK, "Order paid successfully", order, nil)
}
