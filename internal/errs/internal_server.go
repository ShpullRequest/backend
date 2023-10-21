package errs

import "net/http"

type internalServer struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func InternalServer(message string) *internalServer {
	return &internalServer{
		Code:    http.StatusInternalServerError,
		Message: message,
	}
}
