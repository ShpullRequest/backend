package errs

import "net/http"

type badGateway struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func NewBadGateway(message string) *badGateway {
	return &badGateway{
		Code:    http.StatusBadGateway,
		Message: message,
	}
}
