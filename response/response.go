package response

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// SuccessResponse holds response body for success response
type SuccessResponse struct {
	Data    interface{} `json:"data,omitempty"`
	Message string      `json:"message,omitempty"`
	Errors  []ErrorInfo `json:"errors,omitempty"`
	Meta    MetaInfo    `json:"meta"`
}

// ErrorResponse holds data for error response
type ErrorResponse struct {
	Errors []ErrorInfo `json:"errors"`
	Meta   MetaInfo    `json:"meta"`
}

// MetaInfo holds meta data
type MetaInfo struct {
	HTTPStatus int         `json:"http_status"`
	Offset     int         `json:"offset,omitempty"`
	Limit      int         `json:"limit,omitempty"`
	Total      int64       `json:"total,omitempty"`
	Sort       string      `json:"sort,omitempty"`
	Facets     interface{} `json:"facets,omitempty"`
}

// ErrorInfo holds error detail
type ErrorInfo struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
	Field   string `json:"field,omitempty"`
}

// CustomError holds data for customized error
type CustomError struct {
	Message  string
	Field    string
	Code     int
	HTTPCode int
}

var (
	// DefaultError represents Unexpected Internal Server error
	DefaultError = CustomError{
		Message:  "Internal Server Error. Please try again in a few minuts",
		Code:     9001,
		HTTPCode: http.StatusInternalServerError,
	}

	// UserUnauthorizedError represents User unauthorized error
	UserUnauthorizedError = CustomError{
		Message:  "You can't access this page",
		Code:     9002,
		HTTPCode: http.StatusForbidden,
	}

	// BadRequestError represents invalid request body or parameter error
	BadRequestError = CustomError{
		Message:  "Request Body or Parameter is not valid",
		Code:     9003,
		HTTPCode: http.StatusBadRequest,
	}
)

func (c CustomError) Error() string {
	return c.Message
}

// Error is a function to convert error to string.
// It exists to satisfy error interface
func (resp ErrorResponse) Error() string {
	msg := fmt.Sprintf("(%d) ", resp.Meta.HTTPStatus)
	for idx, err := range resp.Errors {
		if idx > 0 {
			msg += ", "
		}
		msg += err.Message
	}
	return msg
}

func (err ErrorInfo) Error() string {
	msg := fmt.Sprintf("%s (%d) ", err.Message, err.Code)
	if err.Field != "" {
		msg += fmt.Sprintf(", Field: %s", err.Field)
	}
	return msg
}

// BuildSuccess is a function to create success SuccessResponse
func BuildSuccess(data interface{}, meta MetaInfo) SuccessResponse {
	return SuccessResponse{
		Data: data,
		Meta: meta,
	}
}

// BuildError is a function to create error ErrorResponse
func BuildError(errors []error) ErrorResponse {
	if len(errors) == 0 {
		errors = []error{DefaultError}
	}

	return BuildErrors(errors)
}

// BuildErrors is a function to create error ErrorResponse
func BuildErrors(errors []error) ErrorResponse {
	var (
		ce         CustomError
		ok         bool
		errorInfos []ErrorInfo
	)

	for _, err := range errors {
		ce, ok = err.(CustomError)
		if !ok {
			ce = DefaultError
		}

		errorInfo := ErrorInfo{
			Code:    ce.Code,
			Field:   ce.Field,
			Message: ce.Message,
		}

		errorInfos = append(errorInfos, errorInfo)
	}

	return ErrorResponse{
		Errors: errorInfos,
		Meta: MetaInfo{
			HTTPStatus: ce.HTTPCode,
		},
	}
}

func BuildErrorAndStatus(err error, fieldName string) (ErrorResponse, int) {
	if ce, ok := err.(CustomError); ok {
		return BuildError([]error{ce}), ce.HTTPCode
	}

	return BuildError([]error{DefaultError}), DefaultError.HTTPCode
}

// Write is a function to write data in json format
func Write(w http.ResponseWriter, result interface{}, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(result)
}
