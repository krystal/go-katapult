package core

type OperatingSystem struct {
	ID    string      `json:"id,omitempty"`
	Name  string      `json:"name,omitempty"`
	Badge *Attachment `json:"badge,omitempty"`
}
