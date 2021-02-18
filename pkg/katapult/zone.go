package katapult

import "strings"

type Zone struct {
	ID         string      `json:"id,omitempty"`
	Name       string      `json:"name,omitempty"`
	Permalink  string      `json:"permalink,omitempty"`
	DataCenter *DataCenter `json:"data_center,omitempty"`
}

// NewZoneLookup takes a string that is a Zone ID or Permalink, returning a
// empty *Zone struct with either the ID or Permalink field populated with the
// given value. This struct is suitable as input to other methods which accept a
// *Zone as input.
func NewZoneLookup(
	idOrPermalink string,
) (lr *Zone, f FieldName) {
	if strings.HasPrefix(idOrPermalink, "zone_") {
		return &Zone{ID: idOrPermalink}, IDField
	}

	return &Zone{Permalink: idOrPermalink}, PermalinkField
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
