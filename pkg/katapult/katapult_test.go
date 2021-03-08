package katapult

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/url"
	"reflect"
	"testing"
	"time"

	"github.com/augurysys/timestamp"
	"github.com/krystal/go-katapult/internal/codec"
	"github.com/krystal/go-katapult/internal/golden"
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

	if golden.Update() {
		golden.Set(t, buf.Bytes())
	}

	g := golden.Get(t)
	assert.Equal(t, string(g), buf.String(), "encoding does not match golden")

	want := decoded
	if isNil(want) {
		want = input
	}

	got := reflect.New(reflect.TypeOf(want).Elem()).Interface()
	err = c.Decode(bytes.NewBuffer(g), got)
	require.NoError(t, err, "decoding golden failed")
	assert.Equal(t, want, got,
		"decoding from golden does not match expected object",
	)
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

func isNil(i interface{}) bool {
	if i == nil {
		return true
	}

	switch reflect.TypeOf(i).Kind() {
	case reflect.Ptr, reflect.Map, reflect.Array, reflect.Chan, reflect.Slice:
		return reflect.ValueOf(i).IsNil()
	}

	return false
}
