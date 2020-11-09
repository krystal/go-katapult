package katapult

import (
	"bytes"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"path/filepath"
	"reflect"
	"testing"
	"time"

	"github.com/augurysys/timestamp"
	"github.com/krystal/go-katapult/internal/codec"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	updateGoldenFlag = flag.Bool("update-golden", false, "update .golden files")

	fixtureInvalidAPITokenErr = "invalid_api_token: The API token provided " +
		"was not valid (it may not exist or have expired)"
	fixtureInvalidAPITokenResponseError = &ResponseError{
		Code: "invalid_api_token",
		Description: "The API token provided was not valid " +
			"(it may not exist or have expired)",
		Detail: json.RawMessage(`{}`),
	}

	fixturePermissionDeniedErr = "permission_denied: The authenticated " +
		"identity is not permitted to perform this action"
	fixturePermissionDeniedResponseError = &ResponseError{
		Code: "permission_denied",
		Description: "The authenticated identity is not permitted to perform " +
			"this action",
		//nolint:lll
		Detail: json.RawMessage(`{
      "details": "Additional information regarding the reason why permission was denied"
    }`),
	}

	fixtureValidationErrorErr = "validation_error: A validation error " +
		"occurred with the object that was being created/updated/deleted"
	fixtureValidationErrorResponseError = &ResponseError{
		Code: "validation_error",
		Description: "A validation error occurred with the object that was " +
			"being created/updated/deleted",
		Detail: json.RawMessage(`{
      "errors": [
        "Failed reticulating 3-dimensional splines",
        "Failed preparing captive simulators"
      ]
    }`,
		),
	}

	fixtureObjectInTrashErr = "object_in_trash: The object found is in the " +
		"trash and therefore cannot be manipulated through the API. It " +
		"should be restored in order to run this operation."
	fixtureObjectInTrashResponseError = &ResponseError{
		Code: "object_in_trash",
		Description: "The object found is in the trash and therefore cannot " +
			"be manipulated through the API. It should be restored in order " +
			"to run this operation.",
		Detail: json.RawMessage(`{}`),
	}
)

func goldenFile(t *testing.T) string {
	return filepath.Join("testdata", filepath.FromSlash(t.Name())+".golden")
}

func getGolden(t *testing.T) []byte {
	gp := goldenFile(t)
	g, err := ioutil.ReadFile(gp)
	if err != nil {
		t.Fatalf("failed reading .golden: %s", err)
	}

	return g
}

func updateGolden(t *testing.T, got []byte) {
	gp := goldenFile(t)
	dir := filepath.Dir(gp)

	t.Logf("updating .golden file: %s", gp)

	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatalf("failed to update .golden directory: %s", err)
	}

	if err := ioutil.WriteFile(gp, got, 0o644); err != nil { //nolint:gosec
		t.Fatalf("failed to update .golden file: %s", err)
	}
}

func timestampPtr(unixtime int64) *timestamp.Timestamp {
	ts := timestamp.Timestamp(time.Unix(unixtime, 0).UTC())

	return &ts
}

func strictUmarshal(r io.Reader, v interface{}) error {
	d := json.NewDecoder(r)
	d.DisallowUnknownFields()

	return d.Decode(v)
}

func fixture(name string) []byte {
	file := fmt.Sprintf("fixtures/%s.json", name)
	c, err := ioutil.ReadFile(file)
	if err != nil {
		panic(err)
	}

	return c
}

func testJSONMarshaling(t *testing.T, input interface{}) {
	testCustomJSONMarshaling(t, input, nil)
}

func testCustomJSONMarshaling(
	t *testing.T,
	input interface{},
	decoded interface{},
) {
	c := &codec.JSON{}

	buf := &bytes.Buffer{}
	err := c.Encode(input, buf)
	require.NoError(t, err, "encoding failed")

	if *updateGoldenFlag {
		updateGolden(t, buf.Bytes())
	}

	g := getGolden(t)
	assert.Equal(t, string(g), buf.String(),
		"encoding does not match golden")

	if decoded != nil {
		got := reflect.New(reflect.TypeOf(decoded).Elem()).Interface()
		err = c.Decode(bytes.NewBuffer(g), got)
		require.NoError(t, err, "decoding golden failed")
		assert.Equal(t, decoded, got,
			"decoding from golden does not match expected object")
	} else {
		got := reflect.New(reflect.TypeOf(input).Elem()).Interface()
		err = c.Decode(bytes.NewBuffer(g), got)
		require.NoError(t, err, "decoding golden failed")
		assert.Equal(t, input, got,
			"decoding from golden does not match expected object")
	}
}

var testUserAgent = "go-katapult/test"

// prepareTestClient creates a test HTTP server for mock API responses, and
// creates a Katapult client configured to talk to the mock server.
func prepareTestClient() (
	client *Client,
	mux *http.ServeMux,
	serverURL string,
	teardown func(),
) {
	mux = http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(
			os.Stderr,
			"FAIL: Request for unhandled request in test server received:",
		)
		fmt.Fprintln(os.Stderr)
		fmt.Fprintln(os.Stderr, "\t"+r.URL.String())
		w.WriteHeader(http.StatusNotImplemented)
		fmt.Fprint(w, "")
	})

	server := httptest.NewServer(mux)
	url, err := url.Parse(server.URL)
	if err != nil {
		log.Fatalf("test failed, invalid URL: %s", err.Error())
	}

	client, err = NewClient(&Config{
		BaseURL:   url,
		UserAgent: testUserAgent,
	})
	if err != nil {
		log.Fatalf("test client setup failure: %s", err)
	}

	return client, mux, url.String(), server.Close
}

var (
	testDefaultBaseURL   = &url.URL{Scheme: "https", Host: "api.katapult.io"}
	testDefaultUserAgent = "go-katapult"
	testDefaultTimeout   = time.Second * 60
)

func TestNewClient(t *testing.T) {
	type args struct {
		config *Config
	}
	tests := []struct {
		name   string
		args   args
		errStr string
		errIs  []error
	}{
		{
			name: "nil config",
			args: args{
				config: nil,
			},
		},
		{
			name: "empty config",
			args: args{
				config: &Config{},
			},
		},
		{
			name: "config with UserAgent",
			args: args{
				config: &Config{
					UserAgent: "Terraform/0.13.5 terraform-provider-katapult",
				},
			},
		},
		{
			name: "config with BaseURL",
			args: args{
				config: &Config{
					BaseURL: &url.URL{Scheme: "http", Host: "localhost:3001"},
				},
			},
		},
		{
			name: "config with BaseURL where Scheme is empty",
			args: args{
				config: &Config{
					BaseURL: &url.URL{Host: "localhost:3001"},
				},
			},
			errStr: "client: base URL scheme is empty",
			errIs:  []error{Err, ErrClient},
		},
		{
			name: "config with BaseURL where Host is empty",
			args: args{
				config: &Config{
					BaseURL: &url.URL{Scheme: "http"},
				},
			},
			errStr: "client: base URL host is empty",
			errIs:  []error{Err, ErrClient},
		},
		{
			name: "config with HTTPClient",
			args: args{
				config: &Config{
					HTTPClient: &http.Client{Timeout: time.Second * 13},
				},
			},
		},
		{
			name: "config with HTTPTransport",
			args: args{
				config: &Config{
					Transport: &http.Transport{MaxIdleConns: 42},
				},
			},
		},
		{
			name: "config with HTTPClient and HTTPTransport",
			args: args{
				config: &Config{
					HTTPClient: &http.Client{Timeout: time.Second * 13},
					Transport:  &http.Transport{MaxIdleConns: 42},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, err := NewClient(tt.args.config)

			if tt.errStr != "" || len(tt.errIs) > 0 {
				assert.Nil(t, c)

				if tt.errStr != "" {
					assert.EqualError(t, err, tt.errStr)
				}
				for _, e := range tt.errIs {
					assert.True(t, errors.Is(err, e))
				}
			} else {
				assert.IsType(t, new(codec.JSON), c.apiClient.codec)

				if tt.args.config != nil && tt.args.config.UserAgent != "" {
					assert.Equal(t,
						tt.args.config.UserAgent, c.apiClient.UserAgent,
					)
				} else {
					assert.Equal(t, testDefaultUserAgent, c.apiClient.UserAgent)
				}

				if tt.args.config != nil && tt.args.config.BaseURL != nil {
					assert.Equal(t, tt.args.config.BaseURL, c.apiClient.BaseURL)
				} else {
					assert.Equal(t, testDefaultBaseURL, c.apiClient.BaseURL)
				}

				if tt.args.config != nil && tt.args.config.HTTPClient != nil {
					assert.Equal(t,
						tt.args.config.HTTPClient, c.apiClient.httpClient,
					)
					assert.Equal(t,
						tt.args.config.HTTPClient.Timeout,
						c.apiClient.httpClient.Timeout,
					)
				} else {
					assert.IsType(t, new(http.Client), c.apiClient.httpClient)
					assert.Equal(t,
						testDefaultTimeout, c.apiClient.httpClient.Timeout,
					)
				}

				if tt.args.config != nil &&
					tt.args.config.Transport != nil {
					assert.Equal(t,
						tt.args.config.Transport,
						c.apiClient.httpClient.Transport,
					)
				}
			}
		})
	}
}

func TestClient_UserAgent(t *testing.T) {
	tests := []struct {
		name  string
		agent string
	}{
		{
			name: "default user agent",
		},
		{
			name:  "custom user agent",
			agent: "katapult-cli",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, err := NewClient(nil)
			require.NoError(t, err)

			if tt.agent != "" {
				c.apiClient.UserAgent = tt.agent
			}

			got := c.UserAgent()

			if tt.agent != "" {
				assert.Equal(t, tt.agent, got)
			} else {
				assert.Equal(t, testDefaultUserAgent, got)
			}
		})
	}
}

func TestClient_SetUserAgent(t *testing.T) {
	type args struct {
		agent string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "non-empty",
			args: args{
				agent: "katapult-cli",
			},
		},
		{
			name: "empty",
			args: args{
				agent: "",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, err := NewClient(nil)
			require.NoError(t, err)

			c.SetUserAgent(tt.args.agent)

			assert.Equal(t, tt.args.agent, c.apiClient.UserAgent)
		})
	}
}

func TestClient_BaseURL(t *testing.T) {
	tests := []struct {
		name string
		url  *url.URL
	}{
		{
			name: "default",
			url:  nil,
		},
		{
			name: "custom",
			url:  &url.URL{Scheme: "http", Host: "127.0.0.1:3000"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, err := NewClient(nil)
			require.NoError(t, err)

			if tt.url != nil {
				c.apiClient.BaseURL = tt.url
			}

			got := c.BaseURL()

			if tt.url != nil {
				assert.Equal(t, tt.url, got)
			} else {
				assert.Equal(t, testDefaultBaseURL, got)
			}
		})
	}
}

func TestClient_SetBaseURL(t *testing.T) {
	type args struct {
		url *url.URL
	}
	tests := []struct {
		name   string
		args   args
		errStr string
		errIs  []error
	}{
		{
			name: "non-empty",
			args: args{
				url: &url.URL{Scheme: "http", Host: "127.0.0.1:3000"},
			},
		},
		{
			name: "empty scheme",
			args: args{
				url: &url.URL{Host: "127.0.0.1:3000"},
			},
			errStr: "client: base URL scheme is empty",
			errIs:  []error{Err, ErrClient},
		},
		{
			name: "empty host",
			args: args{
				url: &url.URL{Scheme: "http"},
			},
			errStr: "client: base URL host is empty",
			errIs:  []error{Err, ErrClient},
		},
		{
			name: "nil",
			args: args{
				url: nil,
			},
			errStr: "client: base URL cannot be nil",
			errIs:  []error{Err, ErrClient},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, err := NewClient(nil)
			require.NoError(t, err)

			err = c.SetBaseURL(tt.args.url)

			if tt.errStr != "" || len(tt.errIs) > 0 {
				assert.Equal(t, testDefaultBaseURL, c.apiClient.BaseURL)

				if tt.errStr != "" {
					assert.EqualError(t, err, tt.errStr)
				}
				for _, e := range tt.errIs {
					assert.True(t, errors.Is(err, e))
				}
			} else {
				if tt.args.url != nil {
					assert.Equal(t, tt.args.url, c.apiClient.BaseURL)
				} else {
					assert.Equal(t, testDefaultBaseURL, c.apiClient.BaseURL)
				}
			}
		})
	}
}

func TestClient_HTTPClient(t *testing.T) {
	tests := []struct {
		name       string
		httpClient *http.Client
	}{
		{
			name: "default",
		},
		{
			name:       "custom",
			httpClient: &http.Client{Timeout: time.Second * 93},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, err := NewClient(nil)
			require.NoError(t, err)

			if tt.httpClient != nil {
				c.apiClient.httpClient = tt.httpClient
			}

			got := c.HTTPClient()

			if tt.httpClient != nil {
				assert.Equal(t, tt.httpClient, got)
			} else {
				assert.IsType(t, new(http.Client), got)
				assert.Equal(t, testDefaultTimeout, got.Timeout)
			}
		})
	}
}

func TestClient_SetHTTPClient(t *testing.T) {
	type args struct {
		httpClient *http.Client
	}
	tests := []struct {
		name   string
		args   args
		errStr string
		errIs  []error
	}{
		{
			name: "custom",
			args: args{
				httpClient: &http.Client{Timeout: time.Second * 83},
			},
		},
		{
			name: "nil",
			args: args{
				httpClient: nil,
			},
			errStr: "client: http client cannot be nil",
			errIs:  []error{Err, ErrClient},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, err := NewClient(nil)
			require.NoError(t, err)

			original := c.apiClient.httpClient
			err = c.SetHTTPClient(tt.args.httpClient)

			if tt.errStr != "" || len(tt.errIs) > 0 {
				assert.Equal(t, original, c.apiClient.httpClient)

				if tt.errStr != "" {
					assert.EqualError(t, err, tt.errStr)
				}
				for _, e := range tt.errIs {
					assert.True(t, errors.Is(err, e))
				}
			} else {
				assert.Equal(t, tt.args.httpClient, c.apiClient.httpClient)
			}
		})
	}
}

func TestClient_Transport(t *testing.T) {
	tests := []struct {
		name          string
		httpClient    *http.Client
		httpTransport http.RoundTripper
		expected      http.RoundTripper
	}{
		{
			name:          "default",
			httpClient:    &http.Client{},
			httpTransport: nil,
			expected:      nil,
		},
		{
			name:          "custom",
			httpClient:    &http.Client{},
			httpTransport: &http.Transport{MaxConnsPerHost: 949},
			expected:      &http.Transport{MaxConnsPerHost: 949},
		},
		{
			name:          "nil http client",
			httpClient:    nil,
			httpTransport: &http.Transport{MaxConnsPerHost: 949},
			expected:      nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, err := NewClient(nil)
			require.NoError(t, err)

			c.apiClient.httpClient = tt.httpClient

			if tt.httpTransport != nil && c.apiClient.httpClient != nil {
				c.apiClient.httpClient.Transport = tt.httpTransport
			}

			got := c.Transport()

			assert.Equal(t, tt.expected, got)
		})
	}
}

func TestClient_SetTransport(t *testing.T) {
	type args struct {
		httpTransport http.RoundTripper
	}
	tests := []struct {
		name   string
		args   args
		errStr string
		errIs  []error
	}{
		{
			name: "custom",
			args: args{
				httpTransport: &http.Transport{MaxIdleConnsPerHost: 9438},
			},
		},
		{
			name: "nil",
			args: args{
				httpTransport: nil,
			},
			errStr: "client: http transport cannot be nil",
			errIs:  []error{Err, ErrClient},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, err := NewClient(nil)
			require.NoError(t, err)

			original := c.apiClient.httpClient.Transport
			err = c.SetTransport(tt.args.httpTransport)

			if tt.errStr != "" || len(tt.errIs) > 0 {
				assert.Equal(t, original, c.apiClient.httpClient.Transport)

				if tt.errStr != "" {
					assert.EqualError(t, err, tt.errStr)
				}
				for _, e := range tt.errIs {
					assert.True(t, errors.Is(err, e))
				}
			} else {
				assert.Equal(t,
					tt.args.httpTransport, c.apiClient.httpClient.Transport,
				)
			}
		})
	}
}
