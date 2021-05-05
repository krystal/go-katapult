package core

import (
	"context"
	"encoding/json"
	"net/http"
	"net/url"
	"testing"

	"github.com/krystal/go-katapult"
	"github.com/stretchr/testify/assert"
)

var (
	fixtureNetworkNotFoundErr = "network_not_found: No network was found " +
		"matching any of the criteria provided in the arguments"
	fixtureNetworkNotFoundResponseError = &katapult.ResponseError{
		Code: "network_not_found",
		Description: "No network was found matching any of the criteria " +
			"provided in the arguments",
		Detail: json.RawMessage(`{}`),
	}

	fixtureNetworkFull = &Network{
		ID:         "netw_zDW7KYAeqqfRfVag",
		Name:       "Public Network",
		Permalink:  "public",
		DataCenter: &DataCenter{ID: "id2"},
	}
	fixtureNetworkNoID = &Network{
		Name:       fixtureNetworkFull.Name,
		Permalink:  fixtureNetworkFull.Permalink,
		DataCenter: fixtureNetworkFull.DataCenter,
	}
	fixtureNetworkNoLookupField = &Network{
		Name:       fixtureNetworkFull.Name,
		DataCenter: fixtureNetworkFull.DataCenter,
	}
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

func TestNewNetworkLookup(t *testing.T) {
	type args struct {
		idOrPermalink string
	}
	tests := []struct {
		name  string
		args  args
		want  *Network
		field FieldName
	}{
		{
			name:  "empty string",
			args:  args{idOrPermalink: ""},
			want:  &Network{},
			field: PermalinkField,
		},
		{
			name:  "netw_ prefixed ID",
			args:  args{idOrPermalink: "netw_UoGX2x12BlVK0CAo"},
			want:  &Network{ID: "netw_UoGX2x12BlVK0CAo"},
			field: IDField,
		},
		{
			name:  "permalink",
			args:  args{idOrPermalink: "country-city-1-public"},
			want:  &Network{Permalink: "country-city-1-public"},
			field: PermalinkField,
		},
		{
			name:  "random text",
			args:  args{idOrPermalink: "JRdJ017AgV4WYbkv"},
			want:  &Network{Permalink: "JRdJ017AgV4WYbkv"},
			field: PermalinkField,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, field := NewNetworkLookup(tt.args.idOrPermalink)

			assert.Equal(t, tt.want, got)
			assert.Equal(t, tt.field, field)
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
			obj:  fixtureNetworkFull,
			want: &Network{ID: "netw_zDW7KYAeqqfRfVag"},
		},
		{
			name: "no ID",
			obj:  fixtureNetworkNoID,
			want: &Network{Permalink: "public"},
		},
		{
			name: "no ID or Permalink",
			obj:  fixtureNetworkNoLookupField,
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

func TestNetwork_queryValues(t *testing.T) {
	tests := []struct {
		name string
		obj  *Network
	}{
		{
			name: "nil",
			obj:  nil,
		},
		{
			name: "empty",
			obj:  &Network{},
		},
		{
			name: "full",
			obj:  fixtureNetworkFull,
		},
		{
			name: "no ID",
			obj:  fixtureNetworkNoID,
		},
		{
			name: "no ID or Permalink",
			obj:  fixtureNetworkNoLookupField,
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
		org *Organization
	}
	tests := []struct {
		name       string
		args       args
		wantNets   []*Network
		wantVnets  []*VirtualNetwork
		wantQuery  *url.Values
		errStr     string
		errResp    *katapult.ResponseError
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
			rm, mux, _, teardown := prepareTestClient(t)
			defer teardown()
			c := NewNetworksClient(rm)

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

			got1, got2, resp, err := c.List(tt.args.ctx, tt.args.org)

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

func TestNetworksClient_Get(t *testing.T) {
	// Correlates to fixtures/network_get.json
	network := &Network{
		ID:        "netw_zDW7KYAeqqfRfVag",
		Name:      "Public Network",
		Permalink: "public",
	}

	type args struct {
		ctx           context.Context
		idOrPermalink string
	}
	tests := []struct {
		name       string
		args       args
		want       *Network
		wantQuery  *url.Values
		errStr     string
		errResp    *katapult.ResponseError
		respStatus int
		respBody   []byte
	}{
		{
			name: "by ID",
			args: args{
				ctx:           context.Background(),
				idOrPermalink: "netw_zDW7KYAeqqfRfVag",
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
				ctx:           context.Background(),
				idOrPermalink: "public",
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
				ctx:           context.Background(),
				idOrPermalink: "netw_nopethisbegone",
			},
			errStr:     fixtureNetworkNotFoundErr,
			errResp:    fixtureNetworkNotFoundResponseError,
			respStatus: http.StatusNotFound,
			respBody:   fixture("network_not_found_error"),
		},
		{
			name: "non-existent network by Permalink",
			args: args{
				ctx:           context.Background(),
				idOrPermalink: "public",
			},
			errStr:     fixtureNetworkNotFoundErr,
			errResp:    fixtureNetworkNotFoundResponseError,
			respStatus: http.StatusNotFound,
			respBody:   fixture("network_not_found_error"),
		},
		{
			name: "nil context",
			args: args{
				ctx:           nil,
				idOrPermalink: "netw_zDW7KYAeqqfRfVag",
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

					if tt.wantQuery != nil {
						assert.Equal(t, *tt.wantQuery, r.URL.Query())
					}

					w.WriteHeader(tt.respStatus)
					_, _ = w.Write(tt.respBody)
				},
			)

			got, resp, err := c.Get(
				tt.args.ctx, tt.args.idOrPermalink,
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

					qs := url.Values{
						"network[id]": []string{tt.args.id},
					}
					assert.Equal(t, qs, r.URL.Query())

					w.WriteHeader(tt.respStatus)
					_, _ = w.Write(tt.respBody)
				},
			)

			got, resp, err := c.GetByID(tt.args.ctx, tt.args.id)

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

					qs := url.Values{
						"network[permalink]": []string{tt.args.permalink},
					}
					assert.Equal(t, qs, r.URL.Query())

					w.WriteHeader(tt.respStatus)
					_, _ = w.Write(tt.respBody)
				},
			)

			got, resp, err := c.GetByPermalink(
				tt.args.ctx, tt.args.permalink,
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
