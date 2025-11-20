package httphandler

import (
	"encoding/json"
	"microservice-challenge/package/errors"
	"microservice-challenge/package/log"
	"microservice-challenge/package/pagination"
	"microservice-challenge/package/response"
	"microservice-challenge/services/inventory/model"
	inventoryservice "microservice-challenge/services/inventory/service/inventory"
	"net/http"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

const (
	maxRequestBodySize = 1 << 20
)

type Handler struct {
	service *inventoryservice.Service
	logger  log.Logger
}

func NewHandler(service *inventoryservice.Service, logger log.Logger) *Handler {
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

func (h *Handler) ListItems(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	limit, offset := pagination.GetLimitOffset(r)

	items, err := h.service.ListItems(ctx, limit, offset)
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

	item, err := h.service.GetItemByID(ctx, id)
	if err != nil {
		h.logger.Error(ctx, "failed to get item", zap.Error(err))
		response.SendErrorResponse(w, err)
		return
	}

	response.SendSuccessResponse(w, http.StatusOK, "Item retrieved successfully", item, nil)
}

func (h *Handler) CreateItem(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req model.CreateItemRequest
	if err := h.parseAndValidateRequest(w, r, &req); err != nil {
		return
	}

	item, err := h.service.CreateItem(ctx, req)
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

	var req model.UpdateItemRequest
	if err := h.parseAndValidateRequest(w, r, &req); err != nil {
		return
	}

	item, err := h.service.UpdateItem(ctx, id, req)
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

	if err := h.service.DeleteItem(ctx, id); err != nil {
		h.logger.Error(ctx, "failed to delete item", zap.Error(err))
		response.SendErrorResponse(w, err)
		return
	}

	response.SendSuccessResponse(w, http.StatusOK, "Item deleted successfully", nil, nil)
}

func (h *Handler) GetStock(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	itemID := chi.URLParam(r, "item_id")

	stock, err := h.service.GetStockByItemID(ctx, itemID)
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

	var req model.AdjustStockRequest
	if err := h.parseAndValidateRequest(w, r, &req); err != nil {
		return
	}

	stock, err := h.service.AdjustStock(ctx, itemID, req.Quantity)
	if err != nil {
		h.logger.Error(ctx, "failed to adjust stock", zap.Error(err))
		response.SendErrorResponse(w, err)
		return
	}

	response.SendSuccessResponse(w, http.StatusOK, "Stock adjusted successfully", stock, nil)
}
