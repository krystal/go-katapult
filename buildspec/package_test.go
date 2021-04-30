package buildspec

import (
	"encoding/xml"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPackage_Marshaling(t *testing.T) {
	tests := []struct {
		name    string
		obj     *Package
		decoded *Package
	}{
		{
			name: "empty",
			obj:  &Package{},
		},
		{
			name: "by ID",
			obj:  &Package{ID: "vmpkg_m7mV5O0MafbDFp2n"},
		},
		{
			name: "by Permalink",
			obj:  &Package{Permalink: "rock-3"},
		},
		{
			name: "with ID and Permalink",
			obj: &Package{
				ID:        "vmpkg_m7mV5O0MafbDFp2n",
				Permalink: "rock-3",
			},
			decoded: &Package{ID: "vmpkg_m7mV5O0MafbDFp2n"},
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

func TestPackage_UnmarshalXML_InvalidByAttr(t *testing.T) {
	err := xml.Unmarshal(
		[]byte(`<Package by="other">foo</Package>`),
		&Package{},
	)

	assert.EqualError(t, err, `parse_xml: Package by="other" is not supported`)
	assert.True(t, errors.Is(err, ErrParseXML))
}
