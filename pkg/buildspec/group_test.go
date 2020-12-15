package buildspec

import (
	"encoding/xml"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGroup_Marshaling(t *testing.T) {
	tests := []struct {
		name    string
		obj     *Group
		decoded *Group
	}{
		{
			name: "empty",
			obj:  &Group{},
		},
		{
			name: "by ID",
			obj:  &Group{ID: "vmgrp_dZDXXLw7e54Ep6CG"},
		},
		{
			name: "by Name",
			obj:  &Group{Name: "Web Servers"},
		},
		{
			name: "with ID and Name",
			obj: &Group{
				ID:   "vmgrp_dZDXXLw7e54Ep6CG",
				Name: "Web Servers",
			},
			decoded: &Group{ID: "vmgrp_dZDXXLw7e54Ep6CG"},
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

func TestGroup_UnmarshalXML_InvalidByAttr(t *testing.T) {
	err := xml.Unmarshal(
		[]byte(`<Group by="other">foo</Group>`),
		&Group{},
	)

	assert.EqualError(t, err,
		`parse_xml: Group by="other" is not supported`,
	)
	assert.True(t, errors.Is(err, ErrParseXML))
}
