package buildspec

import (
	"encoding/xml"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSubnet_Marshaling(t *testing.T) {
	tests := []struct {
		name    string
		obj     *Subnet
		decoded *Subnet
	}{
		{
			name: "empty",
			obj:  &Subnet{},
		},
		{
			name: "by ID",
			obj:  &Subnet{ID: "sbnt_xxhvuhr3dsvEHcM5"},
		},
		{
			name: "by Address",
			obj:  &Subnet{Address: "148.213.112.119"},
		},
		{
			name: "with ID and Address",
			obj: &Subnet{
				ID:      "sbnt_xxhvuhr3dsvEHcM5",
				Address: "148.213.112.119",
			},
			decoded: &Subnet{ID: "sbnt_xxhvuhr3dsvEHcM5"},
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

func TestSubnet_UnmarshalXML_InvalidByAttr(t *testing.T) {
	err := xml.Unmarshal(
		[]byte(`<Subnet by="other">foo</Subnet>`),
		&Subnet{},
	)

	assert.EqualError(t, err,
		`parse_xml: Subnet by="other" is not supported`,
	)
	assert.True(t, errors.Is(err, ErrParseXML))
}
