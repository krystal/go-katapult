package buildspec

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"errors"
	"io"
	"reflect"
	"strings"
	"testing"

	"github.com/jimeh/undent"
	"github.com/krystal/go-katapult/internal/golden"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"
)

//
// Helpers
//

type badWriter struct {
	err error
}

func (s *badWriter) Write(p []byte) (int, error) {
	return 0, s.err
}

func testJSONMarshaling(t *testing.T, v interface{}) {
	marshaled, err := json.MarshalIndent(v, "", "  ")
	require.NoError(t, err, "json encoding failed")

	if golden.Update() {
		golden.Set(t, marshaled)
	}

	g := golden.Get(t)
	assert.Equal(t, string(g), string(marshaled),
		"json encoding does not match golden",
	)

	got := reflect.New(reflect.TypeOf(v).Elem()).Interface()
	err = json.Unmarshal(g, got)
	require.NoError(t, err, "json decoding golden failed")
	assert.Equal(t, v, got,
		"json decoding from golden does not match expected object",
	)
}

func testXMLMarshaling(t *testing.T, v interface{}) {
	testCustomXMLMarshaling(t, v, nil)
}

func testCustomXMLMarshaling(
	t *testing.T,
	v interface{},
	decoded interface{},
) {
	marshaled, err := xml.MarshalIndent(v, "", "  ")
	require.NoError(t, err, "xml encoding failed")

	if golden.Update() {
		golden.Set(t, marshaled)
	}

	g := golden.Get(t)
	assert.Equal(t, string(g), string(marshaled),
		"xml encoding does not match golden",
	)

	want := decoded
	if isNil(want) {
		want = v
	}

	got := reflect.New(reflect.TypeOf(want).Elem()).Interface()
	err = xml.Unmarshal(g, got)
	require.NoError(t, err, "xml decoding golden failed")
	assert.Equal(t, want, got,
		"xml decoding from golden does not match expected object",
	)
}

func testYAMLMarshaling(t *testing.T, v interface{}) {
	buf := bytes.Buffer{}
	enc := yaml.NewEncoder(&buf)
	enc.SetIndent(2)
	err := enc.Encode(v)
	require.NoError(t, err, "yaml encoding failed")
	marshaled := buf.Bytes()

	if golden.Update() {
		golden.Set(t, marshaled)
	}

	g := golden.Get(t)
	assert.Equal(t, string(g), string(marshaled),
		"yaml encoding does not match golden",
	)

	got := reflect.New(reflect.TypeOf(v).Elem()).Interface()
	err = yaml.Unmarshal(g, got)
	require.NoError(t, err, "yaml decoding golden failed")
	assert.Equal(t, v, got,
		"yaml decoding from golden does not match expected object",
	)
}

func isNil(i interface{}) bool {
	if i == nil {
		return true
	}

	switch reflect.TypeOf(i).Kind() { //nolint:exhaustive
	case reflect.Ptr, reflect.Map, reflect.Array, reflect.Chan, reflect.Slice:
		return reflect.ValueOf(i).IsNil()
	default:
		return false
	}
}

//
// Tests
//

func TestFromJSON(t *testing.T) {
	tests := []struct {
		name   string
		json   string
		want   *VirtualMachineSpec
		errIs  error
		errStr string
	}{
		{
			name:  "empty string",
			json:  ``,
			want:  &VirtualMachineSpec{},
			errIs: io.EOF,
		},
		{
			name: "empty object",
			json: `{}`,
			want: &VirtualMachineSpec{},
		},
		{
			name: "basic VirtualMachineSpec",
			json: undent.String(`
				{
					"name": "web-3"
				}`,
			),
			want: &VirtualMachineSpec{Name: "web-3"},
		},
		{
			name: "invalid attribute",
			json: undent.String(`
				{
					"rocket_fuel": "maybe"
				}`,
			),
			want:   &VirtualMachineSpec{},
			errStr: `json: unknown field "rocket_fuel"`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := FromJSON(strings.NewReader(tt.json))

			if tt.errIs != nil {
				assert.True(t, errors.Is(err, tt.errIs))
			}

			if tt.errStr != "" {
				assert.EqualError(t, err, tt.errStr)
			}

			if tt.errIs == nil && tt.errStr == "" {
				assert.NoError(t, err)
			}

			assert.Equal(t, tt.want, got)
		})
	}
}

func TestFromXML(t *testing.T) {
	tests := []struct {
		name   string
		xml    string
		want   *VirtualMachineSpec
		errIs  error
		errStr string
	}{
		{
			name:  "empty string",
			xml:   ``,
			want:  &VirtualMachineSpec{},
			errIs: io.EOF,
		},
		{
			name: "empty XML",
			xml: undent.String(`
				<?xml version="1.0" encoding="UTF-8"?>`,
			),
			want:  &VirtualMachineSpec{},
			errIs: io.EOF,
		},
		{
			name: "empty VirtualMachineSpec",
			xml: undent.String(`
				<?xml version="1.0" encoding="UTF-8"?>
				<VirtualMachineSpec></VirtualMachineSpec>`,
			),
			want: &VirtualMachineSpec{},
		},
		{
			name: "basic VirtualMachineSpec",
			xml: undent.String(`
				<?xml version="1.0" encoding="UTF-8"?>
				<VirtualMachineSpec>
					<Name>web-3</Name>
				</VirtualMachineSpec>`,
			),
			want: &VirtualMachineSpec{Name: "web-3"},
		},
		{
			name: "missing XML header",
			xml: undent.String(`
				<VirtualMachineSpec>
					<Name>web-3</Name>
				</VirtualMachineSpec>`,
			),
			want: &VirtualMachineSpec{Name: "web-3"},
		},
		{
			name: "invalid child element",
			xml: undent.String(`
				<?xml version="1.0" encoding="UTF-8"?>
				<VirtualMachineSpec>
					<RocketFuel>maybe</RocketFuel>
				</VirtualMachineSpec>`,
			),
			want: &VirtualMachineSpec{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := FromXML(strings.NewReader(tt.xml))

			if tt.errIs != nil {
				assert.True(t, errors.Is(err, tt.errIs))
			}

			if tt.errStr != "" {
				assert.EqualError(t, err, tt.errStr)
			}

			if tt.errIs == nil && tt.errStr == "" {
				assert.NoError(t, err)
			}

			assert.Equal(t, tt.want, got)
		})
	}
}

func TestFromYAML(t *testing.T) {
	tests := []struct {
		name   string
		yaml   string
		want   *VirtualMachineSpec
		errIs  error
		errStr string
	}{
		{
			name:  "empty string",
			yaml:  ``,
			want:  &VirtualMachineSpec{},
			errIs: io.EOF,
		},
		{
			name: "empty object",
			yaml: `{}`,
			want: &VirtualMachineSpec{},
		},
		{
			name: "basic VirtualMachineSpec",
			yaml: undent.String(`
				name: web-3`,
			),
			want: &VirtualMachineSpec{Name: "web-3"},
		},
		{
			name: "invalid attribute",
			yaml: undent.String(`
				rocket_fuel: maybe`,
			),
			want: &VirtualMachineSpec{},
			errStr: undent.String(`
				yaml: unmarshal errors:
				  line 1: field rocket_fuel not found in type buildspec.VirtualMachineSpec`, //nolint:lll
			),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := FromYAML(strings.NewReader(tt.yaml))

			if tt.errIs != nil {
				assert.True(t, errors.Is(err, tt.errIs))
			}

			if tt.errStr != "" {
				assert.EqualError(t, err, tt.errStr)
			}

			if tt.errIs == nil && tt.errStr == "" {
				assert.NoError(t, err)
			}

			assert.Equal(t, tt.want, got)
		})
	}
}
