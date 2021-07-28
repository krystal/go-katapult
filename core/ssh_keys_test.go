package core

import "testing"

func Test_AuthSSHKey_JSONMarshaling(t *testing.T) {
	tests := []struct {
		name string
		obj  *AuthSSHKey
	}{
		{
			name: "empty",
			obj:  &AuthSSHKey{},
		},
		{
			name: "full",
			obj: &AuthSSHKey{
				ID:          "a",
				Name:        "b",
				Fingerprint: "c",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testJSONMarshaling(t, tt.obj)
		})
	}
}

func Test_sshKeyResponseBody_JSONMarshaling(t *testing.T) {
	tests := []struct {
		name string
		obj  *sshKeysResponseBody
	}{
		{
			name: "empty",
			obj:  &sshKeysResponseBody{},
		},
		{
			name: "ssh_key",
			obj: &sshKeysResponseBody{
				SSHKey: &AuthSSHKey{},
			},
		},
		{
			name: "ssh_keys",
			obj: &sshKeysResponseBody{
				SSHKeys: []*AuthSSHKey{{}},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testJSONMarshaling(t, tt.obj)
		})
	}
}

func Test_AuthSSHKeyProperties_JSONMarshaling(t *testing.T) {
	tests := []struct {
		name string
		obj  *AuthSSHKey
	}{
		{
			name: "empty",
			obj:  &AuthSSHKey{},
		},
		{
			name: "full",
			obj: &AuthSSHKey{
				ID:          "a",
				Name:        "b",
				Fingerprint: "c",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testJSONMarshaling(t, tt.obj)
		})
	}
}
