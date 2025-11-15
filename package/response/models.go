package response

import "microservice-challenge/package/errors"

type Response struct {
	Status  int            `json:"status,omitempty"`
	Message string         `json:"message,omitempty"`
	Data    any            `json:"data,omitempty"`
	Meta    any            `json:"meta,omitempty"`
	Error   *ErrorResponse `json:"error,omitempty"`
}

type ErrorResponse struct {
	Type    errors.ErrorType `json:"type"`
	Message string           `json:"message,omitempty"`
	Details []FieldError     `json:"details"`
}

type FieldError struct {
	Title       string `json:"title,omitempty"`
	Description string `json:"description,omitempty"`
}

type SuccessResponse struct {
	Data any `json:"data"`
}

type ValidationErrorResponse struct {
	Error *ValidationErrorDetail `json:"error"`
}

type ValidationErrorDetail struct {
	Type    errors.ErrorType `json:"type" example:"validation"`
	Message string           `json:"message" example:"Validation failed"`
	Details []FieldError     `json:"details"`
}

type SimpleErrorResponse struct {
	Error *SimpleErrorDetail `json:"error"`
}

type SimpleErrorDetail struct {
	Type    errors.ErrorType `json:"type" example:"unauthorized"`
	Message string           `json:"message" example:"unauthorized"`
}
