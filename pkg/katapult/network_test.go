package katapult

import (
	"context"
	"net/http"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNetwork_JSONMarshaling(t *testing.T) {
	tests := []struct {
		name string
		obj  *Network
	}{
		{
			name: "empty",
			obj:  &Network{},
		},
		{
			name: "full",
			obj: &Network{
				ID:         "netw_zDW7KYAeqqfRfVag",
				Name:       "Public Network",
				Permalink:  "public",
				DataCenter: &DataCenter{ID: "id2"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testJSONMarshaling(t, tt.obj)
		})
	}
}

func TestNetwork_lookupReference(t *testing.T) {
	tests := []struct {
		name string
		obj  *Network
		want *Network
	}{
		{
			name: "nil",
			obj:  nil,
			want: nil,
		},
		{
			name: "empty",
			obj:  &Network{},
			want: &Network{},
		},
		{
			name: "full",
			obj: &Network{
				ID:         "netw_zDW7KYAeqqfRfVag",
				Name:       "Public Network",
				Permalink:  "public",
				DataCenter: &DataCenter{ID: "id2"},
			},
			want: &Network{ID: "netw_zDW7KYAeqqfRfVag"},
		},
		{
			name: "no ID",
			obj: &Network{
				Name:       "Public Network",
				Permalink:  "public",
				DataCenter: &DataCenter{ID: "id2"},
			},
			want: &Network{Permalink: "public"},
		},
		{
			name: "no ID or Permalink",
			obj: &Network{
				Name:       "Public Network",
				DataCenter: &DataCenter{ID: "id2"},
			},
			want: &Network{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.obj.lookupReference()

			assert.Equal(t, tt.want, got)
		})
	}
}

func TestVirtualNetwork_JSONMarshaling(t *testing.T) {
	tests := []struct {
		name string
		obj  *VirtualNetwork
	}{
		{
			name: "empty",
			obj:  &VirtualNetwork{},
		},
		{
			name: "full",
			obj: &VirtualNetwork{
				ID:         "id1",
				Name:       "name",
				DataCenter: &DataCenter{ID: "id2"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testJSONMarshaling(t, tt.obj)
		})
	}
}

func Test_networksResponseBody_JSONMarshaling(t *testing.T) {
	tests := []struct {
		name string
		obj  *networksResponseBody
	}{
		{
			name: "empty",
			obj:  &networksResponseBody{},
		},
		{
			name: "full",
			obj: &networksResponseBody{
				Networks:        []*Network{{ID: "id1"}},
				VirtualNetworks: []*VirtualNetwork{{ID: "id2"}},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testJSONMarshaling(t, tt.obj)
		})
	}
}

func TestNetworksClient_List(t *testing.T) {
	type args struct {
		ctx context.Context
		org *Organization
	}
	tests := []struct {
		name       string
		args       args
		wantNets   []*Network
		wantVnets  []*VirtualNetwork
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
			wantNets: []*Network{
				{
					ID:   "netw_zDW7KYAeqqfRfVag",
					Name: "Public Network",
				},
				{
					ID:   "netw_t7Rbyvr6ahqpDohR",
					Name: "Private Network",
				},
			},
			wantVnets: []*VirtualNetwork{
				{
					ID:   "vnet_1erVCx7A5Y09WknB",
					Name: "Make-Believe Network",
				},
			},
			wantQuery: &url.Values{
				"organization[id]": []string{"org_O648YDMEYeLmqdmn"},
			},
			respStatus: http.StatusOK,
			respBody:   fixture("networks_list"),
		},
		{
			name: "by organization SubDomain",
			args: args{
				ctx: context.Background(),
				org: &Organization{SubDomain: "acme"},
			},
			wantNets: []*Network{
				{
					ID:   "netw_zDW7KYAeqqfRfVag",
					Name: "Public Network",
				},
				{
					ID:   "netw_t7Rbyvr6ahqpDohR",
					Name: "Private Network",
				},
			},
			wantVnets: []*VirtualNetwork{
				{
					ID:   "vnet_1erVCx7A5Y09WknB",
					Name: "Make-Believe Network",
				},
			},
			wantQuery: &url.Values{
				"organization[sub_domain]": []string{"acme"},
			},
			respStatus: http.StatusOK,
			respBody:   fixture("networks_list"),
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
				org: &Organization{ID: "org_nopethisbegone"},
			},
			errStr:     fixtureOrganizationNotFoundErr,
			errResp:    fixtureOrganizationNotFoundResponseError,
			respStatus: http.StatusNotFound,
			respBody:   fixture("organization_not_found_error"),
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
				"/core/v1/organizations/_/available_networks",
				func(w http.ResponseWriter, r *http.Request) {
					assert.Equal(t, "GET", r.Method)
					assertEmptyFieldSpec(t, r)
					assertAuthorization(t, r)

					if tt.wantQuery != nil {
						assert.Equal(t, *tt.wantQuery, r.URL.Query())
					} else {
						assert.Equal(t,
							*tt.args.org.queryValues(), r.URL.Query(),
						)
					}

					w.WriteHeader(tt.respStatus)
					_, _ = w.Write(tt.respBody)
				},
			)

			got1, got2, resp, err := c.Networks.List(tt.args.ctx, tt.args.org)

			if tt.respStatus != 0 {
				assert.Equal(t, tt.respStatus, resp.StatusCode)
			}

			if tt.errStr == "" {
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, err, tt.errStr)
			}

			if tt.wantNets != nil {
				assert.Equal(t, tt.wantNets, got1)
			}

			if tt.wantVnets != nil {
				assert.Equal(t, tt.wantVnets, got2)
			}

			if tt.errResp != nil {
				assert.Equal(t, tt.errResp, resp.Error)
			}
		})
	}
}
