package katapult

import "testing"

func TestIPAddress_JSONMarshaling(t *testing.T) {
	tests := []struct {
		name string
		obj  *IPAddress
	}{
		{
			name: "empty",
			obj:  &IPAddress{},
		},
		{
			name: "full",
			obj: &IPAddress{
				ID:              "id1",
				Address:         "address",
				ReverseDNS:      "reverse_dns",
				VIP:             true,
				AddressWithMask: "address_with_mask",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testJSONMarshaling(t, tt.obj)
		})
	}
}
