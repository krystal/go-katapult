package core

type ISO struct {
	ID              string           `json:"id,omitempty"`
	Name            string           `json:"name,omitempty"`
	OperatingSystem *OperatingSystem `json:"operating_system,omitempty"`
}
