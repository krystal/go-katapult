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

func TestCertificatesService_List(t *testing.T) {
	// Correlates to fixtures/certificates_list*.json
	certificateList := []*Certificate{
		{
			ID:           "cert_Xr8jREhulOP3UJoM",
			Name:         "test1.example.com",
			Issuer:       "lets_encrypt",
			State:        "issued",
			ExpiresAt:    timestampPtr(1611139536),
			LastIssuedAt: timestampPtr(1603190706),
		},
		{
			ID:           "cert_HJxL4lqK5o7Qy3mM",
			Name:         "test-2.example.com",
			Issuer:       "custom",
			State:        "issued",
			ExpiresAt:    timestampPtr(1610016353),
			LastIssuedAt: timestampPtr(1602067569),
		},
		{
			ID:           "cert_BJz8pI5zjmABRsE0",
			Name:         "test-3.example.com",
			Issuer:       "self_signed",
			State:        "issued",
			ExpiresAt:    timestampPtr(1667472488),
			LastIssuedAt: timestampPtr(1602326866),
		},
	}

	tests := []struct {
		name       string
		org        string
		opts       *ListOptions
		certs      []*Certificate
		pagination *Pagination
		err        string
		errResp    *ResponseError
		respStatus int
		respBody   []byte
	}{
		{
			name:  "fetch list of certificates",
			org:   "org_O648YDMEYeLmqdmn",
			certs: certificateList,
			pagination: &Pagination{
				CurrentPage: 1,
				TotalPages:  1,
				Total:       3,
				PerPage:     30,
				LargeSet:    false,
			},
			respStatus: http.StatusOK,
			respBody:   fixture("certificates_list"),
		},
		{
			name:  "fetch page 1 of certificates list",
			org:   "org_O648YDMEYeLmqdmn",
			opts:  &ListOptions{Page: 1, PerPage: 2},
			certs: certificateList[0:2],
			pagination: &Pagination{
				CurrentPage: 1,
				TotalPages:  2,
				Total:       3,
				PerPage:     2,
				LargeSet:    false,
			},
			respStatus: http.StatusOK,
			respBody:   fixture("certificates_list_page_1"),
		},
		{
			name:  "fetch page 2 of certificates list",
			org:   "org_O648YDMEYeLmqdmn",
			opts:  &ListOptions{Page: 2, PerPage: 2},
			certs: certificateList[2:],
			pagination: &Pagination{
				CurrentPage: 2,
				TotalPages:  2,
				Total:       3,
				PerPage:     2,
				LargeSet:    false,
			},
			respStatus: http.StatusOK,
			respBody:   fixture("certificates_list_page_2"),
		},
		{
			name:       "invalid API token response",
			org:        "org_O648YDMEYeLmqdmn",
			err:        fixtureInvalidAPITokenErr,
			errResp:    fixtureInvalidAPITokenResponseError,
			respStatus: http.StatusForbidden,
			respBody:   fixture("invalid_api_token_error"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, mux, _, teardown := setup()
			defer teardown()

			mux.HandleFunc(
				fmt.Sprintf("/core/v1/organizations/%s/certificates", tt.org),
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

			got, resp, err := c.Certificates.List(
				context.Background(), tt.org, tt.opts,
			)

			assert.Equal(t, tt.respStatus, resp.StatusCode)

			if tt.err == "" {
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, err, tt.err)
			}

			if tt.certs != nil {
				assert.Equal(t, tt.certs, got)
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

func TestCertificatesService_Get(t *testing.T) {
	// Correlates to fixtures/certificate_get.json
	certificateGet := &Certificate{
		ID:                  "cert_Xr8jREhulOP3UJoM",
		Name:                "test-1.example.com",
		AdditionalNames:     []string{"test1.domain.com"},
		Issuer:              "self_signed",
		State:               "issued",
		CreatedAt:           timestampPtr(1611139536),
		ExpiresAt:           timestampPtr(1611139536),
		LastIssuedAt:        timestampPtr(1603190706),
		IssueError:          "",
		AuthorizationMethod: "",
		CertificateAPIURL: "https://certificates.katapult.io/" +
			"cert_Xr8jREhulOP3UJoM/" +
			"l1XAqAqcuhERLEna4UPvwLJWAj7EtLUYFu67iEgU",
		Certificate: "-----BEGIN CERTIFICATE-----\n" +
			"YllvaVFUdjJmaVFnN2Z6VndzYWk4dm53RDI4M0h" +
			"WV3ByeE1NRHN4VDdqOHVCWm56Y3E2UzZVWk1u\n" +
			"VTExZlQwakpld0g4aWtBM1VUdHExU0FxeDhMVUt" +
			"QREhncUFkQUNPOVVtVkZ5SG9Dd2JKZUNTelUy\n" +
			"TmtveGYxRk45OG1VS0I=\n" +
			"-----END CERTIFICATE-----\n",
		Chain: "-----BEGIN CERTIFICATE-----\n" +
			"bEtZYTNHTTF0OFBzSEs0bjhvWlNKejdLMjF3enB" +
			"DdjdEQVhtNDlXajExTDBDeHlPSzZMNGpSb2Fi\n" +
			"MkI5YUhNS0xaTHhDZmFJUXVHUTIxZjFsZkRvWjl" +
			"EaU16TUE0RnhJelVqR0pFMjZ0dmU5ZmdhbUQ4\n" +
			"Y2hmZ3ZXdm11YmFyUXQ=\n" +
			"-----END CERTIFICATE-----\n",
		PrivateKey: "-----BEGIN RSA PRIVATE KEY-----\n" +
			"N2c2dUQ3NVM5NUhYZzNOQUZzUUdmMkc5cnR2ejI" +
			"0U1BxYW9Wd3M4STFnNGxJRlJUSjFGMzRWV2FY\n" +
			"cDFDSmd2RlVyVFU5TDROZHhnQ1VzWFdKV1FqMXJ" +
			"EQzBuZzB3SVpSQ3gxcTRnYmlFdEl1YWJLSUxt\n" +
			"ZjNYTHRSVkxlTkZRbmY=\n" +
			"-----END RSA PRIVATE KEY-----\n",
	}

	tests := []struct {
		name       string
		id         string
		expected   *Certificate
		err        string
		errResp    *ResponseError
		respStatus int
		respBody   []byte
	}{
		{
			name:       "specific Certificate",
			id:         "cert_Xr8jREhulOP3UJoM",
			expected:   certificateGet,
			respStatus: http.StatusOK,
			respBody:   fixture("certificate_get"),
		},
		{
			name: "non-existent Certificate",
			id:   "org_nopethisbegone",
			err: "certificate_not_found: No certificate was found " +
				"matching any of the criteria provided in the arguments",
			errResp: &ResponseError{
				Code: "certificate_not_found",
				Description: "No certificate was found matching any of the " +
					"criteria provided in the arguments",
				Detail: json.RawMessage(`{}`),
			},
			respStatus: http.StatusNotFound,
			respBody:   fixture("certificate_not_found_error"),
		},
		{
			name:       "non-existent Organization",
			id:         "org_nopethisbegone",
			err:        fixtureOrganizationNotFoundErr,
			errResp:    fixtureOrganizationNotFoundResponseError,
			respStatus: http.StatusNotFound,
			respBody:   fixture("organization_not_found_error"),
		},
		{
			name:       "suspended Organization",
			id:         "acme",
			err:        fixtureOrganizationSuspendedErr,
			errResp:    fixtureOrganizationSuspendedResponseError,
			respStatus: http.StatusForbidden,
			respBody:   fixture("organization_suspended_error"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, mux, _, teardown := setup()
			defer teardown()

			mux.HandleFunc(fmt.Sprintf("/core/v1/certificates/%s", tt.id),
				func(w http.ResponseWriter, r *http.Request) {
					assert.Equal(t, "GET", r.Method)
					assert.Equal(t, "", r.Header.Get("X-Field-Spec"))

					w.WriteHeader(tt.respStatus)
					_, _ = w.Write(tt.respBody)
				},
			)

			got, resp, err := c.Certificates.Get(context.Background(), tt.id)

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
