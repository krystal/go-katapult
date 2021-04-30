package core

import "testing"

func TestCurrency_JSONMarshaling(t *testing.T) {
	tests := []struct {
		name string
		obj  *Currency
	}{
		{
			name: "empty",
			obj:  &Currency{},
		},
		{
			name: "full",
			obj: &Currency{
				ID:      "id",
				Name:    "name",
				ISOCode: "iso",
				Symbol:  "symbol",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testJSONMarshaling(t, tt.obj)
		})
	}
}
