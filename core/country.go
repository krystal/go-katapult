package core

type Country struct {
	ID       string `json:"id,omitempty"`
	Name     string `json:"name,omitempty"`
	ISOCode2 string `json:"iso_code2,omitempty"`
	ISOCode3 string `json:"iso_code3,omitempty"`
	TimeZone string `json:"time_zone,omitempty"`
	EU       bool   `json:"eu,omitempty"`
}

type CountryState struct {
	ID      string   `json:"id,omitempty"`
	Name    string   `json:"name,omitempty"`
	Code    string   `json:"code,omitempty"`
	Country *Country `json:"country,omitempty"`
}
