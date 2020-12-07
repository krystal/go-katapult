package katapult

import (
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestListOptions_queryValues(t *testing.T) {
	tests := []struct {
		name string
		obj  *ListOptions
		want *url.Values
	}{
		{
			name: "nil *ListOptions",
			obj:  nil,
			want: &url.Values{},
		},
		{
			name: "empty *ListOptions",
			obj:  &ListOptions{},
			want: &url.Values{},
		},
		{
			name: "zero'd values",
			obj:  &ListOptions{Page: 0, PerPage: 0},
			want: &url.Values{},
		},
		{
			name: "non-zero Page value",
			obj:  &ListOptions{Page: 3},
			want: &url.Values{"page": []string{"3"}},
		},
		{
			name: "non-zero PerPage value",
			obj:  &ListOptions{PerPage: 15},
			want: &url.Values{"per_page": []string{"15"}},
		},
		{
			name: "non-zero Page and PerPage values",
			obj:  &ListOptions{Page: 5, PerPage: 15},
			want: &url.Values{
				"page":     []string{"5"},
				"per_page": []string{"15"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.obj.queryValues()

			assert.Equal(t, tt.want, got)
		})
	}
}
