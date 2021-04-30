package test

import (
	"bytes"
	"github.com/krystal/go-katapult/internal/codec"
	"github.com/krystal/go-katapult/internal/golden"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"reflect"
	"testing"
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
