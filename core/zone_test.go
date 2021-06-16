package core

import (
	"encoding/json"
	"testing"

	"github.com/krystal/go-katapult"
	"github.com/stretchr/testify/assert"
)

var (
	fixtureZoneNotFoundErr = "katapult: not_found: zone_not_found: No zone " +
		"was found matching any of the criteria provided in the arguments"
	fixtureZoneNotFoundResponseError = &katapult.ResponseError{
		Code: "zone_not_found",
		Description: "No zone was found matching any of the criteria " +
			"provided in the arguments",
		Detail: json.RawMessage(`{}`),
	}
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

func TestZone_Ref(t *testing.T) {
	tests := []struct {
		name string
		obj  *Zone
		want ZoneRef
	}{
		{
			name: "empty",
			obj:  &Zone{},
			want: ZoneRef{},
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
			want: ZoneRef{ID: "zone_kY2sPRG24sJVRM2U"},
		},
		{
			name: "no ID or Permalink",
			obj: &Zone{
				Name: "North West",
				DataCenter: &DataCenter{
					ID: "id4",
				},
			},
			want: ZoneRef{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.obj.Ref()

			assert.Equal(t, tt.want, got)
		})
	}
}
