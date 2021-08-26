package katapult

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

//go:generate go run github.com/krystal/go-katapult/tools/codegen -t errors -p katapult -o . -i (Apia|Rapid)\/.* -f ./schemas/core/v1.json

var (
	// Err is the top-most parent of any error returned by katapult.
	Err = errors.New("katapult")

	// ErrRequest is returned when there's an issue with building the request.
	ErrRequest = fmt.Errorf("%w: request", Err)

	// ErrConfig is returned when there's a configuration related issue.
	ErrConfig = fmt.Errorf("%w: config", Err)

	// ErrResponse is the parent of all response API related errors.
	ErrResponse = fmt.Errorf("%w", Err)

	// ErrResourceNotFound is a parent error of all resource-specific not found
	// errors.
	ErrResourceNotFound = fmt.Errorf("%w", ErrNotFound)

	// ErrRouteNotFound indicates the API endpoint called does not exist. This
	// generally will mean go-katapult is very old and needs to be updated.
	ErrRouteNotFound = fmt.Errorf("%w: route_not_found", ErrResponse)

	// ErrUnexpectedResponse is returned if the response body did not contain
	// expected data.
	ErrUnexpectedResponse = fmt.Errorf("%w: unexpected_response", ErrResponse)

	// ErrUnknown is returned if the response error could not be understood.
	ErrUnknown = fmt.Errorf("%w: unknown_error", ErrResponse)
)

// HTTP status-based errors. These may be returned directly, or act as the
// parent error for a more specific error.
var (
	ErrBadRequest          = fmt.Errorf("%w: bad_request", ErrResponse)
	ErrUnauthorized        = fmt.Errorf("%w: unauthorized", ErrResponse)
	ErrForbidden           = fmt.Errorf("%w", ErrUnauthorized)
	ErrNotFound            = fmt.Errorf("%w: not_found", ErrResponse)
	ErrNotAcceptable       = fmt.Errorf("%w: not_acceptable", ErrResponse)
	ErrConflict            = fmt.Errorf("%w: conflict", ErrResponse)
	ErrUnprocessableEntity = fmt.Errorf("%w: unprocessable_entity", ErrResponse)
	ErrTooManyRequests     = fmt.Errorf("%w: too_many_requests", ErrResponse)

	ErrInternalServerError = fmt.Errorf(
		"%w: internal_server_error", ErrResponse,
	)
	ErrBadGateway         = fmt.Errorf("%w: bad_gateway", ErrResponse)
	ErrServiceUnavailable = fmt.Errorf("%w: service_unavailable", ErrResponse)
	ErrGatewayTimeout     = fmt.Errorf("%w: gateway_timeout", ErrResponse)
)

type responseErrorBody struct {
	Error *ResponseError `json:"error,omitempty"`
}

type ResponseError struct {
	parent error

	Code        string          `json:"code,omitempty"`
	Description string          `json:"description,omitempty"`
	Detail      json.RawMessage `json:"detail,omitempty"`
}

func NewResponseError(
	httpStatus int,
	code string,
	description string,
	rawDetail json.RawMessage,
) *ResponseError {
	var parent error
	switch httpStatus {
	case http.StatusBadRequest:
		parent = ErrBadRequest
	case http.StatusUnauthorized:
		parent = ErrUnauthorized
	case http.StatusForbidden:
		parent = ErrForbidden
	case http.StatusNotFound:
		parent = ErrResourceNotFound
	case http.StatusNotAcceptable:
		parent = ErrNotAcceptable
	case http.StatusConflict:
		parent = ErrConflict
	case http.StatusUnprocessableEntity:
		parent = ErrUnprocessableEntity
	case http.StatusTooManyRequests:
		parent = ErrTooManyRequests
	case http.StatusInternalServerError:
		parent = ErrInternalServerError
	case http.StatusBadGateway:
		parent = ErrBadGateway
	case http.StatusServiceUnavailable:
		parent = ErrServiceUnavailable
	case http.StatusGatewayTimeout:
		parent = ErrGatewayTimeout
	default:
		parent = ErrUnknown
	}

	return &ResponseError{
		parent:      parent,
		Code:        code,
		Description: description,
		Detail:      rawDetail,
	}
}

func (s *ResponseError) Error() string {
	out := s.Code
	if s.Description != "" {
		out += ": " + s.Description
	}

	// When RawDetail is not a empty JSON object ("{}" ), we prettify and
	// include it in the error.
	if len(s.Detail) > 2 {
		buf := &bytes.Buffer{}
		_ = json.Indent(buf, s.Detail, "", "  ")
		out += " -- " + buf.String()
	}

	return out
}

func (s *ResponseError) Is(target error) bool {
	return errors.Is(s.parent, target)
}

func (s *ResponseError) Unwrap() error {
	return s.parent
}

// CommonError handles common logic shared between all API-based error types.
type CommonError struct {
	parent error

	Code        string `json:"code,omitempty"`
	Description string `json:"description,omitempty"`
}

func NewCommonError(parent error, code, description string) CommonError {
	return CommonError{
		parent:      parent,
		Code:        code,
		Description: description,
	}
}

func (s *CommonError) BaseError() string {
	if s.parent != nil {
		return s.parent.Error()
	}

	return s.Code
}

func (s *CommonError) Error() string {
	out := s.BaseError()
	if s.Description != "" {
		out += ": " + s.Description
	}

	return out
}

func (s *CommonError) Is(target error) bool {
	return errors.Is(s.parent, target)
}

func (s *CommonError) Unwrap() error {
	return s.parent
}
