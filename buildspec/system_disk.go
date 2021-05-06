package buildspec

import "encoding/xml"

type SystemDisk struct {
	Name           string          `json:"name,omitempty" yaml:"name,omitempty"`
	Size           int             `json:"size,omitempty" yaml:"size,omitempty"`
	Speed          string          `json:"speed,omitempty" yaml:"speed,omitempty"`
	IOProfile      *DiskIOProfile  `json:"io_profile,omitempty" yaml:"io_profile,omitempty"`
	FileSystemType string          `json:"file_system_type,omitempty" yaml:"file_system_type,omitempty"`
	BackupPolicies []*BackupPolicy `json:"backup_policies,omitempty" yaml:"backup_policies,omitempty"`
}

func (s *SystemDisk) MarshalXML(
	e *xml.Encoder,
	start xml.StartElement,
) error {
	x := &xmlSystemDisk{
		Name:           s.Name,
		Size:           s.Size,
		Speed:          s.Speed,
		IOProfile:      s.IOProfile,
		FileSystemType: s.FileSystemType,
	}

	if len(s.BackupPolicies) > 0 {
		x.BackupPolicies = &xmlBackupPolicies{
			BackupPolicies: s.BackupPolicies,
		}
	}

	return e.EncodeElement(x, start)
}

func (s *SystemDisk) UnmarshalXML(
	d *xml.Decoder,
	start xml.StartElement,
) error {
	x := &xmlSystemDisk{}
	_ = d.DecodeElement(x, &start)

	v := SystemDisk{
		Name:           x.Name,
		Size:           x.Size,
		Speed:          x.Speed,
		IOProfile:      x.IOProfile,
		FileSystemType: x.FileSystemType,
	}

	if x.BackupPolicies != nil {
		v.BackupPolicies = x.BackupPolicies.BackupPolicies
	}

	*s = v

	return nil
}

type xmlSystemDisk struct {
	Name           string             `xml:",omitempty"`
	Size           int                `xml:",omitempty"`
	Speed          string             `xml:",omitempty"`
	IOProfile      *DiskIOProfile     `xml:",omitempty"`
	FileSystemType string             `xml:",omitempty"`
	BackupPolicies *xmlBackupPolicies `xml:",omitempty"`
}

type xmlSystemDisks struct {
	SystemDisks []*SystemDisk `xml:"Disk,omitempty"`
}
