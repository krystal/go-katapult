package buildspec

import (
	"encoding/xml"
	"errors"
	"testing"

	"github.com/jimeh/undent"
	"github.com/stretchr/testify/assert"
)

func TestDiskTemplate_Marshaling(t *testing.T) {
	tests := []struct {
		name    string
		obj     *DiskTemplate
		decoded *DiskTemplate
	}{
		{
			name: "empty",
			obj:  &DiskTemplate{},
		},
		{
			name: "by ID",
			obj:  &DiskTemplate{ID: "dtpl_rlinMl51Lb1uvTez"},
		},
		{
			name: "by Permalink",
			obj:  &DiskTemplate{Permalink: "ubuntu-18-04"},
		},
		{
			name: "with ID and Permalink",
			obj: &DiskTemplate{
				ID:        "dtpl_rlinMl51Lb1uvTez",
				Permalink: "ubuntu-18-04",
			},
			decoded: &DiskTemplate{ID: "dtpl_rlinMl51Lb1uvTez"},
		},
		{
			name: "full",
			obj: &DiskTemplate{
				ID:      "dtpl_rlinMl51Lb1uvTez",
				Version: 4,
				Options: []*DiskTemplateOption{
					{Key: "foo", Value: "bar"},
					{Key: "hello", Value: "world"},
				},
			},
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

func TestDiskTemplate_UnmarshalXML_InvalidByAttr(t *testing.T) {
	err := xml.Unmarshal(
		[]byte(undent.String(`
			<DiskTemplate>
				<DiskTemplate by="other">foo</DiskTemplate>
			</DiskTemplate>`,
		)),
		&DiskTemplate{},
	)

	assert.EqualError(t, err,
		`parse_xml: DiskTemplate by="other" is not supported`,
	)
	assert.True(t, errors.Is(err, ErrParseXML))
}
