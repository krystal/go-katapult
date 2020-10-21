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
	defaultBaseURL = "https://api.katapult.io/core/"
	apiVersion     = "v1"
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

	DataCenters   *DataCentersService
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
	c.common.apiVersion = apiVersion

	c.DataCenters = &DataCentersService{&c.common}
	c.Organizations = &OrganizationsService{&c.common}

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

func (c *Client) Do(
	req *http.Request,
	v interface{},
) (*Response, error) {
	ctx := req.Context()

	resp, err := c.client.Do(req)
	if err != nil {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}

		return nil, err
	}
	defer resp.Body.Close()

	response := newResponse(resp)
	if resp.StatusCode/100 == 2 {
		if v != nil && resp.StatusCode != 204 {
			if w, ok := v.(io.Writer); ok {
				_, err = io.Copy(w, resp.Body)
			} else {
				err = c.codec.Decode(resp.Body, v)
			}
		}

		return response, err
	}

	responseBody := &ErrorResponseBody{}
	err = c.codec.Decode(resp.Body, responseBody)
	if err != nil {
		return response, err
	}

	response.Error = responseBody.Error
	if response.Error == nil {
		return response, errors.New("unexpected response")
	}

	return response, fmt.Errorf("%s: %s",
		response.Error.Code,
		response.Error.Description,
	)
}
