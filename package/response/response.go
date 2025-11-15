package response

import (
	"encoding/json"
	"microservice-challenge/package/errors"
	"net/http"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

func SendSuccessResponse(w http.ResponseWriter, statusCode int, message string, data, meta any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	if err := json.NewEncoder(w).Encode(Response{
		Status:  statusCode,
		Message: message,
		Data:    data,
		Meta:    meta,
	}); err != nil {
		w.WriteHeader(errors.ErrorMap[errors.ErrUnexpected])
		json.NewEncoder(w).Encode(Response{
			Status: errors.ErrorMap[errors.ErrUnexpected],
			Error: &ErrorResponse{
				Type:    errors.ErrorTypeUnknown,
				Message: errors.ErrUnexpected.Error(),
			},
		})
		return
	}
}

func SendErrorResponse(w http.ResponseWriter, err error) {
	w.Header().Set("Content-Type", "application/json")

	if ve, ok := err.(validation.Errors); ok {
		fieldErr := ErrorFields(ve)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(Response{
			Status: http.StatusBadRequest,
			Error: &ErrorResponse{
				Type:    errors.ErrorTypeValidation,
				Message: "invalid input data",
				Details: fieldErr,
			},
		})
		return
	}

	statusCode, ok := errors.ErrorMap[err]
	if !ok {
		w.WriteHeader(errors.ErrorMap[errors.ErrUnexpected])
		json.NewEncoder(w).Encode(Response{
			Status: errors.ErrorMap[errors.ErrUnexpected],
			Error: &ErrorResponse{
				Type:    errors.ErrorTypeUnknown,
				Message: errors.ErrUnexpected.Error(),
			},
		})
		return
	}

	w.WriteHeader(statusCode)
	if err := json.NewEncoder(w).Encode(Response{
		Status: statusCode,
		Error: &ErrorResponse{
			Type:    errors.ErrorTypeMap[err],
			Message: err.Error(),
		},
	}); err != nil {
		w.WriteHeader(errors.ErrorMap[errors.ErrUnexpected])
		json.NewEncoder(w).Encode(Response{
			Status: errors.ErrorMap[errors.ErrUnexpected],
			Error: &ErrorResponse{
				Type:    errors.ErrorTypeUnknown,
				Message: errors.ErrUnexpected.Error(),
			},
		})
	}
}

func ErrorFields(err error) []FieldError {
	var errs []FieldError

	if data, ok := err.(validation.Errors); ok {
		for i, v := range data {
			nestedErrors := ErrorFields(v)
			if len(nestedErrors) > 0 {
				errs = append(errs, nestedErrors...)
			} else {
				errs = append(errs, FieldError{
					Title:       i,
					Description: v.Error(),
				})
			}
		}

		return errs
	}

	return nil
}
