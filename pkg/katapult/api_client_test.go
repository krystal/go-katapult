package katapult

import (
	"bytes"
	"context"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/krystal/go-katapult/internal/codec"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type testCtxKey int

type testResponseBody struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func Test_apiClient_NewRequestWithContext(t *testing.T) {
	type reqBody struct {
		Name string `json:"name"`
	}
	type args struct {
		ctx    context.Context
		method string
		url    *url.URL
		body   interface{}
	}
	tests := []struct {
		name         string
		args         args
		baseURL      *url.URL
		codec        *codec.JSON
		expectedBody string
		errStr       string
	}{
		{
			name: "request without body",
			args: args{
				ctx: context.WithValue(
					context.Background(), testCtxKey(0), "bar",
				),
				method: "GET",
				url:    &url.URL{Path: "v1/data_centers"},
			},
		},
		{
			name: "request with body",
			args: args{
				ctx: context.WithValue(
					context.Background(), testCtxKey(2), "bye",
				),
				method: "PATCH",
				url: &url.URL{
					Path: "v1/file_storage_volumes/fsv_SOIPKzqLkyPan28",
				},
				body: &reqBody{Name: "Other Vol"},
			},
			expectedBody: `{"name":"Other Vol"}`,
		},
		{
			name: "nil context",
			args: args{
				ctx:    nil,
				method: "GET",
				url:    &url.URL{Path: "bbq"},
			},
			errStr: "net/http: nil Context",
		},
		{
			name: "invalid method",
			args: args{
				ctx:    context.Background(),
				method: "foo bar",
				url:    &url.URL{Path: "bbq"},
			},
			errStr: "net/http: invalid method \"foo bar\"",
		},
		{
			name: "invalid body",
			args: args{
				ctx:    context.Background(),
				method: "GET",
				url:    &url.URL{Path: "bbq"},
				body:   make(chan int),
			},
			errStr: "json: unsupported type: chan int",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.baseURL == nil {
				tt.baseURL = testDefaultBaseURL
			}
			if tt.codec == nil {
				tt.codec = &codec.JSON{}
			}

			c := &apiClient{BaseURL: tt.baseURL, codec: tt.codec}

			got, err := c.NewRequestWithContext(
				tt.args.ctx, tt.args.method, tt.args.url, tt.args.body,
			)

			if tt.errStr != "" {
				assert.EqualError(t, err, tt.errStr)
			} else {
				expectedURL := tt.baseURL.ResolveReference(tt.args.url)

				assert.NoError(t, err)
				assert.Equal(t, tt.args.ctx, got.Context())
				assert.Equal(t, tt.args.method, got.Method)
				assert.Equal(t, expectedURL.String(), got.URL.String())
				assert.Equal(t,
					c.UserAgent,
					got.Header.Get("User-Agent"),
				)
				assert.Equal(t,
					c.codec.Accept(),
					got.Header.Get("Accept"),
				)

				if tt.args.body != nil {
					assert.Equal(t,
						c.codec.ContentType(),
						got.Header.Get("Content-Type"),
					)

					body, err := ioutil.ReadAll(got.Body)
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

func Test_apiClient_Do(t *testing.T) {
	tests := []struct {
		name       string
		ctx        *context.Context
		reqBody    string
		v          interface{}
		expected   interface{}
		errStr     string
		errResp    *ResponseError
		respStatus int
		respBody   []byte
		respDelay  time.Duration
	}{
		{
			name:       "struct body with JSON tags",
			v:          &testResponseBody{},
			expected:   &testResponseBody{ID: "foo", Name: "bar"},
			respStatus: http.StatusOK,
			respBody:   []byte(`{"id":"foo","name":"bar"}`),
		},
		{
			name:       "io.Writer body",
			v:          &strings.Builder{},
			expected:   `{"id":"foo"}`,
			respStatus: http.StatusOK,
			respBody:   []byte(`{"id":"foo"}`),
		},
		{
			name:       "request body is submitted to the remote server",
			reqBody:    `hello world`,
			respStatus: http.StatusOK,
			respBody:   []byte(`hi`),
		},
		{
			name:       "response body is ignored for HTTP 204 responses",
			v:          &strings.Builder{},
			expected:   "",
			respBody:   []byte(`hi`),
			respStatus: http.StatusNoContent,
		},
		{
			name:       "request times out",
			errStr:     "Get \"{{baseURL}}/bar\": context deadline exceeded",
			respStatus: http.StatusOK,
			respDelay:  10,
		},
		{
			name:       "response is an error",
			errStr:     fixtureInvalidAPITokenErr,
			errResp:    fixtureInvalidAPITokenResponseError,
			respStatus: http.StatusForbidden,
			respBody:   fixture("invalid_api_token_error"),
		},
		{
			name:       "response is an error of invalid JSON",
			errStr:     "unexpected EOF",
			respStatus: http.StatusForbidden,
			respBody:   []byte(`{"error":{`),
		},
		{
			name:       "response is an error without error info",
			errStr:     "unexpected response",
			respStatus: http.StatusForbidden,
			respBody:   []byte(`{"hello":"world"}`),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, mux, baseURL, teardown := prepareTestClient()
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
					_, _ = w.Write(tt.respBody)
				},
			)

			req, err := http.NewRequestWithContext(
				ctx,
				method,
				c.apiClient.BaseURL.String()+"/bar",
				strings.NewReader(tt.reqBody),
			)
			require.NoError(t, err)

			got, err := c.apiClient.Do(req, tt.v)

			if tt.errResp != nil {
				assert.Equal(t, tt.errResp, got.Error)
			}

			if tt.errStr != "" {
				tt.errStr = strings.ReplaceAll(
					tt.errStr, "{{baseURL}}", baseURL,
				)
				assert.EqualError(t, err, tt.errStr)
			} else {
				assert.Equal(t, tt.respStatus, got.StatusCode)

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
