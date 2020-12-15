package buildspec

import (
	"encoding/xml"
	"fmt"
)

type User struct {
	ID           string `json:"id,omitempty" yaml:"id,omitempty"`
	EmailAddress string `json:"email_address,omitempty" yaml:"email_address,omitempty"`
}

func (s *User) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	x := &xmlUser{}

	switch {
	case s.ID != "":
		x.Value = s.ID
	case s.EmailAddress != "":
		x.By = "email_address"
		x.Value = s.EmailAddress
	}

	return e.EncodeElement(x, start)
}

func (s *User) UnmarshalXML(
	d *xml.Decoder,
	start xml.StartElement,
) error {
	x := &xmlUser{}
	_ = d.DecodeElement(&x, &start)

	v := User{}

	switch {
	case x.By == "":
		v.ID = x.Value
	case x.By == "email_address":
		v.EmailAddress = x.Value
	default:
		return fmt.Errorf(
			`%w: User by="%s" is not supported`, ErrParseXML, x.By,
		)
	}

	*s = v

	return nil
}

type xmlUser xmlLookupElem

type xmlUsers struct {
	All   string  `xml:"all,attr,omitempty"`
	Users []*User `xml:"User,omitempty"`
}
