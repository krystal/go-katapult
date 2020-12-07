package buildspec

import "testing"

func Test_xmlSSHKeys_Marshaling(t *testing.T) {
	tests := []struct {
		name string
		obj  *xmlSSHKeys
	}{
		{
			name: "empty",
			obj:  &xmlSSHKeys{},
		},
		{
			name: "full",
			obj: &xmlSSHKeys{
				SSHKeys: []string{
					"ssh-rsa saAfjbp8ADNgGRDi19oKuYVIBTHNC66A",
					"ssh-rsa ouQEux6pnfmyxblOgGQNRpmqIweAO3Qs",
					"ssh-rsa sxyCHXCRBeZliR6EZDmRZ9A82oHkh1nq",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run("xml_"+tt.name, func(t *testing.T) {
			testXMLMarshaling(t, tt.obj)
		})
	}
}
