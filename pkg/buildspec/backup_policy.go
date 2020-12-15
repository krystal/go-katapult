package buildspec

type BackupPolicy struct {
	Retention int       `xml:",omitempty" json:"retention,omitempty" yaml:"retention,omitempty"`
	Schedule  *Schedule `xml:",omitempty" json:"schedule,omitempty" yaml:"schedule,omitempty"`
}

type xmlBackupPolicies struct {
	BackupPolicies []*BackupPolicy `xml:"BackupPolicy,omitempty"`
}
