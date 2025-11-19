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

// ListItems godoc
// @Summary      List items
// @Description  Get a paginated list of items
// @Tags         inventory
// @Accept       json
// @Produce      json
// @Param        limit query int false "Limit" default(10)
// @Param        offset query int false "Offset" default(0)
// @Success      200 {object} response.SuccessResponse{data=[]model.Item}
// @Failure      401 {object} response.SimpleErrorResponse
// @Failure      403 {object} response.SimpleErrorResponse
// @Failure      500 {object} response.SimpleErrorResponse
// @Router       /items [get]
// @Security     BearerAuth
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

// GetItem godoc
// @Summary      Get item by ID
// @Description  Get a single item by its ID
// @Tags         inventory
// @Accept       json
// @Produce      json
// @Param        id path string true "Item ID"
// @Success      200 {object} response.SuccessResponse{data=model.Item}
// @Failure      401 {object} response.SimpleErrorResponse
// @Failure      404 {object} response.SimpleErrorResponse
// @Failure      500 {object} response.SimpleErrorResponse
// @Router       /items/{id} [get]
// @Security     BearerAuth
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

// CreateItem godoc
// @Summary      Create a new item
// @Description  Create a new item with name, description, SKU, and unit price
// @Tags         inventory
// @Accept       json
// @Produce      json
// @Param        request body model.CreateItemRequest true "Item creation request"
// @Success      201 {object} response.SuccessResponse{data=model.Item}
// @Failure      400 {object} response.ValidationErrorResponse
// @Failure      401 {object} response.SimpleErrorResponse
// @Failure      403 {object} response.SimpleErrorResponse
// @Failure      409 {object} response.SimpleErrorResponse
// @Failure      500 {object} response.SimpleErrorResponse
// @Router       /items [post]
// @Security     BearerAuth
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

	item, err := h.service.CreateItem(ctx, req)
	if err != nil {
		h.logger.Error(ctx, "failed to create item", zap.Error(err))
		response.SendErrorResponse(w, err)
		return
	}

	response.SendSuccessResponse(w, http.StatusCreated, "Item created successfully", item, nil)
}

// UpdateItem godoc
// @Summary      Update an item
// @Description  Update an existing item's details
// @Tags         inventory
// @Accept       json
// @Produce      json
// @Param        id path string true "Item ID"
// @Param        request body model.UpdateItemRequest true "Item update request"
// @Success      200 {object} response.SuccessResponse{data=model.Item}
// @Failure      400 {object} response.ValidationErrorResponse
// @Failure      401 {object} response.SimpleErrorResponse
// @Failure      403 {object} response.SimpleErrorResponse
// @Failure      404 {object} response.SimpleErrorResponse
// @Failure      500 {object} response.SimpleErrorResponse
// @Router       /items/{id} [put]
// @Security     BearerAuth
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

	item, err := h.service.UpdateItem(ctx, id, req)
	if err != nil {
		h.logger.Error(ctx, "failed to update item", zap.Error(err))
		response.SendErrorResponse(w, err)
		return
	}

	response.SendSuccessResponse(w, http.StatusOK, "Item updated successfully", item, nil)
}

// DeleteItem godoc
// @Summary      Delete an item
// @Description  Delete an item by its ID
// @Tags         inventory
// @Accept       json
// @Produce      json
// @Param        id path string true "Item ID"
// @Success      200 {object} response.SuccessResponse
// @Failure      401 {object} response.SimpleErrorResponse
// @Failure      403 {object} response.SimpleErrorResponse
// @Failure      404 {object} response.SimpleErrorResponse
// @Failure      500 {object} response.SimpleErrorResponse
// @Router       /items/{id} [delete]
// @Security     BearerAuth
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

// GetStock godoc
// @Summary      Get stock by item ID
// @Description  Get the current stock quantity for a specific item
// @Tags         inventory
// @Accept       json
// @Produce      json
// @Param        item_id path string true "Item ID"
// @Success      200 {object} response.SuccessResponse{data=model.Stock}
// @Failure      401 {object} response.SimpleErrorResponse
// @Failure      404 {object} response.SimpleErrorResponse
// @Failure      500 {object} response.SimpleErrorResponse
// @Router       /items/{item_id}/stock [get]
// @Security     BearerAuth
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

// AdjustStock godoc
// @Summary      Adjust stock quantity
// @Description  Adjust the stock quantity for a specific item (positive to increase, negative to decrease)
// @Tags         inventory
// @Accept       json
// @Produce      json
// @Param        item_id path string true "Item ID"
// @Param        request body model.AdjustStockRequest true "Stock adjustment request"
// @Success      200 {object} response.SuccessResponse{data=model.Stock}
// @Failure      400 {object} response.ValidationErrorResponse
// @Failure      401 {object} response.SimpleErrorResponse
// @Failure      403 {object} response.SimpleErrorResponse
// @Failure      404 {object} response.SimpleErrorResponse
// @Failure      500 {object} response.SimpleErrorResponse
// @Router       /items/{item_id}/stock [put]
// @Security     BearerAuth
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

	stock, err := h.service.AdjustStock(ctx, itemID, req.Quantity)
	if err != nil {
		h.logger.Error(ctx, "failed to adjust stock", zap.Error(err))
		response.SendErrorResponse(w, err)
		return
	}

	response.SendSuccessResponse(w, http.StatusOK, "Stock adjusted successfully", stock, nil)
}
