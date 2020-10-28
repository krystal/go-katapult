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

	if err := os.MkdirAll(dir, 0755); err != nil {
		t.Fatalf("failed to update .golden directory: %s", err)
	}

	if err := ioutil.WriteFile(gp, got, 0644); err != nil { //nolint:gosec
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

func testJSONMarshaling(t *testing.T, v interface{}) {
	c := &codec.JSON{}

	buf := &bytes.Buffer{}
	err := c.Encode(v, buf)
	require.NoError(t, err, "encoding failed")

	if *updateGoldenFlag {
		updateGolden(t, buf.Bytes())
	}

	g := getGolden(t)

	assert.NoError(t, err)
	assert.Equal(t, buf.Bytes(), g,
		"encoding does not match golden")

	got := reflect.New(reflect.TypeOf(v).Elem()).Interface()
	err = c.Decode(bytes.NewBuffer(g), got)
	require.NoError(t, err, "decoding golden failed")
	assert.Equal(t, v, got)
}

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

	client = NewClient(nil)
	client.SetBaseURL(url)

	return client, mux, url.String(), server.Close
}

type customTestHTTPClient struct{}

func (s *customTestHTTPClient) Do(*http.Request) (*http.Response, error) {
	return nil, errors.New("nope")
}

var (
	testDefaultBaseURL   = &url.URL{Scheme: "https", Host: "api.katapult.io"}
	testDefaultUserAgent = "go-katapult"
)

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

			assert.Equal(t, testDefaultBaseURL, c.apiClient.BaseURL)
			assert.Equal(t, testDefaultUserAgent, c.apiClient.UserAgent)
			assert.IsType(t, new(codec.JSON), c.apiClient.codec)

			if tt.httpClient == nil {
				assert.Implements(t, new(HTTPClient), c.apiClient.httpClient)
			} else {
				assert.Equal(t, tt.httpClient, c.apiClient.httpClient)
			}
		})
	}
}

func TestClient_BaseURL(t *testing.T) {
	tests := []struct {
		name string
		url  *url.URL
	}{
		{name: "default base URL"},
		{
			name: "custom base URL",
			url:  &url.URL{Scheme: "http", Host: "127.0.0.1:3000"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := NewClient(nil)

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
	tests := []struct {
		name string
		url  *url.URL
	}{
		{
			name: "custom base URL",
			url:  &url.URL{Scheme: "http", Host: "127.0.0.1:3000"},
		},
		{
			name: "nil",
			url:  nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := NewClient(nil)

			c.SetBaseURL(tt.url)

			if tt.url != nil {
				assert.Equal(t, tt.url, c.apiClient.BaseURL)
			} else {
				assert.Equal(t, testDefaultBaseURL, c.apiClient.BaseURL)
			}
		})
	}
}

func TestClient_UserAgent(t *testing.T) {
	tests := []struct {
		name  string
		agent string
	}{
		{name: "default user agent"},
		{
			name:  "custom user agent",
			agent: "katapult-cli",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := NewClient(nil)

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
	tests := []struct {
		name  string
		agent string
	}{
		{
			name:  "custom user agent",
			agent: "katapult-cli",
		},
		{
			name:  "empty user agent",
			agent: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := NewClient(nil)

			c.SetUserAgent(tt.agent)

			if tt.agent != "" {
				assert.Equal(t, tt.agent, c.apiClient.UserAgent)
			} else {
				assert.Equal(t, testDefaultUserAgent, c.apiClient.UserAgent)
			}
		})
	}
}
