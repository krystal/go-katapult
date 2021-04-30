package buildspec

import (
	"encoding/xml"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIPAddress_Marshaling(t *testing.T) {
	tests := []struct {
		name    string
		obj     *IPAddress
		decoded *IPAddress
	}{
		{
			name: "empty",
			obj:  &IPAddress{},
		},
		{
			name: "by ID",
			obj:  &IPAddress{ID: "ip_Hb8WpvV9qRMznHwZ"},
		},
		{
			name: "by Address",
			obj:  &IPAddress{Address: "21.124.234.68"},
		},
		{
			name: "with ID and Address",
			obj: &IPAddress{
				ID:      "ip_Hb8WpvV9qRMznHwZ",
				Address: "21.124.234.68",
			},
			decoded: &IPAddress{ID: "ip_Hb8WpvV9qRMznHwZ"},
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

func TestIPAddress_UnmarshalXML_InvalidByAttr(t *testing.T) {
	err := xml.Unmarshal(
		[]byte(`<IPAddress by="other">foo</IPAddress>`),
		&IPAddress{},
	)

	assert.EqualError(t, err,
		`parse_xml: IPAddress by="other" is not supported`,
	)
	assert.True(t, errors.Is(err, ErrParseXML))
}
