package errors

import (
	"database/sql"
	"errors"
	"net/http"
)

type ErrorType string

const (
	ErrorTypeValidation      ErrorType = "validation"
	ErrorTypeNotFound        ErrorType = "not_found"
	ErrorTypeConflict        ErrorType = "conflict"
	ErrorTypeUnauthorized    ErrorType = "unauthorized"
	ErrorTypeForbidden       ErrorType = "forbidden"
	ErrorTypeRateLimit       ErrorType = "rate_limit"
	ErrorTypeInternal        ErrorType = "internal"
	ErrorTypeTimeout         ErrorType = "timeout"
	ErrorTypeUnavailable     ErrorType = "unavailable"
	ErrorTypeDatabase        ErrorType = "database"
	ErrorTypeExternal        ErrorType = "external_service"
	ErrorTypeNetwork         ErrorType = "network"
	ErrorTypeBadRequest      ErrorType = "bad_request"
	ErrorTypeInvalidFormat   ErrorType = "invalid_format"
	ErrorTypePayloadTooLarge ErrorType = "payload_too_large"
	ErrorTypeUnknown         ErrorType = "unknown"
)

var (
	ErrNoRows              = sql.ErrNoRows
	ErrUnexpected          = errors.New("unexpected error occurred. please try again later")
	ErrUnauthorized        = errors.New("unauthorized")
	ErrInternalServerError = errors.New("internal server error. please try again later")
	ErrBadRequest          = errors.New("bad request")
	ErrInvalidData         = errors.New("invalid input data")
	ErrNotFound            = errors.New("resource not found")
	ErrConflict            = errors.New("resource conflict")
	ErrRequestTimeout      = errors.New("request timeout")
	ErrForbidden           = errors.New("forbidden")
	ErrInvalidToken        = errors.New("invalid or expired token")
	ErrTokenExpired        = errors.New("token has expired")
)

var ErrorMap = map[error]int{
	ErrBadRequest:          http.StatusBadRequest,
	ErrInvalidData:         http.StatusBadRequest,
	ErrUnexpected:          http.StatusInternalServerError,
	ErrUnauthorized:        http.StatusUnauthorized,
	ErrInternalServerError: http.StatusInternalServerError,
	ErrNotFound:            http.StatusNotFound,
	ErrConflict:            http.StatusConflict,
	ErrRequestTimeout:      http.StatusRequestTimeout,
	ErrForbidden:           http.StatusForbidden,
	ErrInvalidToken:        http.StatusBadRequest,
	ErrTokenExpired:        http.StatusBadRequest,
}

var ErrorTypeMap = map[error]ErrorType{
	ErrBadRequest:          ErrorTypeBadRequest,
	ErrInternalServerError: ErrorTypeInternal,
	ErrInvalidData:         ErrorTypeBadRequest,
	ErrUnauthorized:        ErrorTypeUnauthorized,
	ErrUnexpected:          ErrorTypeInternal,
	ErrNotFound:            ErrorTypeNotFound,
	ErrConflict:            ErrorTypeConflict,
	ErrRequestTimeout:      ErrorTypeTimeout,
	ErrForbidden:           ErrorTypeForbidden,
	ErrInvalidToken:        ErrorTypeBadRequest,
	ErrTokenExpired:        ErrorTypeBadRequest,
}
