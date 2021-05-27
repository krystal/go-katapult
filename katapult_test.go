package katapult

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/jimeh/undent"
	"github.com/krystal/go-katapult/internal/codec"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var testDefaultBaseURL = &url.URL{Scheme: "https", Host: "api.katapult.io"}

func TestClient_NewRequestWithContext(t *testing.T) {
	type testCtxKey int
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
		name       string
		args       args
		apiKey     string
		baseURL    *url.URL
		codec      *codec.JSON
		wantedAuth string
		wantedBody string
		errStr     string
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
			name: "request with auth",
			args: args{
				ctx: context.WithValue(
					context.Background(), testCtxKey(0), "bar",
				),
				method: "GET",
				url:    &url.URL{Path: "v1/data_centers"},
			},
			apiKey:     "xyzzy",
			wantedAuth: "Bearer xyzzy",
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
			wantedBody: `{"name":"Other Vol"}`,
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

			c := &Client{
				BaseURL: tt.baseURL,
				Codec:   tt.codec,
				APIKey:  tt.apiKey,
			}

			got, err := c.NewRequestWithContext(
				tt.args.ctx, tt.args.method, tt.args.url, tt.args.body,
			)

			if tt.errStr != "" {
				assert.EqualError(t, err, tt.errStr)
			} else {
				wantedURL := tt.baseURL.ResolveReference(tt.args.url)

				assert.NoError(t, err)
				assert.Equal(t, tt.args.ctx, got.Context())
				assert.Equal(t, tt.args.method, got.Method)
				assert.Equal(t, wantedURL.String(), got.URL.String())
				assert.Equal(t,
					c.UserAgent,
					got.Header.Get("User-Agent"),
				)
				assert.Equal(t,
					c.Codec.Accept(),
					got.Header.Get("Accept"),
				)
				assert.Equal(t,
					tt.wantedAuth,
					got.Header.Get("Authorization"),
				)

				if tt.args.body != nil {
					assert.Equal(t,
						c.Codec.ContentType(),
						got.Header.Get("Content-Type"),
					)

					body, err := ioutil.ReadAll(got.Body)
					assert.NoError(t, err)
					assert.Equal(t,
						tt.wantedBody,
						string(bytes.TrimSpace(body)),
					)
				}
			}
		})
	}
}

func TestClient_Do(t *testing.T) {
	type respBody struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	}
	tests := []struct {
		name       string
		ctx        *context.Context
		reqBody    string
		v          interface{}
		want       interface{}
		errStr     string
		errResp    *ResponseError
		respStatus int
		respBody   []byte
		respDelay  time.Duration
	}{
		{
			name:       "struct body with JSON tags",
			v:          &respBody{},
			want:       &respBody{ID: "foo", Name: "bar"},
			respStatus: http.StatusOK,
			respBody:   []byte(`{"id":"foo","name":"bar"}`),
		},
		{
			name:       "io.Writer body",
			v:          &strings.Builder{},
			want:       `{"id":"foo"}`,
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
			want:       "",
			respBody:   []byte(`hi`),
			respStatus: http.StatusNoContent,
		},
		{
			name:       "request times out",
			v:          &respBody{},
			want:       &Response{Response: &http.Response{}},
			errStr:     "Get \"{{baseURL}}/bar\": context deadline exceeded",
			respStatus: http.StatusOK,
			respDelay:  10,
		},
		{
			name: "response is an error without details",
			errStr: "error_without_details: This is an error without " +
				"details",
			errResp: &ResponseError{
				Code:        "error_without_details",
				Description: "This is an error without details",
				Detail:      json.RawMessage("{}"),
			},
			respStatus: http.StatusForbidden,
			respBody: undent.Bytes(`
				{
					"error": {
						"code": "error_without_details",
						"description": "This is an error without details",
						"detail": {}
					}
				}`,
			),
		},
		{
			name: "response is an error with details",
			errStr: "error_with_details: This is an error with " +
				"details -- " +
				"{\n  \"errors\": [\n    \"hello\",\n    \"world\"\n  ]\n}",
			errResp: &ResponseError{
				Code:        "error_with_details",
				Description: "This is an error with details",
				Detail:      json.RawMessage(`{"errors": ["hello","world"]}`),
			},
			respStatus: http.StatusForbidden,
			respBody: undent.Bytes(`
				{
					"error": {
						"code": "error_with_details",
						"description": "This is an error with details",
						"detail": {"errors": ["hello","world"]}
					}
				}`,
			),
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
			m := http.NewServeMux()
			server := httptest.NewServer(m)
			defer server.Close()

			url, err := url.Parse(server.URL)
			if err != nil {
				t.Fatalf("test failed, invalid URL: %s", err.Error())
			}

			c, err := New(
				WithBaseURL(url),
			)
			if err != nil {
				t.Fatalf("failed to setup katapult client: %s", err)
			}

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

			m.HandleFunc("/bar",
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
				c.BaseURL.String()+"/bar",
				strings.NewReader(tt.reqBody),
			)
			require.NoError(t, err)

			got, err := c.Do(req, tt.v)

			if tt.errResp != nil {
				assert.Equal(t, tt.errResp, got.Error)
			}

			if tt.errStr != "" {
				tt.errStr = strings.ReplaceAll(
					tt.errStr, "{{baseURL}}", server.URL,
				)
				assert.EqualError(t, err, tt.errStr)
				if tt.want != nil {
					assert.Equal(t, tt.want, got)
				}
			} else {
				assert.Equal(t, tt.respStatus, got.StatusCode)

				switch v := tt.v.(type) {
				case *strings.Builder:
					assert.Equal(t, tt.want, v.String())
				default:
					assert.Equal(t, tt.want, tt.v)
				}
			}
		})
	}
}

func TestNew(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name          string
		opts          []Opt
		wantAPIKey    string
		wantUserAgent string
		wantBaseURL   *url.URL
		wantTimeout   time.Duration
		wantErr       string
	}{
		{
			name:          "defaults",
			wantAPIKey:    "",
			wantUserAgent: DefaultUserAgent,
			wantBaseURL:   DefaultURL,
			wantTimeout:   DefaultTimeout,
		},
		{
			name: "options specified",
			opts: []Opt{
				WithAPIKey("xyzzy"),
				WithUserAgent("skynet"),
			},
			wantAPIKey:    "xyzzy",
			wantUserAgent: "skynet",
			wantBaseURL:   DefaultURL,
			wantTimeout:   DefaultTimeout,
		},
		{
			name: "err propagates",
			opts: []Opt{
				func(c *Client) error {
					return errors.New("tribbles in the vents")
				},
			},
			wantErr: "tribbles in the vents",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, err := New(tt.opts...)
			if tt.wantErr != "" {
				assert.EqualError(t, err, tt.wantErr)

				return
			}

			assert.NoError(t, err)
			require.NotNil(t, c)

			assert.Equal(t, tt.wantAPIKey, c.APIKey)
			assert.Equal(t, tt.wantUserAgent, c.UserAgent)
			assert.Equal(t, tt.wantBaseURL, c.BaseURL)
			assert.Equal(t, tt.wantTimeout, c.HTTPClient.Timeout)
		})
	}
}

func TestWithTimeout(t *testing.T) {
	c := &Client{HTTPClient: http.DefaultClient}
	timeout := 42 * time.Second
	err := WithTimeout(timeout)(c)
	assert.NoError(t, err)
	assert.Equal(t, timeout, c.HTTPClient.Timeout)
}

func TestWithHTTPClient(t *testing.T) {
	hc := &http.Client{Timeout: 42 * time.Second}
	c := &Client{HTTPClient: http.DefaultClient}
	err := WithHTTPClient(hc)(c)
	assert.NoError(t, err)
	assert.Equal(t, hc, c.HTTPClient)
}

func TestWithUserAgent(t *testing.T) {
	c := &Client{}
	ua := "roger_moore/0.0.7"
	err := WithUserAgent(ua)(c)
	assert.NoError(t, err)
	assert.Equal(t, ua, c.UserAgent)
}

func TestWithAPIKey(t *testing.T) {
	c := &Client{}
	key := "extremely_very_secret_secret"
	err := WithAPIKey(key)(c)
	assert.NoError(t, err)
	assert.Equal(t, key, c.APIKey)
}

func TestWithBaseURL(t *testing.T) {
	tests := []struct {
		name    string
		wantErr string
		baseURL *url.URL
	}{
		{
			name:    "nil causes error",
			baseURL: nil,
			wantErr: "katapult: base URL cannot be nil",
		},
		{
			name:    "empty scheme causes error",
			baseURL: &url.URL{Scheme: "", Host: "google.com"},
			wantErr: "katapult: base URL scheme is empty",
		},
		{
			name:    "empty host causes error",
			baseURL: &url.URL{Scheme: "https", Host: ""},
			wantErr: "katapult: base URL host is empty",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Client{}
			err := WithBaseURL(tt.baseURL)(c)
			if tt.wantErr != "" {
				assert.EqualError(t, err, tt.wantErr)

				return
			}

			assert.NotNil(t, c)
			assert.Equal(t, tt.baseURL, c.BaseURL)
		})
	}
}
