package buildspec //nolint:dupl

import (
	"encoding/xml"
	"fmt"
)

type Subnet struct {
	ID      string `json:"id,omitempty" yaml:"id,omitempty"`
	Address string `json:"address,omitempty" yaml:"address,omitempty"`
}

func (s *Subnet) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	x := &xmlSubnet{}

	switch {
	case s.ID != "":
		x.Value = s.ID
	case s.Address != "":
		x.By = address
		x.Value = s.Address
	}

	return e.EncodeElement(x, start)
}

func (s *Subnet) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	x := &xmlSubnet{}
	_ = d.DecodeElement(x, &start)

	v := Subnet{}

	switch {
	case x.By == "":
		v.ID = x.Value
	case x.By == address:
		v.Address = x.Value
	default:
		return fmt.Errorf(
			`%w: Subnet by="%s" is not supported`, ErrParseXML, x.By,
		)
	}

	*s = v

	return nil
}

type xmlSubnet xmlLookupElem
