package v1

import "github.com/tweethub/backend/cmd/statistics/storage"

const (
	InvalidRequestMsg      = "invalid request"
	NotFoundMsg            = "not found"
	InternalServerErrorMsg = "internal server error"
)

// RelevanceResponse represents relevance response.
type RelevanceResponse = storage.Relevance

// ErrorResponse represents http error response.
type ErrorResponse struct {
	Message string `json:"message"`
}

// NewErrorResponse returns custom error response.
func NewErrorResponse(msg string) *ErrorResponse {
	return &ErrorResponse{
		Message: msg,
	}
}

// NewInvalidRequestError returns invalid request error response.
func NewInvalidRequestError() *ErrorResponse {
	return &ErrorResponse{
		Message: InvalidRequestMsg,
	}
}

// NewNotFoundError returns not found error response.
func NewNotFoundError() *ErrorResponse {
	return &ErrorResponse{
		Message: NotFoundMsg,
	}
}

// NewInternalServerError returns internal server error response.
func NewInternalServerError() *ErrorResponse {
	return &ErrorResponse{
		Message: InternalServerErrorMsg,
	}
}