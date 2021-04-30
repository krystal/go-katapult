package core

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	fixtureZoneNotFoundErr = "zone_not_found: No zone was found matching " +
		"any of the criteria provided in the arguments"
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

func TestNewZoneLookup(t *testing.T) {
	type args struct {
		idOrPermalink string
	}
	tests := []struct {
		name  string
		args  args
		want  *Zone
		field FieldName
	}{
		{
			name:  "empty string",
			args:  args{idOrPermalink: ""},
			want:  &Zone{},
			field: PermalinkField,
		},
		{
			name:  "zone_ prefixed ID",
			args:  args{idOrPermalink: "zone_NeK95mFtiSXfUtW2"},
			want:  &Zone{ID: "zone_NeK95mFtiSXfUtW2"},
			field: IDField,
		},
		{
			name:  "permalink",
			args:  args{idOrPermalink: "city-zone-1"},
			want:  &Zone{Permalink: "city-zone-1"},
			field: PermalinkField,
		},
		{
			name:  "random text",
			args:  args{idOrPermalink: "1UKzPq0izsQsNCsd"},
			want:  &Zone{Permalink: "1UKzPq0izsQsNCsd"},
			field: PermalinkField,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, field := NewZoneLookup(tt.args.idOrPermalink)

			assert.Equal(t, tt.want, got)
			assert.Equal(t, tt.field, field)
		})
	}
}

func TestZone_lookupReference(t *testing.T) {
	tests := []struct {
		name string
		obj  *Zone
		want *Zone
	}{
		{
			name: "nil",
			obj:  nil,
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
			got := tt.obj.lookupReference()

			assert.Equal(t, tt.want, got)
		})
	}
}
