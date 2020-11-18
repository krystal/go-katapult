package katapult

import "testing"

func TestISO_JSONMarshaling(t *testing.T) {
	tests := []struct {
		name string
		obj  *ISO
	}{
		{
			name: "empty",
			obj:  &ISO{},
		},
		{
			name: "full",
			obj: &ISO{
				ID:              "id1",
				Name:            "name",
				OperatingSystem: &OperatingSystem{ID: "id2"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testJSONMarshaling(t, tt.obj)
		})
	}
}
