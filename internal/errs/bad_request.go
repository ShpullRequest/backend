package errs

import "net/http"

type badRequest struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func NewBadRequest(message string) *badRequest {
	return &badRequest{
		Code:    http.StatusBadRequest,
		Message: message,
	}
}
