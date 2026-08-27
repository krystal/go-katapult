package buildspec

type CloudInit struct {
	UserData string `xml:",omitempty" json:"user_data,omitempty" yaml:"user_data,omitempty"`
}
