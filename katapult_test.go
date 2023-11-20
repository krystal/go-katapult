package katapult

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/jimeh/undent"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type testErrorWriter struct {
	err error
}

func (s *testErrorWriter) Write(_ []byte) (int, error) {
	return 0, s.err
}

type testErrorHTTPClient struct {
	err error
}

func (s *testErrorHTTPClient) Do(_ *http.Request) (*http.Response, error) {
	return nil, s.err
}

//nolint:gocyclo
func TestClient_Do(t *testing.T) {
	testAPIKey := "7b6eb137-2ce3-4959-9b81-d7aca1428fe1" //nolint:gosec
	type testCtxKey int

	type reqBody struct {
		Hello string `json:"hello,omitempty"`
		World string `json:"world,omitempty"`
	}
	type respBody struct {
		Foo string `json:"foo,omitempty"`
		Bar string `json:"bar,omitempty"`
	}

	type wantReq struct {
		method string
		url    *url.URL
		noAuth bool
		header http.Header
		body   string
	}
	type resp struct {
		status  int
		header  http.Header
		body    string
		delay   time.Duration
		timeout time.Duration
	}
	type wantResp struct {
		status     int
		header     http.Header
		pagination *Pagination
		error      *ResponseError
	}

	type fields struct {
		HTTPClient HTTPClient
		APIKey     *string
		UserAgent  *string
	}
	type args struct {
		ctx     context.Context
		request *Request
		v       interface{}
	}
	tests := []struct {
		name     string
		fields   fields
		args     args
		resp     *resp
		want     interface{}
		wantResp *wantResp
		wantReq  *wantReq
		wantErr  string
	}{
		{
			name: "request without body",
			args: args{
				ctx: context.WithValue(
					context.Background(), testCtxKey(0), "bar",
				),
				request: &Request{
					Method: "GET",
					URL:    &url.URL{Path: "/core/v1/data_centers"},
				},
				v: &respBody{},
			},
			resp: &resp{
				body: `{"foo":"foz","bar":"baz"}`,
			},
			want: &respBody{Foo: "foz", Bar: "baz"},
			wantReq: &wantReq{
				method: "GET",
				url:    &url.URL{Path: "/core/v1/data_centers"},
				header: http.Header{
					"Accept":         []string{"application/json"},
					"Authorization":  []string{"Bearer " + testAPIKey},
					"Content-Length": []string(nil),
					"User-Agent":     []string{DefaultUserAgent},
				},
			},
			wantResp: &wantResp{
				status: http.StatusOK,
			},
		},
		{
			name: "request with struct body",
			args: args{
				ctx: context.WithValue(
					context.Background(), testCtxKey(0), "bar",
				),
				request: &Request{
					Method: "POST",
					URL:    &url.URL{Path: "/core/v1/load_balancers"},
					Body:   &reqBody{Hello: "hi", World: "globe"},
				},
				v: &respBody{},
			},
			resp: &resp{
				body: `{"foo":"foz","bar":"baz"}`,
			},
			want: &respBody{Foo: "foz", Bar: "baz"},
			wantReq: &wantReq{
				method: "POST",
				url:    &url.URL{Path: "/core/v1/load_balancers"},
				header: http.Header{
					"Content-Length": []string{"31"},
					"Content-Type":   []string{"application/json"},
				},
				body: `{"hello":"hi","world":"globe"}` + "\n",
			},
			wantResp: &wantResp{
				status: http.StatusOK,
			},
		},
		{
			name: "request with struct body containing HTML in field values",
			args: args{
				ctx: context.WithValue(
					context.Background(), testCtxKey(0), "bar",
				),
				request: &Request{
					Method: "POST",
					URL:    &url.URL{Path: "/core/v1/load_balancers"},
					Body:   &reqBody{Hello: "<b>hi</b>", World: "<i>globe</i>"},
				},
				v: &respBody{},
			},
			resp: &resp{
				body: `{"foo":"foz","bar":"baz"}`,
			},
			want: &respBody{Foo: "foz", Bar: "baz"},
			wantReq: &wantReq{
				method: "POST",
				url:    &url.URL{Path: "/core/v1/load_balancers"},
				header: http.Header{
					"Content-Length": []string{"45"},
					"Content-Type":   []string{"application/json"},
				},
				body: `{"hello":"<b>hi</b>","world":"<i>globe</i>"}` + "\n",
			},
			wantResp: &wantResp{
				status: http.StatusOK,
			},
		},
		{
			name: "request with custom io.Reader body",
			args: args{
				ctx: context.WithValue(
					context.Background(), testCtxKey(0), "bar",
				),
				request: &Request{
					Method:      "POST",
					URL:         &url.URL{Path: "/core/v1/load_balancers"},
					ContentType: "text/csv",
					Body:        bytes.NewBufferString("foo,bar\nyes,no"),
				},
				v: &respBody{},
			},
			resp: &resp{
				body: `{"foo":"sey","bar":"on"}`,
			},
			want: &respBody{Foo: "sey", Bar: "on"},
			wantReq: &wantReq{
				method: "POST",
				url:    &url.URL{Path: "/core/v1/load_balancers"},
				header: http.Header{
					"Content-Length": []string{"14"},
					"Content-Type":   []string{"text/csv"},
				},
				body: "foo,bar\nyes,no",
			},
			wantResp: &wantResp{
				status: http.StatusOK,
			},
		},
		{
			name: "request with body of invalid type for json marshaling",
			args: args{
				ctx: context.WithValue(
					context.Background(), testCtxKey(0), "bar",
				),
				request: &Request{
					Method: "POST",
					URL:    &url.URL{Path: "/core/v1/load_balancers"},
					Body:   make(chan int),
				},
				v: &respBody{},
			},
			wantErr: "json: unsupported type: chan int",
		},
		{
			name: "request with custom headers",
			args: args{
				ctx: context.WithValue(
					context.Background(), testCtxKey(0), "bar",
				),
				request: &Request{
					Method: "GET",
					URL:    &url.URL{Path: "/core/v1/data_centers"},
					Header: http.Header{
						"Ignore-This": []string{"hello", "world"},
						"X-Client":    []string{"Awesome App"},
					},
				},
				v: &respBody{},
			},
			wantReq: &wantReq{
				method: "GET",
				url:    &url.URL{Path: "/core/v1/data_centers"},
				header: http.Header{
					"Ignore-This": []string{"hello", "world"},
					"X-Client":    []string{"Awesome App"},
				},
			},
			wantResp: &wantResp{
				status: http.StatusOK,
			},
		},
		{
			name: "v is a io.Writer",
			args: args{
				ctx: context.WithValue(
					context.Background(), testCtxKey(0), "bar",
				),
				request: &Request{
					Method: "GET",
					URL:    &url.URL{Path: "/core/v1/stats/csv"},
				},
				v: &bytes.Buffer{},
			},
			resp: &resp{
				body: "foo,bar\nsey,on",
			},
			want: "foo,bar\nsey,on",
			wantReq: &wantReq{
				method: "GET",
				url:    &url.URL{Path: "/core/v1/stats/csv"},
			},
		},
		{
			name: "v is a io.Writer which errors on write",
			args: args{
				ctx: context.WithValue(
					context.Background(), testCtxKey(0), "bar",
				),
				request: &Request{
					Method: "GET",
					URL:    &url.URL{Path: "/core/v1/stats/csv"},
				},
				v: &testErrorWriter{err: errors.New("writer is broken")},
			},
			resp: &resp{
				body: "foo,bar\nsey,on",
			},
			wantReq: &wantReq{
				method: "GET",
				url:    &url.URL{Path: "/core/v1/stats/csv"},
			},
			wantErr: "writer is broken",
		},
		{
			name: "HTTPClient.Do() error",
			fields: fields{
				HTTPClient: &testErrorHTTPClient{
					err: errors.New("HTTP failure"),
				},
			},
			args: args{
				ctx: context.WithValue(
					context.Background(), testCtxKey(0), "bar",
				),
				request: &Request{
					Method: "GET",
					URL:    &url.URL{Path: "/core/v1/stats/csv"},
				},
				v: &respBody{},
			},
			wantErr: "HTTP failure",
		},
		{
			name: "unauthenticated request",
			args: args{
				ctx: context.WithValue(
					context.Background(), testCtxKey(0), "bar",
				),
				request: &Request{
					Method: "GET",
					URL:    &url.URL{Path: "/public/v1/stats"},
					NoAuth: true,
				},
				v: &respBody{},
			},
			resp: &resp{
				body: `{"foo":"foz"}`,
			},
			want: &respBody{Foo: "foz"},
			wantReq: &wantReq{
				method: "GET",
				url:    &url.URL{Path: "/public/v1/stats"},
				noAuth: true,
				header: http.Header{
					"Authorization": []string(nil),
				},
			},
			wantResp: &wantResp{
				status: http.StatusOK,
			},
		},
		{
			name: "authenticated request without API key",
			fields: fields{
				APIKey: strPtr(""),
			},
			args: args{
				ctx: context.WithValue(
					context.Background(), testCtxKey(0), "bar",
				),
				request: &Request{
					Method: "GET",
					URL:    &url.URL{Path: "/core/v1/stats"},
				},
				v: &respBody{},
			},
			wantErr: "katapult: request: no API key available for " +
				"authenticated request: GET /core/v1/stats",
		},
		{
			name: "response has custom headers",
			args: args{
				ctx: context.WithValue(
					context.Background(), testCtxKey(0), "bar",
				),
				request: &Request{
					Method: "GET",
					URL:    &url.URL{Path: "/core/v1/data_centers"},
				},
				v: &respBody{},
			},
			resp: &resp{
				body: `{"foo":"foz","bar":"baz"}`,
				header: http.Header{
					"X-RateLimit-Permitted": []string{"100"},
					"X-RateLimit-Remaining": []string{"99"},
					"X-Hello":               []string{"Hi", "Hey"},
				},
			},
			want: &respBody{Foo: "foz", Bar: "baz"},
			wantReq: &wantReq{
				method: "GET",
				url:    &url.URL{Path: "/core/v1/data_centers"},
			},
			wantResp: &wantResp{
				status: http.StatusOK,
				header: http.Header{
					"X-Ratelimit-Permitted": []string{"101"},
					"X-Ratelimit-Remaining": []string{"99"},
					"X-Hello":               []string{"Hi", "Hey"},
				},
			},
		},
		{
			name: "response is an error without details",
			args: args{
				ctx: context.WithValue(
					context.Background(), testCtxKey(0), "bar",
				),
				request: &Request{
					Method: "GET",
					URL:    &url.URL{Path: "/core/v1/stats"},
				},
				v: &respBody{},
			},
			resp: &resp{
				status: http.StatusForbidden,
				body: undent.String(`
					{
						"error": {
							"code": "error_without_details",
							"description": "This is an error without details",
							"detail": {}
						}
					}`,
				),
			},
			wantReq: &wantReq{
				method: "GET",
				url:    &url.URL{Path: "/core/v1/stats"},
			},
			wantResp: &wantResp{
				error: &ResponseError{
					Code:        "error_without_details",
					Description: "This is an error without details",
					Detail:      json.RawMessage("{}"),
				},
			},
			wantErr: "error_without_details: This is an error without " +
				"details",
		},
		{
			name: "response is an error with details",
			args: args{
				ctx: context.WithValue(
					context.Background(), testCtxKey(0), "bar",
				),
				request: &Request{
					Method: "GET",
					URL:    &url.URL{Path: "/core/v1/stats"},
				},
				v: &respBody{},
			},
			resp: &resp{
				status: http.StatusForbidden,
				body: undent.String(`
					{
						"error": {
							"code": "error_with_details",
							"description": "This is an error with details",
							"detail": {"errors": ["hello","world"]}
						}
					}`,
				),
			},
			wantReq: &wantReq{
				method: "GET",
				url:    &url.URL{Path: "/core/v1/stats"},
			},
			wantResp: &wantResp{
				error: &ResponseError{
					Code:        "error_with_details",
					Description: "This is an error with details",
					Detail: json.RawMessage(
						`{"errors": ["hello","world"]}`,
					),
				},
			},
			wantErr: "error_with_details: This is an error with " +
				"details -- " +
				"{\n  \"errors\": [\n    \"hello\",\n    \"world\"\n  ]\n}",
		},
		{
			name: "response is an error with no error info at all",
			args: args{
				ctx: context.WithValue(
					context.Background(), testCtxKey(0), "bar",
				),
				request: &Request{
					Method: "GET",
					URL:    &url.URL{Path: "/core/v1/stats"},
				},
				v: &respBody{},
			},
			resp: &resp{
				status: http.StatusForbidden,
				body:   `{"code": "something wrong"}`,
			},
			wantReq: &wantReq{
				method: "GET",
				url:    &url.URL{Path: "/core/v1/stats"},
			},
			wantErr: "katapult: unexpected_response",
		},
		{
			name: "response body is invalid JSON",
			args: args{
				ctx: context.WithValue(
					context.Background(), testCtxKey(0), "bar",
				),
				request: &Request{
					Method: "GET",
					URL:    &url.URL{Path: "/core/v1/stats"},
				},
				v: &respBody{},
			},
			resp: &resp{
				body: `{"stats":{`,
			},
			wantReq: &wantReq{
				method: "GET",
				url:    &url.URL{Path: "/core/v1/stats"},
			},
			wantErr: "unexpected EOF",
		},
		{
			name: "response is an error with a invalid JSON body",
			args: args{
				ctx: context.WithValue(
					context.Background(), testCtxKey(0), "bar",
				),
				request: &Request{
					Method: "GET",
					URL:    &url.URL{Path: "/core/v1/stats"},
				},
				v: &respBody{},
			},
			resp: &resp{
				status: http.StatusForbidden,
				body:   `{"error":{`,
			},
			wantReq: &wantReq{
				method: "GET",
				url:    &url.URL{Path: "/core/v1/stats"},
			},
			wantErr: "katapult: unexpected_response",
		},
		{
			name: "context timeout",
			args: args{
				ctx: context.WithValue(
					context.Background(), testCtxKey(0), "bar",
				),
				request: &Request{
					Method: "GET",
					URL:    &url.URL{Path: "/core/v1/stats"},
				},
				v: &respBody{},
			},
			resp: &resp{
				body:    `{"foo":"sey","bar":"on"}`,
				delay:   100 * time.Millisecond,
				timeout: 50 * time.Millisecond,
			},
			wantReq: &wantReq{
				method: "GET",
				url:    &url.URL{Path: "/core/v1/stats"},
			},
			wantErr: "Get \"{{ServerURL}}/core/v1/stats\": " +
				"context deadline exceeded",
		},
		{
			name: "nil context",
			args: args{
				ctx: nil,
				request: &Request{
					Method: "GET",
					URL:    &url.URL{Path: "/core/v1/stats"},
				},
			},
			wantErr: "net/http: nil Context",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mux := http.NewServeMux()
			mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
				assert.FailNowf(
					t, "Unhandled request", "%s %s", r.Method, r.URL.String(),
				)
				w.WriteHeader(http.StatusNotImplemented)
				fmt.Fprint(w, "")
			})
			server := httptest.NewServer(mux)
			defer server.Close()

			u, err := url.Parse(server.URL)
			require.NoError(t, err)

			c, err := New(WithBaseURL(u))
			require.NoError(t, err)
			require.NotNil(t, c)

			if tt.fields.HTTPClient != nil {
				_ = WithHTTPClient(tt.fields.HTTPClient)(c)
			}

			apiKey := testAPIKey
			if tt.fields.APIKey != nil {
				apiKey = *tt.fields.APIKey
			}
			_ = WithAPIKey(apiKey)(c)

			if tt.fields.UserAgent != nil {
				_ = WithUserAgent(*tt.fields.UserAgent)(c)
			}

			if tt.wantReq != nil && tt.wantReq.url != nil {
				mux.HandleFunc(tt.wantReq.url.Path,
					func(w http.ResponseWriter, r *http.Request) {
						assert.Equal(t, tt.wantReq.method, r.Method)
						assert.Equal(t, tt.wantReq.url.Query(), r.URL.Query())
						assert.Equal(t, tt.wantReq.url.Fragment, r.URL.Fragment)

						if tt.wantReq.noAuth {
							assert.Equal(t, "", r.Header.Get("Authorization"))
						} else {
							assert.Equal(t,
								"Bearer "+apiKey, r.Header.Get("Authorization"),
							)
						}

						if len(tt.wantReq.header) > 0 {
							for k := range tt.wantReq.header {
								assert.Equalf(t,
									tt.wantReq.header.Values(k),
									r.Header.Values(k),
									"request header: %s", k,
								)
							}
						}

						receivedReqBody, _ := ioutil.ReadAll(r.Body)
						assert.Equal(t,
							tt.wantReq.body, string(receivedReqBody),
						)

						var status int
						var body string
						if tt.resp != nil {
							status = tt.resp.status
							body = tt.resp.body

							if tt.resp.delay != 0 {
								time.Sleep(tt.resp.delay)
							}
						}

						if status == 0 {
							status = 200
						}

						if tt.wantResp != nil {
							for k := range tt.wantResp.header {
								for _, v := range tt.wantResp.header.Values(k) {
									w.Header().Add(k, v)
								}
							}
						}

						w.WriteHeader(status)

						if body != "" {
							_, _ = w.Write([]byte(body))
						}
					},
				)
			}

			ctx := tt.args.ctx
			if tt.resp != nil && tt.resp.timeout > 0 {
				var cancel context.CancelFunc
				ctx, cancel = context.WithTimeout(ctx, tt.resp.timeout)
				defer cancel()
			}

			resp, err := c.Do(ctx, tt.args.request, tt.args.v)

			if tt.wantErr != "" {
				wantErr := strings.ReplaceAll(
					tt.wantErr, "{{ServerURL}}", server.URL,
				)
				assert.EqualError(t, err, wantErr)
			}

			if tt.wantResp != nil && resp != nil {
				if resp.Response != nil {
					if tt.wantResp.status != 0 {
						assert.Equal(t, tt.wantResp.status, resp.StatusCode)
					}
					if len(tt.wantResp.header) > 0 && resp.Response != nil {
						for k := range tt.wantResp.header {
							assert.Equal(t,
								tt.wantResp.header.Get(k), resp.Header.Get(k),
							)
						}
					}
				}
				if tt.wantResp.pagination != nil {
					assert.Equal(t, tt.wantResp.pagination, resp.Pagination)
				}
				if tt.wantResp.error != nil {
					assert.Equal(t, tt.wantResp.error, resp.Error)
				}
			}

			if tt.want != nil {
				if r, ok := tt.args.v.(io.Reader); ok {
					b, err := ioutil.ReadAll(r)
					require.NoError(t, err)
					assert.Equal(t, tt.want, string(b))
				} else {
					assert.Equal(t, tt.want, tt.args.v)
				}
			}
		})
	}
}

func TestNew(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		opts          []Option
		wantAPIKey    string
		wantUserAgent string
		wantBaseURL   *url.URL
		wantTimeout   time.Duration
		wantErr       string
	}{
		{
			name:          "defaults",
			wantAPIKey:    "",
			wantUserAgent: "go-katapult",
			wantBaseURL:   &url.URL{Scheme: "https", Host: "api.katapult.io"},
			wantTimeout:   time.Second * 60,
		},
		{
			name: "options specified",
			opts: []Option{
				WithAPIKey("xyzzy"),
				WithUserAgent("skynet"),
			},
			wantAPIKey:    "xyzzy",
			wantUserAgent: "skynet",
			wantBaseURL:   &url.URL{Scheme: "https", Host: "api.katapult.io"},
			wantTimeout:   time.Second * 60,
		},
		{
			name: "err propagates",
			opts: []Option{
				func(c *Client) error {
					return errors.New("tribbles in the vents")
				},
			},
			wantErr: "tribbles in the vents",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

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
		})
	}
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
