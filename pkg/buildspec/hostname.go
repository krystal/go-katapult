package buildspec

type xmlHostname struct {
	Hostname *xmlHostnameValue `xml:",omitempty"`
}

type xmlHostnameValue struct {
	Type  string `xml:"type,attr,omitempty"`
	Value string `xml:",chardata"`
}
