// Package testclient contains a fake Client used for testing. Its method
// signature matches that of katapult.Client, allowing you to swap it out if
// you're using a interface to target katapult.Client.
package testclient

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/krystal/go-katapult"
	"github.com/mitchellh/copystructure"
)

var Err = errors.New("testclient")

// Client is a fake client intended for use in tests. It has the same method
// signature as *katapult.Client allowing it satisfy interfaces which describe
// *katapult.Client.
type Client struct {
	// Calls is a of calls individual calls to the Do method with the arguments
	// it was called with. Useful for asserting that expected calls with correct
	// arguments were performed.
	Calls []DoCall

	// Responders is a list of Responder functions which can choose to respond
	// to any specific Do call.
	Responders []Responder

	// Convenience fields populated by the very first call to Do. To access
	// subsequent call arguments, use the Calls field.
	Ctx     context.Context
	Request *katapult.Request
	V       interface{}
}

// New creates a new *Client for testing purposes, with the given resp, err, and
// v added to a default responder which will respond to any call to Do.
func New(
	resp *katapult.Response,
	err error,
	v interface{},
) *Client {
	return &Client{
		Responders: []Responder{
			NewAnyResponder(resp, err, v),
		},
	}
}

// Do pretends to perform a request, instead delegating to it's list of
// responders, going through them in reverse order until one indicates it's a
// match, and then uses the return values from the responder, as the return
// values and struct modification source for the call to Do.
func (s *Client) Do(
	ctx context.Context,
	req *katapult.Request,
	v interface{},
) (*katapult.Response, error) {
	call, err := newDoCall(ctx, req, v)
	if err != nil {
		return nil, fmt.Errorf("%w: %s", Err, err.Error())
	}

	// if it's the first call to Do, assign the convenience argument fields.
	if len(s.Calls) == 0 {
		s.Ctx = call.Ctx
		s.Request = call.Request
		s.V = call.V
	}

	s.Calls = append(s.Calls, *call)

	// Iterate through responders in reverse order, allowing latest defined to
	// take precedence.
	for i := len(s.Responders) - 1; i >= 0; i-- {
		f := s.Responders[i]
		ok, retResp, retErr, retV := f(s, ctx, req, v)
		if !ok {
			continue
		}
		err := s.marshalTransfer(retV, v)
		if err != nil {
			return nil, fmt.Errorf("%w: %s", Err, err.Error())
		}
		retResp = s.ensureResponse(retResp)

		return retResp, retErr
	}

	return &katapult.Response{
		Response: &http.Response{StatusCode: http.StatusNotFound},
	}, nil
}

// marshalTransfer transfers the value from source to target by going through
// json marshal/unmarshal. This avoids the need for any type of reflection, and
// also ensures the struct in use can actually be de-serialized from JSON.
func (s *Client) marshalTransfer(source, target interface{}) error {
	b, err := json.Marshal(source)
	if err != nil {
		return err
	}

	return json.Unmarshal(b, target)
}

// ensureResponse will always return a *katapult.Response with at least a 200 OK
// StatusCode, even if passed nil.
func (s *Client) ensureResponse(
	resp *katapult.Response,
) *katapult.Response {
	// Ensure response is not nil.
	if resp == nil {
		return nil
	}
	// Populate embedded http.Response.
	if resp.Response == nil {
		resp.Response = &http.Response{}
	}
	// Assume 200 OK if status code is empty.
	if resp.Response.StatusCode == 0 {
		resp.Response.StatusCode = http.StatusOK
	}

	return resp
}

// DoCall represents a single call to Do, containing all arguments passed to it.
type DoCall struct {
	Ctx     context.Context
	Request *katapult.Request
	V       interface{}
}

// newDoCall creates a new call object, duplicating the given req and v
// arguments. This ensures we can later compare the state of given arguments
// based on their state at the time of call.
func newDoCall(
	ctx context.Context,
	req *katapult.Request,
	v interface{},
) (*DoCall, error) {
	copy, err := copystructure.Copy(&DoCall{
		Request: req,
		V:       v,
	})
	if err != nil {
		return nil, err
	}

	call := copy.(*DoCall)
	call.Ctx = ctx

	return call, nil
}

// Responder is a function that returns fake return values for calls to Do. This
// allows custom functions to inspect the given context, request, and v
// interface to determine if the responder in question should respond, and with
// what.
//
// The first return value indicates if the responder matched or not. To move on
// to the next responder, this should should be false.
type Responder func(
	tc *Client,
	ctx context.Context,
	req *katapult.Request,
	v interface{},
) (ok bool, resp *katapult.Response, err error, retV interface{})

// NewAnyResponder returns a Responder function which returns the given
// arguments for call to (*Client) Do.
func NewAnyResponder(
	resp *katapult.Response,
	err error,
	v interface{},
) Responder {
	return func(
		_ *Client,
		_ context.Context,
		_ *katapult.Request,
		_ interface{},
	) (bool, *katapult.Response, error, interface{}) {
		return true, resp, err, v
	}
}
