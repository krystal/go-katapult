package katapult

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDataCentersService_List(t *testing.T) {
	tests := []struct {
		name       string
		dcs        []*DataCenter
		err        string
		errResp    *ErrorResponse
		respStatus int
		respBody   string
	}{
		{
			name: "fetch list of data centers",
			dcs: []*DataCenter{
				{
					ID:        "loc_25d48761871e4bf",
					Name:      "Shirebury",
					Permalink: "shirebury",
					Country: &Country{
						ID:   "ctry_2f2dc89a5956437",
						Name: "United Kingdom",
					},
				},
				{
					ID:        "loc_a2417980b9874c0",
					Name:      "New Town",
					Permalink: "newtown",
					Country: &Country{
						ID:   "ctry_9a989e68e0ad866",
						Name: "USA",
					},
				},
			},
			respStatus: http.StatusOK,
			respBody: `{
  "data_centers": [
    {
      "id": "loc_25d48761871e4bf",
      "name": "Shirebury",
      "permalink": "shirebury",
      "country": {
        "id": "ctry_2f2dc89a5956437",
        "name": "United Kingdom"
      }
    },
    {
      "id": "loc_a2417980b9874c0",
      "name": "New Town",
      "permalink": "newtown",
      "country": {
        "id": "ctry_9a989e68e0ad866",
        "name": "USA"
      }
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

			mux.HandleFunc("/core/v1/data_centers",
				func(w http.ResponseWriter, r *http.Request) {
					assert.Equal(t, "GET", r.Method)
					w.WriteHeader(tt.respStatus)
					fmt.Fprint(w, tt.respBody)
				},
			)

			dcs, resp, err := c.DataCenters.List(context.Background())

			if tt.err == "" {
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, err, tt.err)
			}

			if tt.dcs != nil {
				assert.Equal(t, tt.dcs, dcs)
			}

			if tt.errResp != nil {
				assert.Equal(t, tt.errResp, resp.Error)
			}
		})
	}
}

func TestDataCentersService_Get(t *testing.T) {
	tests := []struct {
		name       string
		id         string
		expected   *DataCenter
		err        string
		errResp    *ErrorResponse
		respStatus int
		respBody   string
	}{
		{
			name: "specific Data Center",
			id:   "loc_a2417980b9874c0",
			expected: &DataCenter{
				ID:        "loc_a2417980b9874c0",
				Name:      "New Town",
				Permalink: "newtown",
				Country: &Country{
					ID:   "ctry_9a989e68e0ad866",
					Name: "USA",
				},
			},
			respStatus: http.StatusOK,
			respBody: `{
  "data_center": {
    "id": "loc_a2417980b9874c0",
    "name": "New Town",
    "permalink": "newtown",
    "country": {
      "id": "ctry_9a989e68e0ad866",
      "name": "USA"
    }
  }
}`,
		},
		{
			name: "non-existent Data Center",
			id:   "loc_nopethisbegone",
			err: "data_center_not_found: No data centers was found matching " +
				"any of the criteria provided in the arguments",
			errResp: &ErrorResponse{
				Code: "data_center_not_found",
				Description: "No data centers was found matching any of the " +
					"criteria provided in the arguments",
				Detail: json.RawMessage(`{}`),
			},
			respStatus: http.StatusNotFound,
			//nolint:lll
			respBody: `{
  "error": {
    "code": "data_center_not_found",
    "description": "No data centers was found matching any of the criteria provided in the arguments",
    "detail": {}
  }
}`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, mux, _, teardown := setup()
			defer teardown()

			mux.HandleFunc(fmt.Sprintf("/core/v1/data_centers/%s", tt.id),
				func(w http.ResponseWriter, r *http.Request) {
					assert.Equal(t, "GET", r.Method)
					w.WriteHeader(tt.respStatus)
					fmt.Fprint(w, tt.respBody)
				},
			)

			dc, resp, err := c.DataCenters.Get(context.Background(), tt.id)

			if tt.err == "" {
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, err, tt.err)
			}

			if tt.expected != nil {
				assert.Equal(t, tt.expected, dc)
			}

			if tt.errResp != nil {
				assert.Equal(t, tt.errResp, resp.Error)
			}
		})
	}
}
