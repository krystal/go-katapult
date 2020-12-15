package buildspec

import "testing"

func TestIPAddressAllocation_Marshaling(t *testing.T) {
	tests := []struct {
		name    string
		obj     *IPAddressAllocation
		decoded *IPAddressAllocation
	}{
		{
			name: "empty",
			obj:  &IPAddressAllocation{},
		},
		{
			name: "new IPAddress",
			obj: &IPAddressAllocation{
				Type:    NewIPAddressAllocation,
				Version: 4,
			},
		},
		{
			name: "new IPAddress with Subnet",
			obj: &IPAddressAllocation{
				Type:    NewIPAddressAllocation,
				Version: 4,
				Subnet:  &Subnet{ID: "sbnt_xxhvuhr3dsvEHcM5"},
			},
		},
		{
			name: "existing IPAddress",
			obj: &IPAddressAllocation{
				Type:      ExistingIPAddressAllocation,
				IPAddress: &IPAddress{ID: "ip_Hb8WpvV9qRMznHwZ"},
			},
		},
	}
	for _, tt := range tests {
		t.Run("json_"+tt.name, func(t *testing.T) {
			testJSONMarshaling(t, tt.obj)
		})
		t.Run("xml_"+tt.name, func(t *testing.T) {
			testCustomXMLMarshaling(t, tt.obj, tt.decoded)
		})
		t.Run("yaml_"+tt.name, func(t *testing.T) {
			testYAMLMarshaling(t, tt.obj)
		})
	}
}
