package katapult

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	fixtureDiskTemplateNotFoundErr = "disk_template_not_found: No disk " +
		"template was found matching any of the criteria provided in the " +
		"arguments"
	fixtureDiskTemplateNotFoundResponseError = &ResponseError{
		Code: "disk_template_not_found",
		Description: "No disk template was found matching any of the " +
			"criteria provided in the arguments",
		Detail: json.RawMessage(`{}`),
	}
)

func TestDiskTemplate_JSONMarshaling(t *testing.T) {
	tests := []struct {
		name string
		obj  *DiskTemplate
	}{
		{
			name: "empty",
			obj:  &DiskTemplate{},
		},
		{
			name: "full",
			obj: &DiskTemplate{
				ID:              "dtpl_ytP13XD5DE1RdSL9",
				Name:            "Ubuntu 18.04 Server",
				Description:     "A clean installation of Ubuntu 18.04 server",
				Permalink:       "templates/ubuntu-18-04",
				Universal:       true,
				LatestVersion:   &DiskTemplateVersion{ID: "id2"},
				OperatingSystem: &OperatingSystem{ID: "id3"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testJSONMarshaling(t, tt.obj)
		})
	}
}

func TestDiskTemplate_LookupReference(t *testing.T) {
	tests := []struct {
		name string
		obj  *DiskTemplate
		want *DiskTemplate
	}{
		{
			name: "nil",
			obj:  (*DiskTemplate)(nil),
			want: nil,
		},
		{
			name: "empty",
			obj:  &DiskTemplate{},
			want: &DiskTemplate{},
		},
		{
			name: "full",
			obj: &DiskTemplate{
				ID:              "dtpl_ytP13XD5DE1RdSL9",
				Name:            "Ubuntu 18.04 Server",
				Description:     "A clean installation of Ubuntu 18.04 server",
				Permalink:       "templates/ubuntu-18-04",
				Universal:       true,
				LatestVersion:   &DiskTemplateVersion{ID: "id2"},
				OperatingSystem: &OperatingSystem{ID: "id3"},
			},
			want: &DiskTemplate{ID: "dtpl_ytP13XD5DE1RdSL9"},
		},
		{
			name: "no ID",
			obj: &DiskTemplate{
				Name:            "Ubuntu 18.04 Server",
				Description:     "A clean installation of Ubuntu 18.04 server",
				Permalink:       "templates/ubuntu-18-04",
				Universal:       true,
				LatestVersion:   &DiskTemplateVersion{ID: "id2"},
				OperatingSystem: &OperatingSystem{ID: "id3"},
			},
			want: &DiskTemplate{Permalink: "templates/ubuntu-18-04"},
		},
		{
			name: "no ID or Permalink",
			obj: &DiskTemplate{
				Name:            "Ubuntu 18.04 Server",
				Description:     "A clean installation of Ubuntu 18.04 server",
				Universal:       true,
				LatestVersion:   &DiskTemplateVersion{ID: "id2"},
				OperatingSystem: &OperatingSystem{ID: "id3"},
			},
			want: &DiskTemplate{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.obj.LookupReference()

			assert.Equal(t, tt.want, got)
		})
	}
}

func TestDiskTemplateVersion_JSONMarshaling(t *testing.T) {
	tests := []struct {
		name string
		obj  *DiskTemplateVersion
	}{
		{
			name: "empty",
			obj:  &DiskTemplateVersion{},
		},
		{
			name: "full",
			obj: &DiskTemplateVersion{
				ID:       "id2",
				Number:   398,
				Stable:   true,
				SizeInGB: 434,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testJSONMarshaling(t, tt.obj)
		})
	}
}

func TestDiskTemplateOption_JSONMarshaling(t *testing.T) {
	tests := []struct {
		name string
		obj  *DiskTemplateOption
	}{
		{
			name: "empty",
			obj:  &DiskTemplateOption{},
		},
		{
			name: "full",
			obj: &DiskTemplateOption{
				Key:   "hello",
				Value: "world",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testJSONMarshaling(t, tt.obj)
		})
	}
}

func Test_diskTemplateResponseBody_JSONMarshaling(t *testing.T) {
	tests := []struct {
		name string
		obj  *diskTemplateResponseBody
	}{
		{
			name: "empty",
			obj:  &diskTemplateResponseBody{},
		},
		{
			name: "full",
			obj: &diskTemplateResponseBody{
				Pagination:    &Pagination{CurrentPage: 42},
				DiskTemplates: []*DiskTemplate{{ID: "id2"}},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testJSONMarshaling(t, tt.obj)
		})
	}
}

func TestDiskTemplateListOptions_Values(t *testing.T) {
	tests := []struct {
		name string
		opts *DiskTemplateListOptions
		want *url.Values
	}{
		{
			name: "nil *DiskTemplateListOptions",
			opts: nil,
			want: &url.Values{},
		},
		{
			name: "empty *DiskTemplateListOptions",
			opts: &DiskTemplateListOptions{},
			want: &url.Values{},
		},
		{
			name: "zero'd values",
			opts: &DiskTemplateListOptions{Page: 0, PerPage: 0},
			want: &url.Values{},
		},
		{
			name: "non-zero Page value",
			opts: &DiskTemplateListOptions{Page: 3},
			want: &url.Values{"page": []string{"3"}},
		},
		{
			name: "non-zero PerPage value",
			opts: &DiskTemplateListOptions{PerPage: 15},
			want: &url.Values{"per_page": []string{"15"}},
		},
		{
			name: "non-zero IncludeUniversal value",
			opts: &DiskTemplateListOptions{IncludeUniversal: true},
			want: &url.Values{"include_universal": []string{"true"}},
		},
		{
			name: "non-zero Page, PerPage, and IncludeUniversal values",
			opts: &DiskTemplateListOptions{
				IncludeUniversal: true,
				Page:             5,
				PerPage:          15,
			},
			want: &url.Values{
				"include_universal": []string{"true"},
				"page":              []string{"5"},
				"per_page":          []string{"15"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.opts.Values()

			assert.Equal(t, tt.want, got)
		})
	}
}

func TestDiskTemplatesClient_List(t *testing.T) {
	// Correlates to fixtures/disk_templates_list*.json
	certificateList := []*DiskTemplate{
		{
			ID:   "dtpl_YCTIgR4rE2fSgbW0",
			Name: "CentOS 8.0",
		},
		{
			ID:   "dtpl_KXGG3fOWbJqvZvoq",
			Name: "Debian 10",
		},
		{
			ID:          "dtpl_ytP13XD5DE1RdSL9",
			Name:        "Ubuntu 18.04 Server",
			Description: "A clean installation of Ubuntu 18.04 server",
		},
	}

	type args struct {
		ctx   context.Context
		orgID string
		opts  *DiskTemplateListOptions
	}
	tests := []struct {
		name       string
		args       args
		expected   []*DiskTemplate
		pagination *Pagination
		errStr     string
		errResp    *ResponseError
		respStatus int
		respBody   []byte
	}{
		{
			name: "disk templates",
			args: args{
				ctx:   context.Background(),
				orgID: "org_O648YDMEYeLmqdmn",
			},
			expected: certificateList,
			pagination: &Pagination{
				CurrentPage: 1,
				TotalPages:  1,
				Total:       3,
				PerPage:     30,
				LargeSet:    false,
			},
			respStatus: http.StatusOK,
			respBody:   fixture("disk_templates_list"),
		},
		{
			name: "include universal",
			args: args{
				ctx:   context.Background(),
				orgID: "org_O648YDMEYeLmqdmn",
				opts:  &DiskTemplateListOptions{IncludeUniversal: true},
			},
			expected: certificateList,
			pagination: &Pagination{
				CurrentPage: 1,
				TotalPages:  1,
				Total:       3,
				PerPage:     30,
				LargeSet:    false,
			},
			respStatus: http.StatusOK,
			respBody:   fixture("disk_templates_list"),
		},
		{
			name: "page 1 of disk templates",
			args: args{
				ctx:   context.Background(),
				orgID: "org_O648YDMEYeLmqdmn",
				opts:  &DiskTemplateListOptions{Page: 1, PerPage: 2},
			},
			expected: certificateList[0:2],
			pagination: &Pagination{
				CurrentPage: 1,
				TotalPages:  2,
				Total:       3,
				PerPage:     2,
				LargeSet:    false,
			},
			respStatus: http.StatusOK,
			respBody:   fixture("disk_templates_list_page_1"),
		},
		{
			name: "page 2 of disk templates",
			args: args{
				ctx:   context.Background(),
				orgID: "org_O648YDMEYeLmqdmn",
				opts:  &DiskTemplateListOptions{Page: 2, PerPage: 2},
			},
			expected: certificateList[2:],
			pagination: &Pagination{
				CurrentPage: 2,
				TotalPages:  2,
				Total:       3,
				PerPage:     2,
				LargeSet:    false,
			},
			respStatus: http.StatusOK,
			respBody:   fixture("disk_templates_list_page_2"),
		},
		{
			name: "invalid API token response",
			args: args{
				ctx:   context.Background(),
				orgID: "org_O648YDMEYeLmqdmn",
			},
			errStr:     fixtureInvalidAPITokenErr,
			errResp:    fixtureInvalidAPITokenResponseError,
			respStatus: http.StatusForbidden,
			respBody:   fixture("invalid_api_token_error"),
		},
		{
			name: "non-existent organization",
			args: args{
				ctx:   context.Background(),
				orgID: "org_O648YDMEYeLmqdmn",
			},
			errStr:     fixtureOrganizationNotFoundErr,
			errResp:    fixtureOrganizationNotFoundResponseError,
			respStatus: http.StatusNotFound,
			respBody:   fixture("organization_not_found_error"),
		},
		{
			name: "suspended organization",
			args: args{
				ctx:   context.Background(),
				orgID: "org_O648YDMEYeLmqdmn",
			},
			errStr:     fixtureOrganizationSuspendedErr,
			errResp:    fixtureOrganizationSuspendedResponseError,
			respStatus: http.StatusForbidden,
			respBody:   fixture("organization_suspended_error"),
		},
		{
			name: "nil context",
			args: args{
				ctx:   nil,
				orgID: "org_O648YDMEYeLmqdmn",
			},
			errStr: "net/http: nil Context",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, mux, _, teardown := prepareTestClient()
			defer teardown()

			mux.HandleFunc(
				fmt.Sprintf(
					"/core/v1/organizations/%s/disk_templates", tt.args.orgID,
				),
				func(w http.ResponseWriter, r *http.Request) {
					assert.Equal(t, "GET", r.Method)
					assertEmptyFieldSpec(t, r)
					assertAuthorization(t, r)
					if tt.args.opts != nil {
						assert.Equal(t, *tt.args.opts.Values(), r.URL.Query())
					}

					w.WriteHeader(tt.respStatus)
					_, _ = w.Write(tt.respBody)
				},
			)

			got, resp, err := c.DiskTemplates.List(
				tt.args.ctx, tt.args.orgID, tt.args.opts,
			)

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

			if tt.pagination != nil {
				assert.Equal(t, tt.pagination, resp.Pagination)
			}

			if tt.errResp != nil {
				assert.Equal(t, tt.errResp, resp.Error)
			}
		})
	}
}
