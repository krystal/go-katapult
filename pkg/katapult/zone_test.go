package katapult

import "testing"

func TestZone_JSONMarshaling(t *testing.T) {
	tests := []struct {
		name string
		obj  *Zone
	}{
		{
			name: "empty",
			obj:  &Zone{},
		},
		{
			name: "full",
			obj: &Zone{
				ID:         "id",
				Name:       "name",
				Permalink:  "permalink",
				DataCenter: &DataCenter{ID: "id4"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testJSONMarshaling(t, tt.obj)
		})
	}
}
