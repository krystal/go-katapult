package katapult

type Zone struct {
	ID         string      `json:"id,omitempty"`
	Name       string      `json:"name,omitempty"`
	Permalink  string      `json:"permalink,omitempty"`
	DataCenter *DataCenter `json:"data_center,omitempty"`
}

func (s *Zone) lookupReference() *Zone {
	if s == nil {
		return nil
	}

	lr := &Zone{ID: s.ID}
	if lr.ID == "" {
		lr.Permalink = s.Permalink
	}

	return lr
}
