package buildspec

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"io"

	"gopkg.in/yaml.v3"
)

// VirtualMachineSpec is the top-level object representing a virtual machine
// build specification. It can be created manually, or parsed from XML, JSON, or
// YAML, and can output itself as XML, JSON, and YAML.
type VirtualMachineSpec struct {
	Zone              *Zone               `json:"zone,omitempty" yaml:"zone,omitempty"`
	DataCenter        *DataCenter         `json:"data_center,omitempty" yaml:"data_center,omitempty"`
	Resources         *Resources          `json:"resources,omitempty" yaml:"resources,omitempty"`
	DiskTemplate      *DiskTemplate       `json:"disk_template,omitempty" yaml:"disk_template,omitempty"`
	SystemDisks       []*SystemDisk       `json:"system_disks,omitempty" yaml:"system_disks,omitempty"`
	SharedDisks       []*SharedDisk       `json:"shared_disks,omitempty" yaml:"shared_disks,omitempty"`
	NetworkInterfaces []*NetworkInterface `json:"network_interfaces,omitempty" yaml:"network_interfaces,omitempty"`
	Hostname          string              `json:"hostname,omitempty" yaml:"hostname,omitempty"`
	Name              string              `json:"name,omitempty" yaml:"name,omitempty"`
	Description       string              `json:"description,omitempty" yaml:"description,omitempty"`
	Group             *Group              `json:"group,omitempty" yaml:"group,omitempty"`
	AuthorizedKeys    *AuthorizedKeys     `json:"authorized_keys,omitempty" yaml:"authorized_keys,omitempty"`
	BackupPolicies    []*BackupPolicy     `json:"backup_policies,omitempty" yaml:"backup_policies,omitempty"`
	Tags              []string            `json:"tags,omitempty" yaml:"tags,omitempty"`
	ISO               string              `json:"iso,omitempty" yaml:"iso,omitempty"`
}

// JSON returns the build spec in JSON format as a byte slice.
func (s *VirtualMachineSpec) JSON() ([]byte, error) {
	buf := &bytes.Buffer{}
	err := s.WriteJSON(buf)

	return buf.Bytes(), err
}

// JSONIndent returns the build spec in JSON format as a byte slice.
func (s *VirtualMachineSpec) JSONIndent(prefix, indent string) ([]byte, error) {
	buf := &bytes.Buffer{}
	enc := json.NewEncoder(buf)
	enc.SetIndent(prefix, indent)

	err := enc.Encode(s)

	return buf.Bytes(), err
}

// WriteJSON writes the build spec in JSON format to given io.Writer.
func (s *VirtualMachineSpec) WriteJSON(w io.Writer) error {
	enc := json.NewEncoder(w)

	return enc.Encode(s)
}

// XML returns the build spec in XML format as a byte slice.
func (s *VirtualMachineSpec) XML() ([]byte, error) {
	buf := &bytes.Buffer{}
	err := s.WriteXML(buf)

	return buf.Bytes(), err
}

// XMLIndent returns the build spec in XML format as a byte slice.
func (s *VirtualMachineSpec) XMLIndent(prefix, indent string) ([]byte, error) {
	buf := &bytes.Buffer{}
	enc := xml.NewEncoder(buf)
	enc.Indent(prefix, indent)

	err := enc.Encode(s)

	return buf.Bytes(), err
}

// WriteXML writes the build spec in XML format to given io.Writer.
func (s *VirtualMachineSpec) WriteXML(w io.Writer) error {
	_, err := w.Write([]byte(xml.Header))
	if err != nil {
		return err
	}
	enc := xml.NewEncoder(w)

	return enc.Encode(s)
}

// YAML returns the build spec in YAML format as a byte slice.
func (s *VirtualMachineSpec) YAML() ([]byte, error) {
	buf := &bytes.Buffer{}
	err := s.WriteYAML(buf)

	return buf.Bytes(), err
}

// WriteYAML writes the build spec in YAML format to given io.Writer.
func (s *VirtualMachineSpec) WriteYAML(w io.Writer) error {
	enc := yaml.NewEncoder(w)
	enc.SetIndent(2)

	return enc.Encode(s)
}

func (s *VirtualMachineSpec) MarshalXML(
	e *xml.Encoder,
	start xml.StartElement,
) error {
	x := &xmlVirtualMachineSpec{
		Zone:           s.Zone,
		DataCenter:     s.DataCenter,
		Resources:      s.Resources,
		DiskTemplate:   s.DiskTemplate,
		Description:    s.Description,
		Group:          s.Group,
		AuthorizedKeys: s.AuthorizedKeys,
		ISO:            s.ISO,
	}

	if len(s.SystemDisks) > 0 {
		x.SystemDisks = &xmlSystemDisks{
			SystemDisks: s.SystemDisks,
		}
	}

	if len(s.SharedDisks) > 0 {
		x.SharedDisks = &xmlSharedDisks{
			SharedDisks: s.SharedDisks,
		}
	}

	if len(s.NetworkInterfaces) > 0 {
		x.NetworkInterfaces = &xmlNetworkInterfaces{
			NetworkInterfaces: s.NetworkInterfaces,
		}
	}

	hostname := &xmlHostnameValue{Value: s.Hostname}
	if hostname.Value == "" {
		hostname.Type = "random"
	}
	x.Hostname = &xmlHostname{Hostname: hostname}

	if s.Name != "" {
		x.Name = &xmlName{Value: s.Name}
	}

	if len(s.BackupPolicies) > 0 {
		x.BackupPolicies = &xmlBackupPolicies{
			BackupPolicies: s.BackupPolicies,
		}
	}

	if len(s.Tags) > 0 {
		x.Tags = &xmlTags{
			Tags: s.Tags,
		}
	}

	return e.EncodeElement(x, start)
}

func (s *VirtualMachineSpec) UnmarshalXML(
	d *xml.Decoder,
	start xml.StartElement,
) error {
	x := &xmlVirtualMachineSpec{}
	err := d.DecodeElement(x, &start)
	if err != nil {
		return err
	}

	v := VirtualMachineSpec{
		Zone:           x.Zone,
		DataCenter:     x.DataCenter,
		Resources:      x.Resources,
		DiskTemplate:   x.DiskTemplate,
		Description:    x.Description,
		Group:          x.Group,
		AuthorizedKeys: x.AuthorizedKeys,
		ISO:            x.ISO,
	}

	if x.SystemDisks != nil {
		v.SystemDisks = x.SystemDisks.SystemDisks
	}

	if x.SharedDisks != nil {
		v.SharedDisks = x.SharedDisks.SharedDisks
	}

	if x.NetworkInterfaces != nil {
		v.NetworkInterfaces = x.NetworkInterfaces.NetworkInterfaces
	}

	if x.Hostname != nil && x.Hostname.Hostname != nil {
		v.Hostname = x.Hostname.Hostname.Value
	}

	if x.Name != nil {
		v.Name = x.Name.Value
	}

	if x.BackupPolicies != nil {
		v.BackupPolicies = x.BackupPolicies.BackupPolicies
	}

	if x.Tags != nil {
		v.Tags = x.Tags.Tags
	}

	*s = v

	return nil
}

type xmlVirtualMachineSpec struct {
	Zone              *Zone                 `xml:",omitempty"`
	DataCenter        *DataCenter           `xml:",omitempty"`
	Resources         *Resources            `xml:",omitempty"`
	DiskTemplate      *DiskTemplate         `xml:",omitempty"`
	SystemDisks       *xmlSystemDisks       `xml:",omitempty"`
	SharedDisks       *xmlSharedDisks       `xml:",omitempty"`
	NetworkInterfaces *xmlNetworkInterfaces `xml:",omitempty"`
	Hostname          *xmlHostname          `xml:",omitempty"`
	Name              *xmlName              `xml:",omitempty"`
	Description       string                `xml:",omitempty"`
	Group             *Group                `xml:",omitempty"`
	AuthorizedKeys    *AuthorizedKeys       `xml:",omitempty"`
	BackupPolicies    *xmlBackupPolicies    `xml:",omitempty"`
	Tags              *xmlTags              `xml:",omitempty"`
	ISO               string                `xml:",omitempty"`
}
