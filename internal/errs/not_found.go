package errs

import "net/http"

type notFound struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func NewNotFound(message string) *notFound {
	return &notFound{
		Code:    http.StatusNotFound,
		Message: message,
	}
}
