package buildspec

type Resources struct {
	Package  *Package `xml:",omitempty" json:"package,omitempty" yaml:"package,omitempty"`
	Memory   int      `xml:",omitempty" json:"memory,omitempty" yaml:"memory,omitempty"`
	CPUCores int      `xml:",omitempty" json:"cpu_cores,omitempty" yaml:"cpu_cores,omitempty"`
}
