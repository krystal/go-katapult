package katapult

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNetworksService_List(t *testing.T) {
	tests := []struct {
		name       string
		orgID      string
		nets       []*Network
		vnets      []*VirtualNetwork
		err        string
		errResp    *ResponseError
		respStatus int
		respBody   []byte
	}{
		{
			name:  "fetch list of networks",
			orgID: "org_O648YDMEYeLmqdmn",
			nets: []*Network{
				{
					ID:   "netw_zDW7KYAeqqfRfVag",
					Name: "Public Network",
					DataCenter: &DataCenter{
						ID:        "loc_25d48761871e4bf",
						Name:      "Shirebury",
						Permalink: "shirebury",
					},
				},
				{
					ID:   "netw_t7Rbyvr6ahqpDohR",
					Name: "Private Network",
					DataCenter: &DataCenter{
						ID:        "loc_25d48761871e4bf",
						Name:      "Shirebury",
						Permalink: "shirebury",
					},
				},
			},
			vnets: []*VirtualNetwork{
				{
					ID:   "netw_1erVCx7A5Y09WknB",
					Name: "Make-Believe Network",
					DataCenter: &DataCenter{
						ID:        "loc_25d48761871e4bf",
						Name:      "Shirebury",
						Permalink: "shirebury",
					},
				},
			},
			respStatus: http.StatusOK,
			respBody:   fixture("networks_list"),
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
			orgID:      "org_nopethisbegone",
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
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, mux, _, teardown := setup()
			defer teardown()

			mux.HandleFunc(
				fmt.Sprintf(
					"/core/v1/organizations/%s/available_networks",
					tt.orgID,
				),
				func(w http.ResponseWriter, r *http.Request) {
					assert.Equal(t, "GET", r.Method)
					assert.Equal(t, "", r.Header.Get("X-Field-Spec"))

					w.WriteHeader(tt.respStatus)
					_, _ = w.Write(tt.respBody)
				},
			)

			nets, vnets, resp, err := c.Networks.List(
				context.Background(), tt.orgID,
			)

			if tt.err == "" {
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, err, tt.err)
			}

			if tt.nets != nil {
				assert.Equal(t, tt.nets, nets)
			}

			if tt.vnets != nil {
				assert.Equal(t, tt.vnets, vnets)
			}

			if tt.errResp != nil {
				assert.Equal(t, tt.errResp, resp.Error)
			}
		})
	}
}
