package buildspec

import (
	"encoding/xml"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDiskIOProfile_Marshaling(t *testing.T) {
	tests := []struct {
		name    string
		obj     *DiskIOProfile
		decoded *DiskIOProfile
	}{
		{
			name: "empty",
			obj:  &DiskIOProfile{},
		},
		{
			name: "by ID",
			obj:  &DiskIOProfile{ID: "diop_xPlNw7iDmrGOnPRA"},
		},
		{
			name: "by Permalink",
			obj:  &DiskIOProfile{Permalink: "ssd"},
		},
		{
			name: "with ID and Permalink",
			obj: &DiskIOProfile{
				ID:        "diop_xPlNw7iDmrGOnPRA",
				Permalink: "ssd",
			},
			decoded: &DiskIOProfile{ID: "diop_xPlNw7iDmrGOnPRA"},
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

func TestDiskIOProfile_UnmarshalXML_InvalidByAttr(t *testing.T) {
	err := xml.Unmarshal(
		[]byte(`<DiskIOProfile by="other">foo</DiskIOProfile>`),
		&DiskIOProfile{},
	)

	assert.EqualError(t, err,
		`parse_xml: DiskIOProfile by="other" is not supported`,
	)
	assert.True(t, errors.Is(err, ErrParseXML))
}
