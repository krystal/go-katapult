package buildspec //nolint:dupl

import (
	"encoding/xml"
	"fmt"
)

type Network struct {
	ID        string `json:"id,omitempty" yaml:"id,omitempty"`
	Permalink string `json:"permalink,omitempty" yaml:"permalink,omitempty"`
}

func (s *Network) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	x := &xmlNetwork{}

	switch {
	case s.ID != "":
		x.Value = s.ID
	case s.Permalink != "":
		x.By = permalink
		x.Value = s.Permalink
	}

	return e.EncodeElement(x, start)
}

func (s *Network) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	x := &xmlNetwork{}
	_ = d.DecodeElement(x, &start)

	v := Network{}

	switch {
	case x.By == "":
		v.ID = x.Value
	case x.By == permalink:
		v.Permalink = x.Value
	default:
		return fmt.Errorf(
			`%w: Network by="%s" is not supported`, ErrParseXML, x.By,
		)
	}

	*s = v

	return nil
}

type xmlNetwork xmlLookupElem
