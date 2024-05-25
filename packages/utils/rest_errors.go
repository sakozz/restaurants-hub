package rest_errors

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-playground/validator/v10"
)

type ValidationErrs map[string][]interface{}
type RestErr interface {
	Message() string
	Status() int
	Error() string
	Causes() ValidationErrs
}

type restErr struct {
	ErrMessage string         `json:"message"`
	ErrStatus  int            `json:"status"`
	ErrError   string         `json:"error"`
	ErrCauses  ValidationErrs `json:"causes"`
}

func (e restErr) Error() string {
	return fmt.Sprintf("message: %s - status: %d - error: %s - causes: %v",
		e.ErrMessage, e.ErrStatus, e.ErrError, e.ErrCauses)
}

func (e restErr) Message() string {
	return e.ErrMessage
}

func (e restErr) Status() int {
	return e.ErrStatus
}

func (e restErr) Causes() ValidationErrs {
	return e.ErrCauses
}

func NewValidationError(causes *ValidationErrs) RestErr {
	return restErr{
		ErrMessage: "Validation error",
		ErrStatus:  http.StatusUnprocessableEntity,
		ErrError:   "validation_error",
		ErrCauses:  *causes,
	}
}

func StructValidationErrors(errs error) *ValidationErrs {
	causes := ValidationErrs{}
	for _, err := range errs.(validator.ValidationErrors) {
		causes[err.StructField()] = append(causes[err.StructField()], formattedValidationErrors(err))
	}
	return &causes
}

func formattedValidationErrors(err validator.FieldError) interface{} {
	switch err.ActualTag() {
	case "required":
		return map[string]interface{}{
			"error": "required",
		}
	case "email":
		return map[string]interface{}{
			"error": "invalid_email_format",
		}
	case "min":
		minLength, _ := strconv.ParseInt(err.Param(), 10, 64)
		return map[string]interface{}{
			"error":    "min_length_required",
			"expected": minLength,
			"provided": len(err.Value().(string)),
		}
	case "max":
		maxLength, _ := strconv.ParseInt(err.Param(), 10, 64)
		return map[string]interface{}{
			"error":    "max_length_exceeded",
			"expected": maxLength,
			"provided": len(err.Value().(string)),
		}
	default:
		return err.Tag()
	}
}

func FormattedDbValidationError(attr string, constraint string) interface{} {
	switch constraint {
	case "uniqueness":
		return map[string]interface{}{
			"error": "must_be_unique",
		}
	case "not_found":
		return map[string]interface{}{
			"error": "record_not_found",
		}
	default:
		return "validation_error"
	}
}

func NewRestError(message string, status int, err string, causes ValidationErrs) RestErr {
	return restErr{
		ErrMessage: message,
		ErrStatus:  status,
		ErrError:   err,
		ErrCauses:  causes,
	}
}

func NewRestErrorFromBytes(bytes []byte) (RestErr, error) {
	var apiErr RestErr
	if err := json.Unmarshal(bytes, &apiErr); err != nil {
		return nil, errors.New("invalid json")
	}
	return apiErr, nil
}

func NewBadRequestError(message string) RestErr {
	return restErr{
		ErrMessage: message,
		ErrStatus:  http.StatusBadRequest,
		ErrError:   "bad_request",
	}
}

func NewNotFoundError(message string) RestErr {
	return restErr{
		ErrMessage: message,
		ErrStatus:  http.StatusNotFound,
		ErrError:   "not_found",
	}
}

func InvalidError(message string) RestErr {
	return restErr{
		ErrMessage: message,
		ErrStatus:  http.StatusUnprocessableEntity,
		ErrError:   "invalid_record",
	}
}

func NewUnauthorizedError(message string) RestErr {
	return restErr{
		ErrMessage: message,
		ErrStatus:  http.StatusUnauthorized,
		ErrError:   "unauthorized",
	}
}

func NewInternalServerError(err error) RestErr {
	if err != nil {
		fmt.Println("Server Error: ", err)
	}
	return restErr{
		ErrMessage: "Internal server error",
		ErrStatus:  http.StatusInternalServerError,
		ErrError:   "internal_server_error",
	}
}
