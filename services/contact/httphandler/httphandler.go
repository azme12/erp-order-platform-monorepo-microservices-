package httphandler

import (
	"encoding/json"
	"microservice-challenge/package/errors"
	"microservice-challenge/package/log"
	"microservice-challenge/package/pagination"
	"microservice-challenge/package/response"
	"microservice-challenge/services/contact/model"
	"microservice-challenge/services/contact/usecase/contact"
	"net/http"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

type Handler struct {
	usecase *contact.Usecase
	logger  log.Logger
}

func NewHandler(usecase *contact.Usecase, logger log.Logger) *Handler {
	return &Handler{
		usecase: usecase,
		logger:  logger,
	}
}

// @Summary List customers
// @Description Get a list of customers. Supports both pagination styles: page/size or limit/offset
// @Tags customers
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param page query int false "Page number (1-based)" default(1)
// @Param size query int false "Page size" default(10)
// @Param limit query int false "Limit (alternative to size)" default(10)
// @Param offset query int false "Offset (alternative to page)" default(0)
// @Success 200 {array} model.Customer
// @Failure 400 {string} string "Bad Request"
// @Failure 401 {string} string "Unauthorized"
// @Router /customers [get]
func (h *Handler) ListCustomers(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	limit, offset := pagination.GetLimitOffset(r)

	customers, err := h.usecase.ListCustomers(ctx, limit, offset)
	if err != nil {
		h.logger.Error(ctx, "failed to list customers", zap.Error(err))
		response.SendErrorResponse(w, err)
		return
	}

	response.SendSuccessResponse(w, http.StatusOK, "Customers retrieved successfully", customers, nil)
}

// @Summary Get customer by ID
// @Description Get a customer by ID
// @Tags customers
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Customer ID"
// @Success 200 {object} model.Customer
// @Failure 400 {string} string "Bad Request"
// @Failure 401 {string} string "Unauthorized"
// @Failure 404 {string} string "Not Found"
// @Router /customers/{id} [get]
func (h *Handler) GetCustomer(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := chi.URLParam(r, "id")

	customer, err := h.usecase.GetCustomerByID(ctx, id)
	if err != nil {
		h.logger.Error(ctx, "failed to get customer", zap.Error(err))
		response.SendErrorResponse(w, err)
		return
	}

	response.SendSuccessResponse(w, http.StatusOK, "Customer retrieved successfully", customer, nil)
}

// @Summary Create customer
// @Description Create a new customer
// @Tags customers
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body model.CreateCustomerRequest true "Create customer request"
// @Success 201 {object} model.Customer
// @Failure 400 {string} string "Bad Request"
// @Failure 401 {string} string "Unauthorized"
// @Failure 409 {string} string "Conflict"
// @Router /customers [post]
func (h *Handler) CreateCustomer(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	r.Body = http.MaxBytesReader(w, r.Body, 1<<20)

	var req model.CreateCustomerRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.SendErrorResponse(w, errors.ErrBadRequest)
		return
	}

	if err := req.Validate(); err != nil {
		response.SendErrorResponse(w, err)
		return
	}

	customer, err := h.usecase.CreateCustomer(ctx, req)
	if err != nil {
		h.logger.Error(ctx, "failed to create customer", zap.Error(err))
		response.SendErrorResponse(w, err)
		return
	}

	response.SendSuccessResponse(w, http.StatusCreated, "Customer created successfully", customer, nil)
}

// @Summary Update customer
// @Description Update a customer
// @Tags customers
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Customer ID"
// @Param request body model.UpdateCustomerRequest true "Update customer request"
// @Success 200 {object} model.Customer
// @Failure 400 {string} string "Bad Request"
// @Failure 401 {string} string "Unauthorized"
// @Failure 404 {string} string "Not Found"
// @Failure 409 {string} string "Conflict"
// @Router /customers/{id} [put]
func (h *Handler) UpdateCustomer(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := chi.URLParam(r, "id")
	r.Body = http.MaxBytesReader(w, r.Body, 1<<20)

	var req model.UpdateCustomerRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.SendErrorResponse(w, errors.ErrBadRequest)
		return
	}

	if err := req.Validate(); err != nil {
		response.SendErrorResponse(w, err)
		return
	}

	customer, err := h.usecase.UpdateCustomer(ctx, id, req)
	if err != nil {
		h.logger.Error(ctx, "failed to update customer", zap.Error(err))
		response.SendErrorResponse(w, err)
		return
	}

	response.SendSuccessResponse(w, http.StatusOK, "Customer updated successfully", customer, nil)
}

// @Summary Delete customer
// @Description Delete a customer
// @Tags customers
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Customer ID"
// @Success 200 {string} string "Customer deleted successfully"
// @Failure 400 {string} string "Bad Request"
// @Failure 401 {string} string "Unauthorized"
// @Failure 403 {string} string "Forbidden"
// @Failure 404 {string} string "Not Found"
// @Router /customers/{id} [delete]
func (h *Handler) DeleteCustomer(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := chi.URLParam(r, "id")

	if err := h.usecase.DeleteCustomer(ctx, id); err != nil {
		h.logger.Error(ctx, "failed to delete customer", zap.Error(err))
		response.SendErrorResponse(w, err)
		return
	}

	response.SendSuccessResponse(w, http.StatusOK, "Customer deleted successfully", nil, nil)
}

// @Summary List vendors
// @Description Get a list of vendors. Supports both pagination styles: page/size or limit/offset
// @Tags vendors
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param page query int false "Page number (1-based)" default(1)
// @Param size query int false "Page size" default(10)
// @Param limit query int false "Limit (alternative to size)" default(10)
// @Param offset query int false "Offset (alternative to page)" default(0)
// @Success 200 {array} model.Vendor
// @Failure 400 {string} string "Bad Request"
// @Failure 401 {string} string "Unauthorized"
// @Router /vendors [get]
func (h *Handler) ListVendors(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	limit, offset := pagination.GetLimitOffset(r)

	vendors, err := h.usecase.ListVendors(ctx, limit, offset)
	if err != nil {
		h.logger.Error(ctx, "failed to list vendors", zap.Error(err))
		response.SendErrorResponse(w, err)
		return
	}

	response.SendSuccessResponse(w, http.StatusOK, "Vendors retrieved successfully", vendors, nil)
}

// @Summary Get vendor by ID
// @Description Get a vendor by ID
// @Tags vendors
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Vendor ID"
// @Success 200 {object} model.Vendor
// @Failure 400 {string} string "Bad Request"
// @Failure 401 {string} string "Unauthorized"
// @Failure 404 {string} string "Not Found"
// @Router /vendors/{id} [get]
func (h *Handler) GetVendor(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := chi.URLParam(r, "id")

	vendor, err := h.usecase.GetVendorByID(ctx, id)
	if err != nil {
		h.logger.Error(ctx, "failed to get vendor", zap.Error(err))
		response.SendErrorResponse(w, err)
		return
	}

	response.SendSuccessResponse(w, http.StatusOK, "Vendor retrieved successfully", vendor, nil)
}

// @Summary Create vendor
// @Description Create a new vendor
// @Tags vendors
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body model.CreateVendorRequest true "Create vendor request"
// @Success 201 {object} model.Vendor
// @Failure 400 {string} string "Bad Request"
// @Failure 401 {string} string "Unauthorized"
// @Failure 409 {string} string "Conflict"
// @Router /vendors [post]
func (h *Handler) CreateVendor(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	r.Body = http.MaxBytesReader(w, r.Body, 1<<20)

	var req model.CreateVendorRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.SendErrorResponse(w, errors.ErrBadRequest)
		return
	}

	if err := req.Validate(); err != nil {
		response.SendErrorResponse(w, err)
		return
	}

	vendor, err := h.usecase.CreateVendor(ctx, req)
	if err != nil {
		h.logger.Error(ctx, "failed to create vendor", zap.Error(err))
		response.SendErrorResponse(w, err)
		return
	}

	response.SendSuccessResponse(w, http.StatusCreated, "Vendor created successfully", vendor, nil)
}

// @Summary Update vendor
// @Description Update a vendor
// @Tags vendors
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Vendor ID"
// @Param request body model.UpdateVendorRequest true "Update vendor request"
// @Success 200 {object} model.Vendor
// @Failure 400 {string} string "Bad Request"
// @Failure 401 {string} string "Unauthorized"
// @Failure 404 {string} string "Not Found"
// @Failure 409 {string} string "Conflict"
// @Router /vendors/{id} [put]
func (h *Handler) UpdateVendor(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := chi.URLParam(r, "id")
	r.Body = http.MaxBytesReader(w, r.Body, 1<<20)

	var req model.UpdateVendorRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.SendErrorResponse(w, errors.ErrBadRequest)
		return
	}

	if err := req.Validate(); err != nil {
		response.SendErrorResponse(w, err)
		return
	}

	vendor, err := h.usecase.UpdateVendor(ctx, id, req)
	if err != nil {
		h.logger.Error(ctx, "failed to update vendor", zap.Error(err))
		response.SendErrorResponse(w, err)
		return
	}

	response.SendSuccessResponse(w, http.StatusOK, "Vendor updated successfully", vendor, nil)
}

// @Summary Delete vendor
// @Description Delete a vendor
// @Tags vendors
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Vendor ID"
// @Success 200 {string} string "Vendor deleted successfully"
// @Failure 400 {string} string "Bad Request"
// @Failure 401 {string} string "Unauthorized"
// @Failure 403 {string} string "Forbidden"
// @Failure 404 {string} string "Not Found"
// @Router /vendors/{id} [delete]
func (h *Handler) DeleteVendor(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := chi.URLParam(r, "id")

	if err := h.usecase.DeleteVendor(ctx, id); err != nil {
		h.logger.Error(ctx, "failed to delete vendor", zap.Error(err))
		response.SendErrorResponse(w, err)
		return
	}

	response.SendSuccessResponse(w, http.StatusOK, "Vendor deleted successfully", nil, nil)
}
