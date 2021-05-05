package core

import (
	"context"
	"net/http"
	"net/url"
	"testing"

	"github.com/krystal/go-katapult"
	"github.com/stretchr/testify/assert"
)

func TestClient_NetworkSpeedProfiles(t *testing.T) {
	c := New(&fakeRequestMaker{})

	assert.IsType(t, &NetworkSpeedProfilesClient{}, c.NetworkSpeedProfiles)
}

func TestNetworkSpeedProfile_JSONMarshaling(t *testing.T) {
	tests := []struct {
		name string
		obj  *NetworkSpeedProfile
	}{
		{
			name: "empty",
			obj:  &NetworkSpeedProfile{},
		},
		{
			name: "full",
			obj: &NetworkSpeedProfile{
				ID:                  "nsp_CReSzkaCt01kWoi7",
				Name:                "1 Gbps",
				UploadSpeedInMbit:   100,
				DownloadSpeedInMbit: 1000,
				Permalink:           "1gbps",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testJSONMarshaling(t, tt.obj)
		})
	}
}

func TestNewNetworkSpeedProfileLookup(t *testing.T) {
	type args struct {
		idOrPermalink string
	}
	tests := []struct {
		name  string
		args  args
		want  *NetworkSpeedProfile
		field FieldName
	}{
		{
			name:  "empty string",
			args:  args{idOrPermalink: ""},
			want:  &NetworkSpeedProfile{},
			field: PermalinkField,
		},
		{
			name:  "nsp_ prefixed ID",
			args:  args{idOrPermalink: "nsp_wEyUfJ74ZQu2KmZr"},
			want:  &NetworkSpeedProfile{ID: "nsp_wEyUfJ74ZQu2KmZr"},
			field: IDField,
		},
		{
			name:  "permalink",
			args:  args{idOrPermalink: "10gbps"},
			want:  &NetworkSpeedProfile{Permalink: "10gbps"},
			field: PermalinkField,
		},
		{
			name:  "random text",
			args:  args{idOrPermalink: "kKAvlqM1FVEn3NAG"},
			want:  &NetworkSpeedProfile{Permalink: "kKAvlqM1FVEn3NAG"},
			field: PermalinkField,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, field := NewNetworkSpeedProfileLookup(tt.args.idOrPermalink)

			assert.Equal(t, tt.want, got)
			assert.Equal(t, tt.field, field)
		})
	}
}

func TestNetworkSpeedProfile_lookupReference(t *testing.T) {
	tests := []struct {
		name string
		obj  *NetworkSpeedProfile
		want *NetworkSpeedProfile
	}{
		{
			name: "nil",
			obj:  nil,
			want: nil,
		},
		{
			name: "empty",
			obj:  &NetworkSpeedProfile{},
			want: &NetworkSpeedProfile{},
		},
		{
			name: "full",
			obj: &NetworkSpeedProfile{
				ID:                  "nsp_CReSzkaCt01kWoi7",
				Name:                "1 Gbps",
				UploadSpeedInMbit:   100,
				DownloadSpeedInMbit: 1000,
				Permalink:           "1gbps",
			},
			want: &NetworkSpeedProfile{ID: "nsp_CReSzkaCt01kWoi7"},
		},
		{
			name: "no ID",
			obj: &NetworkSpeedProfile{
				Name:                "1 Gbps",
				UploadSpeedInMbit:   100,
				DownloadSpeedInMbit: 1000,
				Permalink:           "1gbps",
			},
			want: &NetworkSpeedProfile{Permalink: "1gbps"},
		},
		{
			name: "no ID or Permalink",
			obj: &NetworkSpeedProfile{
				Name:                "1 Gbps",
				UploadSpeedInMbit:   100,
				DownloadSpeedInMbit: 1000,
			},
			want: &NetworkSpeedProfile{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.obj.lookupReference()

			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_networkSpeedProfileResponseBody_JSONMarshaling(t *testing.T) {
	tests := []struct {
		name string
		obj  *networkSpeedProfileResponseBody
	}{
		{
			name: "empty",
			obj:  &networkSpeedProfileResponseBody{},
		},
		{
			name: "full",
			obj: &networkSpeedProfileResponseBody{
				Pagination: &katapult.Pagination{
					CurrentPage: 1,
					PerPage:     40,
				},
				NetworkSpeedProfiles: []*NetworkSpeedProfile{
					{
						ID:   "nsp_CReSzkaCt01kWoi7",
						Name: "1 Gbps",
					},
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

func TestNetworkSpeedTestsClient_List(t *testing.T) {
	type args struct {
		ctx  context.Context
		org  *Organization
		opts *ListOptions
	}
	tests := []struct {
		name           string
		args           args
		want           []*NetworkSpeedProfile
		wantQuery      *url.Values
		wantPagination *katapult.Pagination
		errStr         string
		errResp        *katapult.ResponseError
		respStatus     int
		respBody       []byte
	}{
		{
			name: "by organization ID",
			args: args{
				ctx: context.Background(),
				org: &Organization{ID: "org_O648YDMEYeLmqdmn"},
			},
			want: []*NetworkSpeedProfile{
				{
					ID:        "nsp_H3Mknnus3dtDIbIc",
					Name:      "1 Gbps",
					Permalink: "1gbps",
				},
				{
					ID:        "nsp_m2yvaph9SoMFbupJ",
					Name:      "100 Mbps",
					Permalink: "100mbps",
				},
				{
					ID:        "nsp_m2yvaph9SoMFbupJ",
					Name:      "10 Mbps",
					Permalink: "10mbps",
				},
			},
			wantQuery: &url.Values{
				"organization[id]": []string{"org_O648YDMEYeLmqdmn"},
			},
			wantPagination: &katapult.Pagination{
				CurrentPage: 1,
				TotalPages:  1,
				Total:       3,
				PerPage:     30,
				LargeSet:    false,
			},
			respStatus: http.StatusOK,
			respBody:   fixture("network_speed_profiles_list"),
		},
		{
			name: "by organization SubDomain",
			args: args{
				ctx: context.Background(),
				org: &Organization{SubDomain: "acme"},
			},
			want: []*NetworkSpeedProfile{
				{
					ID:        "nsp_H3Mknnus3dtDIbIc",
					Name:      "1 Gbps",
					Permalink: "1gbps",
				},
				{
					ID:        "nsp_m2yvaph9SoMFbupJ",
					Name:      "100 Mbps",
					Permalink: "100mbps",
				},
				{
					ID:        "nsp_m2yvaph9SoMFbupJ",
					Name:      "10 Mbps",
					Permalink: "10mbps",
				},
			},
			wantQuery: &url.Values{
				"organization[sub_domain]": []string{"acme"},
			},
			wantPagination: &katapult.Pagination{
				CurrentPage: 1,
				TotalPages:  1,
				Total:       3,
				PerPage:     30,
				LargeSet:    false,
			},
			respStatus: http.StatusOK,
			respBody:   fixture("network_speed_profiles_list"),
		},
		{
			name: "page 1",
			args: args{
				ctx:  context.Background(),
				org:  &Organization{ID: "org_O648YDMEYeLmqdmn"},
				opts: &ListOptions{Page: 1, PerPage: 2},
			},
			want: []*NetworkSpeedProfile{
				{
					ID:        "nsp_H3Mknnus3dtDIbIc",
					Name:      "1 Gbps",
					Permalink: "1gbps",
				},
				{
					ID:        "nsp_m2yvaph9SoMFbupJ",
					Name:      "100 Mbps",
					Permalink: "100mbps",
				},
			},
			wantQuery: &url.Values{
				"organization[id]": []string{"org_O648YDMEYeLmqdmn"},
				"page":             []string{"1"},
				"per_page":         []string{"2"},
			},
			wantPagination: &katapult.Pagination{
				CurrentPage: 1,
				TotalPages:  2,
				Total:       3,
				PerPage:     2,
				LargeSet:    false,
			},
			respStatus: http.StatusOK,
			respBody:   fixture("network_speed_profiles_list_page_1"),
		},
		{
			name: "page 2",
			args: args{
				ctx:  context.Background(),
				org:  &Organization{ID: "org_O648YDMEYeLmqdmn"},
				opts: &ListOptions{Page: 2, PerPage: 2},
			},
			want: []*NetworkSpeedProfile{
				{
					ID:        "nsp_m2yvaph9SoMFbupJ",
					Name:      "10 Mbps",
					Permalink: "10mbps",
				},
			},
			wantQuery: &url.Values{
				"organization[id]": []string{"org_O648YDMEYeLmqdmn"},
				"page":             []string{"2"},
				"per_page":         []string{"2"},
			},
			wantPagination: &katapult.Pagination{
				CurrentPage: 2,
				TotalPages:  2,
				Total:       3,
				PerPage:     2,
				LargeSet:    false,
			},
			respStatus: http.StatusOK,
			respBody:   fixture("network_speed_profiles_list_page_2"),
		},
		{
			name: "invalid API token response",
			args: args{
				ctx: context.Background(),
				org: &Organization{ID: "org_O648YDMEYeLmqdmn"},
			},
			errStr:     fixtureInvalidAPITokenErr,
			errResp:    fixtureInvalidAPITokenResponseError,
			respStatus: http.StatusForbidden,
			respBody:   fixture("invalid_api_token_error"),
		},
		{
			name: "non-existent organization",
			args: args{
				ctx: context.Background(),
				org: &Organization{ID: "org_nopethisbegone"},
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
				org: &Organization{ID: "org_O648YDMEYeLmqdmn"},
			},
			errStr:     fixtureOrganizationSuspendedErr,
			errResp:    fixtureOrganizationSuspendedResponseError,
			respStatus: http.StatusForbidden,
			respBody:   fixture("organization_suspended_error"),
		},
		{
			name: "nil organization",
			args: args{
				ctx: context.Background(),
				org: nil,
			},
			errStr:     fixtureOrganizationNotFoundErr,
			errResp:    fixtureOrganizationNotFoundResponseError,
			respStatus: http.StatusNotFound,
			respBody:   fixture("organization_not_found_error"),
		},
		{
			name: "nil context",
			args: args{
				ctx: nil,
				org: &Organization{ID: "org_O648YDMEYeLmqdmn"},
			},
			errStr: "net/http: nil Context",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rm, mux, _, teardown := prepareTestClient(t)
			defer teardown()
			c := NewNetworkSpeedProfilesClient(rm)

			mux.HandleFunc(
				"/core/v1/organizations/_/network_speed_profiles",
				func(w http.ResponseWriter, r *http.Request) {
					assert.Equal(t, "GET", r.Method)
					assertEmptyFieldSpec(t, r)
					assertAuthorization(t, r)

					if tt.wantQuery != nil {
						assert.Equal(t, *tt.wantQuery, r.URL.Query())
					} else {
						assert.Equal(t,
							*tt.args.org.queryValues(), r.URL.Query(),
						)
					}

					w.WriteHeader(tt.respStatus)
					_, _ = w.Write(tt.respBody)
				},
			)

			got, resp, err := c.List(
				tt.args.ctx, tt.args.org, tt.args.opts,
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

			if tt.wantPagination != nil {
				assert.Equal(t, tt.wantPagination, resp.Pagination)
			}

			if tt.errResp != nil {
				assert.Equal(t, tt.errResp, resp.Error)
			}
		})
	}
}
