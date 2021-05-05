package core

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"testing"
	"time"

	"github.com/krystal/go-katapult"

	"github.com/augurysys/timestamp"
	"github.com/krystal/go-katapult/internal/golden"
	"github.com/krystal/go-katapult/internal/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

//
// Helpers
//

func boolPtr(b bool) *bool {
	return &b
}

var (
	truePtr  = boolPtr(true)
	falsePtr = boolPtr(false)
)

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
	test.CustomJSONMarshaling(t, input, nil)
}

func testQueryableEncoding(t *testing.T, obj queryable) {
	qs := obj.queryValues()
	queryStr := qs.Encode()

	if golden.Update() {
		golden.Set(t, []byte(queryStr))
	}

	g := string(golden.Get(t))
	assert.Equal(t, queryStr, g, "query string does not match golden")

	parsedQuery, err := url.ParseQuery(g)
	require.NoError(t, err, "parsing golden query string failed")
	assert.Equal(t, qs, &parsedQuery, "parsed golden values do not match")
}

// prepareTestClient creates a test HTTP server for mock API responses, and
// creates a Katapult client configured to talk to the mock server.
func prepareTestClient(t *testing.T) (
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
		katapult.WithAPIKey(test.APIKey),
		katapult.WithBaseURL(url),
		katapult.WithUserAgent(test.UserAgent),
	)
	if err != nil {
		t.Fatalf("failed to setup katapult client: %s", err)
	}

	return rm, mux, url.String(), server.Close
}
