package katapult

import (
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestListOptions_Values(t *testing.T) {
	tests := []struct {
		name string
		opts *ListOptions
		want *url.Values
	}{
		{
			name: "nil *ListOptions",
			opts: nil,
			want: &url.Values{},
		},
		{
			name: "empty *ListOptions",
			opts: &ListOptions{},
			want: &url.Values{},
		},
		{
			name: "zero'd values",
			opts: &ListOptions{Page: 0, PerPage: 0},
			want: &url.Values{},
		},
		{
			name: "non-zero Page value",
			opts: &ListOptions{Page: 3},
			want: &url.Values{"page": []string{"3"}},
		},
		{
			name: "non-zero PerPage value",
			opts: &ListOptions{PerPage: 15},
			want: &url.Values{"per_page": []string{"15"}},
		},
		{
			name: "non-zero Page and PerPage values",
			opts: &ListOptions{Page: 5, PerPage: 15},
			want: &url.Values{
				"page":     []string{"5"},
				"per_page": []string{"15"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.opts.Values()

			assert.Equal(t, tt.want, got)
		})
	}
}
