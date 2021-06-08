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
	fixtureOrganizationSuspendedErr = "organization_suspended: " +
		"An organization was found from the arguments provided but it was " +
		"suspended"
	fixtureOrganizationSuspendedResponseError = &katapult.ResponseError{
		Code: "organization_suspended",
		Description: "An organization was found from the arguments " +
			"provided but it was suspended",
		Detail: json.RawMessage(`{}`),
	}

	fixtureOrganizationNotFoundErr = "organization_not_found: " +
		"No organization was found matching any of the criteria provided " +
		"in the arguments"
	fixtureOrganizationNotFoundResponseError = &katapult.ResponseError{
		Code: "organization_not_found",
		Description: "No organization was found matching any of the " +
			"criteria provided in the arguments",
		Detail: json.RawMessage(`{}`),
	}

	fixtureOrganizationNotActivatedErr = "organization_not_activated: " +
		"An organization was found from the arguments provided but it wasn't " +
		"activated yet"
	fixtureOrganizationNotActivatedResponseError = &katapult.ResponseError{
		Code: "organization_not_activated",
		Description: "An organization was found from the arguments provided " +
			"but it wasn't activated yet",
		Detail: json.RawMessage(`{}`),
	}
)

func TestClient_Organizations(t *testing.T) {
	c := New(&testclient.Client{})

	assert.IsType(t, &OrganizationsClient{}, c.Organizations)
}

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
				ID:                   "org_O648YDMEYeLmqdmn",
				Name:                 "ACME Inc.",
				SubDomain:            "acme",
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

func TestOrganization_Ref(t *testing.T) {
	o := Organization{ID: "org_O648YDMEYeLmqdmn"}
	assert.Equal(t, OrganizationRef{ID: "org_O648YDMEYeLmqdmn"}, o.Ref())
}

func TestOrganizationRef_queryValues(t *testing.T) {
	tests := []struct {
		name string
		opts OrganizationRef
	}{
		{
			name: "empty",
			opts: OrganizationRef{},
		},
		{
			name: "both",
			opts: OrganizationRef{
				ID:        "org_O648YDMEYeLmqdmn",
				SubDomain: "acme",
			},
		},
		{
			name: "id alone",
			opts: OrganizationRef{
				ID: "org_O648YDMEYeLmqdmn",
			},
		},
		{
			name: "subdomain alone",
			opts: OrganizationRef{
				SubDomain: "acme",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testQueryableEncoding(t, tt.opts)
		})
	}
}

func Test_organizationCreateManagedRequest_JSONMarshaling(t *testing.T) {
	tests := []struct {
		name string
		obj  *organizationCreateManagedRequest
	}{
		{
			name: "empty",
			obj:  &organizationCreateManagedRequest{},
		},
		{
			name: "full",
			obj: &organizationCreateManagedRequest{
				Organization: OrganizationRef{ID: "org_O648YDMEYeLmqdmn"},
				Name:         "ACME Rockets Inc.",
				SubDomain:    "acme-rockets",
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

func TestOrganizationsClient_List(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name       string
		args       args
		want       []*Organization
		errStr     string
		errResp    *katapult.ResponseError
		respStatus int
		respBody   []byte
	}{
		{
			name: "fetch list of organizations",
			args: args{
				ctx: context.Background(),
			},
			want: []*Organization{
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
			rm, mux, _, teardown := prepareTestClient(t)
			defer teardown()
			c := NewOrganizationsClient(rm)

			mux.HandleFunc("/core/v1/organizations",
				func(w http.ResponseWriter, r *http.Request) {
					assert.Equal(t, "GET", r.Method)
					assertEmptyFieldSpec(t, r)
					assertAuthorization(t, r)

					w.WriteHeader(tt.respStatus)
					_, _ = w.Write(tt.respBody)
				},
			)

			got, resp, err := c.List(tt.args.ctx)

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

func TestOrganizationsClient_Get(t *testing.T) {
	type args struct {
		ctx context.Context
		ref OrganizationRef
	}
	tests := []struct {
		name       string
		args       args
		reqPath    string
		reqQuery   *url.Values
		want       *Organization
		errStr     string
		errResp    *katapult.ResponseError
		respStatus int
		respBody   []byte
	}{
		{
			name: "by ID",
			args: args{
				ctx: context.Background(),
				ref: OrganizationRef{ID: "org_O648YDMEYeLmqdmn"},
			},
			reqQuery: &url.Values{
				"organization[id]": []string{"org_O648YDMEYeLmqdmn"},
			},
			want: &Organization{
				ID:        "org_O648YDMEYeLmqdmn",
				Name:      "ACME Inc.",
				SubDomain: "acme",
			},
			respStatus: http.StatusOK,
			respBody:   fixture("organization_get"),
		},
		{
			name: "by SubDomain",
			args: args{
				ctx: context.Background(),
				ref: OrganizationRef{SubDomain: "acme"},
			},
			reqQuery: &url.Values{
				"organization[sub_domain]": []string{"acme"},
			},
			want: &Organization{
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
				ref: OrganizationRef{ID: "org_O648YDMEYeLmqdmn"},
			},
			reqQuery: &url.Values{
				"organization[id]": []string{"org_O648YDMEYeLmqdmn"},
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
				ref: OrganizationRef{ID: "org_O648YDMEYeLmqdmn"},
			},
			reqQuery: &url.Values{
				"organization[id]": []string{"org_O648YDMEYeLmqdmn"},
			},
			errStr:     fixtureOrganizationSuspendedErr,
			errResp:    fixtureOrganizationSuspendedResponseError,
			respStatus: http.StatusForbidden,
			respBody:   fixture("organization_suspended_error"),
		},
		{
			name: "not activated organization",
			args: args{
				ctx: context.Background(),
				ref: OrganizationRef{ID: "org_O648YDMEYeLmqdmn"},
			},
			reqQuery: &url.Values{
				"organization[id]": []string{"org_O648YDMEYeLmqdmn"},
			},
			errStr:     fixtureOrganizationNotActivatedErr,
			errResp:    fixtureOrganizationNotActivatedResponseError,
			respStatus: http.StatusForbidden,
			respBody:   fixture("organization_not_activated_error"),
		},
		{
			name: "nil context",
			args: args{
				ctx: nil,
				ref: OrganizationRef{ID: "org_O648YDMEYeLmqdmn"},
			},
			errStr: "net/http: nil Context",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rm, mux, _, teardown := prepareTestClient(t)
			defer teardown()
			c := NewOrganizationsClient(rm)

			mux.HandleFunc(
				"/core/v1/organizations/_",
				func(w http.ResponseWriter, r *http.Request) {
					assert.Equal(t, "GET", r.Method)
					assertEmptyFieldSpec(t, r)
					assertAuthorization(t, r)

					if tt.reqQuery != nil {
						assert.Equal(t, *tt.reqQuery, r.URL.Query())
					}

					w.WriteHeader(tt.respStatus)
					_, _ = w.Write(tt.respBody)
				},
			)

			got, resp, err := c.Get(
				tt.args.ctx, tt.args.ref,
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

func TestOrganizationsClient_GetByID(t *testing.T) {
	type args struct {
		ctx context.Context
		id  string
	}
	tests := []struct {
		name       string
		args       args
		want       *Organization
		errStr     string
		errResp    *katapult.ResponseError
		respStatus int
		respBody   []byte
	}{
		{
			name: "organization",
			args: args{
				ctx: context.Background(),
				id:  "org_O648YDMEYeLmqdmn",
			},
			want: &Organization{
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
			name: "not activated organization",
			args: args{
				ctx: context.Background(),
				id:  "org_O648YDMEYeLmqdmn",
			},
			errStr:     fixtureOrganizationNotActivatedErr,
			errResp:    fixtureOrganizationNotActivatedResponseError,
			respStatus: http.StatusForbidden,
			respBody:   fixture("organization_not_activated_error"),
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
			rm, mux, _, teardown := prepareTestClient(t)
			defer teardown()
			c := NewOrganizationsClient(rm)

			mux.HandleFunc("/core/v1/organizations/_",
				func(w http.ResponseWriter, r *http.Request) {
					assert.Equal(t, "GET", r.Method)
					assertEmptyFieldSpec(t, r)
					assertAuthorization(t, r)

					assert.Equal(t, url.Values{
						"organization[id]": []string{tt.args.id},
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
		})
	}
}

func TestOrganizationsClient_GetBySubDomain(t *testing.T) {
	type args struct {
		ctx       context.Context
		subDomain string
	}
	tests := []struct {
		name       string
		args       args
		want       *Organization
		errStr     string
		errResp    *katapult.ResponseError
		respStatus int
		respBody   []byte
	}{
		{
			name: "organization",
			args: args{
				ctx:       context.Background(),
				subDomain: "acme",
			},
			want: &Organization{
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
				ctx:       context.Background(),
				subDomain: "not-here",
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
				subDomain: "acme",
			},
			errStr:     fixtureOrganizationSuspendedErr,
			errResp:    fixtureOrganizationSuspendedResponseError,
			respStatus: http.StatusForbidden,
			respBody:   fixture("organization_suspended_error"),
		},
		{
			name: "not activated organization",
			args: args{
				ctx:       context.Background(),
				subDomain: "acme",
			},
			errStr:     fixtureOrganizationNotActivatedErr,
			errResp:    fixtureOrganizationNotActivatedResponseError,
			respStatus: http.StatusForbidden,
			respBody:   fixture("organization_not_activated_error"),
		},
		{
			name: "nil context",
			args: args{
				ctx:       nil,
				subDomain: "acme",
			},
			errStr: "net/http: nil Context",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rm, mux, _, teardown := prepareTestClient(t)
			defer teardown()
			c := NewOrganizationsClient(rm)

			mux.HandleFunc("/core/v1/organizations/_",
				func(w http.ResponseWriter, r *http.Request) {
					assert.Equal(t, "GET", r.Method)
					assertEmptyFieldSpec(t, r)
					assertAuthorization(t, r)

					qs := url.Values{
						"organization[sub_domain]": []string{tt.args.subDomain},
					}
					assert.Equal(t, qs, r.URL.Query())

					w.WriteHeader(tt.respStatus)
					_, _ = w.Write(tt.respBody)
				},
			)

			got, resp, err := c.GetBySubDomain(
				tt.args.ctx, tt.args.subDomain,
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

func TestOrganizationsClient_CreateManaged(t *testing.T) {
	type args struct {
		ctx    context.Context
		parent OrganizationRef
		args   *OrganizationManagedArguments
	}
	tests := []struct {
		name       string
		args       args
		reqBody    *organizationCreateManagedRequest
		want       *Organization
		errStr     string
		errResp    *katapult.ResponseError
		respStatus int
		respBody   []byte
	}{
		{
			name: "by organization ID",
			args: args{
				ctx: context.Background(),
				parent: OrganizationRef{
					ID: "org_O648YDMEYeLmqdmn",
				},
				args: &OrganizationManagedArguments{
					Name:      "NERV Corp.",
					SubDomain: "nerv",
				},
			},
			reqBody: &organizationCreateManagedRequest{
				Organization: OrganizationRef{ID: "org_O648YDMEYeLmqdmn"},
				Name:         "NERV Corp.",
				SubDomain:    "nerv",
			},
			want: &Organization{
				ID:        "org_TZQHTxMg1G8COlfu",
				Name:      "NERV Corp.",
				SubDomain: "nerv",
			},
			respStatus: http.StatusCreated,
			respBody:   fixture("organization_managed"),
		},
		{
			name: "by organization SubDomain",
			args: args{
				ctx: context.Background(),
				parent: OrganizationRef{
					SubDomain: "acme",
				},
				args: &OrganizationManagedArguments{
					Name:      "NERV Corp.",
					SubDomain: "nerv",
				},
			},
			reqBody: &organizationCreateManagedRequest{
				Organization: OrganizationRef{SubDomain: "acme"},
				Name:         "NERV Corp.",
				SubDomain:    "nerv",
			},
			want: &Organization{
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
				ctx:    context.Background(),
				parent: OrganizationRef{ID: "org_O648YDMEYeLmqdmn"},
				args: &OrganizationManagedArguments{
					Name:      "NERV Corp.",
					SubDomain: "nerv",
				},
			},
			errStr: "organization_limit_reached: The maxmium number of " +
				"organizations that can be created has been reached",
			errResp: &katapult.ResponseError{
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
				ctx:    context.Background(),
				parent: OrganizationRef{ID: "org_nopewhatbye"},
				args: &OrganizationManagedArguments{
					Name:      "NERV Corp.",
					SubDomain: "nerv",
				},
			},
			errStr:     fixtureOrganizationNotFoundErr,
			errResp:    fixtureOrganizationNotFoundResponseError,
			respStatus: http.StatusNotFound,
			respBody:   fixture("organization_not_found_error"),
		},
		{
			name: "suspended organization",
			args: args{
				ctx:    context.Background(),
				parent: OrganizationRef{ID: "org_O648YDMEYeLmqdmn"},
				args: &OrganizationManagedArguments{
					Name:      "NERV Corp.",
					SubDomain: "nerv",
				},
			},
			errStr:     fixtureOrganizationSuspendedErr,
			errResp:    fixtureOrganizationSuspendedResponseError,
			respStatus: http.StatusForbidden,
			respBody:   fixture("organization_suspended_error"),
		},
		{
			name: "not activated organization",
			args: args{
				ctx:    context.Background(),
				parent: OrganizationRef{ID: "org_O648YDMEYeLmqdmn"},
				args: &OrganizationManagedArguments{
					Name:      "NERV Corp.",
					SubDomain: "nerv",
				},
			},
			errStr:     fixtureOrganizationNotActivatedErr,
			errResp:    fixtureOrganizationNotActivatedResponseError,
			respStatus: http.StatusForbidden,
			respBody:   fixture("organization_not_activated_error"),
		},
		{
			name: "validation error for new org details",
			args: args{
				ctx:    context.Background(),
				parent: OrganizationRef{ID: "org_O648YDMEYeLmqdmn"},
				args: &OrganizationManagedArguments{
					Name:      "NERV Corp.",
					SubDomain: "nerv",
				},
			},
			//nolint:lll
			errStr: "validation_error: A validation error occurred with the " +
				"object that was being created/updated/deleted -- " +
				"{\n  \"errors\": [\n    \"Sub domain has already been taken\"\n  ]\n}",
			errResp: &katapult.ResponseError{
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
			name: "nil args",
			args: args{
				ctx:    context.Background(),
				parent: OrganizationRef{ID: "org_O648YDMEYeLmqdmn"},
				args:   nil,
			},
			errStr:     fixtureValidationErrorErr,
			errResp:    fixtureValidationErrorResponseError,
			respStatus: http.StatusUnprocessableEntity,
			respBody:   fixture("validation_error"),
		},
		{
			name: "nil context",
			args: args{
				ctx:    nil,
				parent: OrganizationRef{ID: "org_O648YDMEYeLmqdmn"},
				args: &OrganizationManagedArguments{
					Name:      "NERV Corp.",
					SubDomain: "nerv",
				},
			},
			errStr: "net/http: nil Context",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rm, mux, _, teardown := prepareTestClient(t)
			defer teardown()
			c := NewOrganizationsClient(rm)

			mux.HandleFunc(
				"/core/v1/organizations/_/managed",
				func(w http.ResponseWriter, r *http.Request) {
					assert.Equal(t, "POST", r.Method)
					assertEmptyFieldSpec(t, r)
					assertAuthorization(t, r)

					if tt.reqBody != nil {
						reqBody := &organizationCreateManagedRequest{}
						err := strictUmarshal(r.Body, reqBody)
						assert.NoError(t, err)
						assert.Equal(t, tt.reqBody, reqBody)
					}
					w.WriteHeader(tt.respStatus)
					_, _ = w.Write(tt.respBody)
				},
			)

			got, resp, err := c.CreateManaged(
				tt.args.ctx, tt.args.parent, tt.args.args,
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
