package buildspec

import "testing"

func TestNetworkInterface_Marshaling(t *testing.T) {
	tests := []struct {
		name string
		obj  *NetworkInterface
	}{
		{
			name: "empty",
			obj:  &NetworkInterface{},
		},
		{
			name: "for Network",
			obj: &NetworkInterface{
				Network:      &Network{ID: "netw_17w3MepxvWE4J3Zx"},
				SpeedProfile: &NetworkSpeedProfile{ID: "nsp_bFQhDNAluyp4t2A9"},
				IPAddressAllocations: []*IPAddressAllocation{
					{Type: NewIPAddressAllocation},
					{Type: ExistingIPAddressAllocation},
				},
			},
		},
		{
			name: "for VirtualNetwork",
			obj: &NetworkInterface{
				VirtualNetwork: &VirtualNetwork{ID: "vnet_Cuc45YcBaUhWqx6u"},
			},
		},
	}
	for _, tt := range tests {
		t.Run("json_"+tt.name, func(t *testing.T) {
			testJSONMarshaling(t, tt.obj)
		})
		t.Run("xml_"+tt.name, func(t *testing.T) {
			testXMLMarshaling(t, tt.obj)
		})
		t.Run("yaml_"+tt.name, func(t *testing.T) {
			testYAMLMarshaling(t, tt.obj)
		})
	}
}

func Test_xmlNetworkInterfaces_Marshaling(t *testing.T) {
	tests := []struct {
		name string
		obj  *xmlNetworkInterfaces
	}{
		{
			name: "empty",
			obj:  &xmlNetworkInterfaces{},
		},
		{
			name: "full",
			obj: &xmlNetworkInterfaces{
				NetworkInterfaces: []*NetworkInterface{
					{
						Network: &Network{ID: "netw_17w3MepxvWE4J3Zx"},
						SpeedProfile: &NetworkSpeedProfile{
							ID: "nsp_bFQhDNAluyp4t2A9",
						},
						IPAddressAllocations: []*IPAddressAllocation{
							{Type: NewIPAddressAllocation},
							{Type: ExistingIPAddressAllocation},
						},
					},
					{
						VirtualNetwork: &VirtualNetwork{
							ID: "vnet_Cuc45YcBaUhWqx6u",
						},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run("xml_"+tt.name, func(t *testing.T) {
			testXMLMarshaling(t, tt.obj)
		})
	}
}
