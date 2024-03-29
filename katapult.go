package katapult

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

const (
	DefaultUserAgent = "go-katapult"
	DefaultTimeout   = time.Second * 60
)

var DefaultURL = &url.URL{Scheme: "https", Host: "api.katapult.io"}

type Option func(c *Client) error

func WithHTTPClient(hc HTTPClient) Option {
	return func(c *Client) error {
		c.HTTPClient = hc

		return nil
	}
}

func WithUserAgent(ua string) Option {
	return func(c *Client) error {
		c.UserAgent = ua

		return nil
	}
}

func WithBaseURL(u *url.URL) Option {
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

func WithAPIKey(key string) Option {
	return func(c *Client) error {
		c.APIKey = key

		return nil
	}
}

type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type Client struct {
	HTTPClient HTTPClient

	APIKey    string
	UserAgent string
	BaseURL   *url.URL
}

func New(opts ...Option) (*Client, error) {
	// Define default values for client
	c := &Client{
		HTTPClient: &http.Client{Timeout: DefaultTimeout},
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

func (c *Client) Do(
	ctx context.Context,
	request *Request,
	v interface{},
) (*Response, error) {
	contentType, bodyReader, err := request.bodyContent()
	if err != nil {
		return nil, err
	}

	u := c.BaseURL.ResolveReference(request.URL)
	req, err := http.NewRequestWithContext(
		ctx, request.Method, u.String(), bodyReader,
	)
	if err != nil {
		return nil, err
	}

	if len(request.Header) > 0 {
		for k := range request.Header {
			for _, v := range request.Header.Values(k) {
				req.Header.Add(k, v)
			}
		}
	}

	if !request.NoAuth {
		if c.APIKey == "" {
			return nil, fmt.Errorf(
				"%w: no API key available for authenticated request: %s %s",
				ErrRequest, request.Method, request.URL.Path,
			)
		}
		req.Header.Set(
			"Authorization",
			fmt.Sprintf("Bearer %s", c.APIKey),
		)
	}
	req.Header.Set("User-Agent", c.UserAgent)
	req.Header.Set("Accept", "application/json")
	if contentType != "" {
		req.Header.Set("Content-Type", contentType)
	}

	r, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer r.Body.Close()

	resp := NewResponse(r)
	if resp.StatusCode/100 != 2 {
		return c.handleResponseError(resp)
	}

	if v != nil && resp.StatusCode != 204 {
		if w, ok := v.(io.Writer); ok {
			_, err = io.Copy(w, r.Body)
		} else {
			err = json.NewDecoder(resp.Body).Decode(v)
		}
	}

	return resp, err
}

func (c *Client) handleResponseError(resp *Response) (*Response, error) {
	var body responseErrorBody
	err := json.NewDecoder(resp.Body).Decode(&body)
	if err != nil {
		return resp, ErrUnexpectedResponse
	}

	if body.Error == nil || body.Error.Code == "" {
		return resp, ErrUnexpectedResponse
	}
	resp.Error = body.Error
	respErr := NewResponseError(
		resp.StatusCode,
		body.Error.Code,
		body.Error.Description,
		body.Error.Detail,
	)

	return resp, castResponseError(respErr)
}
