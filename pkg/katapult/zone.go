package katapult

type Zone struct {
	ID         string      `json:"id,omitempty"`
	Name       string      `json:"name,omitempty"`
	Permalink  string      `json:"permalink,omitempty"`
	DataCenter *DataCenter `json:"data_center,omitempty"`
}

// LookupReference returns a new *Zone stripped down to just ID or Permalink
// fields, making it suitable for endpoints which require a reference to a Zone
// by ID or Permalink.
func (s *Zone) LookupReference() *Zone {
	if s == nil {
		return nil
	}

	lr := &Zone{ID: s.ID}
	if lr.ID == "" {
		lr.Permalink = s.Permalink
	}

	return lr
}
