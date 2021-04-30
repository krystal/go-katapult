package core

type Currency struct {
	ID      string `json:"id,omitempty"`
	Name    string `json:"name,omitempty"`
	ISOCode string `json:"iso_code,omitempty"`
	Symbol  string `json:"symbol,omitempty"`
}
