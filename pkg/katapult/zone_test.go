package katapult

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

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
				ID:        "zone_kY2sPRG24sJVRM2U",
				Name:      "North West",
				Permalink: "north-west",
				DataCenter: &DataCenter{
					ID: "id4",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testJSONMarshaling(t, tt.obj)
		})
	}
}

func TestZone_LookupReference(t *testing.T) {
	tests := []struct {
		name string
		obj  *Zone
		want *Zone
	}{
		{
			name: "nil",
			obj:  (*Zone)(nil),
			want: nil,
		},
		{
			name: "empty",
			obj:  &Zone{},
			want: &Zone{},
		},
		{
			name: "full",
			obj: &Zone{
				ID:        "zone_kY2sPRG24sJVRM2U",
				Name:      "North West",
				Permalink: "north-west",
				DataCenter: &DataCenter{
					ID: "id4",
				},
			},
			want: &Zone{ID: "zone_kY2sPRG24sJVRM2U"},
		},
		{
			name: "no ID",
			obj: &Zone{
				Name:      "North West",
				Permalink: "north-west",
				DataCenter: &DataCenter{
					ID: "id4",
				},
			},
			want: &Zone{Permalink: "north-west"},
		},
		{
			name: "no ID or Permalink",
			obj: &Zone{
				Name: "North West",
				DataCenter: &DataCenter{
					ID: "id4",
				},
			},
			want: &Zone{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.obj.LookupReference()

			assert.Equal(t, tt.want, got)
		})
	}
}
