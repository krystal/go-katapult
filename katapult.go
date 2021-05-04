package katapult

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/krystal/go-katapult/internal/codec"
)

const (
	DefaultUserAgent = "go-katapult"
	DefaultTimeout   = time.Second * 60
)

var (
	DefaultURL = &url.URL{Scheme: "https", Host: "api.katapult.io"}
)

func WithTimeout(t time.Duration) Opt {
	return func(c *Client) error {
		c.HTTPClient.Timeout = t

		return nil
	}
}

func WithUserAgent(ua string) Opt {
	return func(c *Client) error {
		c.UserAgent = ua

		return nil
	}
}

func WithBaseURL(u *url.URL) Opt {
	return func(c *Client) error {
		switch {
		case u == nil:
			return fmt.Errorf("katapult: base URL cannot be nil")
		case u.Scheme == "":
			return fmt.Errorf("katapult: base URL scheme is empty")
		case u.Host == "":
			return fmt.Errorf("katapult: base URL host is empty")
		}

		c.BaseURL = u

		return nil
	}
}

func WithAPIKey(key string) Opt {
	return func(c *Client) error {
		c.APIKey = key

		return nil
	}
}

type Opt func(c *Client) error

func New(opts ...Opt) (*Client, error) {
	// Define default values for client
	c := &Client{
		HTTPClient: &http.Client{Timeout: DefaultTimeout},
		Codec:      &codec.JSON{},
		BaseURL:    DefaultURL,
		UserAgent:  DefaultUserAgent,
	}

	// Apply options to created Client
	for _, o := range opts {
		err := o(c)
		if err != nil {
			return nil, err
		}
	}

	return c, nil
}

type Client struct {
	HTTPClient *http.Client
	Codec      codec.Codec

	APIKey    string
	UserAgent string
	BaseURL   *url.URL
}

// NewRequestWithContext returns a http.Request created for sending to the API.
func (c *Client) NewRequestWithContext(
	ctx context.Context,
	method string,
	u *url.URL,
	body interface{},
) (*http.Request, error) {
	var buf io.ReadWriter
	if body != nil {
		buf = &bytes.Buffer{}
		err := c.Codec.Encode(body, buf)
		if err != nil {
			return nil, err
		}
	}

	u = c.BaseURL.ResolveReference(u)
	req, err := http.NewRequestWithContext(ctx, method, u.String(), buf)
	if err != nil {
		return nil, err
	}

	if c.APIKey != "" {
		req.Header.Set(
			"Authorization",
			fmt.Sprintf("Bearer %s", c.APIKey),
		)
	}

	req.Header.Set("User-Agent", c.UserAgent)
	req.Header.Set("Accept", c.Codec.Accept())

	if body != nil {
		req.Header.Set("Content-Type", c.Codec.ContentType())
	}

	return req, nil
}

// Do executes a request, decoding the response body into argument v.
func (c *Client) Do(req *http.Request, v interface{}) (*Response, error) {
	r, err := c.HTTPClient.Do(req)
	if err != nil {
		return NewResponse(nil), err
	}
	defer r.Body.Close()

	resp := NewResponse(r)
	if resp.StatusCode/100 == 2 {
		if v != nil && resp.StatusCode != 204 {
			if w, ok := v.(io.Writer); ok {
				_, err = io.Copy(w, r.Body)
			} else {
				err = c.Codec.Decode(resp.Body, v)
			}
		}

		return resp, err
	}

	return c.handleResponseError(resp)
}

func (c *Client) handleResponseError(resp *Response) (*Response, error) {
	var body responseErrorBody
	err := c.Codec.Decode(resp.Body, &body)
	if err != nil {
		return resp, err
	}

	if body.ErrorInfo == nil {
		return resp, errors.New("unexpected response")
	}
	resp.Error = body.ErrorInfo

	if len(resp.Error.Detail) > 2 {
		buf := &bytes.Buffer{}
		_ = json.Indent(buf, resp.Error.Detail, "", "  ")

		baseErr := fmt.Errorf("%s: %s",
			resp.Error.Code,
			resp.Error.Description,
		)

		return resp, fmt.Errorf("%w -- %s",
			baseErr,
			buf.String(),
		)
	}

	return resp, fmt.Errorf("%s: %s",
		resp.Error.Code,
		resp.Error.Description,
	)
}
