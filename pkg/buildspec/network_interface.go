package buildspec

type NetworkInterface struct {
	Network              *Network               `xml:",omitempty" json:"network,omitempty" yaml:"network,omitempty"`
	VirtualNetwork       *VirtualNetwork        `xml:",omitempty" json:"virtual_network,omitempty" yaml:"virtual_network,omitempty"`
	SpeedProfile         *NetworkSpeedProfile   `xml:",omitempty" json:"speed_profile,omitempty" yaml:"speed_profile,omitempty"`
	IPAddressAllocations []*IPAddressAllocation `xml:"IPAddressAllocation,omitempty" json:"ip_address_allocations,omitempty" yaml:"ip_address_allocations,omitempty"`
}

type xmlNetworkInterfaces struct {
	NetworkInterfaces []*NetworkInterface `xml:"NetworkInterface,omitempty"`
}
