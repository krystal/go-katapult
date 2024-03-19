package core

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"testing"

	"github.com/krystal/go-katapult"
	"github.com/krystal/go-katapult/internal/testclient"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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

	fixtureDNSZoneNotVerifiedErr = "katapult: unprocessable_entity: " +
		"dns_zone_not_verified: The DNS zone could not be verified, check " +
		"the nameservers are set correctly"
	fixtureDNSZoneNotVerifiedResponseError = &katapult.ResponseError{
		Code: "dns_zone_not_verified",
		Description: "The DNS zone could not be verified, check the " +
			"nameservers are set correctly",
		Detail: json.RawMessage(`{}`),
	}

	// Correlates to fixtures/dns_zone_get.json.
	fixtureDNSZone = &DNSZone{
		ID:         "dnszone_k75eFc4UBOgeE5Zy",
		Name:       "test1.example.com",
		DefaultTTL: 3600,
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
				ID:         "dnszone_k75eFc4UBOgeE5Zy",
				Name:       "test1.example.com",
				DefaultTTL: 343,
				Verified:   true,
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

func TestDNSZoneCreateArguments_JSONMarshaling(t *testing.T) {
	tests := []struct {
		name string
		obj  *DNSZoneCreateArguments
	}{
		{
			name: "empty",
			obj:  &DNSZoneCreateArguments{},
		},
		{
			name: "full",
			obj: &DNSZoneCreateArguments{
				Name:       "name",
				DefaultTTL: 493,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testJSONMarshaling(t, tt.obj)
		})
	}
}

func TestDNSZoneUpdateArguments_JSONMarshaling(t *testing.T) {
	tests := []struct {
		name string
		obj  *DNSZoneUpdateArguments
	}{
		{
			name: "empty",
			obj:  &DNSZoneUpdateArguments{},
		},
		{
			name: "full",
			obj: &DNSZoneUpdateArguments{
				Name:       "name",
				DefaultTTL: 493,
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
				Organization: OrganizationRef{ID: "org_QwNl81npdQQGinmt"},
				Properties: &DNSZoneCreateArguments{
					Name: "test1.example.com",
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

func Test_dnsZoneUpdateRequest_JSONMarshaling(t *testing.T) {
	tests := []struct {
		name string
		obj  *dnsZoneUpdateRequest
	}{
		{
			name: "empty",
			obj:  &dnsZoneUpdateRequest{},
		},
		{
			name: "full",
			obj: &dnsZoneUpdateRequest{
				DNSZone:    DNSZoneRef{ID: "dnszone_JH38nIy3murRa1c9"},
				Properties: &DNSZoneUpdateArguments{Name: "test1.example.com"},
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
				Pagination:  &katapult.Pagination{CurrentPage: 934},
				DNSZones:    []*DNSZone{{ID: "id1"}},
				DNSZone:     &DNSZone{ID: "id2"},
				Deleted:     boolPtr(true),
				Nameservers: []string{"ns1.foo.bar", "ns2.foo.bar"},
			},
		},
		{
			name: "deleted false",
			obj: &dnsZoneResponseBody{
				Deleted: boolPtr(false),
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
			ID:         "dnszone_k75eFc4UBOgeE5Zy",
			Name:       "test1.example.com",
			DefaultTTL: 3600,
		},
		{
			ID:         "dnszone_lwz66kyviwCQyqQc",
			Name:       "test-2.example.com",
			DefaultTTL: 3600,
		},
		{
			ID:         "dnszone_qr9KPhSwkGNh7IMb",
			Name:       "test-3.example.com",
			DefaultTTL: 3600,
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
				"/core/v1/organizations/_/dns_zones",
				func(w http.ResponseWriter, r *http.Request) {
					assert.Equal(t, "GET", r.Method)
					assertEmptyFieldSpec(t, r)
					assertAuthorization(t, r)
					assertRequestOptionHeader(t, r)

					qs := queryValues(tt.args.org, tt.args.opts)
					assert.Equal(t, *qs, r.URL.Query())

					w.WriteHeader(tt.respStatus)
					_, _ = w.Write(tt.respBody)
				},
			)

			got, resp, err := c.List(
				tt.args.ctx, tt.args.org, tt.args.opts, testRequestOption,
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

func TestDNSZonesClient_Nameservers(t *testing.T) {
	type args struct {
		ctx context.Context
		org OrganizationRef
	}
	tests := []struct {
		name       string
		args       args
		want       []string
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
			want:       []string{"ns1.foo.bar", "ns2.foo.bar"},
			respStatus: http.StatusOK,
			respBody:   fixture("dns_zone_nameservers"),
		},
		{
			name: "by organization SubDomain",
			args: args{
				ctx: context.Background(),
				org: OrganizationRef{SubDomain: "blackmesa"},
			},
			want:       []string{"ns1.foo.bar", "ns2.foo.bar"},
			respStatus: http.StatusOK,
			respBody:   fixture("dns_zone_nameservers"),
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
				"/core/v1/organizations/_/dns_zones/nameservers",
				func(w http.ResponseWriter, r *http.Request) {
					assert.Equal(t, "GET", r.Method)
					assertEmptyFieldSpec(t, r)
					assertAuthorization(t, r)
					assertRequestOptionHeader(t, r)

					assert.Equal(t, *tt.args.org.queryValues(), r.URL.Query())

					w.WriteHeader(tt.respStatus)
					_, _ = w.Write(tt.respBody)
				},
			)

			got, resp, err := c.Nameservers(
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
				"/core/v1/dns_zones/_",
				func(w http.ResponseWriter, r *http.Request) {
					assert.Equal(t, "GET", r.Method)
					assertEmptyFieldSpec(t, r)
					assertAuthorization(t, r)
					assertRequestOptionHeader(t, r)

					assert.Equal(t, *tt.args.ref.queryValues(), r.URL.Query())

					w.WriteHeader(tt.respStatus)
					_, _ = w.Write(tt.respBody)
				},
			)

			got, resp, err := c.Get(tt.args.ctx, tt.args.ref, testRequestOption)

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

			mux.HandleFunc("/core/v1/dns_zones/_",
				func(w http.ResponseWriter, r *http.Request) {
					assert.Equal(t, "GET", r.Method)
					assertEmptyFieldSpec(t, r)
					assertAuthorization(t, r)
					assertRequestOptionHeader(t, r)

					assert.Equal(t, url.Values{
						"dns_zone[id]": []string{tt.args.id},
					}, r.URL.Query())

					w.WriteHeader(tt.respStatus)
					_, _ = w.Write(tt.respBody)
				},
			)

			got, resp, err := c.GetByID(
				tt.args.ctx,
				tt.args.id,
				testRequestOption,
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

			mux.HandleFunc("/core/v1/dns_zones/_",
				func(w http.ResponseWriter, r *http.Request) {
					assert.Equal(t, "GET", r.Method)
					assertEmptyFieldSpec(t, r)
					assertAuthorization(t, r)
					assertRequestOptionHeader(t, r)

					qs := url.Values{"dns_zone[name]": []string{tt.args.name}}
					assert.Equal(t, qs, r.URL.Query())

					w.WriteHeader(tt.respStatus)
					_, _ = w.Write(tt.respBody)
				},
			)

			got, resp, err := c.GetByName(
				tt.args.ctx,
				tt.args.name,
				testRequestOption,
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

func TestDNSZonesClient_Create(t *testing.T) {
	type args struct {
		ctx  context.Context
		org  OrganizationRef
		args *DNSZoneCreateArguments
	}
	tests := []struct {
		name        string
		orgID       string
		args        args
		wantQuery   *url.Values
		wantReqBody *dnsZoneCreateRequest
		want        *DNSZone
		errStr      string
		errResp     *katapult.ResponseError
		errIs       error
		respStatus  int
		respBody    []byte
	}{
		{
			name: "create a DNS zone by organization ID",
			args: args{
				ctx: context.Background(),
				org: OrganizationRef{ID: "org_O648YDMEYeLmqdmn"},
				args: &DNSZoneCreateArguments{
					Name:       "test-1.com",
					DefaultTTL: 1800,
				},
			},
			wantQuery: &url.Values{
				"organization[id]": []string{"org_O648YDMEYeLmqdmn"},
			},
			wantReqBody: &dnsZoneCreateRequest{
				Properties: &DNSZoneCreateArguments{
					Name:       "test-1.com",
					DefaultTTL: 1800,
				},
			},
			want: &DNSZone{
				ID:         "dnszone_yqflWVIdu5vnirLq",
				Name:       "test-1.com",
				DefaultTTL: 1800,
			},
			respStatus: http.StatusCreated,
			respBody:   fixture("dns_zone_create"),
		},
		{
			name: "create a DNS zone by organization SubDomain",
			args: args{
				ctx: context.Background(),
				org: OrganizationRef{SubDomain: "acme-inc"},
				args: &DNSZoneCreateArguments{
					Name:       "test-1.com",
					DefaultTTL: 1800,
				},
			},
			wantQuery: &url.Values{
				"organization[sub_domain]": []string{"acme-inc"},
			},
			wantReqBody: &dnsZoneCreateRequest{
				Properties: &DNSZoneCreateArguments{
					Name:       "test-1.com",
					DefaultTTL: 1800,
				},
			},
			want: &DNSZone{
				ID:         "dnszone_yqflWVIdu5vnirLq",
				Name:       "test-1.com",
				DefaultTTL: 1800,
			},
			respStatus: http.StatusCreated,
			respBody:   fixture("dns_zone_create"),
		},
		{
			name: "without default TTL",
			args: args{
				ctx:  context.Background(),
				org:  OrganizationRef{ID: "org_O648YDMEYeLmqdmn"},
				args: &DNSZoneCreateArguments{Name: "test-1.com"},
			},
			wantQuery: &url.Values{
				"organization[id]": []string{"org_O648YDMEYeLmqdmn"},
			},
			wantReqBody: &dnsZoneCreateRequest{
				Properties: &DNSZoneCreateArguments{
					Name: "test-1.com",
				},
			},
			respStatus: http.StatusCreated,
			respBody:   fixture("dns_zone_create"),
		},
		{
			name: "without TTL",
			args: args{
				ctx:  context.Background(),
				org:  OrganizationRef{ID: "org_O648YDMEYeLmqdmn"},
				args: &DNSZoneCreateArguments{Name: "test-1.com"},
			},
			respStatus: http.StatusCreated,
			respBody:   fixture("dns_zone_create"),
		},
		{
			name: "non-existent Organization",
			args: args{
				ctx:  context.Background(),
				org:  OrganizationRef{ID: "org_O648YDMEYeLmqdmn"},
				args: &DNSZoneCreateArguments{Name: "test-1.com"},
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
				ctx:  context.Background(),
				org:  OrganizationRef{ID: "org_O648YDMEYeLmqdmn"},
				args: &DNSZoneCreateArguments{Name: "test-1.com"},
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
				ctx:  context.Background(),
				org:  OrganizationRef{ID: "org_O648YDMEYeLmqdmn"},
				args: &DNSZoneCreateArguments{Name: "test-1.com"},
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
				org:  OrganizationRef{ID: "org_O648YDMEYeLmqdmn"},
				args: &DNSZoneCreateArguments{Name: ""},
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
				ctx:  context.Background(),
				org:  OrganizationRef{ID: "org_O648YDMEYeLmqdmn"},
				args: nil,
			},
			wantQuery: &url.Values{
				"organization[id]": []string{"org_O648YDMEYeLmqdmn"},
			},
			wantReqBody: &dnsZoneCreateRequest{},
			errStr:      fixtureValidationErrorErr,
			errResp:     fixtureValidationErrorResponseError,
			errIs:       ErrValidationError,
			respStatus:  http.StatusUnprocessableEntity,
			respBody:    fixture("validation_error"),
		},
		{
			name: "nil context",
			args: args{
				ctx:  nil,
				org:  OrganizationRef{ID: "org_O648YDMEYeLmqdmn"},
				args: &DNSZoneCreateArguments{Name: "hi"},
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
				"/core/v1/organizations/_/dns_zones",
				func(w http.ResponseWriter, r *http.Request) {
					assert.Equal(t, "POST", r.Method)
					assertEmptyFieldSpec(t, r)
					assertAuthorization(t, r)
					assertRequestOptionHeader(t, r)

					assert.Equal(t, *tt.args.org.queryValues(), r.URL.Query())

					if tt.wantReqBody != nil {
						reqBody := &dnsZoneCreateRequest{}
						err := strictUmarshal(r.Body, reqBody)
						assert.NoError(t, err)
						assert.Equal(t, tt.wantReqBody, reqBody)
					}

					w.WriteHeader(tt.respStatus)
					_, _ = w.Write(tt.respBody)
				},
			)

			got, resp, err := c.Create(
				tt.args.ctx, tt.args.org, tt.args.args, testRequestOption,
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

func TestDNSZonesClient_Update(t *testing.T) {
	type args struct {
		ctx  context.Context
		zone DNSZoneRef
		args *DNSZoneUpdateArguments
	}
	tests := []struct {
		name        string
		orgID       string
		args        args
		wantQuery   *url.Values
		wantReqBody *dnsZoneUpdateRequest
		want        *DNSZone
		errStr      string
		errResp     *katapult.ResponseError
		errIs       error
		respStatus  int
		respBody    []byte
	}{
		{
			name: "by ID",
			args: args{
				ctx:  context.Background(),
				zone: DNSZoneRef{ID: "dnszone_lwz66kyviwCQyqQc"},
				args: &DNSZoneUpdateArguments{
					Name:       "test-1.com",
					DefaultTTL: 1842,
				},
			},
			wantQuery: &url.Values{
				"dns_zone[id]": []string{"dnszone_lwz66kyviwCQyqQc"},
			},
			wantReqBody: &dnsZoneUpdateRequest{
				Properties: &DNSZoneUpdateArguments{
					Name:       "test-1.com",
					DefaultTTL: 1842,
				},
			},
			want: &DNSZone{
				ID:         "dnszone_lwz66kyviwCQyqQc",
				Name:       "test-1.com",
				DefaultTTL: 1842,
			},
			respStatus: http.StatusOK,
			respBody:   fixture("dns_zone_update"),
		},
		{
			name: "by Name",
			args: args{
				ctx:  context.Background(),
				zone: DNSZoneRef{Name: "test1.example.com"},
				args: &DNSZoneUpdateArguments{
					Name:       "test-1.com",
					DefaultTTL: 1842,
				},
			},
			wantQuery: &url.Values{
				"dns_zone[name]": []string{"test1.example.com"},
			},
			wantReqBody: &dnsZoneUpdateRequest{
				Properties: &DNSZoneUpdateArguments{
					Name:       "test-1.com",
					DefaultTTL: 1842,
				},
			},
			want: &DNSZone{
				ID:         "dnszone_lwz66kyviwCQyqQc",
				Name:       "test-1.com",
				DefaultTTL: 1842,
			},
			respStatus: http.StatusOK,
			respBody:   fixture("dns_zone_update"),
		},
		{
			name: "without default TTL",
			args: args{
				ctx:  context.Background(),
				zone: DNSZoneRef{ID: "dnszone_lwz66kyviwCQyqQc"},
				args: &DNSZoneUpdateArguments{Name: "test-1.com"},
			},
			wantQuery: &url.Values{
				"dns_zone[id]": []string{"dnszone_lwz66kyviwCQyqQc"},
			},
			wantReqBody: &dnsZoneUpdateRequest{
				Properties: &DNSZoneUpdateArguments{
					Name: "test-1.com",
				},
			},
			respStatus: http.StatusOK,
			respBody:   fixture("dns_zone_update"),
		},
		{
			name: "without name",
			args: args{
				ctx:  context.Background(),
				zone: DNSZoneRef{ID: "dnszone_lwz66kyviwCQyqQc"},
				args: &DNSZoneUpdateArguments{
					DefaultTTL: 1842,
				},
			},
			wantQuery: &url.Values{
				"dns_zone[id]": []string{"dnszone_lwz66kyviwCQyqQc"},
			},
			wantReqBody: &dnsZoneUpdateRequest{
				Properties: &DNSZoneUpdateArguments{
					DefaultTTL: 1842,
				},
			},
			want: &DNSZone{
				ID:         "dnszone_lwz66kyviwCQyqQc",
				Name:       "test-1.com",
				DefaultTTL: 1842,
			},
			respStatus: http.StatusOK,
			respBody:   fixture("dns_zone_update"),
		},
		{
			name: "non-existent Organization",
			args: args{
				ctx:  context.Background(),
				zone: DNSZoneRef{ID: "dnszone_lwz66kyviwCQyqQc"},
				args: &DNSZoneUpdateArguments{Name: "test-1.com"},
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
				ctx:  context.Background(),
				zone: DNSZoneRef{ID: "dnszone_lwz66kyviwCQyqQc"},
				args: &DNSZoneUpdateArguments{Name: "test-1.com"},
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
				ctx:  context.Background(),
				zone: DNSZoneRef{ID: "dnszone_lwz66kyviwCQyqQc"},
				args: &DNSZoneUpdateArguments{Name: "test-1.com"},
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
				args: &DNSZoneUpdateArguments{Name: ""},
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
				ctx:  context.Background(),
				zone: DNSZoneRef{ID: "dnszone_lwz66kyviwCQyqQc"},
				args: nil,
			},
			wantQuery: &url.Values{
				"dns_zone[id]": []string{"dnszone_lwz66kyviwCQyqQc"},
			},
			wantReqBody: &dnsZoneUpdateRequest{},
			errStr:      fixtureValidationErrorErr,
			errResp:     fixtureValidationErrorResponseError,
			errIs:       ErrValidationError,
			respStatus:  http.StatusUnprocessableEntity,
			respBody:    fixture("validation_error"),
		},
		{
			name: "nil context",
			args: args{
				ctx:  nil,
				zone: DNSZoneRef{ID: "dnszone_lwz66kyviwCQyqQc"},
				args: &DNSZoneUpdateArguments{Name: "hi"},
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
				"/core/v1/dns_zones/_",
				func(w http.ResponseWriter, r *http.Request) {
					assert.Equal(t, "PATCH", r.Method)
					assertEmptyFieldSpec(t, r)
					assertAuthorization(t, r)
					assertRequestOptionHeader(t, r)

					assert.Equal(t, *tt.args.zone.queryValues(), r.URL.Query())

					if tt.wantReqBody != nil {
						reqBody := &dnsZoneUpdateRequest{}
						err := strictUmarshal(r.Body, reqBody)
						assert.NoError(t, err)
						assert.Equal(t, tt.wantReqBody, reqBody)
					}

					w.WriteHeader(tt.respStatus)
					_, _ = w.Write(tt.respBody)
				},
			)

			got, resp, err := c.Update(
				tt.args.ctx, tt.args.zone, tt.args.args, testRequestOption,
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
		want       *bool
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
			want: boolPtr(true),
			wantQuery: &url.Values{
				"dns_zone[id]": []string{"dnszone_k75eFc4UBOgeE5Zy"},
			},
			respStatus: http.StatusOK,
			respBody:   fixture("dns_zone_deleted"),
		},
		{
			name: "by Name",
			args: args{
				ctx:  context.Background(),
				zone: DNSZoneRef{Name: "test1.example.com"},
			},
			want: boolPtr(true),
			wantQuery: &url.Values{
				"dns_zone[name]": []string{"test1.example.com"},
			},
			respStatus: http.StatusOK,
			respBody:   fixture("dns_zone_deleted"),
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
				"/core/v1/dns_zones/_",
				func(w http.ResponseWriter, r *http.Request) {
					assert.Equal(t, "DELETE", r.Method)
					assertEmptyFieldSpec(t, r)
					assertAuthorization(t, r)
					assertRequestOptionHeader(t, r)

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

			got, resp, err := c.Delete(
				tt.args.ctx,
				tt.args.zone,
				testRequestOption,
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
		wantQuery  *url.Values
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
			wantQuery: &url.Values{
				"dns_zone[id]": []string{"dnszone_k75eFc4UBOgeE5Zy"},
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
			wantQuery: &url.Values{
				"dns_zone[name]": []string{"test1.example.com"},
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
			name: "DNS zone not verified",
			args: args{
				ctx:  context.Background(),
				zone: DNSZoneRef{ID: "dnszone_k75eFc4UBOgeE5Zy"},
			},
			errStr:     fixtureDNSZoneNotVerifiedErr,
			errResp:    fixtureDNSZoneNotVerifiedResponseError,
			errIs:      ErrDNSZoneNotVerified,
			respStatus: http.StatusUnprocessableEntity,
			respBody:   fixture("dns_zone_not_verified_error"),
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
				"/core/v1/dns_zones/_/verify",
				func(w http.ResponseWriter, r *http.Request) {
					assert.Equal(t, "POST", r.Method)
					assertEmptyFieldSpec(t, r)
					assertAuthorization(t, r)
					assertRequestOptionHeader(t, r)

					if tt.wantQuery != nil {
						assert.Equal(t, *tt.wantQuery, r.URL.Query())
					} else {
						assert.Equal(t,
							*tt.args.zone.queryValues(), r.URL.Query(),
						)
					}

					b, err := io.ReadAll(r.Body)
					require.NoError(t, err)
					defer r.Body.Close()

					assert.Empty(t, b)

					w.WriteHeader(tt.respStatus)
					_, _ = w.Write(tt.respBody)
				},
			)

			got, resp, err := c.Verify(
				tt.args.ctx, tt.args.zone, testRequestOption,
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
