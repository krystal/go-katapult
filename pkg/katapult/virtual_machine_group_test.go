package katapult

import (
	"context"
	"net/http"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestVirtualMachineGroup_JSONMarshaling(t *testing.T) {
	tests := []struct {
		name string
		obj  *VirtualMachineGroup
	}{
		{
			name: "empty",
			obj:  &VirtualMachineGroup{},
		},
		{
			name: "full",
			obj: &VirtualMachineGroup{
				ID:        "id",
				Name:      "name",
				Segregate: true,
				CreatedAt: timestampPtr(934834834),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testJSONMarshaling(t, tt.obj)
		})
	}
}

func TestVirtualMachineGroupsClient_List(t *testing.T) {
	// Correlates to fixtures/virtual_machine_groups_list*.json
	virtualMachineGroupList := []*VirtualMachineGroup{
		{
			ID:        "vmgrp_gsEUFPp3ybVQm5QQ",
			Name:      "vm group 1",
			Segregate: true,
		},
		{
			ID:        "vmgrp_bcfdEFn2viWCm5ve",
			Name:      "vm group 2",
			Segregate: false,
		},
	}

	type args struct {
		ctx context.Context
		org *Organization
	}
	tests := []struct {
		name       string
		args       args
		want       []*VirtualMachineGroup
		wantQuery  *url.Values
		errStr     string
		errResp    *ResponseError
		respStatus int
		respBody   []byte
	}{
		{
			name: "by organization ID",
			args: args{
				ctx: context.Background(),
				org: &Organization{ID: "org_O648YDMEYeLmqdmn"},
			},
			want: virtualMachineGroupList,
			wantQuery: &url.Values{
				"organization[id]": []string{"org_O648YDMEYeLmqdmn"},
			},
			respStatus: http.StatusOK,
			respBody:   fixture("virtual_machine_groups_list"),
		},
		{
			name: "by organization SubDomain",
			args: args{
				ctx: context.Background(),
				org: &Organization{SubDomain: "acme"},
			},
			want: virtualMachineGroupList,
			wantQuery: &url.Values{
				"organization[sub_domain]": []string{"acme"},
			},
			respStatus: http.StatusOK,
			respBody:   fixture("virtual_machine_groups_list"),
		},
		{
			name: "invalid API token response",
			args: args{
				ctx: context.Background(),
				org: &Organization{ID: "org_O648YDMEYeLmqdmn"},
			},
			errStr:     fixtureInvalidAPITokenErr,
			errResp:    fixtureInvalidAPITokenResponseError,
			respStatus: http.StatusForbidden,
			respBody:   fixture("invalid_api_token_error"),
		},
		{
			name: "non-existent organization",
			args: args{
				ctx: context.Background(),
				org: &Organization{ID: "org_O648YDMEYeLmqdmn"},
			},
			errStr:     fixtureOrganizationNotFoundErr,
			errResp:    fixtureOrganizationNotFoundResponseError,
			respStatus: http.StatusNotFound,
			respBody:   fixture("organization_not_found_error"),
		},
		{
			name: "nil organization",
			args: args{
				ctx: context.Background(),
				org: nil,
			},
			errStr:     fixtureOrganizationNotFoundErr,
			errResp:    fixtureOrganizationNotFoundResponseError,
			respStatus: http.StatusNotFound,
			respBody:   fixture("organization_not_found_error"),
		},
		{
			name: "not activated organization",
			args: args{
				ctx: context.Background(),
				org: &Organization{ID: "org_O648YDMEYeLmqdmn"},
			},
			errStr:     fixtureOrganizationNotActivatedErr,
			errResp:    fixtureOrganizationNotActivatedResponseError,
			respStatus: http.StatusForbidden,
			respBody:   fixture("organization_not_activated_error"),
		},
		{
			name: "suspended organization",
			args: args{
				ctx: context.Background(),
				org: &Organization{ID: "org_O648YDMEYeLmqdmn"},
			},
			errStr:     fixtureOrganizationSuspendedErr,
			errResp:    fixtureOrganizationSuspendedResponseError,
			respStatus: http.StatusForbidden,
			respBody:   fixture("organization_suspended_error"),
		},
		{
			name: "permission denied",
			args: args{
				ctx: context.Background(),
				org: &Organization{ID: "org_O648YDMEYeLmqdmn"},
			},
			errStr:     fixturePermissionDeniedErr,
			errResp:    fixturePermissionDeniedResponseError,
			respStatus: http.StatusForbidden,
			respBody:   fixture("permission_denied_error"),
		},
		{
			name: "nil context",
			args: args{
				ctx: nil,
				org: &Organization{ID: "org_O648YDMEYeLmqdmn"},
			},
			errStr: "net/http: nil Context",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, mux, _, teardown := prepareTestClient()
			defer teardown()

			mux.HandleFunc(
				"/core/v1/organizations/_/virtual_machine_groups",
				func(w http.ResponseWriter, r *http.Request) {
					assert.Equal(t, "GET", r.Method)
					assertEmptyFieldSpec(t, r)
					assertAuthorization(t, r)

					if tt.wantQuery != nil {
						assert.Equal(t, *tt.wantQuery, r.URL.Query())
					} else {
						qs := queryValues(tt.args.org)
						assert.Equal(t, *qs, r.URL.Query())
					}

					w.WriteHeader(tt.respStatus)
					_, _ = w.Write(tt.respBody)
				},
			)

			got, resp, err := c.VirtualMachineGroups.List(
				tt.args.ctx, tt.args.org,
			)

			if tt.respStatus != 0 {
				assert.Equal(t, tt.respStatus, resp.StatusCode)
			}

			if tt.errStr == "" {
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, err, tt.errStr)
			}

			if tt.want != nil {
				assert.Equal(t, tt.want, got)
			}

			if tt.errResp != nil {
				assert.Equal(t, tt.errResp, resp.Error)
			}
		})
	}
}
