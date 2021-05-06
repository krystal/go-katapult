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
	fixtureIPAlreadyAllocatedErr = "ip_already_allocated: " +
		"This IP address has already been allocated to another " +
		"virtual machine."
	fixtureIPAlreadyAllocatedResponseError = &katapult.ResponseError{
		Code: "ip_already_allocated",
		Description: "This IP address has already been allocated to another " +
			"virtual machine.",
		Detail: json.RawMessage(`{}`),
	}
	fixtureNoAvailableAddressesErr = "no_available_addresses: We don't have " +
		"any available IPs for that network and address version at the " +
		"moment. Please contact support for assistance."
	fixtureNoAvailableAddressesResponseError = &katapult.ResponseError{
		Code: "no_available_addresses",
		Description: "We don't have any available IPs for that network and " +
			"address version at the moment. Please contact support for " +
			"assistance.",
		Detail: json.RawMessage(`{}`),
	}

	fixtureIPAddressFull = &IPAddress{
		ID:              "ip_Ru4ef2oh6STZEQkC",
		Address:         "218.205.195.217",
		ReverseDNS:      "reverse_dns",
		VIP:             true,
		Label:           "east-3",
		AddressWithMask: "218.205.195.217/24",
		Network:         &Network{ID: "netw_zDW7KYAeqqfRfVag"},
		AllocationID:    "vm_USg3i8oJTG5OdbQM",
		AllocationType:  "VirtualMachine",
	}
	fixtureIPAddressNoID = &IPAddress{
		Address:         fixtureIPAddressFull.Address,
		ReverseDNS:      fixtureIPAddressFull.ReverseDNS,
		VIP:             fixtureIPAddressFull.VIP,
		Label:           fixtureIPAddressFull.Label,
		AddressWithMask: fixtureIPAddressFull.AddressWithMask,
		Network:         fixtureIPAddressFull.Network,
		AllocationID:    fixtureIPAddressFull.AllocationID,
		AllocationType:  fixtureIPAddressFull.AllocationType,
	}
	fixtureIPAddressNoLookupField = &IPAddress{
		ReverseDNS:      fixtureIPAddressFull.ReverseDNS,
		VIP:             fixtureIPAddressFull.VIP,
		Label:           fixtureIPAddressFull.Label,
		AddressWithMask: fixtureIPAddressFull.AddressWithMask,
		Network:         fixtureIPAddressFull.Network,
		AllocationID:    fixtureIPAddressFull.AllocationID,
		AllocationType:  fixtureIPAddressFull.AllocationType,
	}

	fixtureIPAddressNotFoundErr = "ip_address_not_found: No IP addresses " +
		"were found matching any of the criteria provided in the arguments"
	fixtureIPAddressNotFoundResponseError = &katapult.ResponseError{
		Code: "ip_address_not_found",
		Description: "No IP addresses were found matching any of the " +
			"criteria provided in the arguments",
		Detail: json.RawMessage(`{}`),
	}
)

func TestClient_IPAddresses(t *testing.T) {
	c := New(&fakeRequestMaker{})

	assert.IsType(t, &IPAddressesClient{}, c.IPAddresses)
}

func TestIPAddress_JSONMarshaling(t *testing.T) {
	tests := []struct {
		name string
		obj  *IPAddress
	}{
		{
			name: "empty",
			obj:  &IPAddress{},
		},
		{
			name: "full",
			obj:  fixtureIPAddressFull,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testJSONMarshaling(t, tt.obj)
		})
	}
}

func TestNewIPAddressLookup(t *testing.T) {
	type args struct {
		idOrAddress string
	}
	tests := []struct {
		name  string
		args  args
		want  *IPAddress
		field FieldName
	}{
		{
			name:  "empty string",
			args:  args{idOrAddress: ""},
			want:  &IPAddress{},
			field: AddressField,
		},
		{
			name:  "ip_ prefixed ID",
			args:  args{idOrAddress: "ip_robwZQGtT4hnAsx4"},
			want:  &IPAddress{ID: "ip_robwZQGtT4hnAsx4"},
			field: IDField,
		},
		{
			name:  "address",
			args:  args{idOrAddress: "51.130.20.179"},
			want:  &IPAddress{Address: "51.130.20.179"},
			field: AddressField,
		},
		{
			name:  "random text",
			args:  args{idOrAddress: "oiTx8fUUh7f32GSw"},
			want:  &IPAddress{Address: "oiTx8fUUh7f32GSw"},
			field: AddressField,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, field := NewIPAddressLookup(tt.args.idOrAddress)

			assert.Equal(t, tt.want, got)
			assert.Equal(t, tt.field, field)
		})
	}
}

func TestIPAddress_Version(t *testing.T) {
	type fields struct {
		Address string
	}
	tests := []struct {
		name   string
		fields fields
		want   IPVersion
	}{
		{
			name:   "IPv4 basic",
			fields: fields{Address: "192.168.0.1"},
			want:   IPv4,
		},
		{
			name:   "IPv4 with port",
			fields: fields{Address: "192.168.0.1:80"},
			want:   IPv4,
		},
		{
			name:   "IPv6 basic",
			fields: fields{Address: "::FFFF:C0A8:1"},
			want:   IPv6,
		},
		{
			name:   "IPv6 leading zeros",
			fields: fields{Address: "::FFFF:C0A8:0001"},
			want:   IPv6,
		},
		{
			name:   "IPv6 double colon expanded",
			fields: fields{Address: "0000:0000:0000:0000:0000:FFFF:C0A8:1"},
			want:   IPv6,
		},
		{
			name:   "IPv6 with zone info",
			fields: fields{Address: "::FFFF:C0A8:1%1"},
			want:   IPv6,
		},
		{
			name:   "IPv6 IPv4 literal",
			fields: fields{Address: "::FFFF:192.168.0.1"},
			want:   IPv6,
		},
		{
			name:   "IPv6 with port info",
			fields: fields{Address: "[::FFFF:C0A8:1]:80"},
			want:   IPv6,
		},
		{
			name:   "IPv6 with zone and port info",
			fields: fields{Address: "[::FFFF:C0A8:1%1]:80"},
			want:   IPv6,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ip := &IPAddress{Address: tt.fields.Address}

			got := ip.Version()

			assert.Equal(t, tt.want, got)
		})
	}
}

func TestIPAddress_lookupReference(t *testing.T) {
	tests := []struct {
		name string
		obj  *IPAddress
		want *IPAddress
	}{
		{
			name: "nil",
			obj:  nil,
			want: nil,
		},
		{
			name: "empty",
			obj:  &IPAddress{},
			want: &IPAddress{},
		},
		{
			name: "full",
			obj:  fixtureIPAddressFull,
			want: &IPAddress{ID: "ip_Ru4ef2oh6STZEQkC"},
		},
		{
			name: "no ID",
			obj:  fixtureIPAddressNoID,
			want: &IPAddress{Address: "218.205.195.217"},
		},
		{
			name: "no ID or Address",
			obj:  fixtureIPAddressNoLookupField,
			want: &IPAddress{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.obj.lookupReference()

			assert.Equal(t, tt.want, got)
		})
	}
}

func TestIPAddress_queryValues(t *testing.T) {
	tests := []struct {
		name string
		obj  *IPAddress
	}{
		{
			name: "nil",
			obj:  nil,
		},
		{
			name: "empty",
			obj:  &IPAddress{},
		},
		{
			name: "full",
			obj:  fixtureIPAddressFull,
		},
		{
			name: "no ID",
			obj:  fixtureIPAddressNoID,
		},
		{
			name: "no ID or Address",
			obj:  fixtureIPAddressNoLookupField,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testQueryableEncoding(t, tt.obj)
		})
	}
}

func Test_ipAddressCreateRequest_JSONMarshaling(t *testing.T) {
	tests := []struct {
		name string
		obj  *ipAddressCreateRequest
	}{
		{
			name: "empty",
			obj:  &ipAddressCreateRequest{},
		},
		{
			name: "full",
			obj: &ipAddressCreateRequest{
				Organization: &Organization{ID: "org_O648YDMEYeLmqdmn"},
				Network:      &Network{ID: "netw_zDW7KYAeqqfRfVag"},
				Version:      IPv4,
				VIP:          truePtr,
				Label:        "web-east-3",
			},
		},
		{
			name: "false VIP",
			obj:  &ipAddressCreateRequest{VIP: falsePtr},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testJSONMarshaling(t, tt.obj)
		})
	}
}

func Test_ipAddressUpdateRequest_JSONMarshaling(t *testing.T) {
	tests := []struct {
		name string
		obj  *ipAddressUpdateRequest
	}{
		{
			name: "empty",
			obj:  &ipAddressUpdateRequest{},
		},
		{
			name: "full",
			obj: &ipAddressUpdateRequest{
				IPAddress:  &IPAddress{ID: "ip_Ru4ef2oh6STZEQkC"},
				VIP:        truePtr,
				Label:      "web-east-3",
				ReverseDNS: "web-east-3.acme.katapult.cloud",
			},
		},
		{
			name: "false VIP",
			obj:  &ipAddressUpdateRequest{VIP: falsePtr},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testJSONMarshaling(t, tt.obj)
		})
	}
}

func Test_ipAddressesResponseBody_JSONMarshaling(t *testing.T) {
	tests := []struct {
		name string
		obj  *ipAddressesResponseBody
	}{
		{
			name: "empty",
			obj:  &ipAddressesResponseBody{},
		},
		{
			name: "full",
			obj: &ipAddressesResponseBody{
				Pagination:  &katapult.Pagination{CurrentPage: 345},
				IPAddress:   &IPAddress{ID: "id1"},
				IPAddresses: []*IPAddress{{ID: "id2"}},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testJSONMarshaling(t, tt.obj)
		})
	}
}

func TestIPAddressesClient_List(t *testing.T) {
	// Correlates to fixtures/ip_addresses_list*.json
	ipAddressesList := []*IPAddress{
		{
			ID:         "ip_dZLqwQifQFtboHXW",
			Address:    "169.37.118.179",
			ReverseDNS: "bitter-beautiful-mango.acme.katapult.cloud",
		},
		{
			ID:         "ip_fAwrdP9NvW0Z25eE",
			Address:    "95.135.35.113",
			ReverseDNS: "popular-shapely-tank.acme.katapult.cloud",
		},
		{
			ID:         "ip_KDPs2kKBiaFohrsF",
			Address:    "200.175.55.138",
			ReverseDNS: "popular-blue-kumquat.acme.katapult.cloud",
		},
	}

	type args struct {
		ctx  context.Context
		org  *Organization
		opts *ListOptions
	}
	tests := []struct {
		name           string
		args           args
		want           []*IPAddress
		wantQuery      *url.Values
		wantPagination *katapult.Pagination
		errStr         string
		errResp        *katapult.ResponseError
		respStatus     int
		respBody       []byte
	}{
		{
			name: "by organization ID",
			args: args{
				ctx: context.Background(),
				org: &Organization{ID: "org_O648YDMEYeLmqdmn"},
			},
			want: ipAddressesList,
			wantQuery: &url.Values{
				"organization[id]": []string{"org_O648YDMEYeLmqdmn"},
			},
			wantPagination: &katapult.Pagination{
				CurrentPage: 1,
				TotalPages:  1,
				Total:       3,
				PerPage:     30,
				LargeSet:    false,
			},
			respStatus: http.StatusOK,
			respBody:   fixture("ip_addresses_list"),
		},
		{
			name: "by organization SubDomain",
			args: args{
				ctx: context.Background(),
				org: &Organization{SubDomain: "acme"},
			},
			want: ipAddressesList,
			wantQuery: &url.Values{
				"organization[sub_domain]": []string{"acme"},
			},
			wantPagination: &katapult.Pagination{
				CurrentPage: 1,
				TotalPages:  1,
				Total:       3,
				PerPage:     30,
				LargeSet:    false,
			},
			respStatus: http.StatusOK,
			respBody:   fixture("ip_addresses_list"),
		},
		{
			name: "page 1",
			args: args{
				ctx:  context.Background(),
				org:  &Organization{ID: "org_O648YDMEYeLmqdmn"},
				opts: &ListOptions{Page: 1, PerPage: 2},
			},
			want: ipAddressesList[0:2],
			wantQuery: &url.Values{
				"organization[id]": []string{"org_O648YDMEYeLmqdmn"},
				"page":             []string{"1"},
				"per_page":         []string{"2"},
			},
			wantPagination: &katapult.Pagination{
				CurrentPage: 1,
				TotalPages:  2,
				Total:       3,
				PerPage:     2,
				LargeSet:    false,
			},
			respStatus: http.StatusOK,
			respBody:   fixture("ip_addresses_list_page_1"),
		},
		{
			name: "page 2",
			args: args{
				ctx:  context.Background(),
				org:  &Organization{ID: "org_O648YDMEYeLmqdmn"},
				opts: &ListOptions{Page: 2, PerPage: 2},
			},
			want: ipAddressesList[2:],
			wantQuery: &url.Values{
				"organization[id]": []string{"org_O648YDMEYeLmqdmn"},
				"page":             []string{"2"},
				"per_page":         []string{"2"},
			},
			wantPagination: &katapult.Pagination{
				CurrentPage: 2,
				TotalPages:  2,
				Total:       3,
				PerPage:     2,
				LargeSet:    false,
			},
			respStatus: http.StatusOK,
			respBody:   fixture("ip_addresses_list_page_2"),
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
			c := NewIPAddressesClient(rm)

			mux.HandleFunc(
				"/core/v1/organizations/_/ip_addresses",
				func(w http.ResponseWriter, r *http.Request) {
					assert.Equal(t, "GET", r.Method)
					assertEmptyFieldSpec(t, r)
					assertAuthorization(t, r)

					if tt.wantQuery != nil {
						assert.Equal(t, *tt.wantQuery, r.URL.Query())
					} else {
						qs := queryValues(tt.args.org, tt.args.opts)
						assert.Equal(t, *qs, r.URL.Query())
					}

					w.WriteHeader(tt.respStatus)
					_, _ = w.Write(tt.respBody)
				},
			)

			got, resp, err := c.List(
				tt.args.ctx, tt.args.org, tt.args.opts,
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

			if tt.wantPagination != nil {
				assert.Equal(t, tt.wantPagination, resp.Pagination)
			}

			if tt.errResp != nil {
				assert.Equal(t, tt.errResp, resp.Error)
			}
		})
	}
}

func TestIPAddressesClient_Get(t *testing.T) {
	type args struct {
		ctx         context.Context
		idOrAddress string
	}
	tests := []struct {
		name       string
		args       args
		want       *IPAddress
		wantQuery  *url.Values
		errStr     string
		errResp    *katapult.ResponseError
		respStatus int
		respBody   []byte
	}{
		{
			name: "by ID",
			args: args{
				ctx:         context.Background(),
				idOrAddress: "ip_dZLqwQifQFtboHXW",
			},
			want: &IPAddress{
				ID:         "ip_dZLqwQifQFtboHXW",
				Address:    "169.37.118.179",
				ReverseDNS: "bitter-beautiful-mango.acme.katapult.cloud",
			},
			wantQuery: &url.Values{
				"ip_address[id]": []string{"ip_dZLqwQifQFtboHXW"},
			},
			respStatus: http.StatusOK,
			respBody:   fixture("ip_address_get"),
		},
		{
			name: "by Address",
			args: args{
				ctx:         context.Background(),
				idOrAddress: "169.37.118.179",
			},
			want: &IPAddress{
				ID:         "ip_dZLqwQifQFtboHXW",
				Address:    "169.37.118.179",
				ReverseDNS: "bitter-beautiful-mango.acme.katapult.cloud",
			},
			wantQuery: &url.Values{
				"ip_address[address]": []string{"169.37.118.179"},
			},
			respStatus: http.StatusOK,
			respBody:   fixture("ip_address_get"),
		},
		{
			name: "non-existent IP address",
			args: args{
				ctx:         context.Background(),
				idOrAddress: "ip_nopethisbegone",
			},
			errStr:     fixtureIPAddressNotFoundErr,
			errResp:    fixtureIPAddressNotFoundResponseError,
			respStatus: http.StatusNotFound,
			respBody:   fixture("ip_address_not_found_error"),
		},
		{
			name: "empty idOrAddress",
			args: args{
				ctx:         context.Background(),
				idOrAddress: "",
			},
			errStr:     fixtureIPAddressNotFoundErr,
			errResp:    fixtureIPAddressNotFoundResponseError,
			respStatus: http.StatusNotFound,
			respBody:   fixture("ip_address_not_found_error"),
		},
		{
			name: "nil context",
			args: args{
				ctx:         nil,
				idOrAddress: "ip_dZLqwQifQFtboHXW",
			},
			errStr: "net/http: nil Context",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rm, mux, _, teardown := prepareTestClient(t)
			defer teardown()
			c := NewIPAddressesClient(rm)

			mux.HandleFunc(
				"/core/v1/ip_addresses/_",
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
				tt.args.ctx, tt.args.idOrAddress,
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

func TestIPAddressesClient_GetByID(t *testing.T) {
	type args struct {
		ctx context.Context
		id  string
	}
	tests := []struct {
		name       string
		args       args
		want       *IPAddress
		errStr     string
		errResp    *katapult.ResponseError
		respStatus int
		respBody   []byte
	}{
		{
			name: "IP address",
			args: args{
				ctx: context.Background(),
				id:  "ip_dZLqwQifQFtboHXW",
			},
			want: &IPAddress{
				ID:         "ip_dZLqwQifQFtboHXW",
				Address:    "169.37.118.179",
				ReverseDNS: "bitter-beautiful-mango.acme.katapult.cloud",
			},
			respStatus: http.StatusOK,
			respBody:   fixture("ip_address_get"),
		},
		{
			name: "non-existent IP address",
			args: args{
				ctx: context.Background(),
				id:  "ip_nopethisbegone",
			},
			errStr:     fixtureIPAddressNotFoundErr,
			errResp:    fixtureIPAddressNotFoundResponseError,
			respStatus: http.StatusNotFound,
			respBody:   fixture("ip_address_not_found_error"),
		},
		{
			name: "empty ID",
			args: args{
				ctx: context.Background(),
				id:  "",
			},
			errStr:     fixtureIPAddressNotFoundErr,
			errResp:    fixtureIPAddressNotFoundResponseError,
			respStatus: http.StatusNotFound,
			respBody:   fixture("ip_address_not_found_error"),
		},
		{
			name: "nil context",
			args: args{
				ctx: nil,
				id:  "ip_dZLqwQifQFtboHXW",
			},
			errStr: "net/http: nil Context",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rm, mux, _, teardown := prepareTestClient(t)
			defer teardown()
			c := NewIPAddressesClient(rm)

			mux.HandleFunc(
				"/core/v1/ip_addresses/_",
				func(w http.ResponseWriter, r *http.Request) {
					assert.Equal(t, "GET", r.Method)
					assertEmptyFieldSpec(t, r)
					assertAuthorization(t, r)

					qs := url.Values{}
					if tt.args.id != "" {
						qs["ip_address[id]"] = []string{tt.args.id}
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

func TestIPAddressesClient_GetByAddress(t *testing.T) {
	type args struct {
		ctx     context.Context
		address string
	}
	tests := []struct {
		name       string
		args       args
		want       *IPAddress
		errStr     string
		errResp    *katapult.ResponseError
		respStatus int
		respBody   []byte
	}{
		{
			name: "IP address",
			args: args{
				ctx:     context.Background(),
				address: "169.37.118.179",
			},
			want: &IPAddress{
				ID:         "ip_dZLqwQifQFtboHXW",
				Address:    "169.37.118.179",
				ReverseDNS: "bitter-beautiful-mango.acme.katapult.cloud",
			},
			respStatus: http.StatusOK,
			respBody:   fixture("ip_address_get"),
		},
		{
			name: "non-existent IP address",
			args: args{
				ctx:     context.Background(),
				address: "153.225.225.79",
			},
			errStr:     fixtureIPAddressNotFoundErr,
			errResp:    fixtureIPAddressNotFoundResponseError,
			respStatus: http.StatusNotFound,
			respBody:   fixture("ip_address_not_found_error"),
		},
		{
			name: "empty Address",
			args: args{
				ctx:     context.Background(),
				address: "",
			},
			errStr:     fixtureIPAddressNotFoundErr,
			errResp:    fixtureIPAddressNotFoundResponseError,
			respStatus: http.StatusNotFound,
			respBody:   fixture("ip_address_not_found_error"),
		},
		{
			name: "nil context",
			args: args{
				ctx:     nil,
				address: "169.37.118.179",
			},
			errStr: "net/http: nil Context",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rm, mux, _, teardown := prepareTestClient(t)
			defer teardown()
			c := NewIPAddressesClient(rm)

			mux.HandleFunc(
				"/core/v1/ip_addresses/_",
				func(w http.ResponseWriter, r *http.Request) {
					assert.Equal(t, "GET", r.Method)
					assertEmptyFieldSpec(t, r)
					assertAuthorization(t, r)

					qs := url.Values{}
					if tt.args.address != "" {
						qs["ip_address[address]"] = []string{tt.args.address}
					}
					assert.Equal(t, qs, r.URL.Query())

					w.WriteHeader(tt.respStatus)
					_, _ = w.Write(tt.respBody)
				},
			)

			got, resp, err := c.GetByAddress(
				tt.args.ctx, tt.args.address,
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

func TestIPAddressesClient_Create(t *testing.T) {
	ipArgs := &IPAddressCreateArguments{
		Network: &Network{ID: "netw_zDW7KYAeqqfRfVag"},
		Version: IPv4,
	}

	type args struct {
		ctx  context.Context
		org  *Organization
		args *IPAddressCreateArguments
	}
	tests := []struct {
		name       string
		args       args
		reqBody    *ipAddressCreateRequest
		want       *IPAddress
		errStr     string
		errResp    *katapult.ResponseError
		respStatus int
		respBody   []byte
	}{
		{
			name: "IPv4 address by organization ID",
			args: args{
				ctx: context.Background(),
				org: &Organization{
					ID:        "org_O648YDMEYeLmqdmn",
					Name:      "ACME Inc.",
					SubDomain: "acme",
				},
				args: &IPAddressCreateArguments{
					Network: &Network{
						ID:        "netw_zDW7KYAeqqfRfVag",
						Name:      "Public Network",
						Permalink: "public",
					},
					Version: IPv4,
				},
			},
			reqBody: &ipAddressCreateRequest{
				Organization: &Organization{ID: "org_O648YDMEYeLmqdmn"},
				Network:      &Network{ID: "netw_zDW7KYAeqqfRfVag"},
				Version:      IPv4,
			},
			want: &IPAddress{
				ID:      "ip_68u3d61zpezcp1Sf",
				Address: "101.240.4.249",
				Network: &Network{ID: "netw_zDW7KYAeqqfRfVag"},
			},
			respStatus: http.StatusCreated,
			respBody:   fixture("ipv4_address_create"),
		},
		{
			name: "IPv6 address by organization SubDomain",
			args: args{
				ctx: context.Background(),
				org: &Organization{
					Name:      "ACME Inc.",
					SubDomain: "acme",
				},
				args: &IPAddressCreateArguments{
					Network: &Network{
						Name:      "Public IPv6 Network",
						Permalink: "public-v6",
					},
					Version: IPv6,
				},
			},
			reqBody: &ipAddressCreateRequest{
				Organization: &Organization{SubDomain: "acme"},
				Network:      &Network{Permalink: "public-v6"},
				Version:      IPv6,
			},
			want: &IPAddress{
				ID:      "ip_bPKp77kkebXaNOrq",
				Address: "94ef:258c:b165:a9d1:84eb:681f:9f57:23b1",
				Network: &Network{ID: "netw_EPXhiG2BCFtni4c1"},
			},
			respStatus: http.StatusCreated,
			respBody:   fixture("ipv6_address_create"),
		},
		{
			name: "non-existent Organization",
			args: args{
				ctx:  context.Background(),
				org:  &Organization{ID: "org_O648YDMEYeLmqdmn"},
				args: ipArgs,
			},
			errStr:     fixtureOrganizationNotFoundErr,
			errResp:    fixtureOrganizationNotFoundResponseError,
			respStatus: http.StatusNotFound,
			respBody:   fixture("organization_not_found_error"),
		},
		{
			name: "suspended Organization",
			args: args{
				ctx:  context.Background(),
				org:  &Organization{ID: "org_O648YDMEYeLmqdmn"},
				args: ipArgs,
			},
			errStr:     fixtureOrganizationSuspendedErr,
			errResp:    fixtureOrganizationSuspendedResponseError,
			respStatus: http.StatusForbidden,
			respBody:   fixture("organization_suspended_error"),
		},
		{
			name: "non-existent Network",
			args: args{
				ctx:  context.Background(),
				org:  &Organization{ID: "org_O648YDMEYeLmqdmn"},
				args: ipArgs,
			},
			errStr:     fixtureNetworkNotFoundErr,
			errResp:    fixtureNetworkNotFoundResponseError,
			respStatus: http.StatusNotFound,
			respBody:   fixture("network_not_found_error"),
		},
		{
			name: "no available addresses",
			args: args{
				ctx:  context.Background(),
				org:  &Organization{ID: "org_O648YDMEYeLmqdmn"},
				args: ipArgs,
			},
			errStr:     fixtureNoAvailableAddressesErr,
			errResp:    fixtureNoAvailableAddressesResponseError,
			respStatus: http.StatusServiceUnavailable,
			respBody:   fixture("no_available_addresses_error"),
		},
		{
			name: "permission denied",
			args: args{
				ctx:  context.Background(),
				org:  &Organization{ID: "org_O648YDMEYeLmqdmn"},
				args: ipArgs,
			},
			errStr:     fixturePermissionDeniedErr,
			errResp:    fixturePermissionDeniedResponseError,
			respStatus: http.StatusForbidden,
			respBody:   fixture("permission_denied_error"),
		},
		{
			name: "validation error",
			args: args{
				ctx:  context.Background(),
				org:  &Organization{ID: "org_O648YDMEYeLmqdmn"},
				args: ipArgs,
			},
			errStr:     fixtureValidationErrorErr,
			errResp:    fixtureValidationErrorResponseError,
			respStatus: http.StatusUnprocessableEntity,
			respBody:   fixture("validation_error"),
		},
		{
			name: "nil ip address arguments",
			args: args{
				ctx:  context.Background(),
				org:  &Organization{ID: "org_O648YDMEYeLmqdmn"},
				args: nil,
			},
			errStr:     fixtureValidationErrorErr,
			errResp:    fixtureValidationErrorResponseError,
			respStatus: http.StatusUnprocessableEntity,
			respBody:   fixture("validation_error"),
		},
		{
			name: "nil context",
			args: args{
				ctx:  nil,
				org:  &Organization{ID: "org_O648YDMEYeLmqdmn"},
				args: ipArgs,
			},
			errStr: "net/http: nil Context",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rm, mux, _, teardown := prepareTestClient(t)
			defer teardown()
			c := NewIPAddressesClient(rm)

			var netName string
			if tt.args.args != nil && tt.args.args.Network != nil {
				netName = tt.args.args.Network.Name
			}

			mux.HandleFunc(
				"/core/v1/organizations/_/ip_addresses",
				func(w http.ResponseWriter, r *http.Request) {
					assert.Equal(t, "POST", r.Method)
					assertEmptyFieldSpec(t, r)
					assertAuthorization(t, r)

					if tt.reqBody != nil {
						reqBody := &ipAddressCreateRequest{}
						err := strictUmarshal(r.Body, reqBody)
						assert.NoError(t, err)
						assert.Equal(t, tt.reqBody, reqBody)
					}

					w.WriteHeader(tt.respStatus)
					_, _ = w.Write(tt.respBody)
				},
			)

			got, resp, err := c.Create(
				tt.args.ctx, tt.args.org, tt.args.args,
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

			if tt.args.args != nil && tt.args.args.Network != nil {
				// ensure the input IPAddressCreateArguments are not modified
				assert.Equal(t, netName, tt.args.args.Network.Name)
			}
		})
	}
}

func TestIPAddressesClient_Update(t *testing.T) {
	ipArgs := &IPAddressUpdateArguments{
		VIP:        truePtr,
		Label:      "web-east-3",
		ReverseDNS: "web-east-3.acme.katapult.cloud",
	}

	type args struct {
		ctx  context.Context
		ip   *IPAddress
		args *IPAddressUpdateArguments
	}
	tests := []struct {
		name       string
		args       args
		reqBody    *ipAddressUpdateRequest
		want       *IPAddress
		errStr     string
		errResp    *katapult.ResponseError
		respStatus int
		respBody   []byte
	}{
		{
			name: "by ID",
			args: args{
				ctx: context.Background(),
				ip: &IPAddress{
					ID:      "ip_dZLqwQifQFtboHXW",
					Address: "169.37.118.179",
				},
				args: ipArgs,
			},
			reqBody: &ipAddressUpdateRequest{
				IPAddress:  &IPAddress{ID: "ip_dZLqwQifQFtboHXW"},
				VIP:        truePtr,
				Label:      "web-east-3",
				ReverseDNS: "web-east-3.acme.katapult.cloud",
			},
			want: &IPAddress{
				ID:         "ip_dZLqwQifQFtboHXW",
				Address:    "169.37.118.179",
				VIP:        true,
				Label:      "web-east-3",
				ReverseDNS: "web-east-3.acme.katapult.cloud",
			},
			respStatus: http.StatusOK,
			respBody:   fixture("ip_address_updated"),
		},
		{
			name: "by Address",
			args: args{
				ctx:  context.Background(),
				ip:   &IPAddress{Address: "169.37.118.179"},
				args: ipArgs,
			},
			reqBody: &ipAddressUpdateRequest{
				IPAddress:  &IPAddress{Address: "169.37.118.179"},
				VIP:        truePtr,
				Label:      "web-east-3",
				ReverseDNS: "web-east-3.acme.katapult.cloud",
			},
			want: &IPAddress{
				ID:         "ip_dZLqwQifQFtboHXW",
				Address:    "169.37.118.179",
				VIP:        true,
				Label:      "web-east-3",
				ReverseDNS: "web-east-3.acme.katapult.cloud",
			},
			respStatus: http.StatusOK,
			respBody:   fixture("ip_address_updated"),
		},
		{
			name: "non-existent IP address",
			args: args{
				ctx:  context.Background(),
				ip:   &IPAddress{ID: "ip_nopethisdoesnotexist"},
				args: &IPAddressUpdateArguments{},
			},
			errStr:     fixtureIPAddressNotFoundErr,
			errResp:    fixtureIPAddressNotFoundResponseError,
			respStatus: http.StatusNotFound,
			respBody:   fixture("ip_address_not_found_error"),
		},
		{
			name: "permission denied",
			args: args{
				ctx:  context.Background(),
				ip:   &IPAddress{ID: "ip_dZLqwQifQFtboHXW"},
				args: ipArgs,
			},
			errStr:     fixturePermissionDeniedErr,
			errResp:    fixturePermissionDeniedResponseError,
			respStatus: http.StatusForbidden,
			respBody:   fixture("permission_denied_error"),
		},
		{
			name: "validation error",
			args: args{
				ctx:  context.Background(),
				ip:   &IPAddress{ID: "ip_dZLqwQifQFtboHXW"},
				args: ipArgs,
			},
			errStr:     fixtureValidationErrorErr,
			errResp:    fixtureValidationErrorResponseError,
			respStatus: http.StatusUnprocessableEntity,
			respBody:   fixture("validation_error"),
		},
		{
			name: "nil ip address",
			args: args{
				ctx:  context.Background(),
				ip:   nil,
				args: ipArgs,
			},
			errStr:     fixtureValidationErrorErr,
			errResp:    fixtureValidationErrorResponseError,
			respStatus: http.StatusUnprocessableEntity,
			respBody:   fixture("validation_error"),
		},
		{
			name: "nil context",
			args: args{
				ctx:  nil,
				ip:   &IPAddress{ID: "ip_dZLqwQifQFtboHXW"},
				args: ipArgs,
			},
			errStr: "net/http: nil Context",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rm, mux, _, teardown := prepareTestClient(t)
			defer teardown()
			c := NewIPAddressesClient(rm)

			mux.HandleFunc(
				"/core/v1/ip_addresses/_",
				func(w http.ResponseWriter, r *http.Request) {
					assert.Equal(t, "PATCH", r.Method)
					assertEmptyFieldSpec(t, r)
					assertAuthorization(t, r)

					if tt.reqBody != nil {
						reqBody := &ipAddressUpdateRequest{}
						err := strictUmarshal(r.Body, reqBody)
						assert.NoError(t, err)
						assert.Equal(t, tt.reqBody, reqBody)
					}

					w.WriteHeader(tt.respStatus)
					_, _ = w.Write(tt.respBody)
				},
			)

			got, resp, err := c.Update(
				tt.args.ctx, tt.args.ip, tt.args.args,
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

func TestIPAddressesClient_Delete(t *testing.T) {
	type args struct {
		ctx context.Context
		ip  *IPAddress
	}
	tests := []struct {
		name       string
		args       args
		wantQuery  *url.Values
		errStr     string
		errResp    *katapult.ResponseError
		respStatus int
		respBody   []byte
	}{
		{
			name: "by ID",
			args: args{
				ctx: context.Background(),
				ip:  &IPAddress{ID: "ip_dZLqwQifQFtboHXW"},
			},
			wantQuery: &url.Values{
				"ip_address[id]": []string{"ip_dZLqwQifQFtboHXW"},
			},
			respStatus: http.StatusOK,
			respBody:   []byte("{}"),
		},
		{
			name: "by Address",
			args: args{
				ctx: context.Background(),
				ip:  &IPAddress{Address: "169.37.118.179"},
			},
			wantQuery: &url.Values{
				"ip_address[address]": []string{"169.37.118.179"},
			},
			respStatus: http.StatusOK,
			respBody:   []byte("{}"),
		},
		{
			name: "non-existent IP address",
			args: args{
				ctx: context.Background(),
				ip:  &IPAddress{ID: "ip_nopenotfound"},
			},
			errStr:     fixtureIPAddressNotFoundErr,
			errResp:    fixtureIPAddressNotFoundResponseError,
			respStatus: http.StatusNotFound,
			respBody:   fixture("ip_address_not_found_error"),
		},
		{
			name: "permission denied",
			args: args{
				ctx: context.Background(),
				ip:  &IPAddress{ID: "ip_dZLqwQifQFtboHXW"},
			},
			errStr:     fixturePermissionDeniedErr,
			errResp:    fixturePermissionDeniedResponseError,
			respStatus: http.StatusForbidden,
			respBody:   fixture("permission_denied_error"),
		},
		{
			name: "nil IP Address",
			args: args{
				ctx: context.Background(),
				ip:  nil,
			},
			errStr:     fixtureIPAddressNotFoundErr,
			errResp:    fixtureIPAddressNotFoundResponseError,
			respStatus: http.StatusNotFound,
			respBody:   fixture("ip_address_not_found_error"),
		},
		{
			name: "nil context",
			args: args{
				ctx: nil,
				ip:  &IPAddress{ID: "ip_dZLqwQifQFtboHXW"},
			},
			errStr: "net/http: nil Context",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rm, mux, _, teardown := prepareTestClient(t)
			defer teardown()
			c := NewIPAddressesClient(rm)

			mux.HandleFunc(
				"/core/v1/ip_addresses/_",
				func(w http.ResponseWriter, r *http.Request) {
					assert.Equal(t, "DELETE", r.Method)
					assertEmptyFieldSpec(t, r)
					assertAuthorization(t, r)

					if tt.wantQuery != nil {
						assert.Equal(t, *tt.wantQuery, r.URL.Query())
					} else {
						assert.Equal(t,
							*tt.args.ip.queryValues(), r.URL.Query(),
						)
					}

					w.WriteHeader(tt.respStatus)
					_, _ = w.Write(tt.respBody)
				},
			)

			resp, err := c.Delete(tt.args.ctx, tt.args.ip)

			if tt.respStatus != 0 {
				assert.Equal(t, tt.respStatus, resp.StatusCode)
			}

			if tt.errStr == "" {
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, err, tt.errStr)
			}

			if tt.errResp != nil {
				assert.Equal(t, tt.errResp, resp.Error)
			}
		})
	}
}

func TestIPAddressesClient_Unallocate(t *testing.T) {
	type args struct {
		ctx context.Context
		ip  *IPAddress
	}
	tests := []struct {
		name       string
		args       args
		wantQuery  *url.Values
		errStr     string
		errResp    *katapult.ResponseError
		respStatus int
		respBody   []byte
	}{
		{
			name: "by ID",
			args: args{
				ctx: context.Background(),
				ip:  &IPAddress{ID: "ip_dZLqwQifQFtboHXW"},
			},
			wantQuery: &url.Values{
				"ip_address[id]": []string{"ip_dZLqwQifQFtboHXW"},
			},
			respStatus: http.StatusOK,
			respBody:   []byte("{}"),
		},
		{
			name: "by Address",
			args: args{
				ctx: context.Background(),
				ip:  &IPAddress{Address: "169.37.118.179"},
			},
			wantQuery: &url.Values{
				"ip_address[address]": []string{"169.37.118.179"},
			},
			respStatus: http.StatusOK,
			respBody:   fixture("load_balancer_get"),
		},
		{
			name: "non-existent IP address",
			args: args{
				ctx: context.Background(),
				ip:  &IPAddress{ID: "ip_nopenotfound"},
			},
			errStr:     fixtureIPAddressNotFoundErr,
			errResp:    fixtureIPAddressNotFoundResponseError,
			respStatus: http.StatusNotFound,
			respBody:   fixture("ip_address_not_found_error"),
		},
		{
			name: "permission denied",
			args: args{
				ctx: context.Background(),
				ip:  &IPAddress{ID: "ip_dZLqwQifQFtboHXW"},
			},
			errStr:     fixturePermissionDeniedErr,
			errResp:    fixturePermissionDeniedResponseError,
			respStatus: http.StatusForbidden,
			respBody:   fixture("permission_denied_error"),
		},
		{
			name: "nil IP Address",
			args: args{
				ctx: context.Background(),
				ip:  nil,
			},
			errStr:     fixtureIPAddressNotFoundErr,
			errResp:    fixtureIPAddressNotFoundResponseError,
			respStatus: http.StatusNotFound,
			respBody:   fixture("ip_address_not_found_error"),
		},
		{
			name: "nil context",
			args: args{
				ctx: nil,
				ip:  &IPAddress{ID: "ip_dZLqwQifQFtboHXW"},
			},
			errStr: "net/http: nil Context",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rm, mux, _, teardown := prepareTestClient(t)
			defer teardown()
			c := NewIPAddressesClient(rm)

			mux.HandleFunc(
				"/core/v1/ip_addresses/_/unallocate",
				func(w http.ResponseWriter, r *http.Request) {
					assert.Equal(t, "POST", r.Method)
					assertEmptyFieldSpec(t, r)
					assertAuthorization(t, r)

					if tt.wantQuery != nil {
						assert.Equal(t, *tt.wantQuery, r.URL.Query())
					} else {
						assert.Equal(t,
							*tt.args.ip.queryValues(), r.URL.Query(),
						)
					}

					w.WriteHeader(tt.respStatus)
					_, _ = w.Write(tt.respBody)
				},
			)

			resp, err := c.Unallocate(tt.args.ctx, tt.args.ip)

			if tt.respStatus != 0 {
				assert.Equal(t, tt.respStatus, resp.StatusCode)
			}

			if tt.errStr == "" {
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, err, tt.errStr)
			}

			if tt.errResp != nil {
				assert.Equal(t, tt.errResp, resp.Error)
			}
		})
	}
}
