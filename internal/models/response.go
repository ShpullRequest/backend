package models

type Response struct {
	Response interface{} `json:"response"`
}

func NewResponse(v interface{}) *Response {
	return &Response{
		Response: v,
	}
}
