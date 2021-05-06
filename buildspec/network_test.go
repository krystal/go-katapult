package buildspec

import (
	"encoding/xml"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNetwork_Marshaling(t *testing.T) {
	tests := []struct {
		name    string
		obj     *Network
		decoded *Network
	}{
		{
			name: "empty",
			obj:  &Network{},
		},
		{
			name: "by ID",
			obj:  &Network{ID: "netw_17w3MepxvWE4J3Zx"},
		},
		{
			name: "by Permalink",
			obj:  &Network{Permalink: "public"},
		},
		{
			name: "with ID and Permalink",
			obj: &Network{
				ID:        "netw_17w3MepxvWE4J3Zx",
				Permalink: "public",
			},
			decoded: &Network{ID: "netw_17w3MepxvWE4J3Zx"},
		},
	}
	for _, tt := range tests {
		t.Run("json_"+tt.name, func(t *testing.T) {
			testJSONMarshaling(t, tt.obj)
		})
		t.Run("xml_"+tt.name, func(t *testing.T) {
			testCustomXMLMarshaling(t, tt.obj, tt.decoded)
		})
		t.Run("yaml_"+tt.name, func(t *testing.T) {
			testYAMLMarshaling(t, tt.obj)
		})
	}
}

func TestNetwork_UnmarshalXML_InvalidByAttr(t *testing.T) {
	err := xml.Unmarshal(
		[]byte(`<Network by="other">foo</Network>`),
		&Network{},
	)

	assert.EqualError(t, err, `parse_xml: Network by="other" is not supported`)
	assert.True(t, errors.Is(err, ErrParseXML))
}
