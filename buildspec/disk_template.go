package buildspec

import (
	"encoding/xml"
	"fmt"
)

type DiskTemplate struct {
	ID        string                `json:"id,omitempty" yaml:"id,omitempty"`
	Permalink string                `json:"permalink,omitempty" yaml:"permalink,omitempty"`
	Version   int                   `json:"version,omitempty" yaml:"version,omitempty"`
	Options   []*DiskTemplateOption `json:"options,omitempty" yaml:"options,omitempty"`
}

func (s *DiskTemplate) MarshalXML(
	e *xml.Encoder,
	start xml.StartElement,
) error {
	x := &xmlDiskTemplate{}

	switch {
	case s.ID != "":
		x.DiskTemplate = &xmlDiskTemplateLookup{Value: s.ID}
	case s.Permalink != "":
		x.DiskTemplate = &xmlDiskTemplateLookup{
			By:    permalink,
			Value: s.Permalink,
		}
	}

	if s.Version != 0 {
		x.Version = &xmlDiskTemplateVersion{
			By:    "number",
			Value: s.Version,
		}
	}

	x.Option = s.Options

	return e.EncodeElement(x, start)
}

func (s *DiskTemplate) UnmarshalXML(
	d *xml.Decoder,
	start xml.StartElement,
) error {
	x := &xmlDiskTemplate{}
	_ = d.DecodeElement(x, &start)

	v := DiskTemplate{}

	if x.DiskTemplate != nil {
		switch {
		case x.DiskTemplate.By == "":
			v.ID = x.DiskTemplate.Value
		case x.DiskTemplate.By == permalink:
			v.Permalink = x.DiskTemplate.Value
		default:
			return fmt.Errorf(
				`%w: DiskTemplate by="%s" is not supported`,
				ErrParseXML, x.DiskTemplate.By,
			)
		}
	}

	if x.Version != nil {
		v.Version = x.Version.Value
	}

	v.Options = x.Option

	*s = v

	return nil
}

type DiskTemplateOption struct {
	Key   string `xml:"key,attr" json:"key" yaml:"key"`
	Value string `xml:",chardata" json:"value" yaml:"value"`
}

type xmlDiskTemplate struct {
	DiskTemplate *xmlDiskTemplateLookup  `xml:",omitempty"`
	Version      *xmlDiskTemplateVersion `xml:",omitempty"`
	Option       []*DiskTemplateOption   `xml:",omitempty"`
}

type xmlDiskTemplateLookup xmlLookupElem

type xmlDiskTemplateVersion struct {
	By    string `xml:"by,attr,omitempty"`
	Value int    `xml:",chardata"`
}
