package buildspec

import "encoding/xml"

type AuthorizedKeys struct {
	AllUsers   bool     `json:"all_users,omitempty" yaml:"all_users,omitempty"`
	AllSSHKeys bool     `json:"all_ssh_keys,omitempty" yaml:"all_ssh_keys,omitempty"`
	Users      []*User  `json:"users,omitempty" yaml:"users,omitempty"`
	SSHKeys    []string `json:"ssh_keys,omitempty" yaml:"ssh_keys,omitempty"`
}

func (s *AuthorizedKeys) MarshalXML(
	e *xml.Encoder,
	start xml.StartElement,
) error {
	x := xmlAuthorizedKeys{}

	if s.AllUsers {
		x.Users = &xmlUsers{All: "yes"}
	} else if len(s.Users) > 0 {
		x.Users = &xmlUsers{Users: s.Users}
	}

	if s.AllSSHKeys {
		x.SSHKeys = &xmlSSHKeys{All: "yes"}
	} else if len(s.SSHKeys) > 0 {
		x.SSHKeys = &xmlSSHKeys{SSHKeys: s.SSHKeys}
	}

	return e.EncodeElement(x, start)
}

func (s *AuthorizedKeys) UnmarshalXML(
	d *xml.Decoder,
	start xml.StartElement,
) error {
	x := &xmlAuthorizedKeys{}
	_ = d.DecodeElement(x, &start)

	v := AuthorizedKeys{}

	if x.Users != nil {
		v.AllUsers = (x.Users.All == "yes")
		v.Users = x.Users.Users
	}

	if x.SSHKeys != nil {
		v.AllSSHKeys = (x.SSHKeys.All == "yes")
		v.SSHKeys = x.SSHKeys.SSHKeys
	}

	*s = v

	return nil
}

type xmlAuthorizedKeys struct {
	Users   *xmlUsers   `xml:",omitempty"`
	SSHKeys *xmlSSHKeys `xml:",omitempty"`
}
