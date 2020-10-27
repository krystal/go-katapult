package codec

import (
	"bytes"
	"io/ioutil"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestJSON_Encode(t *testing.T) {
	type blob struct {
		Msg   string   `json:"msg,omitempty"`
		Tags  []string `json:"tags,omitempty"`
		Count int      `json:"count,omitempty"`
	}
	tests := []struct {
		name   string
		input  interface{}
		want   string
		errStr string
	}{
		{
			name:  "encode json",
			input: &blob{Msg: "hi", Tags: []string{"foo", "bar"}, Count: 42},
			want:  `{"msg":"hi","tags":["foo","bar"],"count":42}`,
		},
		{
			name:  "don't escape HTML",
			input: &blob{Msg: "hello <b>world</b>"},
			want:  `{"msg":"hello <b>world</b>"}`,
		},
		{
			name:   "invalid value",
			input:  make(chan int),
			errStr: "json: unsupported type: chan int",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &JSON{}

			buf := &bytes.Buffer{}
			err := c.Encode(tt.input, buf)

			if tt.errStr != "" {
				assert.EqualError(t, err, tt.errStr)
			} else {
				assert.NoError(t, err)

				got, err := ioutil.ReadAll(buf)
				require.NoError(t, err)
				assert.Equal(t, tt.want+"\n", string(got))
			}
		})
	}
}

func TestJSON_Decode(t *testing.T) {
	type blob struct {
		Tags  []string `json:"tags,omitempty"`
		Count int      `json:"count,omitempty"`
	}
	tests := []struct {
		name   string
		input  string
		want   *blob
		errStr string
	}{
		{
			name:  "decode json",
			input: `{"tags":["foo","bar"],"count":42}`,
			want:  &blob{Tags: []string{"foo", "bar"}, Count: 42},
		},
		{
			name: "empty input - ignore io.EOF error",
			want: &blob{},
		},
		{
			name:   "malformed json",
			input:  `{"tags":["foo`,
			want:   &blob{},
			errStr: "unexpected EOF",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &JSON{}

			got := &blob{}
			err := c.Decode(strings.NewReader(tt.input), got)

			if tt.errStr != "" {
				assert.EqualError(t, err, tt.errStr)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}

func TestJSON_ContentType(t *testing.T) {
	c := &JSON{}

	assert.Equal(t, "application/json", c.ContentType())
}

func TestJSON_Accept(t *testing.T) {
	c := &JSON{}

	assert.Equal(t, "application/json", c.Accept())
}
