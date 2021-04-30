package core

import (
	"testing"
)

func TestAttachment_JSONMarshaling(t *testing.T) {
	tests := []struct {
		name string
		obj  *Attachment
	}{
		{
			name: "empty",
			obj:  &Attachment{},
		},
		{
			name: "full",
			obj: &Attachment{
				URL:      "url",
				FileName: "name",
				FileType: "type",
				FileSize: 1234,
				Digest:   "dig",
				Token:    "token",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testJSONMarshaling(t, tt.obj)
		})
	}
}
