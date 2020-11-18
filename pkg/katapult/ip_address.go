package katapult

type IPAddress struct {
	ID              string `json:"id,omitempty"`
	Address         string `json:"address,omitempty"`
	ReverseDNS      string `json:"reverse_dns,omitempty"`
	VIP             bool   `json:"vip,omitempty"`
	AddressWithMask string `json:"address_with_mask,omitempty"`
}
