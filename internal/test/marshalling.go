package test

import (
	"bytes"
	"encoding/json"
	"reflect"
	"testing"

	"github.com/krystal/go-katapult/internal/golden"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

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

func CustomJSONMarshaling(
	t *testing.T,
	input interface{},
	decoded interface{},
) {
	var buf bytes.Buffer
	enc := json.NewEncoder(&buf)
	enc.SetEscapeHTML(false)
	err := enc.Encode(input)
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
	gr := bytes.NewBuffer(g)
	dec := json.NewDecoder(gr)
	dec.DisallowUnknownFields()
	err = dec.Decode(got)
	require.NoError(t, err, "decoding golden failed")
	assert.Equal(t, want, got,
		"decoding from golden does not match expected object",
	)
}
