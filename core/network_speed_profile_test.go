package core

import (
	"context"
	"net/http"
	"testing"

	"github.com/krystal/go-katapult"
	"github.com/krystal/go-katapult/internal/testclient"
	"github.com/stretchr/testify/assert"
)

func TestClient_NetworkSpeedProfiles(t *testing.T) {
	c := New(&testclient.Client{})

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

func TestNetworkSpeedProfile_Ref(t *testing.T) {
	tests := []struct {
		name string
		obj  NetworkSpeedProfile
		want NetworkSpeedProfileRef
	}{
		{
			name: "empty",
			obj:  NetworkSpeedProfile{},
			want: NetworkSpeedProfileRef{},
		},
		{
			name: "with id",
			obj: NetworkSpeedProfile{
				ID:                  "nsp_CReSzkaCt01kWoi7",
				Name:                "1 Gbps",
				UploadSpeedInMbit:   100,
				DownloadSpeedInMbit: 1000,
				Permalink:           "1gbps",
			},
			want: NetworkSpeedProfileRef{ID: "nsp_CReSzkaCt01kWoi7"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.obj.Ref()

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
		org  OrganizationRef
		opts *ListOptions
	}
	tests := []struct {
		name           string
		args           args
		want           []*NetworkSpeedProfile
		wantPagination *katapult.Pagination
		errStr         string
		errResp        *katapult.ResponseError
		errIs          error
		respStatus     int
		respBody       []byte
	}{
		{
			name: "by organization ID",
			args: args{
				ctx: context.Background(),
				org: OrganizationRef{ID: "org_O648YDMEYeLmqdmn"},
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
				org: OrganizationRef{SubDomain: "acme"},
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
				org:  OrganizationRef{ID: "org_O648YDMEYeLmqdmn"},
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
				org:  OrganizationRef{ID: "org_O648YDMEYeLmqdmn"},
				opts: &ListOptions{Page: 2, PerPage: 2},
			},
			want: []*NetworkSpeedProfile{
				{
					ID:        "nsp_m2yvaph9SoMFbupJ",
					Name:      "10 Mbps",
					Permalink: "10mbps",
				},
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
				org: OrganizationRef{ID: "org_O648YDMEYeLmqdmn"},
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
				org: OrganizationRef{ID: "org_nopethisbegone"},
			},
			errStr:     fixtureOrganizationNotFoundErr,
			errResp:    fixtureOrganizationNotFoundResponseError,
			errIs:      ErrOrganizationNotFound,
			respStatus: http.StatusNotFound,
			respBody:   fixture("organization_not_found_error"),
		},
		{
			name: "suspended organization",
			args: args{
				ctx: context.Background(),
				org: OrganizationRef{ID: "org_O648YDMEYeLmqdmn"},
			},
			errStr:     fixtureOrganizationSuspendedErr,
			errResp:    fixtureOrganizationSuspendedResponseError,
			errIs:      ErrOrganizationSuspended,
			respStatus: http.StatusForbidden,
			respBody:   fixture("organization_suspended_error"),
		},
		{
			name: "nil context",
			args: args{
				ctx: nil,
				org: OrganizationRef{ID: "org_O648YDMEYeLmqdmn"},
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

					assert.Equal(t,
						*queryValues(tt.args.org, tt.args.opts), r.URL.Query(),
					)

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

			if tt.errIs != nil {
				assert.ErrorIs(t, err, tt.errIs)
			}
		})
	}
}
