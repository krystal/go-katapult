package katapult

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCertificate_JSONMarshaling(t *testing.T) {
	tests := []struct {
		name string
		obj  *Certificate
	}{
		{
			name: "empty",
			obj:  &Certificate{},
		},
		{
			name: "full",
			obj: &Certificate{
				ID:                  "id",
				Name:                "name",
				AdditionalNames:     []string{"a name"},
				Issuer:              "iss",
				State:               "state",
				CreatedAt:           timestampPtr(123),
				ExpiresAt:           timestampPtr(456),
				LastIssuedAt:        timestampPtr(789),
				IssueError:          "isserr",
				AuthorizationMethod: "meth",
				CertificateAPIURL:   "certurl",
				Certificate:         "cert",
				Chain:               "chain",
				PrivateKey:          "key",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testJSONMarshaling(t, tt.obj)
		})
	}
}

func Test_certificatesResponseBody_JSONMarshaling(t *testing.T) {
	tests := []struct {
		name string
		obj  *certificatesResponseBody
	}{
		{
			name: "empty",
			obj:  &certificatesResponseBody{},
		},
		{
			name: "full",
			obj: &certificatesResponseBody{
				Pagination:   &Pagination{CurrentPage: 42},
				Certificate:  &Certificate{ID: "id1"},
				Certificates: []*Certificate{{ID: "id2"}},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testJSONMarshaling(t, tt.obj)
		})
	}
}

func TestCertificatesClient_List(t *testing.T) {
	// Correlates to fixtures/certificates_list*.json
	certificateList := []*Certificate{
		{
			ID:     "cert_Xr8jREhulOP3UJoM",
			Name:   "test1.example.com",
			Issuer: "lets_encrypt",
		},
		{
			ID:     "cert_HJxL4lqK5o7Qy3mM",
			Name:   "test-2.example.com",
			Issuer: "custom",
		},
		{
			ID:     "cert_BJz8pI5zjmABRsE0",
			Name:   "test-3.example.com",
			Issuer: "self_signed",
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
		expected   []*Certificate
		pagination *Pagination
		errStr     string
		errResp    *ResponseError
		respStatus int
		respBody   []byte
	}{
		{
			name: "certificates",
			args: args{
				ctx:   context.Background(),
				orgID: "org_O648YDMEYeLmqdmn",
			},
			expected: certificateList,
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
			name: "page 1 of certificates",
			args: args{
				ctx:   context.Background(),
				orgID: "org_O648YDMEYeLmqdmn",
				opts:  &ListOptions{Page: 1, PerPage: 2},
			},
			expected: certificateList[0:2],
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
			name: "page 2 of certificates",
			args: args{
				ctx:   context.Background(),
				orgID: "org_O648YDMEYeLmqdmn",
				opts:  &ListOptions{Page: 2, PerPage: 2},
			},
			expected: certificateList[2:],
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
				orgID: "org_O648YDMEYeLmqdmn",
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
					"/core/v1/organizations/%s/certificates", tt.args.orgID,
				),
				func(w http.ResponseWriter, r *http.Request) {
					assert.Equal(t, "GET", r.Method)
					assertEmptyFieldSpec(t, r)
					assertAuthorization(t, r)
					if tt.args.opts != nil {
						assert.Equal(t, *tt.args.opts.Values(), r.URL.Query())
					}

					w.WriteHeader(tt.respStatus)
					_, _ = w.Write(tt.respBody)
				},
			)

			got, resp, err := c.Certificates.List(
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

func TestCertificatesClient_Get(t *testing.T) {
	type args struct {
		ctx context.Context
		id  string
	}
	tests := []struct {
		name       string
		args       args
		id         string
		expected   *Certificate
		errStr     string
		errResp    *ResponseError
		respStatus int
		respBody   []byte
	}{
		{
			name: "certificate",
			args: args{
				ctx: context.Background(),
				id:  "cert_Xr8jREhulOP3UJoM",
			},
			expected: &Certificate{
				ID:              "cert_Xr8jREhulOP3UJoM",
				Name:            "test-1.example.com",
				AdditionalNames: []string{"test1.domain.com"},
			},
			respStatus: http.StatusOK,
			respBody:   fixture("certificate_get"),
		},
		{
			name: "non-existent certificate",
			args: args{
				ctx: context.Background(),
				id:  "lb_nopethisbegone",
			},
			errStr: "certificate_not_found: No certificate was found " +
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
			name: "nil context",
			args: args{
				ctx: nil,
				id:  "cert_Xr8jREhulOP3UJoM",
			},
			errStr: "net/http: nil Context",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, mux, _, teardown := prepareTestClient()
			defer teardown()

			mux.HandleFunc(fmt.Sprintf("/core/v1/certificates/%s", tt.args.id),
				func(w http.ResponseWriter, r *http.Request) {
					assert.Equal(t, "GET", r.Method)
					assertEmptyFieldSpec(t, r)
					assertAuthorization(t, r)

					w.WriteHeader(tt.respStatus)
					_, _ = w.Write(tt.respBody)
				},
			)

			got, resp, err := c.Certificates.Get(tt.args.ctx, tt.args.id)

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
