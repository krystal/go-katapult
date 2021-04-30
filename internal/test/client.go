package test

import (
	"fmt"
	"github.com/krystal/go-katapult"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"testing"
)

var (
	APIKey    = "9d7831d8-03f1-4b4c-a1c3-97272ddefe6a"
	UserAgent = "go-katapult/test"
)

// PrepareTestClient creates a test HTTP server for mock API responses, and
// creates a Katapult client configured to talk to the mock server.
func PrepareTestClient(t *testing.T) (
	client *katapult.Client,
	mux *http.ServeMux,
	serverURL string,
	teardown func(),
) {
	t.Helper()

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
		t.Fatalf("test failed, invalid URL: %s", err.Error())
	}

	rm, err := katapult.New(
		katapult.WithAPIKey(APIKey),
		katapult.WithBaseURL(url),
		katapult.WithUserAgent(UserAgent),
	)
	if err != nil {
		t.Fatalf("failed to setup katapult client: %s", err)
	}

	return rm, mux, url.String(), server.Close
}
