package core

import "testing"

func TestOperatingSystem_JSONMarshaling(t *testing.T) {
	tests := []struct {
		name string
		obj  *OperatingSystem
	}{
		{
			name: "empty",
			obj:  &OperatingSystem{},
		},
		{
			name: "full",
			obj: &OperatingSystem{
				ID:    "id1",
				Name:  "name",
				Badge: &Attachment{URL: "url2"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testJSONMarshaling(t, tt.obj)
		})
	}
}
