package katapult

type Attachment struct {
	URL      string `json:"url,omitempty"`
	FileName string `json:"file_name,omitempty"`
	FileType string `json:"file_type,omitempty"`
	FileSize int64  `json:"file_size,omitempty"`
	Digest   string `json:"digest,omitempty"`
	Token    string `json:"token,omitempty"`
}
