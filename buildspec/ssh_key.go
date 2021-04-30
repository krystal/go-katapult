package buildspec

type xmlSSHKeys struct {
	All     string   `xml:"all,attr,omitempty"`
	SSHKeys []string `xml:"SSHKey,omitempty"`
}
