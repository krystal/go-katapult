package buildspec

import "testing"

func TestAuthorizedKeys_Marshaling(t *testing.T) {
	tests := []struct {
		name    string
		obj     *AuthorizedKeys
		decoded *AuthorizedKeys
	}{
		{
			name: "empty",
			obj:  &AuthorizedKeys{},
		},
		{
			name: "AllUsers",
			obj:  &AuthorizedKeys{AllUsers: true},
		},
		{
			name: "AllSSHKeys",
			obj:  &AuthorizedKeys{AllSSHKeys: true},
		},
		{
			name: "Users",
			obj: &AuthorizedKeys{
				Users: []*User{
					{ID: "user_yUfYcKHgU1ywBWzP"},
					{EmailAddress: "jane@doe.com"},
				},
			},
		},
		{
			name: "AllUsers and Users",
			obj: &AuthorizedKeys{
				AllUsers: true,
				Users: []*User{
					{ID: "user_yUfYcKHgU1ywBWzP"},
					{EmailAddress: "jane@doe.com"},
				},
			},
			decoded: &AuthorizedKeys{AllUsers: true},
		},
		{
			name: "SSHKeys",
			obj: &AuthorizedKeys{
				SSHKeys: []string{
					"ssh-rsa saAfjbp8ADNgGRDi19oKuYVIBTHNC66A",
					"ssh-rsa ouQEux6pnfmyxblOgGQNRpmqIweAO3Qs",
					"ssh-rsa sxyCHXCRBeZliR6EZDmRZ9A82oHkh1nq",
				},
			},
		},
		{
			name: "AllSSHKeys and SSHKeys",
			obj: &AuthorizedKeys{
				AllSSHKeys: true,
				SSHKeys: []string{
					"ssh-rsa saAfjbp8ADNgGRDi19oKuYVIBTHNC66A",
					"ssh-rsa ouQEux6pnfmyxblOgGQNRpmqIweAO3Qs",
					"ssh-rsa sxyCHXCRBeZliR6EZDmRZ9A82oHkh1nq",
				},
			},
			decoded: &AuthorizedKeys{AllSSHKeys: true},
		},
	}
	for _, tt := range tests {
		t.Run("json_"+tt.name, func(t *testing.T) {
			testJSONMarshaling(t, tt.obj)
		})
		t.Run("xml_"+tt.name, func(t *testing.T) {
			testCustomXMLMarshaling(t, tt.obj, tt.decoded)
		})
		t.Run("yaml_"+tt.name, func(t *testing.T) {
			testYAMLMarshaling(t, tt.obj)
		})
	}
}
