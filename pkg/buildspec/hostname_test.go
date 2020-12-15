package buildspec

import "testing"

func Test_xmlHostname_Marshaling(t *testing.T) {
	tests := []struct {
		name string
		obj  *xmlHostname
	}{
		{
			name: "empty",
			obj:  &xmlHostname{},
		},
		{
			name: "specific",
			obj: &xmlHostname{
				Hostname: &xmlHostnameValue{Value: "bitter-beautiful-mango"},
			},
		},
		{
			name: "random",
			obj: &xmlHostname{
				Hostname: &xmlHostnameValue{Type: "random"},
			},
		},
	}
	for _, tt := range tests {
		t.Run("xml_"+tt.name, func(t *testing.T) {
			testXMLMarshaling(t, tt.obj)
		})
	}
}
