package katapult

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNetworksService_List(t *testing.T) {
	tests := []struct {
		name       string
		org        string
		nets       []*Network
		vnets      []*VirtualNetwork
		err        string
		errResp    *ErrorResponse
		respStatus int
		respBody   []byte
	}{
		{
			name: "fetch list of networks",
			org:  "org_O648YDMEYeLmqdmn",
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
			name: "invalid API token response",
			org:  "org_O648YDMEYeLmqdmn",
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

			mux.HandleFunc(
				fmt.Sprintf(
					"/core/v1/organizations/%s/available_networks",
					tt.org,
				),
				func(w http.ResponseWriter, r *http.Request) {
					assert.Equal(t, "GET", r.Method)
					assert.Equal(t, "", r.Header.Get("X-Field-Spec"))
					w.WriteHeader(tt.respStatus)
					_, _ = w.Write(tt.respBody)
				},
			)

			nets, vnets, resp, err := c.Networks.List(
				context.Background(), tt.org,
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
