package errs

import "net/http"

type forbidden struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func NewForbidden(message string) *forbidden {
	return &forbidden{
		Code:    http.StatusForbidden,
		Message: message,
	}
}
