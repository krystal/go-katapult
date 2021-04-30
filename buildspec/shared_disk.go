package buildspec

import (
	"encoding/xml"
	"fmt"
)

type SharedDisk struct {
	ID   string `json:"id,omitempty" yaml:"id,omitempty"`
	Name string `json:"name,omitempty" yaml:"name,omitempty"`
}

func (s *SharedDisk) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	x := &xmlSharedDisk{}

	switch {
	case s.ID != "":
		x.Value = s.ID
	case s.Name != "":
		x.By = name
		x.Value = s.Name
	}

	return e.EncodeElement(x, start)
}

func (s *SharedDisk) UnmarshalXML(
	d *xml.Decoder,
	start xml.StartElement,
) error {
	x := &xmlSharedDisk{}
	_ = d.DecodeElement(&x, &start)

	v := SharedDisk{}

	switch {
	case x.By == "":
		v.ID = x.Value
	case x.By == name:
		v.Name = x.Value
	default:
		return fmt.Errorf(
			`%w: SharedDisk by="%s" is not supported`, ErrParseXML, x.By,
		)
	}

	*s = v

	return nil
}

type xmlSharedDisk xmlLookupElem

type xmlSharedDisks struct {
	SharedDisks []*SharedDisk `xml:"Disk,omitempty"`
}
