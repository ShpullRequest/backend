package errs

import "net/http"

type notAuthorized struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func NotAuthorized(message string) *notAuthorized {
	return &notAuthorized{
		Code:    http.StatusUnauthorized,
		Message: message,
	}
}
