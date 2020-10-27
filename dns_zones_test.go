package katapult

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/google/go-querystring/query"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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
	fixtureDNSZoneGet = &DNSZone{
		ID:                 "dnszone_k75eFc4UBOgeE5Zy",
		Name:               "test1.example.com",
		TTL:                3600,
		Verified:           true,
		InfrastructureZone: true,
	}
)

func TestDNSZonesService_List(t *testing.T) {
	// Correlates to fixtures/dns_zones_list*.json
	dnsZonesList := []*DNSZone{
		{
			ID:                 "dnszone_k75eFc4UBOgeE5Zy",
			Name:               "test1.example.com",
			TTL:                3600,
			Verified:           true,
			InfrastructureZone: true,
		},
		{
			ID:                 "dnszone_lwz66kyviwCQyqQc",
			Name:               "test-2.example.com",
			TTL:                3600,
			Verified:           true,
			InfrastructureZone: false,
		},
		{
			ID:                 "dnszone_qr9KPhSwkGNh7IMb",
			Name:               "test-3.example.com",
			TTL:                3600,
			Verified:           true,
			InfrastructureZone: false,
		},
	}

	tests := []struct {
		name       string
		orgID      string
		opts       *ListOptions
		expected   []*DNSZone
		pagination *Pagination
		err        string
		errResp    *ResponseError
		respStatus int
		respBody   []byte
	}{
		{
			name:     "fetch list of dns_zones",
			orgID:    "org_O648YDMEYeLmqdmn",
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
			name:     "fetch page 1 of dns_zones list",
			orgID:    "org_O648YDMEYeLmqdmn",
			opts:     &ListOptions{Page: 1, PerPage: 2},
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
			name:     "fetch page 2 of dns_zones list",
			orgID:    "org_O648YDMEYeLmqdmn",
			opts:     &ListOptions{Page: 2, PerPage: 2},
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
			name:       "invalid API token response",
			orgID:      "org_O648YDMEYeLmqdmn",
			err:        fixtureInvalidAPITokenErr,
			errResp:    fixtureInvalidAPITokenResponseError,
			respStatus: http.StatusForbidden,
			respBody:   fixture("invalid_api_token_error"),
		},
		{
			name:       "non-existent Organization",
			orgID:      "org_O648YDMEYeLmqdmn",
			err:        fixtureOrganizationNotFoundErr,
			errResp:    fixtureOrganizationNotFoundResponseError,
			respStatus: http.StatusNotFound,
			respBody:   fixture("organization_not_found_error"),
		},
		{
			name:       "suspended Organization",
			orgID:      "org_O648YDMEYeLmqdmn",
			err:        fixtureOrganizationSuspendedErr,
			errResp:    fixtureOrganizationSuspendedResponseError,
			respStatus: http.StatusForbidden,
			respBody:   fixture("organization_suspended_error"),
		},
		{
			name:       "permission denied",
			orgID:      "org_O648YDMEYeLmqdmn",
			err:        fixturePermissionDeniedErr,
			errResp:    fixturePermissionDeniedResponseError,
			respStatus: http.StatusForbidden,
			respBody:   fixture("permission_denied_error"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, mux, _, teardown := setup()
			defer teardown()

			mux.HandleFunc(
				fmt.Sprintf("/core/v1/organizations/%s/dns/zones", tt.orgID),
				func(w http.ResponseWriter, r *http.Request) {
					assert.Equal(t, "GET", r.Method)
					assert.Equal(t, "", r.Header.Get("X-Field-Spec"))
					if tt.opts != nil {
						qs, err := query.Values(tt.opts)
						require.NoError(t, err)
						assert.Equal(t, qs, r.URL.Query())
					}

					w.WriteHeader(tt.respStatus)
					_, _ = w.Write(tt.respBody)
				},
			)

			got, resp, err := c.DNSZones.List(
				context.Background(), tt.orgID, tt.opts,
			)

			assert.Equal(t, tt.respStatus, resp.StatusCode)

			if tt.err == "" {
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, err, tt.err)
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

func TestDNSZonesService_Get(t *testing.T) {
	tests := []struct {
		name       string
		id         string
		expected   *DNSZone
		err        string
		errResp    *ResponseError
		respStatus int
		respBody   []byte
	}{
		{
			name:       "specific DNSZone",
			id:         "dnszone_k75eFc4UBOgeE5Zy",
			expected:   fixtureDNSZoneGet,
			respStatus: http.StatusOK,
			respBody:   fixture("dns_zone_get"),
		},
		{
			name:       "non-existent DNSZone",
			id:         "dnszone_k75eFc4UBOgeE5Zy",
			err:        fixtureDNSZoneNotFoundErr,
			errResp:    fixtureDNSZoneNotFoundResponseError,
			respStatus: http.StatusNotFound,
			respBody:   fixture("dns_zone_not_found_error"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, mux, _, teardown := setup()
			defer teardown()

			mux.HandleFunc(fmt.Sprintf("/core/v1/dns/zones/%s", tt.id),
				func(w http.ResponseWriter, r *http.Request) {
					assert.Equal(t, "GET", r.Method)
					assert.Equal(t, "", r.Header.Get("X-Field-Spec"))

					w.WriteHeader(tt.respStatus)
					_, _ = w.Write(tt.respBody)
				},
			)

			got, resp, err := c.DNSZones.Get(context.Background(), tt.id)

			assert.Equal(t, tt.respStatus, resp.StatusCode)

			if tt.err == "" {
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, err, tt.err)
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

func TestDNSZonesService_Create(t *testing.T) {
	type reqBodyDetails struct {
		Name string `json:"name"`
		TTL  int    `json:"ttl,omitempty"`
	}
	type reqBody struct {
		Details          *reqBodyDetails `json:"details"`
		SkipVerification bool            `json:"skip_verification"`
	}
	tests := []struct {
		name       string
		orgID      string
		args       *DNSZoneArguments
		expected   *DNSZone
		err        string
		errResp    *ResponseError
		respStatus int
		respBody   []byte
	}{
		{
			name:  "create a DNS Zone",
			orgID: "org_O648YDMEYeLmqdmn",
			args: &DNSZoneArguments{
				Name: "test-1.com",
				TTL:  1800,
			},
			expected: &DNSZone{
				ID:                 "dnszone_yqflWVIdu5vnirLq",
				Name:               "test-1.com",
				TTL:                1800,
				Verified:           false,
				InfrastructureZone: true,
			},
			respStatus: http.StatusCreated,
			respBody:   fixture("dns_zone_create"),
		},
		{
			name:  "skip verification",
			orgID: "org_O648YDMEYeLmqdmn",
			args: &DNSZoneArguments{
				Name:            "test-1.com",
				TTL:             1800,
				SkipVerfication: true,
			},
			respStatus: http.StatusCreated,
			respBody:   fixture("dns_zone_create"),
		},
		{
			name:       "without TTL",
			orgID:      "org_O648YDMEYeLmqdmn",
			args:       &DNSZoneArguments{Name: "test-1.com"},
			respStatus: http.StatusCreated,
			respBody:   fixture("dns_zone_create"),
		},
		{
			name:       "without TTL",
			orgID:      "org_O648YDMEYeLmqdmn",
			args:       &DNSZoneArguments{Name: "test-1.com"},
			respStatus: http.StatusCreated,
			respBody:   fixture("dns_zone_create"),
		},
		{
			name:       "non-existent Organization",
			orgID:      "org_O648YDMEYeLmqdmn",
			args:       &DNSZoneArguments{Name: "test-1.com"},
			err:        fixtureOrganizationNotFoundErr,
			errResp:    fixtureOrganizationNotFoundResponseError,
			respStatus: http.StatusNotFound,
			respBody:   fixture("organization_not_found_error"),
		},
		{
			name:       "suspended Organization",
			orgID:      "org_O648YDMEYeLmqdmn",
			args:       &DNSZoneArguments{Name: "test-1.com"},
			err:        fixtureOrganizationSuspendedErr,
			errResp:    fixtureOrganizationSuspendedResponseError,
			respStatus: http.StatusForbidden,
			respBody:   fixture("organization_suspended_error"),
		},
		{
			name:       "permission denied",
			orgID:      "org_O648YDMEYeLmqdmn",
			args:       &DNSZoneArguments{Name: "test-1.com"},
			err:        fixturePermissionDeniedErr,
			errResp:    fixturePermissionDeniedResponseError,
			respStatus: http.StatusForbidden,
			respBody:   fixture("permission_denied_error"),
		},
		{
			name:       "validation error",
			orgID:      "org_O648YDMEYeLmqdmn",
			args:       &DNSZoneArguments{Name: ""},
			err:        fixtureValidationErrorErr,
			errResp:    fixtureValidationErrorResponseError,
			respStatus: http.StatusUnprocessableEntity,
			respBody:   fixture("validation_error"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, mux, _, teardown := setup()
			defer teardown()

			mux.HandleFunc(
				fmt.Sprintf("/core/v1/organizations/%s/dns/zones", tt.orgID),
				func(w http.ResponseWriter, r *http.Request) {
					assert.Equal(t, "POST", r.Method)
					assert.Equal(t, "", r.Header.Get("X-Field-Spec"))

					body := &reqBody{}
					err := strictUmarshal(r.Body, body)
					assert.NoError(t, err)
					assert.Equal(t, &reqBody{
						Details: &reqBodyDetails{
							Name: tt.args.Name,
							TTL:  tt.args.TTL,
						},
						SkipVerification: tt.args.SkipVerfication,
					}, body)

					w.WriteHeader(tt.respStatus)
					_, _ = w.Write(tt.respBody)
				},
			)

			got, resp, err := c.DNSZones.Create(
				context.Background(), tt.orgID, tt.args,
			)

			assert.Equal(t, tt.respStatus, resp.StatusCode)

			if tt.err == "" {
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, err, tt.err)
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

func TestDNSZonesService_Delete(t *testing.T) {
	tests := []struct {
		name       string
		id         string
		expected   *DNSZone
		err        string
		errResp    *ResponseError
		respStatus int
		respBody   []byte
	}{
		{
			name:       "specific DNSZone",
			id:         "dnszone_k75eFc4UBOgeE5Zy",
			expected:   fixtureDNSZoneGet,
			respStatus: http.StatusOK,
			respBody:   fixture("dns_zone_get"),
		},
		{
			name:       "non-existent DNSZone",
			id:         "dnszone_k75eFc4UBOgeE5Zy",
			err:        fixtureDNSZoneNotFoundErr,
			errResp:    fixtureDNSZoneNotFoundResponseError,
			respStatus: http.StatusNotFound,
			respBody:   fixture("dns_zone_not_found_error"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, mux, _, teardown := setup()
			defer teardown()

			mux.HandleFunc(fmt.Sprintf("/core/v1/dns/zones/%s", tt.id),
				func(w http.ResponseWriter, r *http.Request) {
					assert.Equal(t, "DELETE", r.Method)
					assert.Equal(t, "", r.Header.Get("X-Field-Spec"))

					w.WriteHeader(tt.respStatus)
					_, _ = w.Write(tt.respBody)
				},
			)

			got, resp, err := c.DNSZones.Delete(context.Background(), tt.id)

			assert.Equal(t, tt.respStatus, resp.StatusCode)

			if tt.err == "" {
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, err, tt.err)
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

func TestDNSZonesService_VerificationDetails(t *testing.T) {
	tests := []struct {
		name       string
		id         string
		expected   *DNSZoneVerificationDetails
		err        string
		errResp    *ResponseError
		respStatus int
		respBody   []byte
	}{
		{
			name: "get details",
			id:   "dnszone_k75eFc4UBOgeE5Zy",
			expected: &DNSZoneVerificationDetails{
				Nameservers: []string{"ns1.katapult.io", "ns2.katapult.io"},
				TXTRecord:   "M0Y0SzE1TzNJTUZPSDRoQUV0TDZ4MEZwckFqbW1FNHI=",
			},
			respStatus: http.StatusOK,
			respBody:   fixture("dns_zone_verification_details"),
		},
		{
			name:       "non-existent DNSZone",
			id:         "dnszone_k75eFc4UBOgeE5Zy",
			err:        fixtureDNSZoneNotFoundErr,
			errResp:    fixtureDNSZoneNotFoundResponseError,
			respStatus: http.StatusNotFound,
			respBody:   fixture("dns_zone_not_found_error"),
		},
		{
			name: "already verified",
			id:   "dnszone_k75eFc4UBOgeE5Zy",
			err: "dns_zone_already_verified: This DNS zone is already " +
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
			name:       "infrastructure DNS Zone cannot be edited",
			id:         "dnszone_k75eFc4UBOgeE5Zy",
			err:        fixtureDNSZoneInfraErr,
			errResp:    fixtureDNSZoneInfraResponseError,
			respStatus: http.StatusForbidden,
			respBody: fixture(
				"dns_zone_infrastructure_zone_cannot_be_edited",
			),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, mux, _, teardown := setup()
			defer teardown()

			mux.HandleFunc(
				fmt.Sprintf(
					"/core/v1/dns/zones/%s/verification_details",
					tt.id,
				),
				func(w http.ResponseWriter, r *http.Request) {
					assert.Equal(t, "GET", r.Method)
					assert.Equal(t, "", r.Header.Get("X-Field-Spec"))

					w.WriteHeader(tt.respStatus)
					_, _ = w.Write(tt.respBody)
				},
			)

			got, resp, err := c.DNSZones.VerificationDetails(
				context.Background(), tt.id,
			)

			assert.Equal(t, tt.respStatus, resp.StatusCode)

			if tt.err == "" {
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, err, tt.err)
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

func TestDNSZonesService_Verify(t *testing.T) {
	tests := []struct {
		name       string
		id         string
		expected   *DNSZone
		err        string
		errResp    *ResponseError
		respStatus int
		respBody   []byte
	}{
		{
			name:       "specific DNSZone",
			id:         "dnszone_k75eFc4UBOgeE5Zy",
			expected:   fixtureDNSZoneGet,
			respStatus: http.StatusOK,
			respBody:   fixture("dns_zone_get"),
		},
		{
			name:       "non-existent DNSZone",
			id:         "dnszone_k75eFc4UBOgeE5Zy",
			err:        fixtureDNSZoneNotFoundErr,
			errResp:    fixtureDNSZoneNotFoundResponseError,
			respStatus: http.StatusNotFound,
			respBody:   fixture("dns_zone_not_found_error"),
		},
		{
			name:       "infrastructure DNS Zone cannot be edited",
			id:         "dnszone_k75eFc4UBOgeE5Zy",
			err:        fixtureDNSZoneInfraErr,
			errResp:    fixtureDNSZoneInfraResponseError,
			respStatus: http.StatusUnprocessableEntity,
			respBody: fixture(
				"dns_zone_infrastructure_zone_cannot_be_edited",
			),
		},
		{
			name:       "validation error",
			id:         "dnszone_k75eFc4UBOgeE5Zy",
			err:        fixtureValidationErrorErr,
			errResp:    fixtureValidationErrorResponseError,
			respStatus: http.StatusUnprocessableEntity,
			respBody:   fixture("validation_error"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, mux, _, teardown := setup()
			defer teardown()

			mux.HandleFunc(fmt.Sprintf("/core/v1/dns/zones/%s/verify", tt.id),
				func(w http.ResponseWriter, r *http.Request) {
					assert.Equal(t, "POST", r.Method)
					assert.Equal(t, "", r.Header.Get("X-Field-Spec"))

					w.WriteHeader(tt.respStatus)
					_, _ = w.Write(tt.respBody)
				},
			)

			got, resp, err := c.DNSZones.Verify(context.Background(), tt.id)

			assert.Equal(t, tt.respStatus, resp.StatusCode)

			if tt.err == "" {
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, err, tt.err)
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

func TestDNSZonesService_UpdateTTL(t *testing.T) {
	type reqBody struct {
		TTL int `json:"ttl"`
	}
	tests := []struct {
		name       string
		id         string
		ttl        int
		expected   *DNSZone
		err        string
		errResp    *ResponseError
		respStatus int
		respBody   []byte
	}{
		{
			name: "update TTL",
			id:   "dnszone_lwz66kyviwCQyqQc",
			ttl:  1842,
			expected: &DNSZone{
				ID:                 "dnszone_lwz66kyviwCQyqQc",
				Name:               "test-2.example.com",
				TTL:                1842,
				Verified:           true,
				InfrastructureZone: false,
			},
			respStatus: http.StatusCreated,
			respBody:   fixture("dns_zone_update_ttl"),
		},
		{
			name:       "high TTL",
			id:         "dnszone_lwz66kyviwCQyqQc",
			ttl:        25200,
			respStatus: http.StatusCreated,
			respBody:   fixture("dns_zone_update_ttl"),
		},
		{
			name:       "permission denied",
			id:         "dnszone_lwz66kyviwCQyqQc",
			ttl:        600,
			err:        fixturePermissionDeniedErr,
			errResp:    fixturePermissionDeniedResponseError,
			respStatus: http.StatusForbidden,
			respBody:   fixture("permission_denied_error"),
		},
		{
			name:       "validation error",
			id:         "dnszone_lwz66kyviwCQyqQc",
			ttl:        600,
			err:        fixtureValidationErrorErr,
			errResp:    fixtureValidationErrorResponseError,
			respStatus: http.StatusUnprocessableEntity,
			respBody:   fixture("validation_error"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, mux, _, teardown := setup()
			defer teardown()

			mux.HandleFunc(
				fmt.Sprintf("/core/v1/dns/zones/%s/update_ttl", tt.id),
				func(w http.ResponseWriter, r *http.Request) {
					assert.Equal(t, "POST", r.Method)
					assert.Equal(t, "", r.Header.Get("X-Field-Spec"))

					body := &reqBody{}
					err := strictUmarshal(r.Body, body)
					assert.NoError(t, err)
					assert.Equal(t, &reqBody{TTL: tt.ttl}, body)

					w.WriteHeader(tt.respStatus)
					_, _ = w.Write(tt.respBody)
				},
			)

			got, resp, err := c.DNSZones.UpdateTTL(
				context.Background(), tt.id, tt.ttl,
			)

			assert.Equal(t, tt.respStatus, resp.StatusCode)

			if tt.err == "" {
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, err, tt.err)
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
