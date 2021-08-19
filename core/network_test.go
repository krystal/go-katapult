package core

import (
	"context"
	"encoding/json"
	"net/http"
	"net/url"
	"testing"

	"github.com/krystal/go-katapult"
	"github.com/krystal/go-katapult/internal/testclient"
	"github.com/stretchr/testify/assert"
)

var (
	fixtureNetworkNotFoundErr = "katapult: not_found: network_not_found: No " +
		"network was found matching any of the criteria provided in the " +
		"arguments"
	fixtureNetworkNotFoundResponseError = &katapult.ResponseError{
		Code: "network_not_found",
		Description: "No network was found matching any of the criteria " +
			"provided in the arguments",
		Detail: json.RawMessage(`{}`),
	}
)

func TestClient_Networks(t *testing.T) {
	c := New(&testclient.Client{})

	assert.IsType(t, &NetworksClient{}, c.Networks)
}

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

func TestNetwork_Ref(t *testing.T) {
	tests := []struct {
		name string
		obj  Network
		want NetworkRef
	}{
		{
			name: "empty",
			obj:  Network{},
			want: NetworkRef{},
		},
		{
			name: "full",
			obj: Network{
				ID:        "netw_zDW7KYAeqqfRfVag",
				Permalink: "public",
			},
			want: NetworkRef{
				ID: "netw_zDW7KYAeqqfRfVag",
			},
		},
		{
			name: "just ID",
			obj: Network{
				ID: "netw_zDW7KYAeqqfRfVag",
			},
			want: NetworkRef{
				ID: "netw_zDW7KYAeqqfRfVag",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.obj.Ref()

			assert.Equal(t, tt.want, got)
		})
	}
}

func TestNetworkRef_queryValues(t *testing.T) {
	tests := []struct {
		name string
		obj  NetworkRef
	}{
		{
			name: "empty",
			obj:  NetworkRef{},
		},
		{
			name: "full",
			obj: NetworkRef{
				ID:        "netw_zDW7KYAeqqfRfVag",
				Permalink: "public",
			},
		},
		{
			name: "just ID",
			obj: NetworkRef{
				ID: "netw_zDW7KYAeqqfRfVag",
			},
		},
		{
			name: "just Permalink",
			obj: NetworkRef{
				Permalink: "public",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testQueryableEncoding(t, tt.obj)
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
				Network:         &Network{ID: "id1"},
				Networks:        []*Network{{ID: "id2"}},
				VirtualNetworks: []*VirtualNetwork{{ID: "id3"}},
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
		org OrganizationRef
	}
	tests := []struct {
		name       string
		args       args
		wantNets   []*Network
		wantVnets  []*VirtualNetwork
		errStr     string
		errResp    *katapult.ResponseError
		errIs      error
		respStatus int
		respBody   []byte
	}{
		{
			name: "by organization ID",
			args: args{
				ctx: context.Background(),
				org: OrganizationRef{ID: "org_O648YDMEYeLmqdmn"},
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
			respStatus: http.StatusOK,
			respBody:   fixture("networks_list"),
		},
		{
			name: "by organization SubDomain",
			args: args{
				ctx: context.Background(),
				org: OrganizationRef{SubDomain: "acme"},
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
			respStatus: http.StatusOK,
			respBody:   fixture("networks_list"),
		},
		{
			name: "invalid API token response",
			args: args{
				ctx: context.Background(),
				org: OrganizationRef{ID: "org_O648YDMEYeLmqdmn"},
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
				org: OrganizationRef{ID: "org_nopethisbegone"},
			},
			errStr:     fixtureOrganizationNotFoundErr,
			errResp:    fixtureOrganizationNotFoundResponseError,
			errIs:      ErrOrganizationNotFound,
			respStatus: http.StatusNotFound,
			respBody:   fixture("organization_not_found_error"),
		},
		{
			name: "suspended organization",
			args: args{
				ctx: context.Background(),
				org: OrganizationRef{ID: "org_O648YDMEYeLmqdmn"},
			},
			errStr:     fixtureOrganizationSuspendedErr,
			errResp:    fixtureOrganizationSuspendedResponseError,
			errIs:      ErrOrganizationSuspended,
			respStatus: http.StatusForbidden,
			respBody:   fixture("organization_suspended_error"),
		},
		{
			name: "nil context",
			args: args{
				ctx: nil,
				org: OrganizationRef{ID: "org_O648YDMEYeLmqdmn"},
			},
			errStr: "net/http: nil Context",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rm, mux, _, teardown := prepareTestClient(t)
			defer teardown()
			c := NewNetworksClient(rm)

			mux.HandleFunc(
				"/core/v1/organizations/_/available_networks",
				func(w http.ResponseWriter, r *http.Request) {
					assert.Equal(t, "GET", r.Method)
					assertEmptyFieldSpec(t, r)
					assertAuthorization(t, r)
					assertRequestOptionHeader(t, r)

					assert.Equal(t,
						*tt.args.org.queryValues(), r.URL.Query(),
					)

					w.WriteHeader(tt.respStatus)
					_, _ = w.Write(tt.respBody)
				},
			)

			got1, got2, resp, err := c.List(
				tt.args.ctx, tt.args.org, testRequestOption,
			)

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

			if tt.errIs != nil {
				assert.ErrorIs(t, err, tt.errIs)
			}
		})
	}
}

func TestNetworksClient_Get(t *testing.T) {
	// Correlates to fixtures/network_get.json
	network := &Network{
		ID:        "netw_zDW7KYAeqqfRfVag",
		Name:      "Public Network",
		Permalink: "public",
	}

	type args struct {
		ctx context.Context
		ref NetworkRef
	}
	tests := []struct {
		name       string
		args       args
		want       *Network
		wantQuery  *url.Values
		errStr     string
		errResp    *katapult.ResponseError
		errIs      error
		respStatus int
		respBody   []byte
	}{
		{
			name: "by ID",
			args: args{
				ctx: context.Background(),
				ref: NetworkRef{ID: "netw_zDW7KYAeqqfRfVag"},
			},
			want: network,
			wantQuery: &url.Values{
				"network[id]": []string{"netw_zDW7KYAeqqfRfVag"},
			},
			respStatus: http.StatusOK,
			respBody:   fixture("network_get"),
		},
		{
			name: "by Permalink",
			args: args{
				ctx: context.Background(),
				ref: NetworkRef{Permalink: "public"},
			},
			want: network,
			wantQuery: &url.Values{
				"network[permalink]": []string{"public"},
			},
			respStatus: http.StatusOK,
			respBody:   fixture("network_get"),
		},
		{
			name: "non-existent network by ID",
			args: args{
				ctx: context.Background(),
				ref: NetworkRef{ID: "netw_nopethisbegone"},
			},
			errStr:     fixtureNetworkNotFoundErr,
			errResp:    fixtureNetworkNotFoundResponseError,
			errIs:      ErrNetworkNotFound,
			respStatus: http.StatusNotFound,
			respBody:   fixture("network_not_found_error"),
		},
		{
			name: "non-existent network by Permalink",
			args: args{
				ctx: context.Background(),
				ref: NetworkRef{Permalink: "public"},
			},
			errStr:     fixtureNetworkNotFoundErr,
			errResp:    fixtureNetworkNotFoundResponseError,
			errIs:      ErrNetworkNotFound,
			respStatus: http.StatusNotFound,
			respBody:   fixture("network_not_found_error"),
		},
		{
			name: "nil context",
			args: args{
				ctx: nil,
				ref: NetworkRef{ID: "netw_zDW7KYAeqqfRfVag"},
			},
			errStr: "net/http: nil Context",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rm, mux, _, teardown := prepareTestClient(t)
			defer teardown()
			c := NewNetworksClient(rm)

			mux.HandleFunc(
				"/core/v1/networks/_",
				func(w http.ResponseWriter, r *http.Request) {
					assert.Equal(t, "GET", r.Method)
					assertEmptyFieldSpec(t, r)
					assertAuthorization(t, r)
					assertRequestOptionHeader(t, r)

					if tt.wantQuery != nil {
						assert.Equal(t, *tt.wantQuery, r.URL.Query())
					}

					w.WriteHeader(tt.respStatus)
					_, _ = w.Write(tt.respBody)
				},
			)

			got, resp, err := c.Get(
				tt.args.ctx, tt.args.ref, testRequestOption,
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

			if tt.errIs != nil {
				assert.ErrorIs(t, err, tt.errIs)
			}
		})
	}
}

func TestNetworksClient_GetByID(t *testing.T) {
	// Correlates to fixtures/network_get.json
	network := &Network{
		ID:        "netw_zDW7KYAeqqfRfVag",
		Name:      "Public Network",
		Permalink: "public",
	}

	type args struct {
		ctx context.Context
		id  string
	}
	tests := []struct {
		name       string
		args       args
		want       *Network
		errStr     string
		errResp    *katapult.ResponseError
		errIs      error
		respStatus int
		respBody   []byte
	}{
		{
			name: "network",
			args: args{
				ctx: context.Background(),
				id:  "netw_zDW7KYAeqqfRfVag",
			},
			want:       network,
			respStatus: http.StatusOK,
			respBody:   fixture("network_get"),
		},
		{
			name: "non-existent network",
			args: args{
				ctx: context.Background(),
				id:  "netw_nopethisbegone",
			},
			errStr:     fixtureNetworkNotFoundErr,
			errResp:    fixtureNetworkNotFoundResponseError,
			errIs:      ErrNetworkNotFound,
			respStatus: http.StatusNotFound,
			respBody:   fixture("network_not_found_error"),
		},
		{
			name: "nil context",
			args: args{
				ctx: nil,
				id:  "netw_zDW7KYAeqqfRfVag",
			},
			errStr: "net/http: nil Context",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rm, mux, _, teardown := prepareTestClient(t)
			defer teardown()
			c := NewNetworksClient(rm)

			mux.HandleFunc("/core/v1/networks/_",
				func(w http.ResponseWriter, r *http.Request) {
					assert.Equal(t, "GET", r.Method)
					assertEmptyFieldSpec(t, r)
					assertAuthorization(t, r)
					assertRequestOptionHeader(t, r)

					qs := url.Values{
						"network[id]": []string{tt.args.id},
					}
					assert.Equal(t, qs, r.URL.Query())

					w.WriteHeader(tt.respStatus)
					_, _ = w.Write(tt.respBody)
				},
			)

			got, resp, err := c.GetByID(tt.args.ctx, tt.args.id, testRequestOption)

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

			if tt.errIs != nil {
				assert.ErrorIs(t, err, tt.errIs)
			}
		})
	}
}

func TestNetworksClient_GetByPermalink(t *testing.T) {
	// Correlates to fixtures/network_get.json
	network := &Network{
		ID:        "netw_zDW7KYAeqqfRfVag",
		Name:      "Public Network",
		Permalink: "public",
	}

	type args struct {
		ctx       context.Context
		permalink string
	}
	tests := []struct {
		name       string
		args       args
		want       *Network
		errStr     string
		errResp    *katapult.ResponseError
		errIs      error
		respStatus int
		respBody   []byte
	}{
		{
			name: "network",
			args: args{
				ctx:       context.Background(),
				permalink: "public",
			},
			want:       network,
			respStatus: http.StatusOK,
			respBody:   fixture("network_get"),
		},
		{
			name: "non-existent network",
			args: args{
				ctx:       context.Background(),
				permalink: "not-here",
			},
			errStr:     fixtureNetworkNotFoundErr,
			errResp:    fixtureNetworkNotFoundResponseError,
			errIs:      ErrNetworkNotFound,
			respStatus: http.StatusNotFound,
			respBody:   fixture("network_not_found_error"),
		},
		{
			name: "nil context",
			args: args{
				ctx:       nil,
				permalink: "public",
			},
			errStr: "net/http: nil Context",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rm, mux, _, teardown := prepareTestClient(t)
			defer teardown()
			c := NewNetworksClient(rm)

			mux.HandleFunc("/core/v1/networks/_",
				func(w http.ResponseWriter, r *http.Request) {
					assert.Equal(t, "GET", r.Method)
					assertEmptyFieldSpec(t, r)
					assertAuthorization(t, r)
					assertRequestOptionHeader(t, r)

					qs := url.Values{
						"network[permalink]": []string{tt.args.permalink},
					}
					assert.Equal(t, qs, r.URL.Query())

					w.WriteHeader(tt.respStatus)
					_, _ = w.Write(tt.respBody)
				},
			)

			got, resp, err := c.GetByPermalink(
				tt.args.ctx, tt.args.permalink, testRequestOption,
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

			if tt.errIs != nil {
				assert.ErrorIs(t, err, tt.errIs)
			}
		})
	}
}
