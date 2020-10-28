package katapult

import "testing"

func TestCountry_JSONMarshaling(t *testing.T) {
	tests := []struct {
		name string
		obj  *Country
	}{
		{
			name: "empty",
			obj:  &Country{},
		},
		{
			name: "full",
			obj: &Country{
				ID:       "id",
				Name:     "name",
				ISOCode2: "iso2",
				ISOCode3: "iso3",
				TimeZone: "tz",
				EU:       true,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testJSONMarshaling(t, tt.obj)
		})
	}
}

func TestCountryState_JSONMarshaling(t *testing.T) {
	tests := []struct {
		name string
		obj  *CountryState
	}{
		{
			name: "empty",
			obj:  &CountryState{},
		},
		{
			name: "full",
			obj: &CountryState{
				ID:   "id",
				Name: "name",
				Code: "code",
				Country: &Country{
					ID: "id2",
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
