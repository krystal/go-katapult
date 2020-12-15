package buildspec

import (
	"encoding/xml"
	"testing"

	"github.com/jimeh/undent"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_xmlName_Marshaling(t *testing.T) {
	tests := []struct {
		name string
		obj  *xmlName
	}{
		{
			name: "empty",
			obj:  &xmlName{},
		},
		{
			name: "value",
			obj:  &xmlName{Value: "web-1"},
		},
	}
	for _, tt := range tests {
		t.Run("xml_"+tt.name, func(t *testing.T) {
			testXMLMarshaling(t, tt.obj)
		})
	}
}

func Test_xmlName_UnmarshalXML(t *testing.T) {
	tests := []struct {
		name string
		xml  string
		want *xmlName
	}{
		{
			name: "empty",
			xml:  `<Name></Name>`,
			want: &xmlName{},
		},
		{
			name: "value",
			xml:  `<Name>database-2</Name>`,
			want: &xmlName{Value: "database-2"},
		},
		{
			name: "empty nested",
			xml: undent.String(`
				<Name>
					<Name></Name>
				</Name>`,
			),
			want: &xmlName{},
		},
		{
			name: "nested value",
			xml: undent.String(`
				<Name>
					<Name>database-2</Name>
				</Name>`,
			),
			want: &xmlName{Value: "database-2"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := &xmlName{}

			err := xml.Unmarshal([]byte(tt.xml), got)
			require.NoError(t, err)

			assert.Equal(t, tt.want, got)
		})
	}
}
