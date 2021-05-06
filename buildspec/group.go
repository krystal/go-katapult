package buildspec

import (
	"encoding/xml"
	"fmt"
)

type Group struct {
	ID   string `json:"id,omitempty" yaml:"id,omitempty"`
	Name string `json:"name,omitempty" yaml:"name,omitempty"`
}

func (s *Group) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	x := &xmlGroup{}

	switch {
	case s.ID != "":
		x.Value = s.ID
	case s.Name != "":
		x.By = name
		x.Value = s.Name
	}

	return e.EncodeElement(x, start)
}

func (s *Group) UnmarshalXML(
	d *xml.Decoder,
	start xml.StartElement,
) error {
	x := &xmlGroup{}
	_ = d.DecodeElement(&x, &start)

	v := Group{}

	switch {
	case x.By == "":
		v.ID = x.Value
	case x.By == name:
		v.Name = x.Value
	default:
		return fmt.Errorf(
			`%w: Group by="%s" is not supported`, ErrParseXML, x.By,
		)
	}

	*s = v

	return nil
}

type xmlGroup xmlLookupElem
