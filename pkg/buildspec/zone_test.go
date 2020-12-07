package buildspec

import (
	"encoding/xml"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestZone_Marshaling(t *testing.T) {
	tests := []struct {
		name    string
		obj     *Zone
		decoded *Zone
	}{
		{
			name: "empty",
			obj:  &Zone{},
		},
		{
			name: "by ID",
			obj:  &Zone{ID: "zone_xmVotL1zwMwo2eXf"},
		},
		{
			name: "by Permalink",
			obj:  &Zone{Permalink: "east-1"},
		},
		{
			name: "with ID and Permalink",
			obj: &Zone{
				ID:        "zone_xmVotL1zwMwo2eXf",
				Permalink: "east-1",
			},
			decoded: &Zone{ID: "zone_xmVotL1zwMwo2eXf"},
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

func TestZone_UnmarshalXML_InvalidByAttr(t *testing.T) {
	err := xml.Unmarshal(
		[]byte(`<Zone by="other">foo</Zone>`),
		&Zone{},
	)

	assert.EqualError(t, err, `parse_xml: Zone by="other" is not supported`)
	assert.True(t, errors.Is(err, ErrParseXML))
}
