package katapult

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"time"

	"github.com/augurysys/timestamp"
)

type testCtxKey int

type testResponseBody struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

var (
	fixtureInvalidAPITokenErr = "invalid_api_token: The API token provided " +
		"was not valid (it may not exist or have expired)"
	fixtureInvalidAPITokenResponseError = &ResponseError{
		Code: "invalid_api_token",
		Description: "The API token provided was not valid " +
			"(it may not exist or have expired)",
		Detail: json.RawMessage(`{}`),
	}
)

// setup creates a test HTTP server for mock API responses, and creates a
// Katapult client configured to talk to the mock server.
func setup() (
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
	url, _ := url.Parse(server.URL + "/")
	client = NewClient(nil)
	client.BaseURL = url

	return client, mux, url.String(), server.Close
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
