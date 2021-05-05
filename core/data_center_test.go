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

var (
	fixtureDataCenterNotFoundErr = "data_center_not_found: No data centers " +
		"was found matching any of the criteria provided in the arguments"
	fixtureDataCenterNotFoundResponseError = &katapult.ResponseError{
		Code: "data_center_not_found",
		Description: "No data centers was found matching any of the " +
			"criteria provided in the arguments",
		Detail: json.RawMessage(`{}`),
	}
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
				ID:        "dc_25d48761871e4bf",
				Name:      "Shirebury",
				Permalink: "shirebury",
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

func TestNewDataCenterLookup(t *testing.T) {
	type args struct {
		idOrPermalink string
	}
	tests := []struct {
		name  string
		args  args
		want  *DataCenter
		field FieldName
	}{
		{
			name:  "empty string",
			args:  args{idOrPermalink: ""},
			want:  &DataCenter{},
			field: PermalinkField,
		},
		{
			name:  "dc_ prefixed ID",
			args:  args{idOrPermalink: "dc_HHwnaBCIwNHqv0aO"},
			want:  &DataCenter{ID: "dc_HHwnaBCIwNHqv0aO"},
			field: IDField,
		},
		{
			name:  "loc_ prefixed ID",
			args:  args{idOrPermalink: "loc_RuHTM4fyzucbYGCK"},
			want:  &DataCenter{ID: "loc_RuHTM4fyzucbYGCK"},
			field: IDField,
		},
		{
			name:  "permalink",
			args:  args{idOrPermalink: "country-city-1"},
			want:  &DataCenter{Permalink: "country-city-1"},
			field: PermalinkField,
		},
		{
			name:  "random text",
			args:  args{idOrPermalink: "dXUt33rNLmbatuAa"},
			want:  &DataCenter{Permalink: "dXUt33rNLmbatuAa"},
			field: PermalinkField,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, field := NewDataCenterLookup(tt.args.idOrPermalink)

			assert.Equal(t, tt.want, got)
			assert.Equal(t, tt.field, field)
		})
	}
}

func TestDataCenter_lookupReference(t *testing.T) {
	tests := []struct {
		name string
		obj  *DataCenter
		want *DataCenter
	}{
		{
			name: "nil",
			obj:  nil,
			want: nil,
		},
		{
			name: "empty",
			obj:  &DataCenter{},
			want: &DataCenter{},
		},
		{
			name: "full",
			obj: &DataCenter{
				ID:        "dc_25d48761871e4bf",
				Name:      "Shirebury",
				Permalink: "shirebury",
				Country: &Country{
					ID: "id2",
				},
			},
			want: &DataCenter{ID: "dc_25d48761871e4bf"},
		},
		{
			name: "no ID",
			obj: &DataCenter{
				Name:      "Shirebury",
				Permalink: "shirebury",
				Country: &Country{
					ID: "id2",
				},
			},
			want: &DataCenter{Permalink: "shirebury"},
		},
		{
			name: "no ID or Permalink",
			obj: &DataCenter{
				Name: "Shirebury",
				Country: &Country{
					ID: "id2",
				},
			},
			want: &DataCenter{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.obj.lookupReference()

			assert.Equal(t, tt.want, got)
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

func TestDataCentersClient_List(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name       string
		args       args
		want       []*DataCenter
		errStr     string
		errResp    *katapult.ResponseError
		respStatus int
		respBody   []byte
	}{
		{
			name: "data centers",
			args: args{
				ctx: context.Background(),
			},
			want: []*DataCenter{
				{
					ID:        "dc_25d48761871e4bf",
					Name:      "Shirebury",
					Permalink: "shirebury",
				},
				{
					ID:        "dc_a2417980b9874c0",
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
			rm, mux, _, teardown := prepareTestClient(t)
			defer teardown()
			c := NewDataCentersClient(rm)

			mux.HandleFunc("/core/v1/data_centers",
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

func TestDataCentersClient_Get(t *testing.T) {
	// Correlates to fixtures/data_center_get.json
	dataCenter := &DataCenter{
		ID:        "dc_a2417980b9874c0",
		Name:      "New Town",
		Permalink: "newtown",
	}

	type args struct {
		ctx           context.Context
		idOrPermalink string
	}
	tests := []struct {
		name       string
		args       args
		reqPath    string
		reqQuery   *url.Values
		want       *DataCenter
		errStr     string
		errResp    *katapult.ResponseError
		respStatus int
		respBody   []byte
	}{
		{
			name: "by ID",
			args: args{
				ctx:           context.Background(),
				idOrPermalink: "dc_a2417980b9874c0",
			},
			reqPath:    "data_centers/dc_a2417980b9874c0",
			want:       dataCenter,
			respStatus: http.StatusOK,
			respBody:   fixture("data_center_get"),
		},
		{
			name: "by legacy ID",
			args: args{
				ctx:           context.Background(),
				idOrPermalink: "loc_a2417980b9874c0",
			},
			reqPath:    "data_centers/loc_a2417980b9874c0",
			want:       dataCenter,
			respStatus: http.StatusOK,
			respBody:   fixture("data_center_get"),
		},
		{
			name: "by Permalink",
			args: args{
				ctx:           context.Background(),
				idOrPermalink: "newtown",
			},
			reqPath: "data_centers/_",
			reqQuery: &url.Values{
				"data_center[permalink]": []string{"newtown"},
			},
			want:       dataCenter,
			respStatus: http.StatusOK,
			respBody:   fixture("data_center_get"),
		},
		{
			name: "non-existent data center",
			args: args{
				ctx:           context.Background(),
				idOrPermalink: "dc_nopethisbegone",
			},
			errStr:     fixtureDataCenterNotFoundErr,
			errResp:    fixtureDataCenterNotFoundResponseError,
			respStatus: http.StatusNotFound,
			respBody:   fixture("data_center_not_found_error"),
		},
		{
			name: "nil context",
			args: args{
				ctx:           nil,
				idOrPermalink: "dc_a2417980b9874c0",
			},
			errStr: "net/http: nil Context",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rm, mux, _, teardown := prepareTestClient(t)
			defer teardown()
			c := NewDataCentersClient(rm)

			path := fmt.Sprintf("data_centers/%s", tt.args.idOrPermalink)
			if tt.reqPath != "" {
				path = tt.reqPath
			}

			mux.HandleFunc(
				"/core/v1/"+path,
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
				tt.args.ctx, tt.args.idOrPermalink,
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

func TestDataCentersClient_GetByID(t *testing.T) {
	// Correlates to fixtures/data_center_get.json
	dataCenter := &DataCenter{
		ID:        "dc_a2417980b9874c0",
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
		want       *DataCenter
		errStr     string
		errResp    *katapult.ResponseError
		respStatus int
		respBody   []byte
	}{
		{
			name: "data center",
			args: args{
				ctx: context.Background(),
				id:  "dc_a2417980b9874c0",
			},
			want:       dataCenter,
			respStatus: http.StatusOK,
			respBody:   fixture("data_center_get"),
		},
		{
			name: "non-existent data center",
			args: args{
				ctx: context.Background(),
				id:  "dc_nopethisbegone",
			},
			errStr:     fixtureDataCenterNotFoundErr,
			errResp:    fixtureDataCenterNotFoundResponseError,
			respStatus: http.StatusNotFound,
			respBody:   fixture("data_center_not_found_error"),
		},
		{
			name: "nil context",
			args: args{
				ctx: nil,
				id:  "dc_a2417980b9874c0",
			},
			errStr: "net/http: nil Context",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rm, mux, _, teardown := prepareTestClient(t)
			defer teardown()
			c := NewDataCentersClient(rm)

			mux.HandleFunc(fmt.Sprintf("/core/v1/data_centers/%s", tt.args.id),
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

func TestDataCentersClient_GetByPermalink(t *testing.T) {
	// Correlates to fixtures/data_center_get.json
	dataCenter := &DataCenter{
		ID:        "dc_a2417980b9874c0",
		Name:      "New Town",
		Permalink: "newtown",
	}

	type args struct {
		ctx       context.Context
		permalink string
	}
	tests := []struct {
		name       string
		args       args
		want       *DataCenter
		errStr     string
		errResp    *katapult.ResponseError
		respStatus int
		respBody   []byte
	}{
		{
			name: "data center",
			args: args{
				ctx:       context.Background(),
				permalink: "newtown",
			},
			want:       dataCenter,
			respStatus: http.StatusOK,
			respBody:   fixture("data_center_get"),
		},
		{
			name: "non-existent data center",
			args: args{
				ctx:       context.Background(),
				permalink: "not-here",
			},
			errStr:     fixtureDataCenterNotFoundErr,
			errResp:    fixtureDataCenterNotFoundResponseError,
			respStatus: http.StatusNotFound,
			respBody:   fixture("data_center_not_found_error"),
		},
		{
			name: "nil context",
			args: args{
				ctx:       nil,
				permalink: "newtown",
			},
			errStr: "net/http: nil Context",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rm, mux, _, teardown := prepareTestClient(t)
			defer teardown()
			c := NewDataCentersClient(rm)

			mux.HandleFunc("/core/v1/data_centers/_",
				func(w http.ResponseWriter, r *http.Request) {
					assert.Equal(t, "GET", r.Method)
					assertEmptyFieldSpec(t, r)
					assertAuthorization(t, r)

					qs := url.Values{
						"data_center[permalink]": []string{tt.args.permalink},
					}
					assert.Equal(t, qs, r.URL.Query())

					w.WriteHeader(tt.respStatus)
					_, _ = w.Write(tt.respBody)
				},
			)

			got, resp, err := c.GetByPermalink(
				tt.args.ctx, tt.args.permalink,
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
