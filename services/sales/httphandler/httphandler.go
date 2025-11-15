package httphandler

import (
	"encoding/json"
	"microservice-challenge/package/errors"
	"microservice-challenge/package/log"
	"microservice-challenge/package/pagination"
	"microservice-challenge/package/response"
	"microservice-challenge/services/sales/model"
	"microservice-challenge/services/sales/usecase/sales"
	"net/http"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

type Handler struct {
	usecase *sales.Usecase
	logger  log.Logger
}

func NewHandler(usecase *sales.Usecase, logger log.Logger) *Handler {
	return &Handler{
		usecase: usecase,
		logger:  logger,
	}
}

// @Summary List orders
// @Description Get a list of sales orders. Supports both pagination styles: page/size or limit/offset
// @Tags orders
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param page query int false "Page number (1-based)" default(1)
// @Param size query int false "Page size" default(10)
// @Param limit query int false "Limit (alternative to size)" default(10)
// @Param offset query int false "Offset (alternative to page)" default(0)
// @Success 200 {array} model.SalesOrder
// @Failure 400 {string} string "Bad Request"
// @Failure 401 {string} string "Unauthorized"
// @Router /orders [get]
func (h *Handler) ListOrders(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	limit, offset := pagination.GetLimitOffset(r)

	orders, err := h.usecase.ListOrders(ctx, limit, offset)
	if err != nil {
		h.logger.Error(ctx, "failed to list orders", zap.Error(err))
		response.SendErrorResponse(w, err)
		return
	}

	response.SendSuccessResponse(w, http.StatusOK, "Orders retrieved successfully", orders, nil)
}

// @Summary Get order by ID
// @Description Get a sales order by ID with items
// @Tags orders
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Order ID"
// @Success 200 {object} model.SalesOrderWithItems
// @Failure 400 {string} string "Bad Request"
// @Failure 401 {string} string "Unauthorized"
// @Failure 404 {string} string "Not Found"
// @Router /orders/{id} [get]
func (h *Handler) GetOrder(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := chi.URLParam(r, "id")

	order, err := h.usecase.GetOrderByID(ctx, id)
	if err != nil {
		h.logger.Error(ctx, "failed to get order", zap.Error(err))
		response.SendErrorResponse(w, err)
		return
	}

	response.SendSuccessResponse(w, http.StatusOK, "Order retrieved successfully", order, nil)
}

// @Summary Create order
// @Description Create a new draft sales order
// @Tags orders
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body model.CreateOrderRequest true "Create order request"
// @Success 201 {object} model.SalesOrderWithItems
// @Failure 400 {string} string "Bad Request"
// @Failure 401 {string} string "Unauthorized"
// @Router /orders [post]
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

	order, err := h.usecase.CreateOrder(ctx, req)
	if err != nil {
		h.logger.Error(ctx, "failed to create order", zap.Error(err))
		response.SendErrorResponse(w, err)
		return
	}

	response.SendSuccessResponse(w, http.StatusCreated, "Order created successfully", order, nil)
}

// @Summary Update order
// @Description Update a draft sales order
// @Tags orders
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Order ID"
// @Param request body model.UpdateOrderRequest true "Update order request"
// @Success 200 {object} model.SalesOrderWithItems
// @Failure 400 {string} string "Bad Request"
// @Failure 401 {string} string "Unauthorized"
// @Failure 404 {string} string "Not Found"
// @Router /orders/{id} [put]
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

	order, err := h.usecase.UpdateOrder(ctx, id, req)
	if err != nil {
		h.logger.Error(ctx, "failed to update order", zap.Error(err))
		response.SendErrorResponse(w, err)
		return
	}

	response.SendSuccessResponse(w, http.StatusOK, "Order updated successfully", order, nil)
}

// @Summary Confirm order
// @Description Confirm a draft sales order and publish event
// @Tags orders
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Order ID"
// @Success 200 {object} model.SalesOrderWithItems
// @Failure 400 {string} string "Bad Request"
// @Failure 401 {string} string "Unauthorized"
// @Failure 403 {string} string "Forbidden"
// @Failure 404 {string} string "Not Found"
// @Router /orders/{id}/confirm [post]
func (h *Handler) ConfirmOrder(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := chi.URLParam(r, "id")

	order, err := h.usecase.ConfirmOrder(ctx, id)
	if err != nil {
		h.logger.Error(ctx, "failed to confirm order", zap.Error(err))
		response.SendErrorResponse(w, err)
		return
	}

	response.SendSuccessResponse(w, http.StatusOK, "Order confirmed successfully", order, nil)
}

// @Summary Pay order
// @Description Mark a confirmed order as paid
// @Tags orders
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Order ID"
// @Success 200 {object} model.SalesOrderWithItems
// @Failure 400 {string} string "Bad Request"
// @Failure 401 {string} string "Unauthorized"
// @Failure 403 {string} string "Forbidden"
// @Failure 404 {string} string "Not Found"
// @Router /orders/{id}/pay [post]
func (h *Handler) PayOrder(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := chi.URLParam(r, "id")

	order, err := h.usecase.PayOrder(ctx, id)
	if err != nil {
		h.logger.Error(ctx, "failed to pay order", zap.Error(err))
		response.SendErrorResponse(w, err)
		return
	}

	response.SendSuccessResponse(w, http.StatusOK, "Order paid successfully", order, nil)
}
