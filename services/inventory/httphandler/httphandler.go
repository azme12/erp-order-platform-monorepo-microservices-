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

// @Summary List items
// @Description Get a list of items. Supports both pagination styles: page/size or limit/offset
// @Tags items
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param page query int false "Page number (1-based)" default(1)
// @Param size query int false "Page size" default(10)
// @Param limit query int false "Limit (alternative to size)" default(10)
// @Param offset query int false "Offset (alternative to page)" default(0)
// @Success 200 {array} model.Item
// @Failure 400 {string} string "Bad Request"
// @Failure 401 {string} string "Unauthorized"
// @Router /items [get]
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

// @Summary Get item by ID
// @Description Get an item by ID
// @Tags items
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Item ID"
// @Success 200 {object} model.Item
// @Failure 400 {string} string "Bad Request"
// @Failure 401 {string} string "Unauthorized"
// @Failure 404 {string} string "Not Found"
// @Router /items/{id} [get]
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

// @Summary Create item
// @Description Create a new item
// @Tags items
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body model.CreateItemRequest true "Create item request"
// @Success 201 {object} model.Item
// @Failure 400 {string} string "Bad Request"
// @Failure 401 {string} string "Unauthorized"
// @Failure 409 {string} string "Conflict"
// @Router /items [post]
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

// @Summary Update item
// @Description Update an item
// @Tags items
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Item ID"
// @Param request body model.UpdateItemRequest true "Update item request"
// @Success 200 {object} model.Item
// @Failure 400 {string} string "Bad Request"
// @Failure 401 {string} string "Unauthorized"
// @Failure 404 {string} string "Not Found"
// @Failure 409 {string} string "Conflict"
// @Router /items/{id} [put]
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

// @Summary Delete item
// @Description Delete an item
// @Tags items
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Item ID"
// @Success 200 {string} string "Item deleted successfully"
// @Failure 400 {string} string "Bad Request"
// @Failure 401 {string} string "Unauthorized"
// @Failure 403 {string} string "Forbidden"
// @Failure 404 {string} string "Not Found"
// @Router /items/{id} [delete]
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

// @Summary Get stock by item ID
// @Description Get stock information for an item
// @Tags stock
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param item_id path string true "Item ID"
// @Success 200 {object} model.Stock
// @Failure 400 {string} string "Bad Request"
// @Failure 401 {string} string "Unauthorized"
// @Failure 404 {string} string "Not Found"
// @Router /items/{item_id}/stock [get]
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

// @Summary Adjust stock
// @Description Adjust stock quantity for an item (positive to increase, negative to decrease)
// @Tags stock
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param item_id path string true "Item ID"
// @Param request body model.AdjustStockRequest true "Adjust stock request"
// @Success 200 {object} model.Stock
// @Failure 400 {string} string "Bad Request"
// @Failure 401 {string} string "Unauthorized"
// @Failure 404 {string} string "Not Found"
// @Router /items/{item_id}/stock [put]
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
