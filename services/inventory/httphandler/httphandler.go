package httphandler

import (
	"encoding/json"
	"microservice-challenge/package/errors"
	"microservice-challenge/package/log"
	"microservice-challenge/package/pagination"
	"microservice-challenge/package/response"
	"microservice-challenge/services/inventory/model"
	"microservice-challenge/services/inventory/usecase/inventory"
	"net/http"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

type Handler struct {
	usecase *inventory.Usecase
	logger  log.Logger
}

func NewHandler(usecase *inventory.Usecase, logger log.Logger) *Handler {
	return &Handler{
		usecase: usecase,
		logger:  logger,
	}
}

func (h *Handler) ListItems(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	limit, offset := pagination.GetLimitOffset(r)

	items, err := h.usecase.ListItems(ctx, limit, offset)
	if err != nil {
		h.logger.Error(ctx, "failed to list items", zap.Error(err))
		response.SendErrorResponse(w, err)
		return
	}

	response.SendSuccessResponse(w, http.StatusOK, "Items retrieved successfully", items, nil)
}

func (h *Handler) GetItem(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := chi.URLParam(r, "id")

	item, err := h.usecase.GetItemByID(ctx, id)
	if err != nil {
		h.logger.Error(ctx, "failed to get item", zap.Error(err))
		response.SendErrorResponse(w, err)
		return
	}

	response.SendSuccessResponse(w, http.StatusOK, "Item retrieved successfully", item, nil)
}

func (h *Handler) CreateItem(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	r.Body = http.MaxBytesReader(w, r.Body, 1<<20)

	var req model.CreateItemRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.SendErrorResponse(w, errors.ErrBadRequest)
		return
	}

	if err := req.Validate(); err != nil {
		response.SendErrorResponse(w, err)
		return
	}

	item, err := h.usecase.CreateItem(ctx, req)
	if err != nil {
		h.logger.Error(ctx, "failed to create item", zap.Error(err))
		response.SendErrorResponse(w, err)
		return
	}

	response.SendSuccessResponse(w, http.StatusCreated, "Item created successfully", item, nil)
}

func (h *Handler) UpdateItem(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := chi.URLParam(r, "id")
	r.Body = http.MaxBytesReader(w, r.Body, 1<<20)

	var req model.UpdateItemRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.SendErrorResponse(w, errors.ErrBadRequest)
		return
	}

	if err := req.Validate(); err != nil {
		response.SendErrorResponse(w, err)
		return
	}

	item, err := h.usecase.UpdateItem(ctx, id, req)
	if err != nil {
		h.logger.Error(ctx, "failed to update item", zap.Error(err))
		response.SendErrorResponse(w, err)
		return
	}

	response.SendSuccessResponse(w, http.StatusOK, "Item updated successfully", item, nil)
}

func (h *Handler) DeleteItem(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := chi.URLParam(r, "id")

	if err := h.usecase.DeleteItem(ctx, id); err != nil {
		h.logger.Error(ctx, "failed to delete item", zap.Error(err))
		response.SendErrorResponse(w, err)
		return
	}

	response.SendSuccessResponse(w, http.StatusOK, "Item deleted successfully", nil, nil)
}

func (h *Handler) GetStock(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	itemID := chi.URLParam(r, "item_id")

	stock, err := h.usecase.GetStockByItemID(ctx, itemID)
	if err != nil {
		h.logger.Error(ctx, "failed to get stock", zap.Error(err))
		response.SendErrorResponse(w, err)
		return
	}

	response.SendSuccessResponse(w, http.StatusOK, "Stock retrieved successfully", stock, nil)
}

func (h *Handler) AdjustStock(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	itemID := chi.URLParam(r, "item_id")
	r.Body = http.MaxBytesReader(w, r.Body, 1<<20)

	var req model.AdjustStockRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.SendErrorResponse(w, errors.ErrBadRequest)
		return
	}

	if err := req.Validate(); err != nil {
		response.SendErrorResponse(w, err)
		return
	}

	stock, err := h.usecase.AdjustStock(ctx, itemID, req.Quantity)
	if err != nil {
		h.logger.Error(ctx, "failed to adjust stock", zap.Error(err))
		response.SendErrorResponse(w, err)
		return
	}

	response.SendSuccessResponse(w, http.StatusOK, "Stock adjusted successfully", stock, nil)
}
