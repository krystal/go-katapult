package katapult

import (
	"fmt"
	"github.com/krystal/go-katapult/internal/codec"
	"net/http"
	"net/url"
	"time"
)

const (
	DefaultUserAgent = "go-katapult"
	DefaultTimeout   = time.Second * 60
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

func WithHTTPClient(h *http.Client) Opt {
	return func(c *Client) error {
		c.HTTPClient = h
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
		BaseURL:    &url.URL{Scheme: "https", Host: "api.katapult.io"},
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
