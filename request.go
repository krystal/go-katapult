package katapult

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

// Request represents a HTTP request to the Katapult API, it is essentially
// similar to http.Request, but stripped down to the bare essentials, with some
// Katapult-specific attributes added.
type Request struct {
	// Method is the HTTP method to perform when making the request.
	Method string

	// URL is the request URL to perform against Katapult's API. Generally the
	// only fields you need to set are Path and RawQuery, as it will be merged
	// with the Client's BaseURL value through its ResolveReference() method.
	URL *url.URL

	// NoAuth instructs Client not to send Authorization header containing the
	// APIKey. This is useful for public endpoints which do not require/use
	// authentication.
	NoAuth bool

	// Header holds request-specific HTTP headers. Client.Do() will set a number
	// of essential headers itself which cannot be customized through
	// Request.Headers.
	Header http.Header

	// ContentType allows sending a custom request body of any mimetype. If set
	// Body must be a io.Reader. If not set Body must be a object which can be
	// serialized with json.Marshal(). Content-Type header is only sent when
	// Body is not nil.
	ContentType string

	// Body can be any object which can be marshaled to JSON through
	// json.Marshal() when ContentType is not set. If ContentType is set, Body
	// must be a io.Reader, or nil.
	//
	// No validation is performed between Method and Body, making it possible to
	// send a body with a method that does not allow it.
	Body interface{}
}

type RequestOption = func(r *Request)

// RequestSetHeader sets a header on the outgoing request. This replaces any
// headers that are currently specified with that key.
func RequestSetHeader(key, value string) RequestOption {
	return func(r *Request) {
		r.Header.Set(key, value)
	}
}

func NewRequest(
	method string,
	u *url.URL,
	body interface{},
	opts ...RequestOption,
) *Request {
	r := &Request{
		Method: method,
		URL:    u,
		Body:   body,
	}

	// Apply options to created Client
	for _, opt := range opts {
		opt(r)
	}

	return r
}

func (r *Request) bodyContent() (string, io.Reader, error) {
	if r.Body == nil {
		return "", nil, nil
	}

	contentType := r.ContentType
	var body io.Reader

	if contentType == "" {
		contentType = "application/json"
		var buf bytes.Buffer
		enc := json.NewEncoder(&buf)
		enc.SetEscapeHTML(false)
		err := enc.Encode(r.Body)
		if err != nil {
			return "", nil, err
		}
		body = &buf
	} else {
		var ok bool
		body, ok = r.Body.(io.Reader)
		if !ok {
			return "", nil, fmt.Errorf(
				"%w: Body must be a io.Reader when ContentType is set",
				ErrRequest,
			)
		}
	}

	return contentType, body, nil
}
