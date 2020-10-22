package katapult

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const (
	defaultBaseURL = "https://api.katapult.io/"
	userAgent      = "go-katapult"
	defaultTimeout = time.Second * 60
)

type HTTPClient interface {
	Do(*http.Request) (*http.Response, error)
}

type Client struct {
	client    HTTPClient
	codec     Codec
	common    service
	BaseURL   *url.URL
	UserAgent string

	Certificates  *CertificatesService
	DataCenters   *DataCentersService
	Networks      *NetworksService
	Organizations *OrganizationsService
}

func NewClient(httpClient HTTPClient) *Client {
	if httpClient == nil {
		httpClient = &http.Client{Timeout: defaultTimeout}
	}

	baseURL, _ := url.Parse(defaultBaseURL)

	c := &Client{
		client:    httpClient,
		codec:     &JSONCodec{},
		BaseURL:   baseURL,
		UserAgent: userAgent,
	}
	c.common.client = c

	c.Certificates = NewCertificatesService(&c.common)
	c.DataCenters = NewDataCentersService(&c.common)
	c.Networks = NewNetworksService(&c.common)
	c.Organizations = NewOrganizationsService(&c.common)

	return c
}

func (c *Client) NewRequestWithContext(
	ctx context.Context,
	method string,
	urlStr string,
	body interface{},
) (*http.Request, error) {
	if !strings.HasSuffix(c.BaseURL.Path, "/") {
		return nil, fmt.Errorf(
			"client BaseURL must have a trailing slash, but %q does not",
			c.BaseURL,
		)
	}

	u, err := c.BaseURL.Parse(urlStr)
	if err != nil {
		return nil, err
	}

	var buf io.ReadWriter
	if body != nil {
		buf = &bytes.Buffer{}
		err = c.codec.Encode(body, buf)
		if err != nil {
			return nil, err
		}
	}

	req, err := http.NewRequestWithContext(ctx, method, u.String(), buf)
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", c.UserAgent)
	req.Header.Set("Accept", c.codec.Accept())

	if body != nil {
		req.Header.Set("Content-Type", c.codec.ContentType())
	}

	return req, nil
}

func (c *Client) Do(req *http.Request, v interface{}) (*Response, error) {
	ctx := req.Context()

	r, err := c.client.Do(req)
	if err != nil {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}

		return nil, err
	}
	defer r.Body.Close()

	resp := newResponse(r)
	if resp.StatusCode/100 == 2 {
		if v != nil && resp.StatusCode != 204 {
			if w, ok := v.(io.Writer); ok {
				_, err = io.Copy(w, r.Body)
			} else {
				err = c.codec.Decode(resp.Body, v)
			}
		}

		return resp, err
	}

	return c.handleErrorResponse(resp)
}

func (c *Client) handleErrorResponse(resp *Response) (*Response, error) {
	var body responseErrorBody
	err := c.codec.Decode(resp.Body, &body)
	if err != nil {
		return resp, err
	}

	if body.ErrorInfo == nil {
		return resp, errors.New("unexpected response")
	}
	resp.Error = body.ErrorInfo

	return resp, fmt.Errorf("%s: %s",
		resp.Error.Code,
		resp.Error.Description,
	)
}

type ListOptions struct {
	Page    int `url:"page,omitempty"`
	PerPage int `url:"per_page,omitempty"`
}
