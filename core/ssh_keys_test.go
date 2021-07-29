package core

import (
	"context"
	"fmt"
	"github.com/krystal/go-katapult"
	"github.com/krystal/go-katapult/internal/test"
	"github.com/krystal/go-katapult/internal/testclient"
	"github.com/stretchr/testify/assert"
	"net/url"
	"testing"
)

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

func TestNewSSHKeysClient(t *testing.T) {
	tc := testclient.New(nil, nil, nil)
	c := NewSSHKeysClient(tc)
	assert.Equal(t, tc, c.client)
	assert.Equal(t, &url.URL{Path: "/core/v1/"}, c.basePath)
}

func TestSSHKeysClient_List(t *testing.T) {
	type args struct {
		ctx  context.Context
		org  OrganizationRef
	}
	tests := []struct {
		name    string
		args    args
		resp    *katapult.Response
		respErr error
		respV   *sshKeysResponseBody
		want    []*AuthSSHKey
		wantReq *katapult.Request
		wantErr string
	}{
		{
			name: "success",
			args: args{
				ctx: context.Background(),
				org: OrganizationRef{ID: "org_O648YDMEYeLmqdmn"},
			},
			resp: &katapult.Response{
				Pagination: &katapult.Pagination{Total: 333},
			},
			respV: &sshKeysResponseBody{
				SSHKeys: []*AuthSSHKey{
					{ID: "ssh_O574YEEEYeLmqdmn"},
				},
			},
			want: []*AuthSSHKey{
				{ID: "ssh_O574YEEEYeLmqdmn"},
			},
			wantReq: &katapult.Request{
				Method: "GET",
				URL: &url.URL{
					Path: "/core/v1/organizations/_/ssh_keys",
					RawQuery: url.Values{
						"organization[id]": []string{
							"org_O648YDMEYeLmqdmn",
						},
					}.Encode(),
				},
			},
		},
		{
			name: "request error",
			args: args{
				ctx:  context.Background(),
				org:  OrganizationRef{ID: "org_O648YDMEYeLmqdmn"},
			},
			respErr: fmt.Errorf("flux capacitor undercharged"),
			wantErr: "flux capacitor undercharged",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tc := testclient.New(tt.resp, tt.respErr, tt.respV)
			c := NewSSHKeysClient(tc)
			ctx := test.Context(tt.args.ctx)

			got, resp, err := c.List(ctx, tt.args.org)

			assert.Equal(t, 1, len(tc.Calls), "only 1 request should be made")
			test.AssertContext(t, ctx, tc.Ctx)

			assert.Equal(t, tt.want, got)

			if tt.resp != nil {
				assert.Equal(t, tt.resp, resp)
			}

			if tt.wantReq != nil {
				assert.Equal(t, tt.wantReq, tc.Request)
			}

			if tt.wantErr == "" {
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, err, tt.wantErr)
			}
		})
	}
}
