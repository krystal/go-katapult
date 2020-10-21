package katapult

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
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

// setup creates a test HTTP server for mock API responses, and creates a
// Katapult client configure to talk to the mock server.
func setup() (
	client *Client,
	mux *http.ServeMux,
	serverURL string,
	teardown func(),
) {
	mux = http.NewServeMux()
	baseURL, err := url.Parse(testDefaultBaseURL)
	if err != nil {
		log.Fatal(err)
	}

	path := baseURL.Path
	if path[len(path)-1:] == "/" {
		path = path[0 : len(path)-1]
	}

	apiHandler := http.NewServeMux()
	apiHandler.Handle(path+"/", http.StripPrefix(path, mux))
	apiHandler.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		fmt.Fprintln(
			os.Stderr,
			"FAIL: Request for unhandled request in test server received:",
		)
		fmt.Fprintln(os.Stderr)
		fmt.Fprintln(os.Stderr, "\t"+req.URL.String())
	})

	server := httptest.NewServer(apiHandler)
	url, _ := url.Parse(server.URL + baseURL.Path)
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
