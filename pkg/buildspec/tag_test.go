package buildspec

import "testing"

func Test_xmlTags_Marshaling(t *testing.T) {
	tests := []struct {
		name string
		obj  *xmlTags
	}{
		{
			name: "empty",
			obj:  &xmlTags{},
		},
		{
			name: "full",
			obj: &xmlTags{
				Tags: []string{"ha", "db", "web"},
			},
		},
	}
	for _, tt := range tests {
		t.Run("xml_"+tt.name, func(t *testing.T) {
			testXMLMarshaling(t, tt.obj)
		})
	}
}
