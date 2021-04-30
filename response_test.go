package katapult

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPagination_JSONMarshaling(t *testing.T) {
	tests := []struct {
		name string
		obj  *Pagination
	}{
		{
			name: "empty",
			obj:  &Pagination{},
		},
		{
			name: "full",
			obj: &Pagination{
				CurrentPage: 5,
				TotalPages:  10,
				Total:       190,
				PerPage:     20,
				LargeSet:    true,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testJSONMarshaling(t, tt.obj)
		})
	}
}

func TestResponseError_JSONMarshaling(t *testing.T) {
	tests := []struct {
		name string
		obj  *ResponseError
	}{
		{
			name: "empty",
			obj:  &ResponseError{},
		},
		{
			name: "full",
			obj: &ResponseError{
				Code:        "code",
				Description: "desc",
				Detail:      json.RawMessage(`[{}]`),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testJSONMarshaling(t, tt.obj)
		})
	}
}

func Test_responseErrorBody_JSONMarshaling(t *testing.T) {
	tests := []struct {
		name string
		obj  *responseErrorBody
	}{
		{
			name: "empty",
			obj:  &responseErrorBody{},
		},
		{
			name: "full",
			obj: &responseErrorBody{
				ErrorInfo: &ResponseError{Code: "nope"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testJSONMarshaling(t, tt.obj)
		})
	}
}

func Test_newResponse(t *testing.T) {
	tests := []struct {
		name string
		r    *http.Response
		want *Response
	}{
		{
			name: "given nil",
			r:    nil,
			want: &Response{Response: &http.Response{}},
		},
		{
			name: "given http.Response",
			r:    &http.Response{StatusCode: http.StatusEarlyHints},
			want: &Response{
				Response: &http.Response{StatusCode: http.StatusEarlyHints},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewResponse(tt.r)

			assert.Equal(t, tt.want, got)
			assert.IsType(t, int(0), got.StatusCode)
		})
	}
}
