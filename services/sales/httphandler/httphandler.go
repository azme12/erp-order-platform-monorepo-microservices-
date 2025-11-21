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

const (
	maxRequestBodySize = 1 << 20
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

func (h *Handler) parseAndValidateRequest(w http.ResponseWriter, r *http.Request, req interface{ Validate() error }) error {
	r.Body = http.MaxBytesReader(w, r.Body, maxRequestBodySize)

	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		response.SendErrorResponse(w, errors.ErrBadRequest)
		return err
	}

	if err := req.Validate(); err != nil {
		response.SendErrorResponse(w, err)
		return err
	}

	return nil
}

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

func (h *Handler) CreateOrder(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req model.CreateOrderRequest
	if err := h.parseAndValidateRequest(w, r, &req); err != nil {
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

func (h *Handler) UpdateOrder(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := chi.URLParam(r, "id")

	var req model.UpdateOrderRequest
	if err := h.parseAndValidateRequest(w, r, &req); err != nil {
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
