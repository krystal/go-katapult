package katapult

import "github.com/augurysys/timestamp"

type TrashObject struct {
	ID        string               `json:"id,omitempty"`
	KeepUntil *timestamp.Timestamp `json:"keep_until,omitempty"`
}
