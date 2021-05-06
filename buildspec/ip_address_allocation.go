package buildspec

type IPAddressAllocationType string

const (
	NewIPAddressAllocation      IPAddressAllocationType = "new"
	ExistingIPAddressAllocation IPAddressAllocationType = "existing"
)

type IPAddressAllocation struct {
	Type      IPAddressAllocationType `xml:"type,attr,omitempty" json:"type,omitempty" yaml:"type,omitempty"`
	IPAddress *IPAddress              `xml:",omitempty" json:"ip_address,omitempty" yaml:"ip_address,omitempty"`
	Version   IPVersion               `xml:",omitempty" json:"version,omitempty" yaml:"version,omitempty"`
	Subnet    *Subnet                 `xml:",omitempty" json:"subnet,omitempty" yaml:"subnet,omitempty"`
}
