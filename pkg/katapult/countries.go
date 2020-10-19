package katapult

type Country struct {
	ID       string `json:"id,omitempty"`
	Name     string `json:"name,omitempty"`
	ISOCode2 string `json:"iso_code2,omitempty"`
	ISOCode3 string `json:"iso_code3,omitempty"`
	TimeZone string `json:"time_zone,omitempty"`
	EU       bool   `json:"eu,omitempty"`
}
