package httphandler

import (
	"encoding/json"
	"microservice-challenge/package/errors"
	"microservice-challenge/package/log"
	"microservice-challenge/package/pagination"
	"microservice-challenge/package/response"
	"microservice-challenge/services/purchase/model"
	"microservice-challenge/services/purchase/usecase/purchase"
	"net/http"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

type Handler struct {
	usecase *purchase.Usecase
	logger  log.Logger
}

func NewHandler(usecase *purchase.Usecase, logger log.Logger) *Handler {
	return &Handler{
		usecase: usecase,
		logger:  logger,
	}
}

func (h *Handler) ListOrders(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	limit, offset := pagination.GetLimitOffset(r)

	orders, err := h.usecase.ListOrders(ctx, limit, offset)
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

	order, err := h.usecase.GetOrderByID(ctx, id)
	if err != nil {
		h.logger.Error(ctx, "failed to get order", zap.Error(err))
		response.SendErrorResponse(w, err)
		return
	}

	response.SendSuccessResponse(w, http.StatusOK, "Purchase order retrieved successfully", order, nil)
}

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

	order, err := h.usecase.CreateOrder(ctx, req)
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

	order, err := h.usecase.UpdateOrder(ctx, id, req)
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

	order, err := h.usecase.ReceiveOrder(ctx, id)
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

	order, err := h.usecase.PayOrder(ctx, id)
	if err != nil {
		h.logger.Error(ctx, "failed to pay order", zap.Error(err))
		response.SendErrorResponse(w, err)
		return
	}

	response.SendSuccessResponse(w, http.StatusOK, "Purchase order paid successfully", order, nil)
}
