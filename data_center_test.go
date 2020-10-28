package katapult

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDataCenter_JSONMarshaling(t *testing.T) {
	tests := []struct {
		name string
		obj  *DataCenter
	}{
		{
			name: "empty",
			obj:  &DataCenter{},
		},
		{
			name: "full",
			obj: &DataCenter{
				ID:        "id",
				Name:      "name",
				Permalink: "permalink",
				Country: &Country{
					ID: "id2",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testJSONMarshaling(t, tt.obj)
		})
	}
}

func Test_dataCentersResponseBody_JSONMarshaling(t *testing.T) {
	tests := []struct {
		name string
		obj  *dataCentersResponseBody
	}{
		{
			name: "empty",
			obj:  &dataCentersResponseBody{},
		},
		{
			name: "full",
			obj: &dataCentersResponseBody{
				DataCenter:  &DataCenter{ID: "id1"},
				DataCenters: []*DataCenter{{ID: "id2"}},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testJSONMarshaling(t, tt.obj)
		})
	}
}

func TestDataCentersResource_List(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name       string
		args       args
		expected   []*DataCenter
		errStr     string
		errResp    *ResponseError
		respStatus int
		respBody   []byte
	}{
		{
			name: "data centers",
			args: args{
				ctx: context.Background(),
			},
			expected: []*DataCenter{
				{
					ID:        "loc_25d48761871e4bf",
					Name:      "Shirebury",
					Permalink: "shirebury",
				},
				{
					ID:        "loc_a2417980b9874c0",
					Name:      "New Town",
					Permalink: "newtown",
				},
			},
			respStatus: http.StatusOK,
			respBody:   fixture("data_centers_list"),
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

			mux.HandleFunc("/core/v1/data_centers",
				func(w http.ResponseWriter, r *http.Request) {
					assert.Equal(t, "GET", r.Method)
					assert.Equal(t, "", r.Header.Get("X-Field-Spec"))

					w.WriteHeader(tt.respStatus)
					_, _ = w.Write(tt.respBody)
				},
			)

			got, resp, err := c.DataCenters.List(tt.args.ctx)

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

func TestDataCentersResource_Get(t *testing.T) {
	// Correlates to fixtures/data_center_get.json
	dataCenter := &DataCenter{
		ID:        "loc_a2417980b9874c0",
		Name:      "New Town",
		Permalink: "newtown",
	}

	type args struct {
		ctx context.Context
		id  string
	}
	tests := []struct {
		name       string
		args       args
		expected   *DataCenter
		errStr     string
		errResp    *ResponseError
		respStatus int
		respBody   []byte
	}{
		{
			name: "data center",
			args: args{
				ctx: context.Background(),
				id:  "loc_a2417980b9874c0",
			},
			expected:   dataCenter,
			respStatus: http.StatusOK,
			respBody:   fixture("data_center_get"),
		},
		{
			name: "non-existent data center",
			args: args{
				ctx: context.Background(),
				id:  "loc_nopethisbegone",
			},
			errStr: "data_center_not_found: No data centers was found " +
				"matching any of the criteria provided in the arguments",
			errResp: &ResponseError{
				Code: "data_center_not_found",
				Description: "No data centers was found matching any of the " +
					"criteria provided in the arguments",
				Detail: json.RawMessage(`{}`),
			},
			respStatus: http.StatusNotFound,
			respBody:   fixture("data_center_not_found_error"),
		},
		{
			name: "nil context",
			args: args{
				ctx: nil,
				id:  "loc_a2417980b9874c0",
			},
			errStr: "net/http: nil Context",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, mux, _, teardown := prepareTestClient()
			defer teardown()

			mux.HandleFunc(fmt.Sprintf("/core/v1/data_centers/%s", tt.args.id),
				func(w http.ResponseWriter, r *http.Request) {
					assert.Equal(t, "GET", r.Method)
					assert.Equal(t, "", r.Header.Get("X-Field-Spec"))

					w.WriteHeader(tt.respStatus)
					_, _ = w.Write(tt.respBody)
				},
			)

			got, resp, err := c.DataCenters.Get(tt.args.ctx, tt.args.id)

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
