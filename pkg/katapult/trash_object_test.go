package katapult

import "testing"

func TestTrashObject_JSONMarshaling(t *testing.T) {
	tests := []struct {
		name string
		obj  *TrashObject
	}{
		{
			name: "empty",
			obj:  &TrashObject{},
		},
		{
			name: "full",
			obj: &TrashObject{
				ID:        "id",
				KeepUntil: timestampPtr(93043),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testJSONMarshaling(t, tt.obj)
		})
	}
}
