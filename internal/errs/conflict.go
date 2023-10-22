package errs

import "net/http"

type conflict struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func NewConflict(message string) *conflict {
	return &conflict{
		Code:    http.StatusConflict,
		Message: message,
	}
}
