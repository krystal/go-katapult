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

	Pagination *Pagination
	Error      *ResponseError
}

func (r *Response) Paginated() bool {
	return r.Pagination.TotalPages > 0
}

type Pagination struct {
	CurrentPage int  `json:"current_page,omitempty"`
	TotalPages  int  `json:"total_pages,omitempty"`
	Total       int  `json:"total,omitempty"`
	PerPage     int  `json:"per_page,omitempty"`
	LargeSet    bool `json:"large_set,omitempty"`
}

type ResponseError struct {
	Code        string          `json:"code,omitempty"`
	Description string          `json:"description,omitempty"`
	Detail      json.RawMessage `json:"detail,omitempty"`
}

type responseErrorBody struct {
	ErrorInfo *ResponseError `json:"error,omitempty"`
}
