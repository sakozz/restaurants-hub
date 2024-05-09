package rest_errors

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-playground/validator/v10"
)

type RestErr interface {
	Message() string
	Status() int
	Error() string
	Causes() []interface{}
	ValidationErrors() map[string][]interface{}
}

type restErr struct {
	ErrMessage     string                   `json:"message"`
	ErrStatus      int                      `json:"status"`
	ErrError       string                   `json:"error"`
	ErrCauses      []interface{}            `json:"causes"`
	ErrValidations map[string][]interface{} `json:"validationErrors"`
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

func (e restErr) Causes() []interface{} {
	return e.ErrCauses
}

func (e restErr) ValidationErrors() map[string][]interface{} {
	return e.ErrValidations
}

func NewValidationError(message string, errs error) RestErr {

	result := restErr{
		ErrMessage: message,
		ErrStatus:  http.StatusUnprocessableEntity,
		ErrError:   "validation_error",
	}
	var causes = map[string][]interface{}{}
	for _, err := range errs.(validator.ValidationErrors) {
		causes[err.StructField()] = append(causes[err.StructField()], formattedValidationErrors(err))
	}

	result.ErrValidations = causes
	return result
}

func formattedValidationErrors(err validator.FieldError) interface{} {
	switch err.ActualTag() {
	case "required":
		return map[string]interface{}{
			"error":    "required",
			"expected": err.Value(),
			"provided": nil,
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

func NewRestError(message string, status int, err string, causes []interface{}, validationErrs map[string][]interface{}) RestErr {
	return restErr{
		ErrMessage:     message,
		ErrStatus:      status,
		ErrError:       err,
		ErrCauses:      causes,
		ErrValidations: validationErrs,
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

func NewInternalServerError(message string, err error) RestErr {
	result := restErr{
		ErrMessage: message,
		ErrStatus:  http.StatusInternalServerError,
		ErrError:   "internal_server_error",
	}
	if err != nil {
		result.ErrCauses = append(result.ErrCauses, err.Error())
	}
	return result
}
