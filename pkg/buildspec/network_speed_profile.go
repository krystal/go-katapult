package buildspec //nolint:dupl

import (
	"encoding/xml"
	"fmt"
)

type NetworkSpeedProfile struct {
	ID        string `json:"id,omitempty" yaml:"id,omitempty"`
	Permalink string `json:"permalink,omitempty" yaml:"permalink,omitempty"`
}

func (s *NetworkSpeedProfile) MarshalXML(
	e *xml.Encoder,
	start xml.StartElement,
) error {
	x := &xmlNetworkSpeedProfile{}

	switch {
	case s.ID != "":
		x.Value = s.ID
	case s.Permalink != "":
		x.By = permalink
		x.Value = s.Permalink
	}

	return e.EncodeElement(x, start)
}

func (s *NetworkSpeedProfile) UnmarshalXML(
	d *xml.Decoder,
	start xml.StartElement,
) error {
	x := &xmlNetworkSpeedProfile{}
	_ = d.DecodeElement(x, &start)

	v := NetworkSpeedProfile{}

	switch {
	case x.By == "":
		v.ID = x.Value
	case x.By == permalink:
		v.Permalink = x.Value
	default:
		return fmt.Errorf(
			`%w: NetworkSpeedProfile by="%s" is not supported`,
			ErrParseXML, x.By,
		)
	}

	*s = v

	return nil
}

type xmlNetworkSpeedProfile xmlLookupElem
