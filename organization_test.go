package katapult

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	fixtureOrganizationSuspendedErr = "organization_suspended: " +
		"An organization was found from the arguments provided but it was " +
		"suspended"
	fixtureOrganizationSuspendedResponseError = &ResponseError{
		Code: "organization_suspended",
		Description: "An organization was found from the arguments " +
			"provided but it was suspended",
		Detail: json.RawMessage(`{}`),
	}

	fixtureOrganizationNotFoundErr = "organization_not_found: " +
		"No organization was found matching any of the criteria provided " +
		"in the arguments"
	fixtureOrganizationNotFoundResponseError = &ResponseError{
		Code: "organization_not_found",
		Description: "No organization was found matching any of the " +
			"criteria provided in the arguments",
		Detail: json.RawMessage(`{}`),
	}
)

func TestOrganization_JSONMarshaling(t *testing.T) {
	tests := []struct {
		name string
		obj  *Organization
	}{
		{
			name: "empty",
			obj:  &Organization{},
		},
		{
			name: "full",
			obj: &Organization{
				ID:                   "Id",
				Name:                 "name",
				SubDomain:            "sub_domain",
				InfrastructureDomain: "infrastructure_domain",
				Personal:             true,
				CreatedAt:            timestampPtr(934933),
				Suspended:            true,
				Managed:              true,
				BillingName:          "billing_name",
				Address1:             "address1",
				Address2:             "address2",
				Address3:             "address3",
				Address4:             "address4",
				Postcode:             "postcode",
				VatNumber:            "vat_number",
				Currency:             &Currency{ID: "id0"},
				Country:              &Country{ID: "id1"},
				CountryState:         &CountryState{ID: "id2"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testJSONMarshaling(t, tt.obj)
		})
	}
}

func Test_organizationsResponseBody_JSONMarshaling(t *testing.T) {
	tests := []struct {
		name string
		obj  *organizationsResponseBody
	}{
		{
			name: "empty",
			obj:  &organizationsResponseBody{},
		},
		{
			name: "full",
			obj: &organizationsResponseBody{
				Organization:  &Organization{ID: "id1"},
				Organizations: []*Organization{{ID: "id2"}},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testJSONMarshaling(t, tt.obj)
		})
	}
}

func TestOrganizationsResource_List(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name       string
		args       args
		orgs       []*Organization
		errStr     string
		errResp    *ResponseError
		respStatus int
		respBody   []byte
	}{
		{
			name: "fetch list of organizations",
			args: args{
				ctx: context.Background()},
			orgs: []*Organization{
				{
					ID:        "org_O648YDMEYeLmqdmn",
					Name:      "ACME Inc.",
					SubDomain: "acme",
				},
				{
					ID:        "org_c0CU62PqQgkON2rZ",
					Name:      "Lex Corp.",
					SubDomain: "lex-corp",
				},
			},
			respStatus: http.StatusOK,
			respBody:   fixture("organizations_list"),
		},
		{
			name: "invalid API token response",
			args: args{
				ctx: context.Background(),
			},
			errStr:     fixtureInvalidAPITokenErr,
			errResp:    fixtureInvalidAPITokenResponseError,
			respStatus: http.StatusForbidden,
			respBody:   fixture("invalid_api_token_error"),
		},
		{
			name: "nil context",
			args: args{
				ctx: nil,
			},
			errStr: "net/http: nil Context",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, mux, _, teardown := prepareTestClient()
			defer teardown()

			mux.HandleFunc("/core/v1/organizations",
				func(w http.ResponseWriter, r *http.Request) {
					assert.Equal(t, "GET", r.Method)
					assert.Equal(t, "", r.Header.Get("X-Field-Spec"))

					w.WriteHeader(tt.respStatus)
					_, _ = w.Write(tt.respBody)
				},
			)

			got, resp, err := c.Organizations.List(tt.args.ctx)

			if tt.respStatus != 0 {
				assert.Equal(t, tt.respStatus, resp.StatusCode)
			}

			if tt.errStr == "" {
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, err, tt.errStr)
			}

			if tt.orgs != nil {
				assert.Equal(t, tt.orgs, got)
			}

			if tt.errResp != nil {
				assert.Equal(t, tt.errResp, resp.Error)
			}
		})
	}
}

func TestOrganizationsResource_Get(t *testing.T) {
	type args struct {
		ctx context.Context
		id  string
	}
	tests := []struct {
		name       string
		args       args
		expected   *Organization
		errStr     string
		errResp    *ResponseError
		respStatus int
		respBody   []byte
	}{
		{
			name: "organization",
			args: args{
				ctx: context.Background(),
				id:  "org_O648YDMEYeLmqdmn",
			},
			expected: &Organization{
				ID:        "org_O648YDMEYeLmqdmn",
				Name:      "ACME Inc.",
				SubDomain: "acme",
			},
			respStatus: http.StatusOK,
			respBody:   fixture("organization_get"),
		},
		{
			name: "non-existent organization",
			args: args{
				ctx: context.Background(),
				id:  "org_nopethisbegone",
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
				id:  "org_O648YDMEYeLmqdmn",
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
				id:  "org_O648YDMEYeLmqdmn",
			},
			errStr: "net/http: nil Context",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, mux, _, teardown := prepareTestClient()
			defer teardown()

			mux.HandleFunc(fmt.Sprintf("/core/v1/organizations/%s", tt.args.id),
				func(w http.ResponseWriter, r *http.Request) {
					assert.Equal(t, "GET", r.Method)
					assert.Equal(t, "", r.Header.Get("X-Field-Spec"))

					w.WriteHeader(tt.respStatus)
					_, _ = w.Write(tt.respBody)
				},
			)

			got, resp, err := c.Organizations.Get(tt.args.ctx, tt.args.id)

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

func TestOrganizationsResource_CreateManaged(t *testing.T) {
	type args struct {
		ctx       context.Context
		parentID  string
		name      string
		subDomain string
	}
	type reqBody struct {
		Name      string `json:"name"`
		SubDomain string `json:"sub_domain"`
	}
	tests := []struct {
		name       string
		args       args
		expected   *Organization
		errStr     string
		errResp    *ResponseError
		respStatus int
		respBody   []byte
	}{
		{
			name: "organization",
			args: args{
				ctx:       context.Background(),
				parentID:  "org_O648YDMEYeLmqdmn",
				name:      "NERV Corp.",
				subDomain: "nerv",
			},
			expected: &Organization{
				ID:        "org_TZQHTxMg1G8COlfu",
				Name:      "NERV Corp.",
				SubDomain: "nerv",
			},
			respStatus: http.StatusCreated,
			respBody:   fixture("organization_managed"),
		},
		{
			name: "managed org limit reached",
			args: args{
				ctx:       context.Background(),
				parentID:  "org_O648YDMEYeLmqdmn",
				name:      "NERV Corp.",
				subDomain: "nerv",
			},
			errStr: "organization_limit_reached: The maxmium number of " +
				"organizations that can be created has been reached",
			errResp: &ResponseError{
				Code: "organization_limit_reached",
				Description: "The maxmium number of organizations that can " +
					"be created has been reached",
				Detail: json.RawMessage(`{}`),
			},
			respStatus: http.StatusUnprocessableEntity,
			respBody:   fixture("organization_limit_reached_error"),
		},
		{
			name: "non-existent organization",
			args: args{
				ctx:       context.Background(),
				parentID:  "org_nopewhatbye",
				name:      "NERV Corp.",
				subDomain: "nerv",
			},
			errStr:     fixtureOrganizationNotFoundErr,
			errResp:    fixtureOrganizationNotFoundResponseError,
			respStatus: http.StatusNotFound,
			respBody:   fixture("organization_not_found_error"),
		},
		{
			name: "suspended organization",
			args: args{
				ctx:       context.Background(),
				parentID:  "org_O648YDMEYeLmqdmn",
				name:      "NERV Corp.",
				subDomain: "nerv",
			},
			errStr:     fixtureOrganizationSuspendedErr,
			errResp:    fixtureOrganizationSuspendedResponseError,
			respStatus: http.StatusForbidden,
			respBody:   fixture("organization_suspended_error"),
		},
		{
			name: "validation error for new org details",
			args: args{
				ctx:       context.Background(),
				parentID:  "org_O648YDMEYeLmqdmn",
				name:      "NERV Corp.",
				subDomain: "acme",
			},
			errStr: "validation_error: A validation error occurred with the " +
				"object that was being created/updated/deleted",
			errResp: &ResponseError{
				Code: "validation_error",
				Description: "A validation error occurred with the object " +
					"that was being created/updated/deleted",
				Detail: json.RawMessage(`{
      "errors": ["Sub domain has already been taken"]
    }`,
				),
			},
			respStatus: http.StatusUnprocessableEntity,
			respBody: fixture(
				"organization_validation_error_sub_domain_taken",
			),
		},
		{
			name: "nil context",
			args: args{
				ctx:       nil,
				parentID:  "org_O648YDMEYeLmqdmn",
				name:      "NERV Corp.",
				subDomain: "acme",
			},
			errStr: "net/http: nil Context",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, mux, _, teardown := prepareTestClient()
			defer teardown()

			mux.HandleFunc(
				fmt.Sprintf("/core/v1/organizations/%s/managed",
					tt.args.parentID,
				),
				func(w http.ResponseWriter, r *http.Request) {
					assert.Equal(t, "POST", r.Method)
					assert.Equal(t, "", r.Header.Get("X-Field-Spec"))

					body := &reqBody{}
					err := strictUmarshal(r.Body, body)
					require.NoError(t, err)
					assert.Equal(t, &reqBody{
						Name:      tt.args.name,
						SubDomain: tt.args.subDomain,
					}, body)

					w.WriteHeader(tt.respStatus)
					_, _ = w.Write(tt.respBody)
				},
			)

			got, resp, err := c.Organizations.CreateManaged(
				tt.args.ctx, tt.args.parentID, tt.args.name, tt.args.subDomain,
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
