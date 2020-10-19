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
		respBody   string
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
			respBody: `{
  "organizations": [
    {
      "id": "org_O648YDMEYeLmqdmn",
      "name": "ACME Inc.",
      "sub_domain": "acme",
      "personal": false,
      "created_at": 1589052170,
      "suspended": false
    },
    {
      "id": "org_c0CU62PqQgkON2rZ",
      "name": "Lex Corp.",
      "sub_domain": "lex-corp",
      "personal": true,
      "created_at": 1542225631,
      "suspended": false
    }
  ]
}`,
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
			//nolint:lll
			respBody: `{
  "error": {
    "code": "invalid_api_token",
    "description": "The API token provided was not valid (it may not exist or have expired)",
    "detail": {}
  }
}`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, mux, _, teardown := setup()
			defer teardown()

			mux.HandleFunc("/v1/organizations",
				func(w http.ResponseWriter, r *http.Request) {
					assert.Equal(t, "GET", r.Method)
					w.WriteHeader(tt.respStatus)
					fmt.Fprint(w, tt.respBody)
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
		respBody   string
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
			respBody: `{
  "organization": {
    "id": "org_O648YDMEYeLmqdmn",
    "name": "ACME Inc.",
    "sub_domain": "acme",
    "infrastructure_domain": "acme.test.kpult.com",
    "personal": false,
    "created_at": 1589052170,
    "suspended": false,
    "managed": false,
    "billing_name": "ACME Inc",
    "address1": "273  Elk Avenue",
    "address2": "Clarklake",
    "address3": "",
    "address4": "",
    "postcode": "49234",
    "vat_number": "GB123456789",
    "currency": {
      "id": "cur_8UFhhlYAcRLf3ua6",
      "name": "United States Dollars",
      "iso_code": "USD",
      "symbol": "$"
    },
    "country": {
      "id": "ctry_V5UmyvGWYlC1pPPg",
      "name": "United States of America",
      "iso_code2": "US",
      "iso_code3": "USA",
      "time_zone": "America/NewYork",
      "eu": false
    },
    "country_state": {
      "id": "ctct_E62qc88s24FD3XIR",
      "name": "Michigan",
      "code": "MI",
      "country": {
        "id": "ctry_V5UmyvGWYlC1pPPg",
        "name": "United States of America",
        "iso_code2": "US",
        "iso_code3": "USA",
        "time_zone": "America/NewYork",
        "eu": false
      }
    }
  }
}`,
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
			//nolint:lll
			respBody: `{
  "error": {
    "code": "organization_not_found",
    "description": "No organization was found matching any of the criteria provided in the arguments",
    "detail": {}
  }
}`,
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
			//nolint:lll
			respBody: `{
  "error": {
    "code": "organization_suspended",
    "description": "An organization was found from the arguments provided but it was suspended",
    "detail": {}
  }
}`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, mux, _, teardown := setup()
			defer teardown()

			mux.HandleFunc(fmt.Sprintf("/v1/organizations/%s", tt.id),
				func(w http.ResponseWriter, r *http.Request) {
					assert.Equal(t, "GET", r.Method)
					w.WriteHeader(tt.respStatus)
					fmt.Fprint(w, tt.respBody)
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
		respBody     string
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
			respBody: `{
  "organization": {
    "id": "org_TZQHTxMg1G8COlfu",
    "name": "NERV Corp.",
    "sub_domain": "nerv",
    "infrastructure_domain": "nerv.test.kpult.com",
    "personal": false,
    "created_at": 1603211871,
    "suspended": false,
    "managed": true,
    "billing_name": "ACME Inc",
    "address1": "273  Elk Avenue",
    "address2": "Clarklake",
    "address3": "",
    "address4": "",
    "postcode": "49234",
    "vat_number": "GB123456789",
    "currency": {
      "id": "cur_8UFhhlYAcRLf3ua6",
      "name": "United States Dollars",
      "iso_code": "USD",
      "symbol": "$"
    },
    "country": {
      "id": "ctry_V5UmyvGWYlC1pPPg",
      "name": "United States of America",
      "iso_code2": "US",
      "iso_code3": "USA",
      "time_zone": "America/NewYork",
      "eu": false
    },
    "country_state": {
      "id": "ctct_E62qc88s24FD3XIR",
      "name": "Michigan",
      "code": "MI",
      "country": {
        "id": "ctry_V5UmyvGWYlC1pPPg",
        "name": "United States of America",
        "iso_code2": "US",
        "iso_code3": "USA",
        "time_zone": "America/NewYork",
        "eu": false
      }
    }
  }
}`,
		},
		{
			name:         "non-existent Organization",
			parent:       &Organization{ID: "org_nopewhatbye"},
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
			//nolint:lll
			respBody: `{
  "error": {
    "code": "organization_limit_reached",
    "description": "The maxmium number of organizations that can be created has been reached",
    "detail": {}
  }
}`,
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
			//nolint:lll
			respBody: `{
  "error": {
    "code": "organization_not_found",
    "description": "No organization was found matching any of the criteria provided in the arguments",
    "detail": {}
  }
}`,
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
			//nolint:lll
			respBody: `{
  "error": {
    "code": "organization_suspended",
    "description": "An organization was found from the arguments provided but it was suspended",
    "detail": {}
  }
}`,
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
      "errors": [
        "Sub domain has already been taken"
      ]
    }`,
				),
			},
			respStatus: http.StatusUnprocessableEntity,
			//nolint:lll
			respBody: `{
  "error": {
    "code": "validation_error",
    "description": "A validation error occurred with the object that was being created/updated/deleted",
    "detail": {
      "errors": [
        "Sub domain has already been taken"
      ]
    }
  }
}`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, mux, _, teardown := setup()
			defer teardown()

			mux.HandleFunc(
				fmt.Sprintf("/v1/organizations/%s/managed", tt.parent.ID),
				func(w http.ResponseWriter, r *http.Request) {
					assert.Equal(t, "POST", r.Method)

					body := &reqBody{}
					err := strictUmarshal(r.Body, body)
					require.NoError(t, err)
					assert.Equal(t, &reqBody{
						Name:      tt.orgName,
						SubDomain: tt.orgSubDomain,
					}, body)

					w.WriteHeader(tt.respStatus)
					fmt.Fprint(w, tt.respBody)
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
