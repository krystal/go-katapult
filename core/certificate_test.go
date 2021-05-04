package core

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"testing"

	"github.com/krystal/go-katapult"
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
				Pagination:   &katapult.Pagination{CurrentPage: 42},
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
		ctx  context.Context
		org  *Organization
		opts *ListOptions
	}
	tests := []struct {
		name           string
		args           args
		want           []*Certificate
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
			want: certificateList,
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
			respBody:   fixture("certificates_list"),
		},
		{
			name: "by organization SubDomain",
			args: args{
				ctx: context.Background(),
				org: &Organization{SubDomain: "acme"},
			},
			want: certificateList,
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
			respBody:   fixture("certificates_list"),
		},
		{
			name: "page 1",
			args: args{
				ctx:  context.Background(),
				org:  &Organization{ID: "org_O648YDMEYeLmqdmn"},
				opts: &ListOptions{Page: 1, PerPage: 2},
			},
			want: certificateList[0:2],
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
			respBody:   fixture("certificates_list_page_1"),
		},
		{
			name: "page 2",
			args: args{
				ctx:  context.Background(),
				org:  &Organization{ID: "org_O648YDMEYeLmqdmn"},
				opts: &ListOptions{Page: 2, PerPage: 2},
			},
			want: certificateList[2:],
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
			respBody:   fixture("certificates_list_page_2"),
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
			rm, mux, _, teardown := katapult.MockClient(t)
			defer teardown()
			c := NewCertificatesClient(rm)

			mux.HandleFunc(
				"/core/v1/organizations/_/certificates",
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

func TestCertificatesClient_Get(t *testing.T) {
	type args struct {
		ctx context.Context
		id  string
	}
	tests := []struct {
		name       string
		args       args
		id         string
		want       *Certificate
		errStr     string
		errResp    *katapult.ResponseError
		respStatus int
		respBody   []byte
	}{
		{
			name: "certificate",
			args: args{
				ctx: context.Background(),
				id:  "cert_Xr8jREhulOP3UJoM",
			},
			want: &Certificate{
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
			errResp: &katapult.ResponseError{
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
			rm, mux, _, teardown := katapult.MockClient(t)
			defer teardown()

			c := NewCertificatesClient(rm)

			mux.HandleFunc(fmt.Sprintf("/core/v1/certificates/%s", tt.args.id),
				func(w http.ResponseWriter, r *http.Request) {
					assert.Equal(t, "GET", r.Method)
					assertEmptyFieldSpec(t, r)
					assertAuthorization(t, r)

					w.WriteHeader(tt.respStatus)
					_, _ = w.Write(tt.respBody)
				},
			)

			got, resp, err := c.Get(tt.args.ctx, tt.args.id)

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

func TestCertificatesClient_GetByID(t *testing.T) {
	type args struct {
		ctx context.Context
		id  string
	}
	tests := []struct {
		name       string
		args       args
		id         string
		want       *Certificate
		errStr     string
		errResp    *katapult.ResponseError
		respStatus int
		respBody   []byte
	}{
		{
			name: "certificate",
			args: args{
				ctx: context.Background(),
				id:  "cert_Xr8jREhulOP3UJoM",
			},
			want: &Certificate{
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
			errResp: &katapult.ResponseError{
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
			rm, mux, _, teardown := katapult.MockClient(t)
			defer teardown()
			c := NewCertificatesClient(rm)

			mux.HandleFunc(fmt.Sprintf("/core/v1/certificates/%s", tt.args.id),
				func(w http.ResponseWriter, r *http.Request) {
					assert.Equal(t, "GET", r.Method)
					assertEmptyFieldSpec(t, r)
					assertAuthorization(t, r)

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
