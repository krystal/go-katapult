package core

import "github.com/augurysys/timestamp"

type Tag struct {
	ID        string               `json:"id,omitempty"`
	Name      string               `json:"name,omitempty"`
	Color     string               `json:"color,omitempty"`
	CreatedAt *timestamp.Timestamp `json:"created_at,omitempty"`
}
