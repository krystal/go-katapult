package core

import (
	"context"
	"encoding/json"
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

func TestClient_DataCenters(t *testing.T) {
	c := New(&fakeRequestMaker{})

	assert.IsType(t, &DataCentersClient{}, c.DataCenters)
}

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

func TestDataCenter_Ref(t *testing.T) {
	tests := []struct {
		name string
		obj  DataCenter
		want DataCenterRef
	}{
		{
			obj:  DataCenter{ID: "dc_25d48761871e4bf"},
			want: DataCenterRef{ID: "dc_25d48761871e4bf"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.obj.Ref()
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestDataCenterRef_queryValues(t *testing.T) {
	tests := []struct {
		name string
		ref  DataCenterRef
		want *url.Values
	}{
		{
			name: "id",
			ref:  DataCenterRef{ID: "dc_25d48761871e4bf"},
			want: &url.Values{
				"data_center[id]": []string{"dc_25d48761871e4bf"},
			},
		},
		{
			name: "permalink",
			ref:  DataCenterRef{Permalink: "central-amazon-jungle"},
			want: &url.Values{
				"data_center[permalink]": []string{"central-amazon-jungle"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(t.Name(), func(t *testing.T) {
			assert.Equal(t, tt.want, tt.ref.queryValues())
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

			mux.HandleFunc("/core/v1/data_centers/_",
				func(w http.ResponseWriter, r *http.Request) {
					assert.Equal(t, "GET", r.Method)
					assertEmptyFieldSpec(t, r)
					assertAuthorization(t, r)

					qs := url.Values{
						"data_center[id]": []string{tt.args.id},
					}
					assert.Equal(t, qs, r.URL.Query())

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
