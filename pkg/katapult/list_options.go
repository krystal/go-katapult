package katapult

import (
	"fmt"
	"net/url"
)

type ListOptions struct {
	Page    int
	PerPage int
}

func (s *ListOptions) Values() *url.Values {
	u := &url.Values{}

	if s != nil {
		if s.Page != 0 {
			u.Set("page", fmt.Sprint(s.Page))
		}

		if s.PerPage != 0 {
			u.Set("per_page", fmt.Sprint(s.PerPage))
		}
	}

	return u
}
