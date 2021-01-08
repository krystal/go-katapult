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

	"github.com/krystal/go-katapult/internal/codec"
)

type apiClient struct {
	httpClient *http.Client
	codec      codec.Codec

	APIKey    string
	UserAgent string
	BaseURL   *url.URL
}

func (c *apiClient) NewRequestWithContext(
	ctx context.Context,
	method string,
	u *url.URL,
	body interface{},
) (*http.Request, error) {
	var buf io.ReadWriter
	if body != nil {
		buf = &bytes.Buffer{}
		err := c.codec.Encode(body, buf)
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
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.APIKey))
	}

	req.Header.Set("User-Agent", c.UserAgent)
	req.Header.Set("Accept", c.codec.Accept())

	if body != nil {
		req.Header.Set("Content-Type", c.codec.ContentType())
	}

	return req, nil
}

func (c *apiClient) Do(req *http.Request, v interface{}) (*Response, error) {
	r, err := c.httpClient.Do(req)
	if err != nil {
		return newResponse(nil), err
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

	return c.handleResponseError(resp)
}

func (c *apiClient) handleResponseError(resp *Response) (*Response, error) {
	var body responseErrorBody
	err := c.codec.Decode(resp.Body, &body)
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
