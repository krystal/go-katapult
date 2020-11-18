package katapult

import "testing"

func TestTag_JSONMarshaling(t *testing.T) {
	tests := []struct {
		name string
		obj  *Tag
	}{
		{
			name: "empty",
			obj:  &Tag{},
		},
		{
			name: "full",
			obj: &Tag{
				ID:        "id1",
				Name:      "name",
				Color:     "color",
				CreatedAt: timestampPtr(3043009),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testJSONMarshaling(t, tt.obj)
		})
	}
}
