package buildspec //nolint:dupl

import (
	"encoding/xml"
	"fmt"
)

type IPAddress struct {
	ID      string `json:"id,omitempty" yaml:"id,omitempty"`
	Address string `json:"address,omitempty" yaml:"address,omitempty"`
}

func (s *IPAddress) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	x := &xmlIPAddress{}

	switch {
	case s.ID != "":
		x.Value = s.ID
	case s.Address != "":
		x.By = address
		x.Value = s.Address
	}

	return e.EncodeElement(x, start)
}

func (s *IPAddress) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	x := &xmlIPAddress{}
	_ = d.DecodeElement(x, &start)

	v := IPAddress{}

	switch {
	case x.By == "":
		v.ID = x.Value
	case x.By == address:
		v.Address = x.Value
	default:
		return fmt.Errorf(
			`%w: IPAddress by="%s" is not supported`, ErrParseXML, x.By,
		)
	}

	*s = v

	return nil
}

type xmlIPAddress xmlLookupElem
