package katapult

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	fixtureDNSZoneNotFoundErr = "dns_zone_not_found: No DNS zone was found " +
		"matching any of the criteria provided in the arguments"
	fixtureDNSZoneNotFoundResponseError = &ResponseError{
		Code: "dns_zone_not_found",
		Description: "No DNS zone was found matching any of the criteria " +
			"provided in the arguments",
		Detail: json.RawMessage(`{}`),
	}

	fixtureDNSZoneInfraErr = "infrastructure_dns_zone_cannot_be_edited: " +
		"Infrastructure DNS zones cannot be edited through the API. " +
		"These are managed exclusively by Katapult."
	fixtureDNSZoneInfraResponseError = &ResponseError{
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
				ID:                 "id",
				Name:               "name",
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

func Test_createDNSZoneRequest_JSONMarshaling(t *testing.T) {
	tests := []struct {
		name string
		obj  *createDNSZoneRequest
	}{
		{
			name: "empty",
			obj:  &createDNSZoneRequest{},
		},
		{
			name: "full",
			obj: &createDNSZoneRequest{
				Details:         &DNSZoneDetails{Name: "name"},
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

func Test_updateDNSZoneTTLRequest_JSONMarshaling(t *testing.T) {
	tests := []struct {
		name string
		obj  *updateDNSZoneTTLRequest
	}{
		{
			name: "empty",
			obj:  &updateDNSZoneTTLRequest{},
		},
		{
			name: "full",
			obj: &updateDNSZoneTTLRequest{
				TTL: 8384,
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
				Pagination: &Pagination{CurrentPage: 934},
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

func TestDNSZonesResource_List(t *testing.T) {
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
		ctx   context.Context
		orgID string
		opts  *ListOptions
	}
	tests := []struct {
		name       string
		args       args
		expected   []*DNSZone
		pagination *Pagination
		errStr     string
		errResp    *ResponseError
		respStatus int
		respBody   []byte
	}{
		{
			name: "DNS zones",
			args: args{
				ctx:   context.Background(),
				orgID: "org_O648YDMEYeLmqdmn",
			},
			expected: dnsZonesList,
			pagination: &Pagination{
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
			name: "page 1 of DNS zones",
			args: args{
				ctx:   context.Background(),
				orgID: "org_O648YDMEYeLmqdmn",
				opts:  &ListOptions{Page: 1, PerPage: 2},
			},
			expected: dnsZonesList[0:2],
			pagination: &Pagination{
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
			name: "page 2 of DNS zones",
			args: args{
				ctx:   context.Background(),
				orgID: "org_O648YDMEYeLmqdmn",
				opts:  &ListOptions{Page: 2, PerPage: 2},
			},
			expected: dnsZonesList[2:],
			pagination: &Pagination{
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
				ctx:   context.Background(),
				orgID: "org_O648YDMEYeLmqdmn",
			},
			errStr:     fixtureInvalidAPITokenErr,
			errResp:    fixtureInvalidAPITokenResponseError,
			respStatus: http.StatusForbidden,
			respBody:   fixture("invalid_api_token_error"),
		},
		{
			name: "non-existent Organization",
			args: args{
				ctx:   context.Background(),
				orgID: "org_O648YDMEYeLmqdmn",
			},
			errStr:     fixtureOrganizationNotFoundErr,
			errResp:    fixtureOrganizationNotFoundResponseError,
			respStatus: http.StatusNotFound,
			respBody:   fixture("organization_not_found_error"),
		},
		{
			name: "suspended Organization",
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
			name: "permission denied",
			args: args{
				ctx:   context.Background(),
				orgID: "org_O648YDMEYeLmqdmn",
			},
			errStr:     fixturePermissionDeniedErr,
			errResp:    fixturePermissionDeniedResponseError,
			respStatus: http.StatusForbidden,
			respBody:   fixture("permission_denied_error"),
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
					"/core/v1/organizations/%s/dns/zones", tt.args.orgID,
				),
				func(w http.ResponseWriter, r *http.Request) {
					assert.Equal(t, "GET", r.Method)
					assert.Equal(t, "", r.Header.Get("X-Field-Spec"))
					if tt.args.opts != nil {
						assert.Equal(t, *tt.args.opts.Values(), r.URL.Query())
					}

					w.WriteHeader(tt.respStatus)
					_, _ = w.Write(tt.respBody)
				},
			)

			got, resp, err := c.DNSZones.List(
				tt.args.ctx, tt.args.orgID, tt.args.opts,
			)

			if tt.respStatus != 0 {
				assert.Equal(t, tt.respStatus, resp.StatusCode)
			}

			if tt.errStr == "" {
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, err, tt.errStr)
			}

			if tt.expected != nil {
				assert.Equal(t, tt.expected, got)
			}

			if tt.pagination != nil {
				assert.Equal(t, tt.pagination, resp.Pagination)
			}

			if tt.errResp != nil {
				assert.Equal(t, tt.errResp, resp.Error)
			}
		})
	}
}

func TestDNSZonesResource_Get(t *testing.T) {
	type args struct {
		ctx context.Context
		id  string
	}
	tests := []struct {
		name       string
		args       args
		expected   *DNSZone
		errStr     string
		errResp    *ResponseError
		respStatus int
		respBody   []byte
	}{
		{
			name: "DNS zone",
			args: args{
				ctx: context.Background(),
				id:  "dnszone_k75eFc4UBOgeE5Zy",
			},
			expected:   fixtureDNSZone,
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
			c, mux, _, teardown := prepareTestClient()
			defer teardown()

			mux.HandleFunc(fmt.Sprintf("/core/v1/dns/zones/%s", tt.args.id),
				func(w http.ResponseWriter, r *http.Request) {
					assert.Equal(t, "GET", r.Method)
					assert.Equal(t, "", r.Header.Get("X-Field-Spec"))

					w.WriteHeader(tt.respStatus)
					_, _ = w.Write(tt.respBody)
				},
			)

			got, resp, err := c.DNSZones.Get(tt.args.ctx, tt.args.id)

			if tt.respStatus != 0 {
				assert.Equal(t, tt.respStatus, resp.StatusCode)
			}

			if tt.errStr == "" {
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, err, tt.errStr)
			}

			if tt.expected != nil {
				assert.Equal(t, tt.expected, got)
			}

			if tt.errResp != nil {
				assert.Equal(t, tt.errResp, resp.Error)
			}
		})
	}
}

func TestDNSZonesResource_GetByName(t *testing.T) {
	type args struct {
		ctx  context.Context
		name string
	}
	tests := []struct {
		name       string
		args       args
		expected   *DNSZone
		errStr     string
		errResp    *ResponseError
		respStatus int
		respBody   []byte
	}{
		{
			name: "DNS zone",
			args: args{
				ctx:  context.Background(),
				name: "test1.example.com",
			},
			expected:   fixtureDNSZone,
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
			c, mux, _, teardown := prepareTestClient()
			defer teardown()

			mux.HandleFunc("/core/v1/dns/zones/_",
				func(w http.ResponseWriter, r *http.Request) {
					assert.Equal(t, "GET", r.Method)
					assert.Equal(t, "", r.Header.Get("X-Field-Spec"))

					qs := url.Values{"dns_zone[name]": []string{tt.args.name}}
					assert.Equal(t, qs, r.URL.Query())

					w.WriteHeader(tt.respStatus)
					_, _ = w.Write(tt.respBody)
				},
			)

			got, resp, err := c.DNSZones.GetByName(tt.args.ctx, tt.args.name)

			if tt.respStatus != 0 {
				assert.Equal(t, tt.respStatus, resp.StatusCode)
			}

			if tt.errStr == "" {
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, err, tt.errStr)
			}

			if tt.expected != nil {
				assert.Equal(t, tt.expected, got)
			}

			if tt.errResp != nil {
				assert.Equal(t, tt.errResp, resp.Error)
			}
		})
	}
}

func TestDNSZonesResource_Create(t *testing.T) {
	type reqBodyDetails struct {
		Name string `json:"name"`
		TTL  int    `json:"ttl,omitempty"`
	}
	type reqBody struct {
		Details          *reqBodyDetails `json:"details"`
		SkipVerification bool            `json:"skip_verification"`
	}
	type args struct {
		ctx      context.Context
		orgID    string
		zoneArgs *DNSZoneArguments
	}
	tests := []struct {
		name       string
		orgID      string
		args       args
		expected   *DNSZone
		errStr     string
		errResp    *ResponseError
		respStatus int
		respBody   []byte
	}{
		{
			name: "create a DNS zone",
			args: args{
				ctx:   context.Background(),
				orgID: "org_O648YDMEYeLmqdmn",
				zoneArgs: &DNSZoneArguments{
					Name: "test-1.com",
					TTL:  1800,
				},
			},
			expected: &DNSZone{
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
				ctx:   context.Background(),
				orgID: "org_O648YDMEYeLmqdmn",
				zoneArgs: &DNSZoneArguments{
					Name:            "test-1.com",
					TTL:             1800,
					SkipVerfication: true,
				},
			},
			respStatus: http.StatusCreated,
			respBody:   fixture("dns_zone_create"),
		},
		{
			name: "without TTL",
			args: args{
				ctx:      context.Background(),
				orgID:    "org_O648YDMEYeLmqdmn",
				zoneArgs: &DNSZoneArguments{Name: "test-1.com"},
			},
			respStatus: http.StatusCreated,
			respBody:   fixture("dns_zone_create"),
		},
		{
			name: "without TTL",
			args: args{
				ctx:      context.Background(),
				orgID:    "org_O648YDMEYeLmqdmn",
				zoneArgs: &DNSZoneArguments{Name: "test-1.com"},
			},
			respStatus: http.StatusCreated,
			respBody:   fixture("dns_zone_create"),
		},
		{
			name: "non-existent Organization",
			args: args{
				ctx:      context.Background(),
				orgID:    "org_O648YDMEYeLmqdmn",
				zoneArgs: &DNSZoneArguments{Name: "test-1.com"},
			},
			errStr:     fixtureOrganizationNotFoundErr,
			errResp:    fixtureOrganizationNotFoundResponseError,
			respStatus: http.StatusNotFound,
			respBody:   fixture("organization_not_found_error"),
		},
		{
			name: "suspended Organization",
			args: args{
				ctx:      context.Background(),
				orgID:    "org_O648YDMEYeLmqdmn",
				zoneArgs: &DNSZoneArguments{Name: "test-1.com"},
			},
			errStr:     fixtureOrganizationSuspendedErr,
			errResp:    fixtureOrganizationSuspendedResponseError,
			respStatus: http.StatusForbidden,
			respBody:   fixture("organization_suspended_error"),
		},
		{
			name: "permission denied",
			args: args{
				ctx:      context.Background(),
				orgID:    "org_O648YDMEYeLmqdmn",
				zoneArgs: &DNSZoneArguments{Name: "test-1.com"},
			},
			errStr:     fixturePermissionDeniedErr,
			errResp:    fixturePermissionDeniedResponseError,
			respStatus: http.StatusForbidden,
			respBody:   fixture("permission_denied_error"),
		},
		{
			name: "validation error",
			args: args{
				ctx:      context.Background(),
				orgID:    "org_O648YDMEYeLmqdmn",
				zoneArgs: &DNSZoneArguments{Name: ""},
			},
			errStr:     fixtureValidationErrorErr,
			errResp:    fixtureValidationErrorResponseError,
			respStatus: http.StatusUnprocessableEntity,
			respBody:   fixture("validation_error"),
		},
		{
			name: "nil zone arguments",
			args: args{
				ctx:      context.Background(),
				orgID:    "org_O648YDMEYeLmqdmn",
				zoneArgs: nil,
			},
			errStr: "nil zone arguments",
		},
		{
			name: "nil context",
			args: args{
				ctx:      nil,
				orgID:    "org_O648YDMEYeLmqdmn",
				zoneArgs: &DNSZoneArguments{Name: "hi"},
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
					"/core/v1/organizations/%s/dns/zones", tt.args.orgID,
				),
				func(w http.ResponseWriter, r *http.Request) {
					assert.Equal(t, "POST", r.Method)
					assert.Equal(t, "", r.Header.Get("X-Field-Spec"))

					body := &reqBody{}
					err := strictUmarshal(r.Body, body)
					assert.NoError(t, err)
					assert.Equal(t, &reqBody{
						Details: &reqBodyDetails{
							Name: tt.args.zoneArgs.Name,
							TTL:  tt.args.zoneArgs.TTL,
						},
						SkipVerification: tt.args.zoneArgs.SkipVerfication,
					}, body)

					w.WriteHeader(tt.respStatus)
					_, _ = w.Write(tt.respBody)
				},
			)

			got, resp, err := c.DNSZones.Create(
				tt.args.ctx, tt.args.orgID, tt.args.zoneArgs,
			)

			if tt.respStatus != 0 {
				assert.Equal(t, tt.respStatus, resp.StatusCode)
			}

			if tt.errStr == "" {
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, err, tt.errStr)
			}

			if tt.expected != nil {
				assert.Equal(t, tt.expected, got)
			}

			if tt.errResp != nil {
				assert.Equal(t, tt.errResp, resp.Error)
			}
		})
	}
}

func TestDNSZonesResource_Delete(t *testing.T) {
	type args struct {
		ctx context.Context
		id  string
	}
	tests := []struct {
		name       string
		args       args
		expected   *DNSZone
		errStr     string
		errResp    *ResponseError
		respStatus int
		respBody   []byte
	}{
		{
			name: "DNS zone",
			args: args{
				ctx: context.Background(),
				id:  "dnszone_k75eFc4UBOgeE5Zy",
			},
			expected:   fixtureDNSZone,
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
			c, mux, _, teardown := prepareTestClient()
			defer teardown()

			mux.HandleFunc(fmt.Sprintf("/core/v1/dns/zones/%s", tt.args.id),
				func(w http.ResponseWriter, r *http.Request) {
					assert.Equal(t, "DELETE", r.Method)
					assert.Equal(t, "", r.Header.Get("X-Field-Spec"))

					w.WriteHeader(tt.respStatus)
					_, _ = w.Write(tt.respBody)
				},
			)

			got, resp, err := c.DNSZones.Delete(tt.args.ctx, tt.args.id)

			if tt.respStatus != 0 {
				assert.Equal(t, tt.respStatus, resp.StatusCode)
			}

			if tt.errStr == "" {
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, err, tt.errStr)
			}

			if tt.expected != nil {
				assert.Equal(t, tt.expected, got)
			}

			if tt.errResp != nil {
				assert.Equal(t, tt.errResp, resp.Error)
			}
		})
	}
}

func TestDNSZonesResource_VerificationDetails(t *testing.T) {
	type args struct {
		ctx context.Context
		id  string
	}
	tests := []struct {
		name       string
		args       args
		expected   *DNSZoneVerificationDetails
		errStr     string
		errResp    *ResponseError
		respStatus int
		respBody   []byte
	}{
		{
			name: "get details",
			args: args{
				ctx: context.Background(),
				id:  "dnszone_k75eFc4UBOgeE5Zy",
			},
			expected: &DNSZoneVerificationDetails{
				Nameservers: []string{"ns1.katapult.io", "ns2.katapult.io"},
				TXTRecord:   "M0Y0SzE1TzNJTUZPSDRoQUV0TDZ4MEZwckFqbW1FNHI=",
			},
			respStatus: http.StatusOK,
			respBody:   fixture("dns_zone_verification_details"),
		},
		{
			name: "non-existent DNS zone",
			args: args{
				ctx: context.Background(),
				id:  "dnszone_k75eFc4UBOgeE5Zy",
			},
			errStr:     fixtureDNSZoneNotFoundErr,
			errResp:    fixtureDNSZoneNotFoundResponseError,
			respStatus: http.StatusNotFound,
			respBody:   fixture("dns_zone_not_found_error"),
		},
		{
			name: "already verified",
			args: args{
				ctx: context.Background(),
				id:  "dnszone_k75eFc4UBOgeE5Zy",
			},
			errStr: "dns_zone_already_verified: This DNS zone is already " +
				"verified, and does not require any verification details",
			errResp: &ResponseError{
				Code: "dns_zone_already_verified",
				Description: "This DNS zone is already verified, and does " +
					"not require any verification details",
				Detail: json.RawMessage(`{}`),
			},
			respStatus: http.StatusUnprocessableEntity,
			respBody:   fixture("dns_zone_already_verified_error"),
		},
		{
			name: "infrastructure DNS zone cannot be edited",
			args: args{
				ctx: context.Background(),
				id:  "dnszone_k75eFc4UBOgeE5Zy",
			},
			errStr:     fixtureDNSZoneInfraErr,
			errResp:    fixtureDNSZoneInfraResponseError,
			respStatus: http.StatusForbidden,
			respBody: fixture(
				"dns_zone_infrastructure_zone_cannot_be_edited",
			),
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
			c, mux, _, teardown := prepareTestClient()
			defer teardown()

			mux.HandleFunc(
				fmt.Sprintf(
					"/core/v1/dns/zones/%s/verification_details",
					tt.args.id,
				),
				func(w http.ResponseWriter, r *http.Request) {
					assert.Equal(t, "GET", r.Method)
					assert.Equal(t, "", r.Header.Get("X-Field-Spec"))

					w.WriteHeader(tt.respStatus)
					_, _ = w.Write(tt.respBody)
				},
			)

			got, resp, err := c.DNSZones.VerificationDetails(
				tt.args.ctx, tt.args.id,
			)

			if tt.respStatus != 0 {
				assert.Equal(t, tt.respStatus, resp.StatusCode)
			}

			if tt.errStr == "" {
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, err, tt.errStr)
			}

			if tt.expected != nil {
				assert.Equal(t, tt.expected, got)
			}

			if tt.errResp != nil {
				assert.Equal(t, tt.errResp, resp.Error)
			}
		})
	}
}

func TestDNSZonesResource_Verify(t *testing.T) {
	type args struct {
		ctx context.Context
		id  string
	}
	tests := []struct {
		name       string
		args       args
		expected   *DNSZone
		errStr     string
		errResp    *ResponseError
		respStatus int
		respBody   []byte
	}{
		{
			name: "DNS zone",
			args: args{
				ctx: context.Background(),
				id:  "dnszone_k75eFc4UBOgeE5Zy",
			},
			expected:   fixtureDNSZone,
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
			respStatus: http.StatusNotFound,
			respBody:   fixture("dns_zone_not_found_error"),
		},
		{
			name: "infrastructure DNS zone cannot be edited",
			args: args{
				ctx: context.Background(),
				id:  "dnszone_k75eFc4UBOgeE5Zy",
			},
			errStr:     fixtureDNSZoneInfraErr,
			errResp:    fixtureDNSZoneInfraResponseError,
			respStatus: http.StatusUnprocessableEntity,
			respBody: fixture(
				"dns_zone_infrastructure_zone_cannot_be_edited",
			),
		},
		{
			name: "validation error",
			args: args{
				ctx: context.Background(),
				id:  "dnszone_k75eFc4UBOgeE5Zy",
			},
			errStr:     fixtureValidationErrorErr,
			errResp:    fixtureValidationErrorResponseError,
			respStatus: http.StatusUnprocessableEntity,
			respBody:   fixture("validation_error"),
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
			c, mux, _, teardown := prepareTestClient()
			defer teardown()

			mux.HandleFunc(
				fmt.Sprintf("/core/v1/dns/zones/%s/verify", tt.args.id),
				func(w http.ResponseWriter, r *http.Request) {
					assert.Equal(t, "POST", r.Method)
					assert.Equal(t, "", r.Header.Get("X-Field-Spec"))

					w.WriteHeader(tt.respStatus)
					_, _ = w.Write(tt.respBody)
				},
			)

			got, resp, err := c.DNSZones.Verify(tt.args.ctx, tt.args.id)

			if tt.respStatus != 0 {
				assert.Equal(t, tt.respStatus, resp.StatusCode)
			}

			if tt.errStr == "" {
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, err, tt.errStr)
			}

			if tt.expected != nil {
				assert.Equal(t, tt.expected, got)
			}

			if tt.errResp != nil {
				assert.Equal(t, tt.errResp, resp.Error)
			}
		})
	}
}

func TestDNSZonesResource_UpdateTTL(t *testing.T) {
	type reqBody struct {
		TTL int `json:"ttl"`
	}
	type args struct {
		ctx context.Context
		id  string
		ttl int
	}
	tests := []struct {
		name       string
		args       args
		expected   *DNSZone
		errStr     string
		errResp    *ResponseError
		respStatus int
		respBody   []byte
	}{
		{
			name: "update TTL",
			args: args{
				ctx: context.Background(),
				id:  "dnszone_lwz66kyviwCQyqQc",
				ttl: 1842,
			},
			expected: &DNSZone{
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
				ctx: context.Background(),
				id:  "dnszone_lwz66kyviwCQyqQc",
				ttl: 25200,
			},
			respStatus: http.StatusCreated,
			respBody:   fixture("dns_zone_update_ttl"),
		},
		{
			name: "permission denied",
			args: args{
				ctx: context.Background(),
				id:  "dnszone_lwz66kyviwCQyqQc",
				ttl: 600,
			},
			errStr:     fixturePermissionDeniedErr,
			errResp:    fixturePermissionDeniedResponseError,
			respStatus: http.StatusForbidden,
			respBody:   fixture("permission_denied_error"),
		},
		{
			name: "validation error",
			args: args{
				ctx: context.Background(),
				id:  "dnszone_lwz66kyviwCQyqQc",
				ttl: 600,
			},
			errStr:     fixtureValidationErrorErr,
			errResp:    fixtureValidationErrorResponseError,
			respStatus: http.StatusUnprocessableEntity,
			respBody:   fixture("validation_error"),
		},
		{
			name: "nil context",
			args: args{
				ctx: nil,
				id:  "dnszone_lwz66kyviwCQyqQc",
				ttl: 600,
			},
			errStr: "net/http: nil Context",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, mux, _, teardown := prepareTestClient()
			defer teardown()

			mux.HandleFunc(
				fmt.Sprintf("/core/v1/dns/zones/%s/update_ttl", tt.args.id),
				func(w http.ResponseWriter, r *http.Request) {
					assert.Equal(t, "POST", r.Method)
					assert.Equal(t, "", r.Header.Get("X-Field-Spec"))

					body := &reqBody{}
					err := strictUmarshal(r.Body, body)
					assert.NoError(t, err)
					assert.Equal(t, &reqBody{TTL: tt.args.ttl}, body)

					w.WriteHeader(tt.respStatus)
					_, _ = w.Write(tt.respBody)
				},
			)

			got, resp, err := c.DNSZones.UpdateTTL(
				tt.args.ctx, tt.args.id, tt.args.ttl,
			)

			if tt.respStatus != 0 {
				assert.Equal(t, tt.respStatus, resp.StatusCode)
			}

			if tt.errStr == "" {
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, err, tt.errStr)
			}

			if tt.expected != nil {
				assert.Equal(t, tt.expected, got)
			}

			if tt.errResp != nil {
				assert.Equal(t, tt.errResp, resp.Error)
			}
		})
	}
}
