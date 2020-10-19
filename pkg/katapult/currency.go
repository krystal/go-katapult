package katapult

type Currency struct {
	ID      string `json:"id,omitempty"`
	Name    string `json:"name,omitempty"`
	IsoCode string `json:"iso_code,omitempty"`
	Symbol  string `json:"symbol,omitempty"`
}
