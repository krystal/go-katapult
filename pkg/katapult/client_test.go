package katapult

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	testDefaultBaseURL = "https://api.katapult.io/core/"
	testAPIVersion     = "v1"
	testUserAgent      = "go-katapult"
)

type customTestHTTPClient struct{}

func (s *customTestHTTPClient) Do(*http.Request) (*http.Response, error) {
	return nil, errors.New("nope")
}

func TestNewClient(t *testing.T) {
	tests := []struct {
		name       string
		httpClient HTTPClient
	}{
		{
			name:       "not given a http.Client",
			httpClient: nil,
		},
		{
			name:       "given a http.Client",
			httpClient: &http.Client{},
		},
		{
			name:       "given a custom HTTPClient",
			httpClient: &customTestHTTPClient{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := NewClient(tt.httpClient)

			assert.Equal(t, testDefaultBaseURL, c.BaseURL.String())
			assert.Equal(t, testUserAgent, c.UserAgent)
			assert.Equal(t, testAPIVersion, c.common.apiVersion)
			assert.IsType(t, new(JSONCodec), c.codec)

			if tt.httpClient == nil {
				assert.Implements(t, new(HTTPClient), c.client)
			} else {
				assert.Equal(t, tt.httpClient, c.client)
			}
		})
	}
}

func TestClient_NewRequestWithContext(t *testing.T) {
	tests := []struct {
		name         string
		ctx          context.Context
		method       string
		baseURL      string
		urlStr       string
		body         interface{}
		expectedBody string
		err          string
	}{
		{
			name: "request without body",
			ctx: context.WithValue(
				context.Background(), testCtxKey(0), "bar",
			),
			method: "GET",
			urlStr: "v1/data_centers",
		},
		{
			name: "request with body",
			ctx: context.WithValue(
				context.Background(), testCtxKey(2), "bye",
			),
			method: "PATCH",
			urlStr: "v1/file_storage_volumes/fsv_SOIPKzqLkyPan28",
			body: struct {
				Name string `json:"name"`
			}{Name: "Other Vol"},
			expectedBody: `{"name":"Other Vol"}`,
		},
		{
			name: "Base URL without trailing slash",
			ctx: context.WithValue(
				context.Background(), testCtxKey(3), "world",
			),
			baseURL: "https://api.katapult.io/core",
			err: `client BaseURL must have a trailing slash, but ` +
				`"https://api.katapult.io/core" does not`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, _, testBaseURL, teardown := setup()
			defer teardown()

			var baseURL *url.URL
			var err error
			if tt.baseURL != "" {
				baseURL, err = url.Parse(tt.baseURL)
				require.NoError(t, err)
				c.BaseURL = baseURL
			} else {
				baseURL, err = url.Parse(testBaseURL)
				require.NoError(t, err)
			}

			req, err := c.NewRequestWithContext(
				tt.ctx, tt.method, tt.urlStr, tt.body,
			)

			if tt.err != "" {
				assert.EqualError(t, err, tt.err)
			} else {
				expectedURL, err := baseURL.Parse(tt.urlStr)
				require.NoError(t, err)

				assert.NoError(t, err)
				assert.Equal(t, tt.ctx, req.Context())
				assert.Equal(t, tt.method, req.Method)
				assert.Equal(t, expectedURL.String(), req.URL.String())

				if tt.body != nil {
					body, err := ioutil.ReadAll(req.Body)
					assert.NoError(t, err)
					assert.Equal(t,
						tt.expectedBody,
						string(bytes.TrimSpace(body)),
					)
				}
			}
		})
	}
}

func TestClient_Do(t *testing.T) {
	tests := []struct {
		name       string
		ctx        *context.Context
		reqBody    string
		v          interface{}
		expected   interface{}
		err        string
		errResp    *ErrorResponse
		respStatus int
		respBody   string
		respDelay  time.Duration
	}{
		{
			name: "body is decoded into v when it is a struct with " +
				"JSON tags",
			v:          &testResponseBody{},
			expected:   &testResponseBody{ID: "foo", Name: "bar"},
			respStatus: http.StatusOK,
			respBody:   `{"id":"foo","name":"bar"}`,
		},
		{
			name:       "body is copied to v when it is a io.Writer",
			v:          &strings.Builder{},
			expected:   `{"id":"foo"}`,
			respStatus: http.StatusOK,
			respBody:   `{"id":"foo"}`,
		},
		{
			name:       "request body is submitted to the remote server",
			reqBody:    `hello world`,
			respStatus: http.StatusOK,
			respBody:   `hi`,
		},
		{
			name:       "response body is ignored when response is HTTP 204",
			v:          &strings.Builder{},
			expected:   "",
			respBody:   `hi`,
			respStatus: http.StatusNoContent,
		},
		{
			name:       "when request times out",
			err:        "context deadline exceeded",
			respStatus: http.StatusOK,
			respDelay:  10,
		},
		{
			name: "response is an error",
			err: "invalid_api_token: The API token provided was not valid " +
				"(it may not exist or have expired)",
			errResp: &ErrorResponse{
				Code: "invalid_api_token",
				Description: "The API token provided was not valid " +
					"(it may not exist or have expired)",
				Detail: json.RawMessage(`{}`),
			},
			respStatus: http.StatusForbidden,
			//nolint:lll
			respBody: `{
  "error": {
    "code": "invalid_api_token",
    "description": "The API token provided was not valid (it may not exist or have expired)",
    "detail": {}
  }
}`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, mux, _, teardown := setup()
			defer teardown()

			method := "GET"
			ctx := context.Background()

			if tt.reqBody != "" {
				method = "POST"
			}
			if tt.ctx != nil {
				ctx = *tt.ctx
			}
			if tt.respDelay != 0 {
				var cancel context.CancelFunc
				ctx, cancel = context.WithTimeout(
					context.Background(),
					(tt.respDelay/2)*time.Millisecond,
				)
				defer cancel()
			}

			mux.HandleFunc("/bar",
				func(w http.ResponseWriter, r *http.Request) {
					assert.Equal(t, method, r.Method)

					receivedReqBody, _ := ioutil.ReadAll(r.Body)
					assert.Equal(t, tt.reqBody, string(receivedReqBody))

					if tt.respDelay != 0 {
						time.Sleep(tt.respDelay * time.Millisecond)
					}

					w.WriteHeader(tt.respStatus)
					fmt.Fprint(w, tt.respBody)
				},
			)

			req, err := http.NewRequestWithContext(
				ctx,
				method,
				c.BaseURL.String()+"bar",
				strings.NewReader(tt.reqBody),
			)
			require.NoError(t, err)

			resp, err := c.Do(req, tt.v)

			if tt.errResp != nil {
				assert.Equal(t, tt.errResp, resp.Error)
			}

			if tt.err != "" {
				assert.EqualError(t, err, tt.err)
			} else {
				assert.Equal(t, tt.respStatus, resp.StatusCode)

				switch v := tt.v.(type) {
				case *strings.Builder:
					assert.Equal(t, tt.expected, v.String())
				default:
					assert.Equal(t, tt.expected, tt.v)
				}
			}
		})
	}
}
