package buildspec

import (
	"encoding/xml"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDataCenter_Marshaling(t *testing.T) {
	tests := []struct {
		name    string
		obj     *DataCenter
		decoded *DataCenter
	}{
		{
			name: "empty",
			obj:  &DataCenter{},
		},
		{
			name: "by ID",
			obj:  &DataCenter{ID: "dc_0KVdXStXduYtcypG"},
		},
		{
			name: "by Name",
			obj:  &DataCenter{Name: "London (UK)"},
		},
		{
			name: "by Permalink",
			obj:  &DataCenter{Permalink: "london"},
		},
		{
			name: "with ID, Name, and Permalink",
			obj: &DataCenter{
				ID:        "dc_0KVdXStXduYtcypG",
				Name:      "London (UK)",
				Permalink: "london",
			},
			decoded: &DataCenter{ID: "dc_0KVdXStXduYtcypG"},
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

func TestDataCenter_UnmarshalXML_InvalidByAttr(t *testing.T) {
	err := xml.Unmarshal(
		[]byte(`<DataCenter by="other">foo</DataCenter>`),
		&DataCenter{},
	)

	assert.EqualError(t, err,
		`parse_xml: DataCenter by="other" is not supported`,
	)
	assert.True(t, errors.Is(err, ErrParseXML))
}
