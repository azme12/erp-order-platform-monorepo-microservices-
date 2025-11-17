package httphandler

import (
	"encoding/json"
	"microservice-challenge/package/errors"
	"microservice-challenge/package/log"
	"microservice-challenge/package/response"
	"microservice-challenge/services/auth/model"
	"microservice-challenge/services/auth/usecase/auth"
	"net/http"

	"go.uber.org/zap"
)

type Handler struct {
	usecase *auth.Usecase
	logger  log.Logger
}

func NewHandler(usecase *auth.Usecase, logger log.Logger) *Handler {
	return &Handler{
		usecase: usecase,
		logger:  logger,
	}
}

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

	user, err := h.usecase.Register(ctx, req)
	if err != nil {
		h.logger.Error(ctx, "failed to register user", zap.Error(err))
		response.SendErrorResponse(w, err)
		return
	}

	response.SendSuccessResponse(w, http.StatusCreated, "User registered successfully", &user, nil)
}

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

	loginResp, err := h.usecase.Login(ctx, req)
	if err != nil {
		h.logger.Error(ctx, "failed to login user", zap.Error(err))
		response.SendErrorResponse(w, err)
		return
	}

	response.SendSuccessResponse(w, http.StatusOK, "Login successful", &loginResp, nil)
}

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

	forgotResp, err := h.usecase.ForgotPassword(ctx, req)
	if err != nil {
		h.logger.Error(ctx, "failed to process forgot password", zap.Error(err))
		response.SendErrorResponse(w, err)
		return
	}

	response.SendSuccessResponse(w, http.StatusOK, forgotResp.Message, &forgotResp, nil)
}

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

	if err := h.usecase.ResetPassword(ctx, req); err != nil {
		h.logger.Error(ctx, "failed to reset password", zap.Error(err))
		response.SendErrorResponse(w, err)
		return
	}

	response.SendSuccessResponse(w, http.StatusOK, "Password reset successfully", nil, nil)
}

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

	tokenResp, err := h.usecase.GenerateServiceToken(ctx, req)
	if err != nil {
		h.logger.Error(ctx, "failed to generate service token", zap.Error(err))
		response.SendErrorResponse(w, err)
		return
	}

	response.SendSuccessResponse(w, http.StatusOK, "Service token generated successfully", &tokenResp, nil)
}
