package katapult

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"testing"
	"time"

	"github.com/krystal/go-katapult/internal/codec"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	fixtureInvalidAPITokenErr = "invalid_api_token: The API token provided " +
		"was not valid (it may not exist or have expired)"
	fixtureInvalidAPITokenResponseError = &ResponseError{
		Code: "invalid_api_token",
		Description: "The API token provided was not valid " +
			"(it may not exist or have expired)",
		Detail: json.RawMessage(`{}`),
	}

	//nolint:lll
	fixturePermissionDeniedErr = "permission_denied: The authenticated " +
		"identity is not permitted to perform this action -- " +
		"{\n  \"details\": \"Additional information regarding the reason why permission was denied\"\n}"
	fixturePermissionDeniedResponseError = &ResponseError{
		Code: "permission_denied",
		Description: "The authenticated identity is not permitted to perform " +
			"this action",
		//nolint:lll
		Detail: json.RawMessage(`{
      "details": "Additional information regarding the reason why permission was denied"
    }`),
	}

	//nolint:lll
	fixtureValidationErrorErr = "validation_error: A validation error " +
		"occurred with the object that was being created/updated/deleted -- " +
		"{\n  \"errors\": [\n    \"Failed reticulating 3-dimensional splines\",\n    \"Failed preparing captive simulators\"\n  ]\n}"
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

	fixtureInvalidArgumentErr = "invalid_argument: The 'X' argument " +
		"is invalid"
	fixtureInvalidArgumentResponseError = &ResponseError{
		Code:        "invalid_argument",
		Description: "The 'X' argument is invalid",
		Detail:      json.RawMessage(`{}`),
	}
)

//
// Helpers
//

func assertFieldSpec(t *testing.T, r *http.Request, spec string) {
	assert.Equal(t, spec, r.Header.Get("X-Field-Spec"))
}

func assertEmptyFieldSpec(t *testing.T, r *http.Request) {
	assertFieldSpec(t, r, "")
}

func assertCustomAuthorization(t *testing.T, r *http.Request, apiKey string) {
	assert.Equal(t,
		fmt.Sprintf("Bearer %s", apiKey), r.Header.Get("Authorization"),
	)
}

func assertAuthorization(t *testing.T, r *http.Request) {
	assertCustomAuthorization(t, r, testAPIKey)
}

var (
	testAPIKey    = "9d7831d8-03f1-4b4c-a1c3-97272ddefe6a"
	testUserAgent = "go-katapult/test"

	testDefaultBaseURL   = &url.URL{Scheme: "https", Host: "api.katapult.io"}
	testDefaultUserAgent = "go-katapult"
	testDefaultTimeout   = time.Second * 60
)

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
		fmt.Fprintf(os.Stderr, "\t%s %s\n\n", r.Method, r.URL.String())

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
		APIKey:    testAPIKey,
		UserAgent: testUserAgent,
	})
	if err != nil {
		log.Fatalf("test client setup failure: %s", err)
	}

	return client, mux, url.String(), server.Close
}

//
// Tests
//

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
			name: "config with APIKey",
			args: args{
				config: &Config{
					APIKey: "016e5c2d-6c21-41e5-a08c-c0a87724fd51",
				},
			},
		},
		{
			name: "config with APIKey",
			args: args{
				config: &Config{
					APIKey: "016e5c2d-6c21-41e5-a08c-c0a87724fd51",
				},
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

				if tt.args.config != nil && tt.args.config.APIKey != "" {
					assert.Equal(t, tt.args.config.APIKey, c.apiClient.APIKey)
				} else {
					assert.Equal(t, "", c.apiClient.APIKey)
				}

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
			}
		})
	}
}

func TestClient_APIKey(t *testing.T) {
	tests := []struct {
		name   string
		apiKey string
	}{
		{
			name:   "non-empty",
			apiKey: "ad8311d3-0e8d-464d-9d1d-c4b12440ebbd",
		},
		{
			name:   "empty",
			apiKey: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, err := NewClient(nil)
			require.NoError(t, err)

			c.apiClient.APIKey = tt.apiKey

			got := c.APIKey()

			assert.Equal(t, tt.apiKey, got)
		})
	}
}

func TestClient_SetAPIKey(t *testing.T) {
	type args struct {
		apiKey string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "non-empty",
			args: args{
				apiKey: "0d297da8-5235-4348-87a0-887be660390b",
			},
		},
		{
			name: "empty",
			args: args{
				apiKey: "",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, err := NewClient(nil)
			require.NoError(t, err)

			c.SetAPIKey(tt.args.apiKey)

			if tt.args.apiKey != "" {
				assert.Equal(t, tt.args.apiKey, c.apiClient.APIKey)
			} else {
				assert.Equal(t, "", c.apiClient.APIKey)
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
