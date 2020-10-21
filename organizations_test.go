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

func TestOrganizationsService_List(t *testing.T) {
	tests := []struct {
		name       string
		orgs       []*Organization
		err        string
		errResp    *ErrorResponse
		respStatus int
		respBody   []byte
	}{
		{
			name: "fetch list of data centers",
			orgs: []*Organization{
				{
					ID:        "org_O648YDMEYeLmqdmn",
					Name:      "ACME Inc.",
					SubDomain: "acme",
					Personal:  false,
					CreatedAt: timestampPtr(1589052170),
					Suspended: false,
				},
				{
					ID:        "org_c0CU62PqQgkON2rZ",
					Name:      "Lex Corp.",
					SubDomain: "lex-corp",
					Personal:  true,
					CreatedAt: timestampPtr(1542225631),
					Suspended: false,
				},
			},
			respStatus: http.StatusOK,
			respBody:   fixture("organizations_list"),
		},
		{
			name: "invalid API token response",
			err: "invalid_api_token: The API token provided was not valid " +
				"(it may not exist or have expired)",
			errResp: &ErrorResponse{
				Code: "invalid_api_token",
				Description: "The API token provided was not valid " +
					"(it may not exist or have expired)",
				Detail: json.RawMessage(`{}`),
			},
			respStatus: http.StatusForbidden,
			respBody:   fixture("invalid_api_token_error"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, mux, _, teardown := setup()
			defer teardown()

			mux.HandleFunc("/core/v1/organizations",
				func(w http.ResponseWriter, r *http.Request) {
					assert.Equal(t, "GET", r.Method)
					assert.Equal(t, "", r.Header.Get("X-Field-Spec"))

					w.WriteHeader(tt.respStatus)
					_, _ = w.Write(tt.respBody)
				},
			)

			orgs, resp, err := c.Organizations.List(context.Background())

			if tt.err == "" {
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, err, tt.err)
			}

			if tt.orgs != nil {
				assert.Equal(t, tt.orgs, orgs)
			}

			if tt.errResp != nil {
				assert.Equal(t, tt.errResp, resp.Error)
			}
		})
	}
}

func TestOrganizationsService_Get(t *testing.T) {
	tests := []struct {
		name       string
		id         string
		expected   *Organization
		err        string
		errResp    *ErrorResponse
		respStatus int
		respBody   []byte
	}{
		{
			name: "specific Organization",
			id:   "org_O648YDMEYeLmqdmn",
			expected: &Organization{
				ID:                   "org_O648YDMEYeLmqdmn",
				Name:                 "ACME Inc.",
				SubDomain:            "acme",
				InfrastructureDomain: "acme.test.kpult.com",
				Personal:             false,
				CreatedAt:            timestampPtr(1589052170),
				Suspended:            false,
				Managed:              false,
				BillingName:          "ACME Inc",
				Address1:             "273  Elk Avenue",
				Address2:             "Clarklake",
				Address3:             "",
				Address4:             "",
				Postcode:             "49234",
				VatNumber:            "GB123456789",
				Currency: &Currency{
					ID:      "cur_8UFhhlYAcRLf3ua6",
					Name:    "United States Dollars",
					IsoCode: "USD",
					Symbol:  "$",
				},
				Country: &Country{
					ID:       "ctry_V5UmyvGWYlC1pPPg",
					Name:     "United States of America",
					ISOCode2: "US",
					ISOCode3: "USA",
					TimeZone: "America/NewYork",
					EU:       false,
				},
				CountryState: &CountryState{
					ID:   "ctct_E62qc88s24FD3XIR",
					Name: "Michigan",
					Code: "MI",
					Country: &Country{
						ID:       "ctry_V5UmyvGWYlC1pPPg",
						Name:     "United States of America",
						ISOCode2: "US",
						ISOCode3: "USA",
						TimeZone: "America/NewYork",
						EU:       false,
					},
				},
			},
			respStatus: http.StatusOK,
			respBody:   fixture("organization_get"),
		},
		{
			name: "non-existent Organization",
			id:   "org_nopethisbegone",
			err: "organization_not_found: No organization was found matching " +
				"any of the criteria provided in the arguments",
			errResp: &ErrorResponse{
				Code: "organization_not_found",
				Description: "No organization was found matching any of the " +
					"criteria provided in the arguments",
				Detail: json.RawMessage(`{}`),
			},
			respStatus: http.StatusNotFound,
			respBody:   fixture("organization_not_found_error"),
		},
		{
			name: "suspended Organization",
			id:   "acme",
			err: "organization_suspended: An organization was found from the " +
				"arguments provided but it was suspended",
			errResp: &ErrorResponse{
				Code: "organization_suspended",
				Description: "An organization was found from the arguments " +
					"provided but it was suspended",
				Detail: json.RawMessage(`{}`),
			},
			respStatus: http.StatusForbidden,
			respBody:   fixture("organization_suspended_error"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, mux, _, teardown := setup()
			defer teardown()

			mux.HandleFunc(fmt.Sprintf("/core/v1/organizations/%s", tt.id),
				func(w http.ResponseWriter, r *http.Request) {
					assert.Equal(t, "GET", r.Method)
					assert.Equal(t, "", r.Header.Get("X-Field-Spec"))

					w.WriteHeader(tt.respStatus)
					_, _ = w.Write(tt.respBody)
				},
			)

			org, resp, err := c.Organizations.Get(context.Background(), tt.id)

			if tt.err == "" {
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, err, tt.err)
			}

			if tt.expected != nil {
				assert.Equal(t, tt.expected, org)
			}

			if tt.errResp != nil {
				assert.Equal(t, tt.errResp, resp.Error)
			}
		})
	}
}

func TestOrganizationsService_CreateManaged(t *testing.T) {
	type reqBody struct {
		Name      string `json:"name"`
		SubDomain string `json:"sub_domain"`
	}
	tests := []struct {
		name         string
		parent       *Organization
		orgName      string
		orgSubDomain string
		expected     *Organization
		err          string
		errResp      *ErrorResponse
		respStatus   int
		respBody     []byte
	}{
		{
			name:         "create a managed organization",
			parent:       &Organization{ID: "org_O648YDMEYeLmqdmn"},
			orgName:      "NERV Corp.",
			orgSubDomain: "nerv",
			expected: &Organization{
				ID:                   "org_TZQHTxMg1G8COlfu",
				Name:                 "NERV Corp.",
				SubDomain:            "nerv",
				InfrastructureDomain: "nerv.test.kpult.com",
				Personal:             false,
				CreatedAt:            timestampPtr(1603211871),
				Suspended:            false,
				Managed:              true,
				BillingName:          "ACME Inc",
				Address1:             "273  Elk Avenue",
				Address2:             "Clarklake",
				Address3:             "",
				Address4:             "",
				Postcode:             "49234",
				VatNumber:            "GB123456789",
				Currency: &Currency{
					ID:      "cur_8UFhhlYAcRLf3ua6",
					Name:    "United States Dollars",
					IsoCode: "USD",
					Symbol:  "$",
				},
				Country: &Country{
					ID:       "ctry_V5UmyvGWYlC1pPPg",
					Name:     "United States of America",
					ISOCode2: "US",
					ISOCode3: "USA",
					TimeZone: "America/NewYork",
					EU:       false,
				},
				CountryState: &CountryState{
					ID:   "ctct_E62qc88s24FD3XIR",
					Name: "Michigan",
					Code: "MI",
					Country: &Country{
						ID:       "ctry_V5UmyvGWYlC1pPPg",
						Name:     "United States of America",
						ISOCode2: "US",
						ISOCode3: "USA",
						TimeZone: "America/NewYork",
						EU:       false,
					},
				},
			},
			respStatus: http.StatusCreated,
			respBody:   fixture("organization_managed"),
		},
		{
			name:         "managed org limit reached",
			parent:       &Organization{ID: "org_O648YDMEYeLmqdmn"},
			orgName:      "NERV Corp.",
			orgSubDomain: "nerv",
			err: "organization_limit_reached: The maxmium number of " +
				"organizations that can be created has been reached",
			errResp: &ErrorResponse{
				Code: "organization_limit_reached",
				Description: "The maxmium number of organizations that can " +
					"be created has been reached",
				Detail: json.RawMessage(`{}`),
			},
			respStatus: http.StatusUnprocessableEntity,
			respBody:   fixture("organization_limit_reached_error"),
		},
		{
			name:         "non-existent Organization",
			parent:       &Organization{ID: "org_nopewhatbye"},
			orgName:      "NERV Corp.",
			orgSubDomain: "nerv",
			err: "organization_not_found: No organization was found matching " +
				"any of the criteria provided in the arguments",
			errResp: &ErrorResponse{
				Code: "organization_not_found",
				Description: "No organization was found matching any of the " +
					"criteria provided in the arguments",
				Detail: json.RawMessage(`{}`),
			},
			respStatus: http.StatusNotFound,
			respBody:   fixture("organization_not_found_error"),
		},
		{
			name:         "suspended Organization",
			parent:       &Organization{ID: "org_O648YDMEYeLmqdmn"},
			orgName:      "NERV Corp.",
			orgSubDomain: "nerv",
			err: "organization_suspended: An organization was found from the " +
				"arguments provided but it was suspended",
			errResp: &ErrorResponse{
				Code: "organization_suspended",
				Description: "An organization was found from the arguments " +
					"provided but it was suspended",
				Detail: json.RawMessage(`{}`),
			},
			respStatus: http.StatusForbidden,
			respBody:   fixture("organization_suspended_error"),
		},
		{
			name:         "validation error for new org details",
			parent:       &Organization{ID: "org_O648YDMEYeLmqdmn"},
			orgName:      "NERV Corp.",
			orgSubDomain: "acme",
			err: "validation_error: A validation error occurred with the " +
				"object that was being created/updated/deleted",
			errResp: &ErrorResponse{
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
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, mux, _, teardown := setup()
			defer teardown()

			mux.HandleFunc(
				fmt.Sprintf("/core/v1/organizations/%s/managed", tt.parent.ID),
				func(w http.ResponseWriter, r *http.Request) {
					assert.Equal(t, "POST", r.Method)
					assert.Equal(t, "", r.Header.Get("X-Field-Spec"))

					body := &reqBody{}
					err := strictUmarshal(r.Body, body)
					require.NoError(t, err)
					assert.Equal(t, &reqBody{
						Name:      tt.orgName,
						SubDomain: tt.orgSubDomain,
					}, body)

					w.WriteHeader(tt.respStatus)
					_, _ = w.Write(tt.respBody)
				},
			)

			org, resp, err := c.Organizations.CreateManaged(
				context.Background(),
				tt.parent, tt.orgName, tt.orgSubDomain,
			)

			if tt.err == "" {
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, err, tt.err)
			}

			if tt.expected != nil {
				assert.Equal(t, tt.expected, org)
			}

			if tt.errResp != nil {
				assert.Equal(t, tt.errResp, resp.Error)
			}
		})
	}
}
