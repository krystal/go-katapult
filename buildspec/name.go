package buildspec

import "encoding/xml"

type xmlName struct {
	Value string `xml:",chardata"`
}

func (s *xmlName) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	x := &xmlNameNested{}
	_ = d.DecodeElement(x, &start)

	if x.Name != nil {
		s.Value = x.Name.Value
	} else {
		s.Value = x.Value
	}

	return nil
}

type xmlNameNested struct {
	Value string              `xml:",chardata"`
	Name  *xmlNameNestedValue `xml:",omitempty"`
}

type xmlNameNestedValue struct {
	Value string `xml:",chardata"`
}
