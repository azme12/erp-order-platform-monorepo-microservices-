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

const (
	maxRequestBodySize = 1 << 20
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

	response.SendSuccessResponse(w, http.StatusOK, "Purchase orders retrieved successfully", orders, nil)
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

	response.SendSuccessResponse(w, http.StatusOK, "Purchase order retrieved successfully", order, nil)
}

func (h *Handler) CreateOrder(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req model.CreatePurchaseOrderRequest
	if err := h.parseAndValidateRequest(w, r, &req); err != nil {
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

func (h *Handler) UpdateOrder(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := chi.URLParam(r, "id")

	var req model.UpdatePurchaseOrderRequest
	if err := h.parseAndValidateRequest(w, r, &req); err != nil {
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
