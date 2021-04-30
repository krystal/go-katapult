package core

import (
	"net/url"
	"strconv"
)

type ListOptions struct {
	Page    int
	PerPage int
}

func (s *ListOptions) queryValues() *url.Values {
	values := &url.Values{}

	if s == nil {
		return values
	}

	if s.Page != 0 {
		values.Set("page", strconv.Itoa(s.Page))
	}

	if s.PerPage != 0 {
		values.Set("per_page", strconv.Itoa(s.PerPage))
	}

	return values
}
