package buildspec

import (
	"testing"
)

func TestResources_Marshaling(t *testing.T) {
	tests := []struct {
		name string
		obj  *Resources
	}{
		{
			name: "empty",
			obj:  &Resources{},
		},
		{
			name: "full",
			obj: &Resources{
				Package:  &Package{ID: "vmpkg_m7mV5O0MafbDFp2n"},
				Memory:   16,
				CPUCores: 4,
			},
		},
	}
	for _, tt := range tests {
		t.Run("json_"+tt.name, func(t *testing.T) {
			testJSONMarshaling(t, tt.obj)
		})
		t.Run("xml_"+tt.name, func(t *testing.T) {
			testXMLMarshaling(t, tt.obj)
		})
		t.Run("yaml_"+tt.name, func(t *testing.T) {
			testYAMLMarshaling(t, tt.obj)
		})
	}
}
