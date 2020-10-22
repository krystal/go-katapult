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

	Error *ResponseError
}

type errorResponseBody struct {
	Error *ResponseError `json:"error,omitempty"`
}

type ResponseError struct {
	Code        string          `json:"code,omitempty"`
	Description string          `json:"description,omitempty"`
	Detail      json.RawMessage `json:"detail,omitempty"`
}
