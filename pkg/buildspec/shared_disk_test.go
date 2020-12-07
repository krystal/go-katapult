package buildspec

import (
	"encoding/xml"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSharedDisk_Marshaling(t *testing.T) {
	tests := []struct {
		name    string
		obj     *SharedDisk
		decoded *SharedDisk
	}{
		{
			name: "empty",
			obj:  &SharedDisk{},
		},
		{
			name: "by ID",
			obj:  &SharedDisk{ID: "disk_gJRNxe3h7zi0Hdh5"},
		},
		{
			name: "by Name",
			obj:  &SharedDisk{Name: "file-uploads"},
		},
		{
			name: "with ID and Name",
			obj: &SharedDisk{
				ID:   "disk_gJRNxe3h7zi0Hdh5",
				Name: "file-uploads",
			},
			decoded: &SharedDisk{ID: "disk_gJRNxe3h7zi0Hdh5"},
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

func TestSharedDisk_UnmarshalXML_InvalidByAttr(t *testing.T) {
	err := xml.Unmarshal(
		[]byte(`<SharedDisk by="other">foo</SharedDisk>`),
		&SharedDisk{},
	)

	assert.EqualError(t, err,
		`parse_xml: SharedDisk by="other" is not supported`,
	)
	assert.True(t, errors.Is(err, ErrParseXML))
}

func Test_xmlSharedDisks_Marshaling(t *testing.T) {
	tests := []struct {
		name string
		obj  *xmlSharedDisks
	}{
		{
			name: "empty",
			obj:  &xmlSharedDisks{},
		},
		{
			name: "full",
			obj: &xmlSharedDisks{
				SharedDisks: []*SharedDisk{
					{ID: "disk_gJRNxe3h7zi0Hdh5"},
					{Name: "image-uploads"},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run("xml_"+tt.name, func(t *testing.T) {
			testXMLMarshaling(t, tt.obj)
		})
	}
}
