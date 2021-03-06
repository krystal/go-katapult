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
	fixtureDNSZoneNotFoundErr = "katapult: not_found: dns_zone_not_found: " +
		"No DNS zone was found matching any of the criteria provided in the " +
		"arguments"
	fixtureDNSZoneNotFoundResponseError = &katapult.ResponseError{
		Code: "dns_zone_not_found",
		Description: "No DNS zone was found matching any of the criteria " +
			"provided in the arguments",
		Detail: json.RawMessage(`{}`),
	}

	fixtureDNSZoneInfraErr = "katapult: unauthorized: " +
		"infrastructure_dns_zone_cannot_be_edited: Infrastructure DNS zones " +
		"cannot be edited through the API. These are managed exclusively by " +
		"Katapult."
	fixtureDNSZoneInfraResponseError = &katapult.ResponseError{
		Code: "infrastructure_dns_zone_cannot_be_edited",
		Description: "Infrastructure DNS zones cannot be edited through the " +
			"API. These are managed exclusively by Katapult.",
		Detail: json.RawMessage(`{}`),
	}

	// Correlates to fixtures/dns_zone_get.json
	fixtureDNSZone = &DNSZone{
		ID:   "dnszone_k75eFc4UBOgeE5Zy",
		Name: "test1.example.com",
		TTL:  3600,
	}
)

func TestClient_DNSZones(t *testing.T) {
	c := New(&testclient.Client{})

	assert.IsType(t, &DNSZonesClient{}, c.DNSZones)
}

func TestDNSZone_JSONMarshaling(t *testing.T) {
	tests := []struct {
		name string
		obj  *DNSZone
	}{
		{
			name: "empty",
			obj:  &DNSZone{},
		},
		{
			name: "full",
			obj: &DNSZone{
				ID:                 "dnszone_k75eFc4UBOgeE5Zy",
				Name:               "test1.example.com",
				TTL:                343,
				Verified:           true,
				InfrastructureZone: true,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testJSONMarshaling(t, tt.obj)
		})
	}
}

func TestDNZZone_Ref(t *testing.T) {
	tests := []struct {
		name string
		obj  DNSZone
		want DNSZoneRef
	}{
		{
			name: "with id",
			obj: DNSZone{
				ID: "dnszone_k75eFc4UBOgeE5Zy",
			},
			want: DNSZoneRef{ID: "dnszone_k75eFc4UBOgeE5Zy"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.obj.Ref())
		})
	}
}

func TestDNSZoneRef_queryValues(t *testing.T) {
	tests := []struct {
		name string
		obj  DNSZoneRef
	}{
		{
			name: "empty",
			obj:  DNSZoneRef{},
		},
		{
			name: "full",
			obj: DNSZoneRef{
				ID:   "dnszone_k75eFc4UBOgeE5Zy",
				Name: "test1.example.com",
			},
		},
		{
			name: "just ID",
			obj: DNSZoneRef{
				ID: "dnszone_k75eFc4UBOgeE5Zy",
			},
		},
		{
			name: "no ID",
			obj: DNSZoneRef{
				Name: "test1.example.com",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testQueryableEncoding(t, tt.obj)
		})
	}
}

func TestDNSZoneVerificationDetails_JSONMarshaling(t *testing.T) {
	tests := []struct {
		name string
		obj  *DNSZoneVerificationDetails
	}{
		{
			name: "empty",
			obj:  &DNSZoneVerificationDetails{},
		},
		{
			name: "full",
			obj: &DNSZoneVerificationDetails{
				Nameservers: []string{"foo", "bar"},
				TXTRecord:   "txt",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testJSONMarshaling(t, tt.obj)
		})
	}
}

func TestDNSZoneArguments_JSONMarshaling(t *testing.T) {
	tests := []struct {
		name string
		obj  *DNSZoneArguments
	}{
		{
			name: "empty",
			obj:  &DNSZoneArguments{},
		},
		{
			name: "full",
			obj: &DNSZoneArguments{
				Name:            "name",
				TTL:             493,
				SkipVerfication: true,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testJSONMarshaling(t, tt.obj)
		})
	}
}

func TestDNSZoneDetails_JSONMarshaling(t *testing.T) {
	tests := []struct {
		name string
		obj  *DNSZoneDetails
	}{
		{
			name: "empty",
			obj:  &DNSZoneDetails{},
		},
		{
			name: "full",
			obj: &DNSZoneDetails{
				Name: "name",
				TTL:  493,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testJSONMarshaling(t, tt.obj)
		})
	}
}

func Test_dnsZoneCreateRequest_JSONMarshaling(t *testing.T) {
	tests := []struct {
		name string
		obj  *dnsZoneCreateRequest
	}{
		{
			name: "empty",
			obj:  &dnsZoneCreateRequest{},
		},
		{
			name: "full",
			obj: &dnsZoneCreateRequest{
				Organization:    OrganizationRef{ID: "org_QwNl81npdQQGinmt"},
				Details:         &DNSZoneDetails{Name: "test1.example.com"},
				SkipVerfication: true,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testJSONMarshaling(t, tt.obj)
		})
	}
}

func Test_dnsZoneUpdateTTLRequest_JSONMarshaling(t *testing.T) {
	tests := []struct {
		name string
		obj  *dnsZoneUpdateTTLRequest
	}{
		{
			name: "empty",
			obj:  &dnsZoneUpdateTTLRequest{},
		},
		{
			name: "full",
			obj: &dnsZoneUpdateTTLRequest{
				DNSZone: DNSZoneRef{ID: "dnszone_gymjA0XKuxJlcQXZ"},
				TTL:     8384,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testJSONMarshaling(t, tt.obj)
		})
	}
}

func Test_dnsZoneResponseBody_JSONMarshaling(t *testing.T) {
	tests := []struct {
		name string
		obj  *dnsZoneResponseBody
	}{
		{
			name: "empty",
			obj:  &dnsZoneResponseBody{},
		},
		{
			name: "full",
			obj: &dnsZoneResponseBody{
				Pagination: &katapult.Pagination{CurrentPage: 934},
				DNSZones:   []*DNSZone{{ID: "id1"}},
				DNSZone:    &DNSZone{ID: "id2"},
				VerificationDetails: &DNSZoneVerificationDetails{
					TXTRecord: "txt",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testJSONMarshaling(t, tt.obj)
		})
	}
}

func TestDNSZonesClient_List(t *testing.T) {
	// Correlates to fixtures/dns_zones_list*.json
	dnsZonesList := []*DNSZone{
		{
			ID:   "dnszone_k75eFc4UBOgeE5Zy",
			Name: "test1.example.com",
			TTL:  3600,
		},
		{
			ID:   "dnszone_lwz66kyviwCQyqQc",
			Name: "test-2.example.com",
			TTL:  3600,
		},
		{
			ID:   "dnszone_qr9KPhSwkGNh7IMb",
			Name: "test-3.example.com",
			TTL:  3600,
		},
	}

	type args struct {
		ctx  context.Context
		org  OrganizationRef
		opts *ListOptions
	}
	tests := []struct {
		name           string
		args           args
		want           []*DNSZone
		wantPagination *katapult.Pagination
		errStr         string
		errResp        *katapult.ResponseError
		errIs          error
		respStatus     int
		respBody       []byte
	}{
		{
			name: "by organization ID",
			args: args{
				ctx: context.Background(),
				org: OrganizationRef{ID: "org_O648YDMEYeLmqdmn"},
			},
			want: dnsZonesList,
			wantPagination: &katapult.Pagination{
				CurrentPage: 1,
				TotalPages:  1,
				Total:       3,
				PerPage:     30,
				LargeSet:    false,
			},
			respStatus: http.StatusOK,
			respBody:   fixture("dns_zones_list"),
		},
		{
			name: "by organization SubDomain",
			args: args{
				ctx: context.Background(),
				org: OrganizationRef{SubDomain: "blackmesa"},
			},
			want: dnsZonesList,
			wantPagination: &katapult.Pagination{
				CurrentPage: 1,
				TotalPages:  1,
				Total:       3,
				PerPage:     30,
				LargeSet:    false,
			},
			respStatus: http.StatusOK,
			respBody:   fixture("dns_zones_list"),
		},
		{
			name: "page 1",
			args: args{
				ctx:  context.Background(),
				org:  OrganizationRef{ID: "org_O648YDMEYeLmqdmn"},
				opts: &ListOptions{Page: 1, PerPage: 2},
			},
			want: dnsZonesList[0:2],

			wantPagination: &katapult.Pagination{
				CurrentPage: 1,
				TotalPages:  2,
				Total:       3,
				PerPage:     2,
				LargeSet:    false,
			},
			respStatus: http.StatusOK,
			respBody:   fixture("dns_zones_list_page_1"),
		},
		{
			name: "page 2",
			args: args{
				ctx:  context.Background(),
				org:  OrganizationRef{ID: "org_O648YDMEYeLmqdmn"},
				opts: &ListOptions{Page: 2, PerPage: 2},
			},
			want: dnsZonesList[2:],

			wantPagination: &katapult.Pagination{
				CurrentPage: 2,
				TotalPages:  2,
				Total:       3,
				PerPage:     2,
				LargeSet:    false,
			},
			respStatus: http.StatusOK,
			respBody:   fixture("dns_zones_list_page_2"),
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
				org: OrganizationRef{ID: "org_O648YDMEYeLmqdmn"},
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
			name: "permission denied",
			args: args{
				ctx: context.Background(),
				org: OrganizationRef{ID: "org_O648YDMEYeLmqdmn"},
			},
			errStr:     fixturePermissionDeniedErr,
			errResp:    fixturePermissionDeniedResponseError,
			errIs:      ErrPermissionDenied,
			respStatus: http.StatusForbidden,
			respBody:   fixture("permission_denied_error"),
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
			c := NewDNSZonesClient(rm)

			mux.HandleFunc(
				"/core/v1/organizations/_/dns/zones",
				func(w http.ResponseWriter, r *http.Request) {
					assert.Equal(t, "GET", r.Method)
					assertEmptyFieldSpec(t, r)
					assertAuthorization(t, r)

					qs := queryValues(tt.args.org, tt.args.opts)
					assert.Equal(t, *qs, r.URL.Query())

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

			if tt.errIs != nil {
				assert.ErrorIs(t, err, tt.errIs)
			}
		})
	}
}

func TestDNSZonesClient_Get(t *testing.T) {
	type args struct {
		ctx context.Context
		ref DNSZoneRef
	}
	tests := []struct {
		name       string
		args       args
		want       *DNSZone
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
				ref: DNSZoneRef{ID: "dnszone_k75eFc4UBOgeE5Zy"},
			},
			want:       fixtureDNSZone,
			respStatus: http.StatusOK,
			respBody:   fixture("dns_zone_get"),
		},
		{
			name: "by Name",
			args: args{
				ctx: context.Background(),
				ref: DNSZoneRef{Name: "test1.example.com"},
			},
			want:       fixtureDNSZone,
			respStatus: http.StatusOK,
			respBody:   fixture("dns_zone_get"),
		},
		{
			name: "non-existent DNS zone",
			args: args{
				ctx: context.Background(),
				ref: DNSZoneRef{Name: "dnszone_k75eFc4UBOgeE5Zy"},
			},
			errStr:     fixtureDNSZoneNotFoundErr,
			errResp:    fixtureDNSZoneNotFoundResponseError,
			errIs:      ErrDNSZoneNotFound,
			respStatus: http.StatusNotFound,
			respBody:   fixture("dns_zone_not_found_error"),
		},
		{
			name: "nil context",
			args: args{
				ctx: nil,
				ref: DNSZoneRef{ID: "dnszone_k75eFc4UBOgeE5Zy"},
			},
			errStr: "net/http: nil Context",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rm, mux, _, teardown := prepareTestClient(t)
			defer teardown()
			c := NewDNSZonesClient(rm)

			mux.HandleFunc(
				"/core/v1/dns/zones/_",
				func(w http.ResponseWriter, r *http.Request) {
					assert.Equal(t, "GET", r.Method)
					assertEmptyFieldSpec(t, r)
					assertAuthorization(t, r)

					assert.Equal(t, *tt.args.ref.queryValues(), r.URL.Query())

					w.WriteHeader(tt.respStatus)
					_, _ = w.Write(tt.respBody)
				},
			)

			got, resp, err := c.Get(tt.args.ctx, tt.args.ref)

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

func TestDNSZonesClient_GetByID(t *testing.T) {
	type args struct {
		ctx context.Context
		id  string
	}
	tests := []struct {
		name       string
		args       args
		want       *DNSZone
		errStr     string
		errResp    *katapult.ResponseError
		errIs      error
		respStatus int
		respBody   []byte
	}{
		{
			name: "DNS zone",
			args: args{
				ctx: context.Background(),
				id:  "dnszone_k75eFc4UBOgeE5Zy",
			},
			want:       fixtureDNSZone,
			respStatus: http.StatusOK,
			respBody:   fixture("dns_zone_get"),
		},
		{
			name: "non-existent DNS zone",
			args: args{
				ctx: context.Background(),
				id:  "dnszone_k75eFc4UBOgeE5Zy",
			},
			errStr:     fixtureDNSZoneNotFoundErr,
			errResp:    fixtureDNSZoneNotFoundResponseError,
			errIs:      ErrDNSZoneNotFound,
			respStatus: http.StatusNotFound,
			respBody:   fixture("dns_zone_not_found_error"),
		},
		{
			name: "nil context",
			args: args{
				ctx: nil,
				id:  "dnszone_k75eFc4UBOgeE5Zy",
			},
			errStr: "net/http: nil Context",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rm, mux, _, teardown := prepareTestClient(t)
			defer teardown()
			c := NewDNSZonesClient(rm)

			mux.HandleFunc("/core/v1/dns/zones/_",
				func(w http.ResponseWriter, r *http.Request) {
					assert.Equal(t, "GET", r.Method)
					assertEmptyFieldSpec(t, r)
					assertAuthorization(t, r)

					assert.Equal(t, url.Values{
						"dns_zone[id]": []string{tt.args.id},
					}, r.URL.Query())

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

			if tt.errIs != nil {
				assert.ErrorIs(t, err, tt.errIs)
			}
		})
	}
}

func TestDNSZonesClient_GetByName(t *testing.T) {
	type args struct {
		ctx  context.Context
		name string
	}
	tests := []struct {
		name       string
		args       args
		want       *DNSZone
		errStr     string
		errResp    *katapult.ResponseError
		errIs      error
		respStatus int
		respBody   []byte
	}{
		{
			name: "DNS zone",
			args: args{
				ctx:  context.Background(),
				name: "test1.example.com",
			},
			want:       fixtureDNSZone,
			respStatus: http.StatusOK,
			respBody:   fixture("dns_zone_get"),
		},
		{
			name: "non-existent DNS zone",
			args: args{
				ctx:  context.Background(),
				name: "test1.examplezzz.com",
			},
			errStr:     fixtureDNSZoneNotFoundErr,
			errResp:    fixtureDNSZoneNotFoundResponseError,
			errIs:      ErrDNSZoneNotFound,
			respStatus: http.StatusNotFound,
			respBody:   fixture("dns_zone_not_found_error"),
		},
		{
			name: "nil context",
			args: args{
				ctx:  nil,
				name: "test1.example.com",
			},
			errStr: "net/http: nil Context",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rm, mux, _, teardown := prepareTestClient(t)
			defer teardown()
			c := NewDNSZonesClient(rm)

			mux.HandleFunc("/core/v1/dns/zones/_",
				func(w http.ResponseWriter, r *http.Request) {
					assert.Equal(t, "GET", r.Method)
					assertEmptyFieldSpec(t, r)
					assertAuthorization(t, r)

					qs := url.Values{"dns_zone[name]": []string{tt.args.name}}
					assert.Equal(t, qs, r.URL.Query())

					w.WriteHeader(tt.respStatus)
					_, _ = w.Write(tt.respBody)
				},
			)

			got, resp, err := c.GetByName(tt.args.ctx, tt.args.name)

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

func TestDNSZonesClient_Create(t *testing.T) {
	type args struct {
		ctx      context.Context
		org      OrganizationRef
		zoneArgs *DNSZoneArguments
	}
	tests := []struct {
		name       string
		orgID      string
		args       args
		reqBody    *dnsZoneCreateRequest
		want       *DNSZone
		errStr     string
		errResp    *katapult.ResponseError
		errIs      error
		respStatus int
		respBody   []byte
	}{
		{
			name: "create a DNS zone",
			args: args{
				ctx: context.Background(),
				org: OrganizationRef{ID: "org_O648YDMEYeLmqdmn"},
				zoneArgs: &DNSZoneArguments{
					Name: "test-1.com",
					TTL:  1800,
				},
			},
			reqBody: &dnsZoneCreateRequest{
				Organization: OrganizationRef{ID: "org_O648YDMEYeLmqdmn"},
				Details: &DNSZoneDetails{
					Name: "test-1.com",
					TTL:  1800,
				},
				SkipVerfication: false,
			},
			want: &DNSZone{
				ID:   "dnszone_yqflWVIdu5vnirLq",
				Name: "test-1.com",
				TTL:  1800,
			},
			respStatus: http.StatusCreated,
			respBody:   fixture("dns_zone_create"),
		},
		{
			name: "skip verification",
			args: args{
				ctx: context.Background(),
				org: OrganizationRef{ID: "org_O648YDMEYeLmqdmn"},
				zoneArgs: &DNSZoneArguments{
					Name:            "test-1.com",
					TTL:             1800,
					SkipVerfication: true,
				},
			},
			reqBody: &dnsZoneCreateRequest{
				Organization: OrganizationRef{ID: "org_O648YDMEYeLmqdmn"},
				Details: &DNSZoneDetails{
					Name: "test-1.com",
					TTL:  1800,
				},
				SkipVerfication: true,
			},
			respStatus: http.StatusCreated,
			respBody:   fixture("dns_zone_create"),
		},
		{
			name: "without TTL",
			args: args{
				ctx:      context.Background(),
				org:      OrganizationRef{ID: "org_O648YDMEYeLmqdmn"},
				zoneArgs: &DNSZoneArguments{Name: "test-1.com"},
			},
			reqBody: &dnsZoneCreateRequest{
				Organization: OrganizationRef{ID: "org_O648YDMEYeLmqdmn"},
				Details: &DNSZoneDetails{
					Name: "test-1.com",
				},
			},
			respStatus: http.StatusCreated,
			respBody:   fixture("dns_zone_create"),
		},
		{
			name: "without TTL",
			args: args{
				ctx:      context.Background(),
				org:      OrganizationRef{ID: "org_O648YDMEYeLmqdmn"},
				zoneArgs: &DNSZoneArguments{Name: "test-1.com"},
			},
			respStatus: http.StatusCreated,
			respBody:   fixture("dns_zone_create"),
		},
		{
			name: "non-existent Organization",
			args: args{
				ctx:      context.Background(),
				org:      OrganizationRef{ID: "org_O648YDMEYeLmqdmn"},
				zoneArgs: &DNSZoneArguments{Name: "test-1.com"},
			},
			errStr:     fixtureOrganizationNotFoundErr,
			errResp:    fixtureOrganizationNotFoundResponseError,
			errIs:      ErrOrganizationNotFound,
			respStatus: http.StatusNotFound,
			respBody:   fixture("organization_not_found_error"),
		},
		{
			name: "suspended Organization",
			args: args{
				ctx:      context.Background(),
				org:      OrganizationRef{ID: "org_O648YDMEYeLmqdmn"},
				zoneArgs: &DNSZoneArguments{Name: "test-1.com"},
			},
			errStr:     fixtureOrganizationSuspendedErr,
			errResp:    fixtureOrganizationSuspendedResponseError,
			errIs:      ErrOrganizationSuspended,
			respStatus: http.StatusForbidden,
			respBody:   fixture("organization_suspended_error"),
		},
		{
			name: "permission denied",
			args: args{
				ctx:      context.Background(),
				org:      OrganizationRef{ID: "org_O648YDMEYeLmqdmn"},
				zoneArgs: &DNSZoneArguments{Name: "test-1.com"},
			},
			errStr:     fixturePermissionDeniedErr,
			errResp:    fixturePermissionDeniedResponseError,
			errIs:      ErrPermissionDenied,
			respStatus: http.StatusForbidden,
			respBody:   fixture("permission_denied_error"),
		},
		{
			name: "validation error",
			args: args{
				ctx:      context.Background(),
				org:      OrganizationRef{ID: "org_O648YDMEYeLmqdmn"},
				zoneArgs: &DNSZoneArguments{Name: ""},
			},
			errStr:     fixtureValidationErrorErr,
			errResp:    fixtureValidationErrorResponseError,
			errIs:      ErrValidationError,
			respStatus: http.StatusUnprocessableEntity,
			respBody:   fixture("validation_error"),
		},
		{
			name: "nil zone arguments",
			args: args{
				ctx:      context.Background(),
				org:      OrganizationRef{ID: "org_O648YDMEYeLmqdmn"},
				zoneArgs: nil,
			},
			reqBody: &dnsZoneCreateRequest{
				Organization: OrganizationRef{ID: "org_O648YDMEYeLmqdmn"},
			}, errStr: fixtureValidationErrorErr,
			errResp:    fixtureValidationErrorResponseError,
			errIs:      ErrValidationError,
			respStatus: http.StatusUnprocessableEntity,
			respBody:   fixture("validation_error"),
		},
		{
			name: "nil context",
			args: args{
				ctx:      nil,
				org:      OrganizationRef{ID: "org_O648YDMEYeLmqdmn"},
				zoneArgs: &DNSZoneArguments{Name: "hi"},
			},
			errStr: "net/http: nil Context",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rm, mux, _, teardown := prepareTestClient(t)
			defer teardown()
			c := NewDNSZonesClient(rm)

			mux.HandleFunc(
				"/core/v1/organizations/_/dns/zones",
				func(w http.ResponseWriter, r *http.Request) {
					assert.Equal(t, "POST", r.Method)
					assertEmptyFieldSpec(t, r)
					assertAuthorization(t, r)

					if tt.reqBody != nil {
						reqBody := &dnsZoneCreateRequest{}
						err := strictUmarshal(r.Body, reqBody)
						assert.NoError(t, err)
						assert.Equal(t, tt.reqBody, reqBody)
					}

					w.WriteHeader(tt.respStatus)
					_, _ = w.Write(tt.respBody)
				},
			)

			got, resp, err := c.Create(
				tt.args.ctx, tt.args.org, tt.args.zoneArgs,
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

func TestDNSZonesClient_Delete(t *testing.T) {
	type args struct {
		ctx  context.Context
		zone DNSZoneRef
	}
	tests := []struct {
		name       string
		args       args
		want       *DNSZone
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
				ctx:  context.Background(),
				zone: DNSZoneRef{ID: "dnszone_k75eFc4UBOgeE5Zy"},
			},
			want: fixtureDNSZone,
			wantQuery: &url.Values{
				"dns_zone[id]": []string{"dnszone_k75eFc4UBOgeE5Zy"},
			},
			respStatus: http.StatusOK,
			respBody:   fixture("dns_zone_get"),
		},
		{
			name: "by Name",
			args: args{
				ctx:  context.Background(),
				zone: DNSZoneRef{Name: "test1.example.com"},
			},
			want: fixtureDNSZone,
			wantQuery: &url.Values{
				"dns_zone[name]": []string{"test1.example.com"},
			},
			respStatus: http.StatusOK,
			respBody:   fixture("dns_zone_get"),
		},
		{
			name: "non-existent DNS zone",
			args: args{
				ctx:  context.Background(),
				zone: DNSZoneRef{ID: "dnszone_k75eFc4UBOgeE5Zy"},
			},
			errStr:     fixtureDNSZoneNotFoundErr,
			errResp:    fixtureDNSZoneNotFoundResponseError,
			errIs:      ErrDNSZoneNotFound,
			respStatus: http.StatusNotFound,
			respBody:   fixture("dns_zone_not_found_error"),
		},
		{
			name: "nil context",
			args: args{
				ctx:  nil,
				zone: DNSZoneRef{ID: "dnszone_k75eFc4UBOgeE5Zy"},
			},
			errStr: "net/http: nil Context",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rm, mux, _, teardown := prepareTestClient(t)
			defer teardown()
			c := NewDNSZonesClient(rm)

			mux.HandleFunc(
				"/core/v1/dns/zones/_",
				func(w http.ResponseWriter, r *http.Request) {
					assert.Equal(t, "DELETE", r.Method)
					assertEmptyFieldSpec(t, r)
					assertAuthorization(t, r)

					if tt.wantQuery != nil {
						assert.Equal(t, *tt.wantQuery, r.URL.Query())
					} else {
						assert.Equal(t,
							*tt.args.zone.queryValues(), r.URL.Query(),
						)
					}

					w.WriteHeader(tt.respStatus)
					_, _ = w.Write(tt.respBody)
				},
			)

			got, resp, err := c.Delete(tt.args.ctx, tt.args.zone)

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

func TestDNSZonesClient_VerificationDetails(t *testing.T) {
	type args struct {
		ctx  context.Context
		zone DNSZoneRef
	}
	tests := []struct {
		name       string
		args       args
		want       *DNSZoneVerificationDetails
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
				ctx:  context.Background(),
				zone: DNSZoneRef{ID: "dnszone_k75eFc4UBOgeE5Zy"},
			},
			want: &DNSZoneVerificationDetails{
				Nameservers: []string{"ns1.katapult.io", "ns2.katapult.io"},
				TXTRecord:   "M0Y0SzE1TzNJTUZPSDRoQUV0TDZ4MEZwckFqbW1FNHI=",
			},
			wantQuery: &url.Values{
				"dns_zone[id]": []string{"dnszone_k75eFc4UBOgeE5Zy"},
			},
			respStatus: http.StatusOK,
			respBody:   fixture("dns_zone_verification_details"),
		},
		{
			name: "by Name",
			args: args{
				ctx:  context.Background(),
				zone: DNSZoneRef{Name: "test1.example.com"},
			},
			want: &DNSZoneVerificationDetails{
				Nameservers: []string{"ns1.katapult.io", "ns2.katapult.io"},
				TXTRecord:   "M0Y0SzE1TzNJTUZPSDRoQUV0TDZ4MEZwckFqbW1FNHI=",
			},
			wantQuery: &url.Values{
				"dns_zone[name]": []string{"test1.example.com"},
			},
			respStatus: http.StatusOK,
			respBody:   fixture("dns_zone_verification_details"),
		},
		{
			name: "non-existent DNS zone",
			args: args{
				ctx:  context.Background(),
				zone: DNSZoneRef{ID: "dnszone_k75eFc4UBOgeE5Zy"},
			},
			errStr:     fixtureDNSZoneNotFoundErr,
			errResp:    fixtureDNSZoneNotFoundResponseError,
			errIs:      ErrDNSZoneNotFound,
			respStatus: http.StatusNotFound,
			respBody:   fixture("dns_zone_not_found_error"),
		},
		{
			name: "already verified",
			args: args{
				ctx:  context.Background(),
				zone: DNSZoneRef{ID: "dnszone_k75eFc4UBOgeE5Zy"},
			},
			errStr: "katapult: unprocessable_entity: " +
				"dns_zone_already_verified: This DNS zone is already " +
				"verified, and does not require any verification details",
			errResp: &katapult.ResponseError{
				Code: "dns_zone_already_verified",
				Description: "This DNS zone is already verified, and does " +
					"not require any verification details",
				Detail: json.RawMessage(`{}`),
			},
			errIs:      ErrDNSZoneAlreadyVerified,
			respStatus: http.StatusUnprocessableEntity,
			respBody:   fixture("dns_zone_already_verified_error"),
		},
		{
			name: "infrastructure DNS zone cannot be edited",
			args: args{
				ctx:  context.Background(),
				zone: DNSZoneRef{ID: "dnszone_k75eFc4UBOgeE5Zy"},
			},
			errStr:     fixtureDNSZoneInfraErr,
			errResp:    fixtureDNSZoneInfraResponseError,
			errIs:      ErrInfrastructureDNSZoneCannotBeEdited,
			respStatus: http.StatusForbidden,
			respBody: fixture(
				"dns_zone_infrastructure_zone_cannot_be_edited",
			),
		},
		{
			name: "nil context",
			args: args{
				ctx:  nil,
				zone: DNSZoneRef{ID: "dnszone_k75eFc4UBOgeE5Zy"},
			},
			errStr: "net/http: nil Context",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rm, mux, _, teardown := prepareTestClient(t)
			defer teardown()
			c := NewDNSZonesClient(rm)

			mux.HandleFunc(
				"/core/v1/dns/zones/_/verification_details",
				func(w http.ResponseWriter, r *http.Request) {
					assert.Equal(t, "GET", r.Method)
					assertEmptyFieldSpec(t, r)
					assertAuthorization(t, r)

					if tt.wantQuery != nil {
						assert.Equal(t, *tt.wantQuery, r.URL.Query())
					} else {
						assert.Equal(t,
							*tt.args.zone.queryValues(), r.URL.Query(),
						)
					}

					w.WriteHeader(tt.respStatus)
					_, _ = w.Write(tt.respBody)
				},
			)

			got, resp, err := c.VerificationDetails(
				tt.args.ctx, tt.args.zone,
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

func TestDNSZonesClient_Verify(t *testing.T) {
	type args struct {
		ctx  context.Context
		zone DNSZoneRef
	}
	tests := []struct {
		name       string
		args       args
		reqBody    *dnsZoneVerifyRequest
		want       *DNSZone
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
				zone: DNSZoneRef{
					ID: "dnszone_k75eFc4UBOgeE5Zy",
				},
			},
			reqBody: &dnsZoneVerifyRequest{
				DNSZone: DNSZoneRef{ID: "dnszone_k75eFc4UBOgeE5Zy"},
			},
			want:       fixtureDNSZone,
			respStatus: http.StatusOK,
			respBody:   fixture("dns_zone_get"),
		},
		{
			name: "by Name",
			args: args{
				ctx: context.Background(),
				zone: DNSZoneRef{
					Name: "test1.example.com",
				},
			},
			reqBody: &dnsZoneVerifyRequest{
				DNSZone: DNSZoneRef{Name: "test1.example.com"},
			},
			want:       fixtureDNSZone,
			respStatus: http.StatusOK,
			respBody:   fixture("dns_zone_get"),
		},
		{
			name: "non-existent DNS zone",
			args: args{
				ctx:  context.Background(),
				zone: DNSZoneRef{ID: "dnszone_k75eFc4UBOgeE5Zy"},
			},
			errStr:     fixtureDNSZoneNotFoundErr,
			errResp:    fixtureDNSZoneNotFoundResponseError,
			errIs:      ErrDNSZoneNotFound,
			respStatus: http.StatusNotFound,
			respBody:   fixture("dns_zone_not_found_error"),
		},
		{
			name: "infrastructure DNS zone cannot be edited",
			args: args{
				ctx:  context.Background(),
				zone: DNSZoneRef{ID: "dnszone_k75eFc4UBOgeE5Zy"},
			},
			errStr:     fixtureDNSZoneInfraErr,
			errResp:    fixtureDNSZoneInfraResponseError,
			errIs:      ErrInfrastructureDNSZoneCannotBeEdited,
			respStatus: http.StatusUnprocessableEntity,
			respBody: fixture(
				"dns_zone_infrastructure_zone_cannot_be_edited",
			),
		},
		{
			name: "validation error",
			args: args{
				ctx:  context.Background(),
				zone: DNSZoneRef{ID: "dnszone_k75eFc4UBOgeE5Zy"},
			},
			errStr:     fixtureValidationErrorErr,
			errResp:    fixtureValidationErrorResponseError,
			errIs:      ErrValidationError,
			respStatus: http.StatusUnprocessableEntity,
			respBody:   fixture("validation_error"),
		},
		{
			name: "nil context",
			args: args{
				ctx:  nil,
				zone: DNSZoneRef{ID: "dnszone_k75eFc4UBOgeE5Zy"},
			},
			errStr: "net/http: nil Context",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rm, mux, _, teardown := prepareTestClient(t)
			defer teardown()
			c := NewDNSZonesClient(rm)

			mux.HandleFunc(
				"/core/v1/dns/zones/_/verify",
				func(w http.ResponseWriter, r *http.Request) {
					assert.Equal(t, "POST", r.Method)
					assertEmptyFieldSpec(t, r)
					assertAuthorization(t, r)

					if tt.reqBody != nil {
						reqBody := &dnsZoneVerifyRequest{}
						err := strictUmarshal(r.Body, reqBody)
						assert.NoError(t, err)
						assert.Equal(t, tt.reqBody, reqBody)
					}

					w.WriteHeader(tt.respStatus)
					_, _ = w.Write(tt.respBody)
				},
			)

			got, resp, err := c.Verify(tt.args.ctx, tt.args.zone)

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

func TestDNSZonesClient_UpdateTTL(t *testing.T) {
	type args struct {
		ctx  context.Context
		zone DNSZoneRef
		ttl  int
	}
	tests := []struct {
		name       string
		args       args
		reqBody    *dnsZoneUpdateTTLRequest
		want       *DNSZone
		errStr     string
		errResp    *katapult.ResponseError
		errIs      error
		respStatus int
		respBody   []byte
	}{
		{
			name: "by ID",
			args: args{
				ctx:  context.Background(),
				zone: DNSZoneRef{ID: "dnszone_lwz66kyviwCQyqQc"},
				ttl:  1842,
			},
			reqBody: &dnsZoneUpdateTTLRequest{
				DNSZone: DNSZoneRef{ID: "dnszone_lwz66kyviwCQyqQc"},
				TTL:     1842,
			},
			want: &DNSZone{
				ID:   "dnszone_lwz66kyviwCQyqQc",
				Name: "test-2.example.com",
				TTL:  1842,
			},
			respStatus: http.StatusCreated,
			respBody:   fixture("dns_zone_update_ttl"),
		},
		{
			name: "by Name",
			args: args{
				ctx:  context.Background(),
				zone: DNSZoneRef{Name: "test1.example.come"},
				ttl:  1842,
			},
			reqBody: &dnsZoneUpdateTTLRequest{
				DNSZone: DNSZoneRef{Name: "test1.example.come"},
				TTL:     1842,
			},
			want: &DNSZone{
				ID:   "dnszone_lwz66kyviwCQyqQc",
				Name: "test-2.example.com",
				TTL:  1842,
			},
			respStatus: http.StatusCreated,
			respBody:   fixture("dns_zone_update_ttl"),
		},
		{
			name: "high TTL",
			args: args{
				ctx:  context.Background(),
				zone: DNSZoneRef{ID: "dnszone_lwz66kyviwCQyqQc"},
				ttl:  25200,
			},
			reqBody: &dnsZoneUpdateTTLRequest{
				DNSZone: DNSZoneRef{ID: "dnszone_lwz66kyviwCQyqQc"},
				TTL:     25200,
			},
			respStatus: http.StatusCreated,
			respBody:   fixture("dns_zone_update_ttl"),
		},
		{
			name: "non-existent DNS zone",
			args: args{
				ctx:  context.Background(),
				zone: DNSZoneRef{ID: "dnszone_lwz66kyviwCQyqQc"},
				ttl:  1842,
			},
			errStr:     fixtureDNSZoneNotFoundErr,
			errResp:    fixtureDNSZoneNotFoundResponseError,
			errIs:      ErrDNSZoneNotFound,
			respStatus: http.StatusNotFound,
			respBody:   fixture("dns_zone_not_found_error"),
		},
		{
			name: "permission denied",
			args: args{
				ctx:  context.Background(),
				zone: DNSZoneRef{ID: "dnszone_lwz66kyviwCQyqQc"},
				ttl:  600,
			},
			errStr:     fixturePermissionDeniedErr,
			errResp:    fixturePermissionDeniedResponseError,
			errIs:      ErrPermissionDenied,
			respStatus: http.StatusForbidden,
			respBody:   fixture("permission_denied_error"),
		},
		{
			name: "validation error",
			args: args{
				ctx:  context.Background(),
				zone: DNSZoneRef{ID: "dnszone_lwz66kyviwCQyqQc"},
				ttl:  600,
			},
			errStr:     fixtureValidationErrorErr,
			errResp:    fixtureValidationErrorResponseError,
			errIs:      ErrValidationError,
			respStatus: http.StatusUnprocessableEntity,
			respBody:   fixture("validation_error"),
		},
		{
			name: "nil context",
			args: args{
				ctx:  nil,
				zone: DNSZoneRef{ID: "dnszone_lwz66kyviwCQyqQc"},
				ttl:  600,
			},
			errStr: "net/http: nil Context",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rm, mux, _, teardown := prepareTestClient(t)
			defer teardown()
			c := NewDNSZonesClient(rm)

			mux.HandleFunc(
				"/core/v1/dns/zones/_/update_ttl",
				func(w http.ResponseWriter, r *http.Request) {
					assert.Equal(t, "POST", r.Method)
					assertEmptyFieldSpec(t, r)
					assertAuthorization(t, r)

					if tt.reqBody != nil {
						reqBody := &dnsZoneUpdateTTLRequest{}
						err := strictUmarshal(r.Body, reqBody)
						assert.NoError(t, err)
						assert.Equal(t, tt.reqBody, reqBody)
					}

					w.WriteHeader(tt.respStatus)
					_, _ = w.Write(tt.respBody)
				},
			)

			got, resp, err := c.UpdateTTL(
				tt.args.ctx, tt.args.zone, tt.args.ttl,
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
