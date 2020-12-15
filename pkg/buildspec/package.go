package buildspec //nolint:dupl

import (
	"encoding/xml"
	"fmt"
)

type Package struct {
	ID        string `json:"id,omitempty" yaml:"id,omitempty"`
	Permalink string `json:"permalink,omitempty" yaml:"permalink,omitempty"`
}

func (s *Package) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	x := &xmlPackage{}

	switch {
	case s.ID != "":
		x.Value = s.ID
	case s.Permalink != "":
		x.By = permalink
		x.Value = s.Permalink
	}

	return e.EncodeElement(x, start)
}

func (s *Package) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	x := &xmlPackage{}
	_ = d.DecodeElement(x, &start)

	v := Package{}

	switch {
	case x.By == "":
		v.ID = x.Value
	case x.By == permalink:
		v.Permalink = x.Value
	default:
		return fmt.Errorf(
			`%w: Package by="%s" is not supported`, ErrParseXML, x.By,
		)
	}

	*s = v

	return nil
}

type xmlPackage xmlLookupElem
