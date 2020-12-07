package buildspec

import (
	"encoding/xml"
	"fmt"
)

type DataCenter struct {
	ID        string `json:"id,omitempty" yaml:"id,omitempty"`
	Name      string `json:"name,omitempty" yaml:"name,omitempty"`
	Permalink string `json:"permalink,omitempty" yaml:"permalink,omitempty"`
}

func (s *DataCenter) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	x := &xmlDataCenter{}

	switch {
	case s.ID != "":
		x.Value = s.ID
	case s.Name != "":
		x.By = name
		x.Value = s.Name
	case s.Permalink != "":
		x.By = permalink
		x.Value = s.Permalink
	}

	return e.EncodeElement(x, start)
}

func (s *DataCenter) UnmarshalXML(
	d *xml.Decoder,
	start xml.StartElement,
) error {
	x := &xmlDataCenter{}
	_ = d.DecodeElement(&x, &start)

	v := DataCenter{}

	switch {
	case x.By == "":
		v.ID = x.Value
	case x.By == name:
		v.Name = x.Value
	case x.By == permalink:
		v.Permalink = x.Value
	default:
		return fmt.Errorf(
			`%w: DataCenter by="%s" is not supported`, ErrParseXML, x.By,
		)
	}

	*s = v

	return nil
}

type xmlDataCenter xmlLookupElem
