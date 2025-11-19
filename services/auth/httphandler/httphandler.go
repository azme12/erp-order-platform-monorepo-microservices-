package httphandler

import (
	"encoding/json"
	"microservice-challenge/package/errors"
	"microservice-challenge/package/log"
	"microservice-challenge/package/response"
	"microservice-challenge/services/auth/model"
	authservice "microservice-challenge/services/auth/service/auth"
	"net/http"

	"go.uber.org/zap"
)

type Handler struct {
	service *authservice.Service
	logger  log.Logger
}

func NewHandler(service *authservice.Service, logger log.Logger) *Handler {
	return &Handler{
		service: service,
		logger:  logger,
	}
}

// Register godoc
// @Summary      Register a new user
// @Description  Register a new user with email, password, and role
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        request body model.RegisterRequest true "Registration request"
// @Success      201 {object} response.SuccessResponse{data=model.User}
// @Failure      400 {object} response.ValidationErrorResponse
// @Failure      409 {object} response.SimpleErrorResponse
// @Failure      500 {object} response.SimpleErrorResponse
// @Router       /register [post]
func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	r.Body = http.MaxBytesReader(w, r.Body, 1<<20)

	var req model.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.SendErrorResponse(w, errors.ErrBadRequest)
		return
	}

	if err := req.Validate(); err != nil {
		response.SendErrorResponse(w, err)
		return
	}

	user, err := h.service.Register(ctx, req)
	if err != nil {
		h.logger.Error(ctx, "failed to register user", zap.Error(err))
		response.SendErrorResponse(w, err)
		return
	}

	response.SendSuccessResponse(w, http.StatusCreated, "User registered successfully", user, nil)
}

// Login godoc
// @Summary      Login user
// @Description  Authenticate user and return access token
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        request body model.LoginRequest true "Login request"
// @Success      200 {object} response.SuccessResponse{data=model.LoginResponse}
// @Failure      400 {object} response.ValidationErrorResponse
// @Failure      401 {object} response.SimpleErrorResponse
// @Failure      500 {object} response.SimpleErrorResponse
// @Router       /login [post]
func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	r.Body = http.MaxBytesReader(w, r.Body, 1<<20)

	var req model.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.SendErrorResponse(w, errors.ErrBadRequest)
		return
	}

	if err := req.Validate(); err != nil {
		response.SendErrorResponse(w, err)
		return
	}

	loginResp, err := h.service.Login(ctx, req)
	if err != nil {
		h.logger.Error(ctx, "failed to login user", zap.Error(err))
		response.SendErrorResponse(w, err)
		return
	}

	response.SendSuccessResponse(w, http.StatusOK, "Login successful", loginResp, nil)
}

// ForgotPassword godoc
// @Summary      Request password reset
// @Description  Generate a password reset token for the user
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        request body model.ForgotPasswordRequest true "Forgot password request"
// @Success      200 {object} response.SuccessResponse{data=model.ForgotPasswordResponse}
// @Failure      400 {object} response.ValidationErrorResponse
// @Failure      404 {object} response.SimpleErrorResponse
// @Failure      500 {object} response.SimpleErrorResponse
// @Router       /forgot-password [post]
func (h *Handler) ForgotPassword(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	r.Body = http.MaxBytesReader(w, r.Body, 1<<20)

	var req model.ForgotPasswordRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.SendErrorResponse(w, errors.ErrBadRequest)
		return
	}

	if err := req.Validate(); err != nil {
		response.SendErrorResponse(w, err)
		return
	}

	forgotResp, err := h.service.ForgotPassword(ctx, req)
	if err != nil {
		h.logger.Error(ctx, "failed to process forgot password", zap.Error(err))
		response.SendErrorResponse(w, err)
		return
	}

	response.SendSuccessResponse(w, http.StatusOK, forgotResp.Message, forgotResp, nil)
}

// ResetPassword godoc
// @Summary      Reset password
// @Description  Reset user password using reset token
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        request body model.ResetPasswordRequest true "Reset password request"
// @Success      200 {object} response.SuccessResponse
// @Failure      400 {object} response.ValidationErrorResponse
// @Failure      401 {object} response.SimpleErrorResponse
// @Failure      500 {object} response.SimpleErrorResponse
// @Router       /reset-password [post]
func (h *Handler) ResetPassword(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	r.Body = http.MaxBytesReader(w, r.Body, 1<<20)

	var req model.ResetPasswordRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.SendErrorResponse(w, errors.ErrBadRequest)
		return
	}

	if err := req.Validate(); err != nil {
		response.SendErrorResponse(w, err)
		return
	}

	if err := h.service.ResetPassword(ctx, req); err != nil {
		h.logger.Error(ctx, "failed to reset password", zap.Error(err))
		response.SendErrorResponse(w, err)
		return
	}

	response.SendSuccessResponse(w, http.StatusOK, "Password reset successfully", nil, nil)
}

// GenerateServiceToken godoc
// @Summary      Generate service token
// @Description  Generate a JWT token for inter-service communication
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        request body model.ServiceTokenRequest true "Service token request"
// @Success      200 {object} response.SuccessResponse{data=model.ServiceTokenResponse}
// @Failure      400 {object} response.ValidationErrorResponse
// @Failure      401 {object} response.SimpleErrorResponse
// @Failure      500 {object} response.SimpleErrorResponse
// @Router       /service-token [post]
func (h *Handler) GenerateServiceToken(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	r.Body = http.MaxBytesReader(w, r.Body, 1<<20)

	var req model.ServiceTokenRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.SendErrorResponse(w, errors.ErrBadRequest)
		return
	}

	if err := req.Validate(); err != nil {
		response.SendErrorResponse(w, err)
		return
	}

	tokenResp, err := h.service.GenerateServiceToken(ctx, req)
	if err != nil {
		h.logger.Error(ctx, "failed to generate service token", zap.Error(err))
		response.SendErrorResponse(w, err)
		return
	}

	response.SendSuccessResponse(w, http.StatusOK, "Service token generated successfully", tokenResp, nil)
}
