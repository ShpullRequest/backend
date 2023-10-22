package models

type ErrorResponse struct {
	Error interface{} `json:"error"`
}

func NewErrorResponse(v interface{}) *ErrorResponse {
	return &ErrorResponse{
		Error: v,
	}
}
