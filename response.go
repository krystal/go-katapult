package katapult

import (
	"net/http"
)

type Response struct {
	*http.Response

	Pagination *Pagination
	Error      *ResponseError
}

func NewResponse(r *http.Response) *Response {
	if r == nil {
		r = &http.Response{}
	}

	return &Response{Response: r}
}

type Pagination struct {
	CurrentPage int  `json:"current_page,omitempty"`
	TotalPages  int  `json:"total_pages,omitempty"`
	Total       int  `json:"total,omitempty"`
	PerPage     int  `json:"per_page,omitempty"`
	LargeSet    bool `json:"large_set,omitempty"`
}
