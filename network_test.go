package katapult

import (
	"context"
	"fmt"
	"net/http"
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
		ctx   context.Context
		orgID string
	}
	tests := []struct {
		name       string
		args       args
		nets       []*Network
		vnets      []*VirtualNetwork
		errStr     string
		errResp    *ResponseError
		respStatus int
		respBody   []byte
	}{
		{
			name: "networks",
			args: args{
				ctx:   context.Background(),
				orgID: "org_O648YDMEYeLmqdmn",
			},
			nets: []*Network{
				{
					ID:   "netw_zDW7KYAeqqfRfVag",
					Name: "Public Network",
				},
				{
					ID:   "netw_t7Rbyvr6ahqpDohR",
					Name: "Private Network",
				},
			},
			vnets: []*VirtualNetwork{
				{
					ID:   "vnet_1erVCx7A5Y09WknB",
					Name: "Make-Believe Network",
				},
			},
			respStatus: http.StatusOK,
			respBody:   fixture("networks_list"),
		},
		{
			name: "invalid API token response",
			args: args{
				ctx:   context.Background(),
				orgID: "org_O648YDMEYeLmqdmn",
			},
			errStr:     fixtureInvalidAPITokenErr,
			errResp:    fixtureInvalidAPITokenResponseError,
			respStatus: http.StatusForbidden,
			respBody:   fixture("invalid_api_token_error"),
		},
		{
			name: "non-existent organization",
			args: args{
				ctx:   context.Background(),
				orgID: "org_nopethisbegone",
			},
			errStr:     fixtureOrganizationNotFoundErr,
			errResp:    fixtureOrganizationNotFoundResponseError,
			respStatus: http.StatusNotFound,
			respBody:   fixture("organization_not_found_error"),
		},
		{
			name: "suspended organization",
			args: args{
				ctx:   context.Background(),
				orgID: "org_O648YDMEYeLmqdmn",
			},
			errStr:     fixtureOrganizationSuspendedErr,
			errResp:    fixtureOrganizationSuspendedResponseError,
			respStatus: http.StatusForbidden,
			respBody:   fixture("organization_suspended_error"),
		},
		{
			name: "nil context",
			args: args{
				ctx:   nil,
				orgID: "org_O648YDMEYeLmqdmn",
			},
			errStr: "net/http: nil Context",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, mux, _, teardown := prepareTestClient()
			defer teardown()

			mux.HandleFunc(
				fmt.Sprintf(
					"/core/v1/organizations/%s/available_networks",
					tt.args.orgID,
				),
				func(w http.ResponseWriter, r *http.Request) {
					assert.Equal(t, "GET", r.Method)
					assert.Equal(t, "", r.Header.Get("X-Field-Spec"))

					w.WriteHeader(tt.respStatus)
					_, _ = w.Write(tt.respBody)
				},
			)

			got1, got2, resp, err := c.Networks.List(tt.args.ctx, tt.args.orgID)

			if tt.respStatus != 0 {
				assert.Equal(t, tt.respStatus, resp.StatusCode)
			}

			if tt.errStr == "" {
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, err, tt.errStr)
			}

			if tt.nets != nil {
				assert.Equal(t, tt.nets, got1)
			}

			if tt.vnets != nil {
				assert.Equal(t, tt.vnets, got2)
			}

			if tt.errResp != nil {
				assert.Equal(t, tt.errResp, resp.Error)
			}
		})
	}
}
