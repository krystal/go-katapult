package buildspec

type xmlLookupElem struct {
	By    string `xml:"by,attr,omitempty"`
	Value string `xml:",chardata"`
}
