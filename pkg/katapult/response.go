package katapult

import (
	"encoding/json"
	"net/http"
)

func newResponse(r *http.Response) *Response {
	return &Response{Response: r}
}

type Response struct {
	*http.Response

	Error *ErrorResponse
}

type ErrorResponseBody struct {
	Error *ErrorResponse `json:"error,omitempty"`
}

type ErrorResponse struct {
	Code        string          `json:"code,omitempty"`
	Description string          `json:"description,omitempty"`
	Detail      json.RawMessage `json:"detail,omitempty"`
}
