package katapult

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/jimeh/undent"
	"github.com/krystal/go-katapult/internal/test"
	"github.com/stretchr/testify/assert"
)

func TestErr(t *testing.T) {
	assert.EqualError(t, Err, "katapult")
}

func TestErrRequest(t *testing.T) {
	assert.EqualError(t, ErrRequest, "katapult: request")
	assert.ErrorIs(t, ErrRequest, Err)
}

func TestErrConfig(t *testing.T) {
	assert.EqualError(t, ErrConfig, "katapult: config")
	assert.ErrorIs(t, ErrConfig, Err)
}

func TestErrResponse(t *testing.T) {
	assert.EqualError(t, ErrResponse, "katapult")
	assert.ErrorIs(t, ErrResponse, Err)
}

func TestErrResourceNotFound(t *testing.T) {
	assert.EqualError(t, ErrResourceNotFound, "katapult: not_found")
	assert.ErrorIs(t, ErrResourceNotFound, ErrNotFound)
}

func TestErrRouteNotFound(t *testing.T) {
	assert.EqualError(t, ErrRouteNotFound, "katapult: route_not_found")
	assert.ErrorIs(t, ErrRouteNotFound, ErrResponse)
}

func TestErrUnexpectedResponse(t *testing.T) {
	assert.EqualError(t, ErrUnexpectedResponse, "katapult: unexpected_response")
	assert.ErrorIs(t, ErrUnexpectedResponse, ErrResponse)
}

func TestErrUnknown(t *testing.T) {
	assert.EqualError(t, ErrUnknown, "katapult: unknown_error")
	assert.ErrorIs(t, ErrUnknown, ErrResponse)
}

func TestErrBadRequest(t *testing.T) {
	assert.EqualError(t, ErrBadRequest, "katapult: bad_request")
	assert.ErrorIs(t, ErrBadRequest, ErrResponse)
}

func TestErrUnauthorized(t *testing.T) {
	assert.EqualError(t, ErrUnauthorized, "katapult: unauthorized")
	assert.ErrorIs(t, ErrUnauthorized, ErrResponse)
}

func TestErrForbidden(t *testing.T) {
	assert.EqualError(t, ErrForbidden, "katapult: unauthorized")
	assert.ErrorIs(t, ErrForbidden, ErrUnauthorized)
}

func TestErrNotFound(t *testing.T) {
	assert.EqualError(t, ErrNotFound, "katapult: not_found")
	assert.ErrorIs(t, ErrNotFound, Err)
}

func TestErrNotAcceptable(t *testing.T) {
	assert.EqualError(t, ErrNotAcceptable, "katapult: not_acceptable")
	assert.ErrorIs(t, ErrNotAcceptable, Err)
}

func TestErrConflict(t *testing.T) {
	assert.EqualError(t, ErrConflict, "katapult: conflict")
	assert.ErrorIs(t, ErrConflict, ErrResponse)
}

func TestErrUnprocessableEntity(t *testing.T) {
	assert.EqualError(t,
		ErrUnprocessableEntity, "katapult: unprocessable_entity",
	)
	assert.ErrorIs(t, ErrUnprocessableEntity, ErrResponse)
}

func TestErrInternalServerError(t *testing.T) {
	assert.EqualError(t,
		ErrInternalServerError, "katapult: internal_server_error",
	)
	assert.ErrorIs(t, ErrInternalServerError, ErrResponse)
}

func TestErrBadGateway(t *testing.T) {
	assert.EqualError(t, ErrBadGateway, "katapult: bad_gateway")
	assert.ErrorIs(t, ErrBadGateway, ErrResponse)
}

func TestErrServiceUnavailable(t *testing.T) {
	assert.EqualError(t, ErrServiceUnavailable, "katapult: service_unavailable")
	assert.ErrorIs(t, ErrServiceUnavailable, ErrResponse)
}

func TestErrGatewayTimeout(t *testing.T) {
	assert.EqualError(t, ErrGatewayTimeout, "katapult: gateway_timeout")
	assert.ErrorIs(t, ErrGatewayTimeout, ErrResponse)
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
				Error: &ResponseError{Code: "nope"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			test.CustomJSONMarshaling(t, tt.obj, nil)
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
			test.CustomJSONMarshaling(t, tt.obj, nil)
		})
	}
}

func TestNewResponseError(t *testing.T) {
	tests := []struct {
		name        string
		httpStatus  int
		code        string
		description string
		rawDetail   json.RawMessage
		parent      error
	}{
		{
			name:        "bad_request",
			httpStatus:  http.StatusBadRequest,
			code:        "bad_request",
			description: "You did: BadRequest",
			rawDetail:   json.RawMessage(`{"code":"bad_request"}`),
			parent:      ErrBadRequest,
		},
		{
			name:        "unauthorized",
			httpStatus:  http.StatusUnauthorized,
			code:        "unauthorized",
			description: "You did: Unauthorized",
			rawDetail:   json.RawMessage(`{"code":"unauthorized"}`),
			parent:      ErrUnauthorized,
		},
		{
			name:        "forbidden",
			httpStatus:  http.StatusForbidden,
			code:        "forbidden",
			description: "You did: Forbidden",
			rawDetail:   json.RawMessage(`{"code":"forbidden"}`),
			parent:      ErrForbidden,
		},
		{
			name:        "resource_not_found",
			httpStatus:  http.StatusNotFound,
			code:        "resource_not_found",
			description: "You did: resource_not_found",
			rawDetail:   json.RawMessage(`{"code":"resource_not_found"}`),
			parent:      ErrResourceNotFound,
		},
		{
			name:        "not_acceptable",
			httpStatus:  http.StatusNotAcceptable,
			code:        "not_acceptable",
			description: "You did: not_acceptable",
			rawDetail:   json.RawMessage(`{"code":"not_acceptable"}`),
			parent:      ErrNotAcceptable,
		},
		{
			name:        "conflict",
			httpStatus:  http.StatusConflict,
			code:        "conflict",
			description: "You did: conflict",
			rawDetail:   json.RawMessage(`{"code":"conflict"}`),
			parent:      ErrConflict,
		},
		{
			name:        "unprocessable_entity",
			httpStatus:  http.StatusUnprocessableEntity,
			code:        "unprocessable_entity",
			description: "You did: unprocessable_entity",
			rawDetail:   json.RawMessage(`{"code":"unprocessable_entity"}`),
			parent:      ErrUnprocessableEntity,
		},
		{
			name:        "too_many_requests",
			httpStatus:  http.StatusTooManyRequests,
			code:        "too_many_requests",
			description: "You did: too_many_requests",
			rawDetail:   json.RawMessage(`{"code":"too_many_requests"}`),
			parent:      ErrTooManyRequests,
		},
		{
			name:        "internal_server_error",
			httpStatus:  http.StatusInternalServerError,
			code:        "internal_server_error",
			description: "You did: internal_server_error",
			rawDetail:   json.RawMessage(`{"code":"internal_server_error"}`),
			parent:      ErrInternalServerError,
		},
		{
			name:        "bad_gateway",
			httpStatus:  http.StatusBadGateway,
			code:        "bad_gateway",
			description: "You did: bad_gateway",
			rawDetail:   json.RawMessage(`{"code":"bad_gateway"}`),
			parent:      ErrBadGateway,
		},
		{
			name:        "service_unavailable",
			httpStatus:  http.StatusServiceUnavailable,
			code:        "service_unavailable",
			description: "You did: service_unavailable",
			rawDetail:   json.RawMessage(`{"code":"service_unavailable"}`),
			parent:      ErrServiceUnavailable,
		},
		{
			name:        "gateway_timeout",
			httpStatus:  http.StatusGatewayTimeout,
			code:        "gateway_timeout",
			description: "You did: gateway_timeout",
			rawDetail:   json.RawMessage(`{"code":"gateway_timeout"}`),
			parent:      ErrGatewayTimeout,
		},
		{
			name:        "unknown",
			httpStatus:  http.StatusTeapot,
			code:        "i_am_a_teapot",
			description: "You did: want_lemon",
			rawDetail:   json.RawMessage(`{"code":"only_have_earl_grey"}`),
			parent:      ErrUnknown,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewResponseError(
				tt.httpStatus,
				tt.code,
				tt.description,
				tt.rawDetail,
			)
			want := &ResponseError{
				parent:      tt.parent,
				Code:        tt.code,
				Description: tt.description,
				Detail:      tt.rawDetail,
			}

			assert.Equal(t, want, got)
		})
	}
}

func TestResponseError_Error(t *testing.T) {
	type fields struct {
		code        string
		description string
		detail      json.RawMessage
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "code only",
			fields: fields{
				code: "disk_not_found",
			},
			want: "disk_not_found",
		},
		{
			name: "code and description",
			fields: fields{
				code:        "disk_not_found",
				description: "Specified disk could not be found",
			},
			want: "disk_not_found: Specified disk could not be found",
		},
		{
			name: "code, description and detail",
			fields: fields{
				code:        "disk_not_found",
				description: "Specified disk could not be found",
				detail:      json.RawMessage(`{"id":"disk_oVsxmtVXYQSmmk8g"}`),
			},
			want: undent.String(`
				disk_not_found: Specified disk could not be found -- {
				  "id": "disk_oVsxmtVXYQSmmk8g"
				}`,
			),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			respErr := &ResponseError{
				Code:        tt.fields.code,
				Description: tt.fields.description,
				Detail:      tt.fields.detail,
			}

			got := respErr.Error()

			assert.Equal(t, tt.want, got)
		})
	}
}

func TestResponseError_Is(t *testing.T) {
	tests := []struct {
		name   string
		parent error
		target error
		want   bool
	}{
		{
			name:   "child",
			parent: ErrNotFound,
			target: ErrResourceNotFound,
			want:   false,
		},
		{
			name:   "nil parent",
			parent: nil,
			target: ErrNotFound,
			want:   false,
		},
		{
			name:   "exact parent",
			parent: ErrNotFound,
			target: ErrNotFound,
			want:   true,
		},
		{
			name:   "grandparent",
			parent: ErrResourceNotFound,
			target: ErrNotFound,
			want:   true,
		},
		{
			name:   "great-grandparent",
			parent: ErrNotFound,
			target: ErrResponse,
			want:   true,
		},
		{
			name:   "great-great-grandparent",
			parent: ErrNotFound,
			target: Err,
			want:   true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			respErr := &ResponseError{parent: tt.parent}

			if tt.want {
				assert.ErrorIs(t, respErr, tt.target)
			} else {
				assert.NotErrorIs(t, respErr, tt.target)
			}
		})
	}
}

func TestResponseError_Unwrap(t *testing.T) {
	tests := []struct {
		name   string
		parent error
		want   error
	}{
		{
			name:   "nil parent",
			parent: nil,
			want:   nil,
		},
		{
			name:   "ErrNotFound",
			parent: ErrNotFound,
			want:   ErrNotFound,
		},
		{
			name:   "ErrInternalServerError",
			parent: ErrInternalServerError,
			want:   ErrInternalServerError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			respErr := &ResponseError{parent: tt.parent}

			got := respErr.Unwrap()

			assert.Equal(t, tt.want, got)
		})
	}
}

func TestNewCommonError(t *testing.T) {
	type args struct {
		parent      error
		code        string
		description string
	}
	tests := []struct {
		name string
		args args
		want CommonError
	}{
		{
			name: "empty",
			args: args{
				parent:      nil,
				code:        "",
				description: "",
			},
			want: CommonError{},
		},
		{
			name: "full",
			args: args{
				parent:      ErrResourceNotFound,
				code:        "disk_not_found",
				description: "Specified disk could not be found",
			},
			want: CommonError{
				parent:      ErrResourceNotFound,
				Code:        "disk_not_found",
				Description: "Specified disk could not be found",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewCommonError(
				tt.args.parent,
				tt.args.code,
				tt.args.description,
			)

			assert.Equal(t, tt.want, got)
		})
	}
}

func TestCommonError_BaseError(t *testing.T) {
	type fields struct {
		parent error
		code   string
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "nil parent",
			fields: fields{
				parent: nil,
				code:   "disk_not_found",
			},
			want: "disk_not_found",
		},
		{
			name: "parent",
			fields: fields{
				parent: ErrUnauthorized,
				code:   "foobar",
			},
			want: ErrUnauthorized.Error(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			comErr := &CommonError{
				parent: tt.fields.parent,
				Code:   tt.fields.code,
			}

			got := comErr.BaseError()

			assert.Equal(t, tt.want, got)
		})
	}
}

func TestCommonError_Error(t *testing.T) {
	type fields struct {
		parent      error
		code        string
		description string
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "code only",
			fields: fields{
				parent:      nil,
				code:        "disk_not_found",
				description: "",
			},
			want: "disk_not_found",
		},
		{
			name: "code and description",
			fields: fields{
				parent:      nil,
				code:        "disk_not_found",
				description: "Specified disk could not be found",
			},
			want: "disk_not_found: Specified disk could not be found",
		},
		{
			name: "parent and code",
			fields: fields{
				parent: ErrUnauthorized,
				code:   "foobar",
			},
			want: ErrUnauthorized.Error(),
		},
		{
			name: "parent, code and description",
			fields: fields{
				parent:      ErrUnauthorized,
				code:        "invalid_api_key",
				description: "Invalid API key",
			},
			want: ErrUnauthorized.Error() + ": Invalid API key",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			comErr := &CommonError{
				parent:      tt.fields.parent,
				Code:        tt.fields.code,
				Description: tt.fields.description,
			}

			got := comErr.Error()

			assert.Equal(t, tt.want, got)
		})
	}
}

func TestCommonError_Is(t *testing.T) {
	tests := []struct {
		name   string
		parent error
		target error
		want   bool
	}{
		{
			name:   "child",
			parent: ErrNotFound,
			target: ErrResourceNotFound,
			want:   false,
		},
		{
			name:   "nil parent",
			parent: nil,
			target: ErrNotFound,
			want:   false,
		},
		{
			name:   "exact parent",
			parent: ErrNotFound,
			target: ErrNotFound,
			want:   true,
		},
		{
			name:   "grandparent",
			parent: ErrResourceNotFound,
			target: ErrNotFound,
			want:   true,
		},
		{
			name:   "great-grandparent",
			parent: ErrNotFound,
			target: ErrResponse,
			want:   true,
		},
		{
			name:   "great-great-grandparent",
			parent: ErrNotFound,
			target: Err,
			want:   true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			respErr := &CommonError{parent: tt.parent}

			if tt.want {
				assert.ErrorIs(t, respErr, tt.target)
			} else {
				assert.NotErrorIs(t, respErr, tt.target)
			}
		})
	}
}

func TestCommonError_Unwrap(t *testing.T) {
	tests := []struct {
		name   string
		parent error
		want   error
	}{
		{
			name:   "nil parent",
			parent: nil,
			want:   nil,
		},
		{
			name:   "ErrNotFound",
			parent: ErrNotFound,
			want:   ErrNotFound,
		},
		{
			name:   "ErrInternalServerError",
			parent: ErrInternalServerError,
			want:   ErrInternalServerError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			respErr := &CommonError{parent: tt.parent}

			got := respErr.Unwrap()

			assert.Equal(t, tt.want, got)
		})
	}
}
