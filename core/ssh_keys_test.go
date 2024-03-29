package core

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"testing"

	"github.com/krystal/go-katapult"
	"github.com/krystal/go-katapult/internal/test"
	"github.com/krystal/go-katapult/internal/testclient"
	"github.com/stretchr/testify/assert"
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
			name: "full",
			obj: &sshKeysResponseBody{
				Pagination: &katapult.Pagination{
					CurrentPage: 33,
				},
				SSHKey:  &AuthSSHKey{ID: "1"},
				SSHKeys: []*AuthSSHKey{{ID: "2"}},
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
		obj  *AuthSSHKeyProperties
	}{
		{
			name: "empty",
			obj:  &AuthSSHKeyProperties{},
		},
		{
			name: "full",
			obj: &AuthSSHKeyProperties{
				Name: "a",
				Key:  "b",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testJSONMarshaling(t, tt.obj)
		})
	}
}

func Test_SSHKeyRef_JSONMarshaling(t *testing.T) {
	tests := []struct {
		name string
		obj  *SSHKeyRef
	}{
		{
			name: "empty",
			obj:  &SSHKeyRef{},
		},
		{
			name: "full",
			obj:  &SSHKeyRef{ID: "a"},
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
		opts *ListOptions
	}
	tests := []struct {
		name     string
		args     args
		resp     *katapult.Response
		respErr  error
		respV    *sshKeysResponseBody
		want     []*AuthSSHKey
		wantReq  *katapult.Request
		wantResp *katapult.Response
		wantErr  string
	}{
		{
			name: "success",
			args: args{
				ctx: context.Background(),
				org: OrganizationRef{ID: "org_O648YDMEYeLmqdmn"},
				opts: &ListOptions{
					Page:    5,
					PerPage: 32,
				},
			},
			resp: &katapult.Response{
				Response: &http.Response{StatusCode: http.StatusOK},
			},
			respV: &sshKeysResponseBody{
				Pagination: &katapult.Pagination{Total: 394},
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
						"page":     []string{"5"},
						"per_page": []string{"32"},
						"organization[id]": []string{
							"org_O648YDMEYeLmqdmn",
						},
					}.Encode(),
				},
			},
			wantResp: &katapult.Response{
				Pagination: &katapult.Pagination{Total: 394},
				Response:   &http.Response{StatusCode: http.StatusOK},
			},
		},
		{
			name: "success with nil options",
			args: args{
				ctx: context.Background(),
				org: OrganizationRef{ID: "org_O648YDMEYeLmqdmn"},
			},
			resp: &katapult.Response{
				Response: &http.Response{StatusCode: http.StatusOK},
			},
			respV: &sshKeysResponseBody{
				Pagination: &katapult.Pagination{Total: 333},
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
			wantResp: &katapult.Response{
				Pagination: &katapult.Pagination{Total: 333},
				Response:   &http.Response{StatusCode: http.StatusOK},
			},
		},
		{
			name: "request error",
			args: args{
				ctx: context.Background(),
				org: OrganizationRef{ID: "org_O648YDMEYeLmqdmn"},
			},
			resp: &katapult.Response{
				Response: &http.Response{
					StatusCode: http.StatusInternalServerError,
				},
			},
			respErr: fmt.Errorf("flux capacitor undercharged"),
			wantResp: &katapult.Response{
				Response: &http.Response{
					StatusCode: http.StatusInternalServerError,
				},
			},
			wantErr: "flux capacitor undercharged",
		},
		{
			name: "request error with nil response",
			args: args{
				ctx: context.Background(),
				org: OrganizationRef{ID: "org_O648YDMEYeLmqdmn"},
			},
			resp:    nil,
			respErr: fmt.Errorf("something is really wrong"),
			wantErr: "something is really wrong",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tc := testclient.New(tt.resp, tt.respErr, tt.respV)
			c := NewSSHKeysClient(tc)
			ctx := test.Context(tt.args.ctx)

			got, resp, err := c.List(
				ctx,
				tt.args.org,
				tt.args.opts,
				testRequestOption,
			)

			assert.Equal(t, 1, len(tc.Calls), "only 1 request should be made")
			test.AssertContext(t, ctx, tc.Ctx)

			assert.Equal(t, tt.want, got)

			if tt.wantResp != nil {
				assert.Equal(t, tt.wantResp, resp)
			}

			if tt.wantReq != nil {
				setWantRequestOptionHeader(tt.wantReq)
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

const sshPublicKey = "ssh-rsa " +
	"AAAAB3NzaC1yc2EAAAADAQABAAABgQC7MZzjBzFWsc6BCcYE2" +
	"EpSo8DOjzhDPb/WndW6QE/G0xM7iqdlezcmQnL3Gw9jtAOI4O" +
	"lNok19v4q8C6ham+1WbX2aGd2labOmKoBVWXIzKyFz9pg2Rs1" +
	"0ZGn+Ly+uJF558rSehSvGJPFmKUagYeBG9c/cwuVMzube0yVb" +
	"tH2CWRs2dMvwhloH5zOh3NMQj/5uBGYMh9uRQKsGHoG8TET08" +
	"VSok3W/CFilSH7jSmaQYziUqJjOLE2hb8ziCzfv/0GhbY5MoJ" +
	"JUZqUdOlGkYgDMR/IVOHxxF93QBvp1AkAzh8RBsvJPajgZHFa" +
	"1lWYJRP7U4TREWuxkpaJrbK3I3AHM74GAfIq76wndoFYJhi5q" +
	"bNgaJjLUJDPPzl8KOcp0Pb5FPqygHWz/K4n1h5SV/LdD0mB48" +
	"7TxeC1NV4XBQQruM5RgfTXSWBW+8W83U0y5h1RNl/Qo9Efo7K" +
	"yc25wCxVT2cWRHr3mxZ98p+JxmFmC1KTdUrM95+B7+Hw9fKYv" +
	"hKz0= user@machine.local"

const sshFingerprint = "SHA256:Ybk7/sbyptVqD87piCCz/XHi" +
	"EKrdvHND2EMDA1qGqRA user@machine.local"

func TestSSHKeysClient_Add(t *testing.T) {
	type args struct {
		ctx  context.Context
		org  OrganizationRef
		args AuthSSHKeyProperties
	}
	tests := []struct {
		name     string
		args     args
		respErr  error
		resp     *katapult.Response
		respV    *sshKeysResponseBody
		want     *AuthSSHKey
		wantReq  *katapult.Request
		wantResp *katapult.Response
		wantErr  string
	}{
		{
			name: "success",
			args: args{
				ctx: context.Background(),
				org: OrganizationRef{ID: "org_O648YDMEYeLmqdmn"},
				args: AuthSSHKeyProperties{
					Name: "test",
					Key:  sshPublicKey,
				},
			},
			resp: &katapult.Response{
				Response: &http.Response{StatusCode: http.StatusOK},
			},
			respV: &sshKeysResponseBody{
				SSHKey: &AuthSSHKey{
					ID:          "testing-id",
					Name:        "test",
					Fingerprint: sshFingerprint,
				},
			},
			want: &AuthSSHKey{
				ID:          "testing-id",
				Name:        "test",
				Fingerprint: sshFingerprint,
			},
			wantReq: &katapult.Request{
				Method: "POST",
				URL: &url.URL{
					Path: "/core/v1/organizations/_/ssh_keys",
					RawQuery: url.Values{
						"organization[id]": []string{
							"org_O648YDMEYeLmqdmn",
						},
					}.Encode(),
				},
				Body: AuthSSHKeyProperties{
					Name: "test",
					Key:  sshPublicKey,
				},
			},
			wantResp: &katapult.Response{
				Response: &http.Response{StatusCode: http.StatusOK},
			},
		},
		{
			name: "request error",
			args: args{
				ctx: context.Background(),
				org: OrganizationRef{ID: "org_O648YDMEYeLmqdmn"},
			},
			resp: &katapult.Response{
				Response: &http.Response{
					StatusCode: http.StatusInternalServerError,
				},
			},
			respErr: fmt.Errorf("flux capacitor undercharged"),
			wantResp: &katapult.Response{
				Response: &http.Response{
					StatusCode: http.StatusInternalServerError,
				},
			},
			wantErr: "flux capacitor undercharged",
		},
		{
			name: "request error with nil response",
			args: args{
				ctx: context.Background(),
				org: OrganizationRef{ID: "org_O648YDMEYeLmqdmn"},
			},
			resp:    nil,
			respErr: fmt.Errorf("something is really wrong"),
			wantErr: "something is really wrong",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tc := testclient.New(tt.resp, tt.respErr, tt.respV)
			c := NewSSHKeysClient(tc)
			ctx := test.Context(tt.args.ctx)

			got, resp, err := c.Add(
				ctx,
				tt.args.org,
				tt.args.args,
				testRequestOption,
			)

			assert.Equal(t, 1, len(tc.Calls), "only 1 request should be made")
			test.AssertContext(t, ctx, tc.Ctx)

			assert.Equal(t, tt.want, got)

			if tt.wantResp != nil {
				assert.Equal(t, tt.wantResp, resp)
			}

			if tt.wantReq != nil {
				setWantRequestOptionHeader(tt.wantReq)
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

func TestSSHKeysClient_Delete(t *testing.T) {
	type args struct {
		ctx context.Context
		ref SSHKeyRef
	}
	tests := []struct {
		name     string
		args     args
		resp     *katapult.Response
		respErr  error
		respV    *sshKeysResponseBody
		want     *AuthSSHKey
		wantReq  *katapult.Request
		wantResp *katapult.Response
		wantErr  string
	}{
		{
			name: "success",
			args: args{
				ctx: context.Background(),
				ref: SSHKeyRef{ID: "ssh_O574YEEEYeLmqdmn"},
			},
			respV: &sshKeysResponseBody{
				SSHKey: &AuthSSHKey{
					ID:          "testing-id",
					Name:        "test",
					Fingerprint: sshFingerprint,
				},
			},
			resp: &katapult.Response{
				Response: &http.Response{StatusCode: http.StatusOK},
			},
			want: &AuthSSHKey{
				ID:          "testing-id",
				Name:        "test",
				Fingerprint: sshFingerprint,
			},
			wantReq: &katapult.Request{
				Method: "DELETE",
				URL: &url.URL{
					Path: "/core/v1/ssh_keys/_",
					RawQuery: url.Values{
						"ssh_key[id]": []string{
							"ssh_O574YEEEYeLmqdmn",
						},
					}.Encode(),
				},
			},
			wantResp: &katapult.Response{
				Response: &http.Response{StatusCode: http.StatusOK},
			},
		},
		{
			name: "request error",
			args: args{
				ctx: context.Background(),
			},
			resp: &katapult.Response{
				Response: &http.Response{
					StatusCode: http.StatusInternalServerError,
				},
			},
			respErr: fmt.Errorf("flux capacitor undercharged"),
			wantResp: &katapult.Response{
				Response: &http.Response{
					StatusCode: http.StatusInternalServerError,
				},
			},
			wantErr: "flux capacitor undercharged",
		},
		{
			name: "request error with nil response",
			args: args{
				ctx: context.Background(),
			},
			resp:    nil,
			respErr: fmt.Errorf("something is really wrong"),
			wantErr: "something is really wrong",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tc := testclient.New(tt.resp, tt.respErr, tt.respV)
			c := NewSSHKeysClient(tc)
			ctx := test.Context(tt.args.ctx)

			got, resp, err := c.Delete(ctx, tt.args.ref, testRequestOption)

			assert.Equal(t, 1, len(tc.Calls), "only 1 request should be made")
			test.AssertContext(t, ctx, tc.Ctx)

			assert.Equal(t, tt.want, got)

			if tt.wantResp != nil {
				assert.Equal(t, tt.wantResp, resp)
			}

			if tt.wantReq != nil {
				setWantRequestOptionHeader(tt.wantReq)
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

func TestClient_SSHKeys(t *testing.T) {
	c := New(&testclient.Client{})

	assert.IsType(t, &SSHKeysClient{}, c.SSHKeys)
}
