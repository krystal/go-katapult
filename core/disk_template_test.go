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
	fixtureDiskTemplateNotFoundErr = "disk_template_not_found: No disk " +
		"template was found matching any of the criteria provided in the " +
		"arguments"
	fixtureDiskTemplateNotFoundResponseError = &katapult.ResponseError{
		Code: "disk_template_not_found",
		Description: "No disk template was found matching any of the " +
			"criteria provided in the arguments",
		Detail: json.RawMessage(`{}`),
	}

	fixtureDiskTemplateFull = &DiskTemplate{
		ID:              "dtpl_ytP13XD5DE1RdSL9",
		Name:            "Ubuntu 18.04 Server",
		Description:     "A clean installation of Ubuntu 18.04 server",
		Permalink:       "templates/ubuntu-18-04",
		Universal:       true,
		LatestVersion:   &DiskTemplateVersion{ID: "id2"},
		OperatingSystem: &OperatingSystem{ID: "id3"},
	}
)

func TestClient_DiskTemplates(t *testing.T) {
	c := New(&fakeRequestMaker{})

	assert.IsType(t, &DiskTemplatesClient{}, c.DiskTemplates)
}

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
			obj:  fixtureDiskTemplateFull,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testJSONMarshaling(t, tt.obj)
		})
	}
}

func TestDiskTemplate_Ref(t *testing.T) {
	tests := []struct {
		name string
		obj  DiskTemplate
		want DiskTemplateRef
	}{
		{
			name: "empty",
			obj:  DiskTemplate{},
			want: DiskTemplateRef{},
		},
		{
			name: "ID",
			obj:  DiskTemplate{ID: "dtpl_ytP13XD5DE1RdSL9"},
			want: DiskTemplateRef{ID: "dtpl_ytP13XD5DE1RdSL9"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.obj.Ref()

			assert.Equal(t, tt.want, got)
		})
	}
}

func TestDiskTemplateRef_queryValues(t *testing.T) {
	tests := []struct {
		name string
		obj  DiskTemplateRef
	}{
		{
			name: "empty",
			obj:  DiskTemplateRef{},
		},
		{
			name: "full",
			obj: DiskTemplateRef{
				ID:        "dtpl_ytP13XD5DE1RdSL9",
				Permalink: "templates/ubuntu-18-04",
			},
		},
		{
			name: "just ID",
			obj: DiskTemplateRef{
				ID: "dtpl_ytP13XD5DE1RdSL9",
			},
		},
		{
			name: "just permalink",
			obj: DiskTemplateRef{
				Permalink: "templates/ubuntu-18-04",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testQueryableEncoding(t, tt.obj)
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
				Pagination:    &katapult.Pagination{CurrentPage: 42},
				DiskTemplate:  &DiskTemplate{ID: "id1"},
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

func TestDiskTemplateListOptions_queryValues(t *testing.T) {
	tests := []struct {
		name string
		obj  *DiskTemplateListOptions
		want *url.Values
	}{
		{
			name: "nil",
			obj:  nil,
			want: &url.Values{},
		},
		{
			name: "empty",
			obj:  &DiskTemplateListOptions{},
			want: &url.Values{},
		},
		{
			name: "full",
			obj: &DiskTemplateListOptions{
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
			got := tt.obj.queryValues()

			assert.Equal(t, tt.want, got)
		})
	}
}

func TestDiskTemplatesClient_List(t *testing.T) {
	// Correlates to fixtures/disk_templates_list*.json
	diskTemplateList := []*DiskTemplate{
		{
			ID:        "dtpl_YCTIgR4rE2fSgbW0",
			Name:      "CentOS 8.0",
			Permalink: "templates/centos-8",
		},
		{
			ID:        "dtpl_KXGG3fOWbJqvZvoq",
			Name:      "Debian 10",
			Permalink: "templates/debian-10",
		},
		{
			ID:        "dtpl_ytP13XD5DE1RdSL9",
			Name:      "Ubuntu 18.04 Server",
			Permalink: "templates/ubuntu-18-04",
		},
	}

	type args struct {
		ctx  context.Context
		org  OrganizationRef
		opts *DiskTemplateListOptions
	}
	tests := []struct {
		name           string
		args           args
		want           []*DiskTemplate
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
				org: OrganizationRef{ID: "org_O648YDMEYeLmqdmn"},
			},
			want: diskTemplateList,
			wantPagination: &katapult.Pagination{
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
			name: "by organization subdomain",
			args: args{
				ctx: context.Background(),
				org: OrganizationRef{SubDomain: "valveinc"},
			},
			want: diskTemplateList,
			wantPagination: &katapult.Pagination{
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
				ctx:  context.Background(),
				org:  OrganizationRef{ID: "org_O648YDMEYeLmqdmn"},
				opts: &DiskTemplateListOptions{IncludeUniversal: true},
			},
			want: diskTemplateList,
			wantPagination: &katapult.Pagination{
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
			name: "page 1",
			args: args{
				ctx:  context.Background(),
				org:  OrganizationRef{ID: "org_O648YDMEYeLmqdmn"},
				opts: &DiskTemplateListOptions{Page: 1, PerPage: 2},
			},
			want: diskTemplateList[0:2],
			wantPagination: &katapult.Pagination{
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
			name: "page 2",
			args: args{
				ctx:  context.Background(),
				org:  OrganizationRef{ID: "org_O648YDMEYeLmqdmn"},
				opts: &DiskTemplateListOptions{Page: 2, PerPage: 2},
			},
			want: diskTemplateList[2:],
			wantPagination: &katapult.Pagination{
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
				org: OrganizationRef{ID: "org_O648YDMEYeLmqdmn"},
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
				org: OrganizationRef{ID: "org_O648YDMEYeLmqdmn"},
			},
			errStr:     fixtureOrganizationSuspendedErr,
			errResp:    fixtureOrganizationSuspendedResponseError,
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
			c := NewDiskTemplatesClient(rm)

			mux.HandleFunc(
				"/core/v1/organizations/_/disk_templates",
				func(w http.ResponseWriter, r *http.Request) {
					assert.Equal(t, "GET", r.Method)
					assertEmptyFieldSpec(t, r)
					assertAuthorization(t, r)

					qs := queryValues(tt.args.org, tt.args.opts)
					assert.Equal(t, *qs, r.URL.Query())

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

func TestDiskTemplatesClient_Get(t *testing.T) {
	// Correlates to fixtures/disk_template_get.json
	diskTemplate := &DiskTemplate{
		ID:        "dtpl_ytP13XD5DE1RdSL9",
		Name:      "Ubuntu 18.04 Server",
		Permalink: "templates/ubuntu-18-04",
	}

	type args struct {
		ctx context.Context
		ref DiskTemplateRef
	}
	tests := []struct {
		name       string
		args       args
		want       *DiskTemplate
		wantQuery  *url.Values
		errStr     string
		errResp    *katapult.ResponseError
		respStatus int
		respBody   []byte
	}{
		{
			name: "by ID",
			args: args{
				ctx: context.Background(),
				ref: DiskTemplateRef{ID: "dtpl_ytP13XD5DE1RdSL9"},
			},
			want: diskTemplate,
			wantQuery: &url.Values{
				"disk_template[id]": []string{"dtpl_ytP13XD5DE1RdSL9"},
			},
			respStatus: http.StatusOK,
			respBody:   fixture("disk_template_get"),
		},
		{
			name: "by Permalink",
			args: args{
				ctx: context.Background(),
				ref: DiskTemplateRef{Permalink: "public"},
			},
			wantQuery: &url.Values{
				"disk_template[permalink]": []string{"public"},
			},
			want:       diskTemplate,
			respStatus: http.StatusOK,
			respBody:   fixture("disk_template_get"),
		},
		{
			name: "non-existent disk template by ID",
			args: args{
				ctx: context.Background(),
				ref: DiskTemplateRef{ID: "dtpl_nopethisbegone"},
			},
			errStr:     fixtureDiskTemplateNotFoundErr,
			errResp:    fixtureDiskTemplateNotFoundResponseError,
			respStatus: http.StatusNotFound,
			respBody:   fixture("disk_template_not_found_error"),
		},
		{
			name: "non-existent disk template by Permalink",
			args: args{
				ctx: context.Background(),
				ref: DiskTemplateRef{Permalink: "templates/darwin-11"},
			},
			errStr:     fixtureDiskTemplateNotFoundErr,
			errResp:    fixtureDiskTemplateNotFoundResponseError,
			respStatus: http.StatusNotFound,
			respBody:   fixture("disk_template_not_found_error"),
		},
		{
			name: "nil context",
			args: args{
				ctx: nil,
				ref: DiskTemplateRef{ID: "dtpl_ytP13XD5DE1RdSL9"},
			},
			errStr: "net/http: nil Context",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rm, mux, _, teardown := prepareTestClient(t)
			defer teardown()
			c := NewDiskTemplatesClient(rm)

			mux.HandleFunc(
				"/core/v1/disk_templates/_",
				func(w http.ResponseWriter, r *http.Request) {
					assert.Equal(t, "GET", r.Method)
					assertEmptyFieldSpec(t, r)
					assertAuthorization(t, r)

					if tt.wantQuery != nil {
						assert.Equal(t, *tt.wantQuery, r.URL.Query())
					}

					w.WriteHeader(tt.respStatus)
					_, _ = w.Write(tt.respBody)
				},
			)

			got, resp, err := c.Get(
				tt.args.ctx, tt.args.ref,
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

func TestDiskTemplatesClient_GetByID(t *testing.T) {
	type args struct {
		ctx context.Context
		id  string
	}
	tests := []struct {
		name       string
		args       args
		want       *DiskTemplate
		errStr     string
		errResp    *katapult.ResponseError
		respStatus int
		respBody   []byte
	}{
		{
			name: "disk template",
			args: args{
				ctx: context.Background(),
				id:  "dtpl_ytP13XD5DE1RdSL9",
			},
			want: &DiskTemplate{
				ID:        "dtpl_ytP13XD5DE1RdSL9",
				Name:      "Ubuntu 18.04 Server",
				Permalink: "templates/ubuntu-18-04",
			},
			respStatus: http.StatusOK,
			respBody:   fixture("disk_template_get"),
		},
		{
			name: "non-existent disk template",
			args: args{
				ctx: context.Background(),
				id:  "dtpl_nopethisbegone",
			},
			errStr:     fixtureDiskTemplateNotFoundErr,
			errResp:    fixtureDiskTemplateNotFoundResponseError,
			respStatus: http.StatusNotFound,
			respBody:   fixture("disk_template_not_found_error"),
		},
		{
			name: "nil context",
			args: args{
				ctx: nil,
				id:  "dtpl_ytP13XD5DE1RdSL9",
			},
			errStr: "net/http: nil Context",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rm, mux, _, teardown := prepareTestClient(t)
			defer teardown()
			c := NewDiskTemplatesClient(rm)

			mux.HandleFunc("/core/v1/disk_templates/_",
				func(w http.ResponseWriter, r *http.Request) {
					assert.Equal(t, "GET", r.Method)
					assertEmptyFieldSpec(t, r)
					assertAuthorization(t, r)

					qs := url.Values{
						"disk_template[id]": []string{tt.args.id},
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

func TestDiskTemplatesClient_GetByPermalink(t *testing.T) {
	type args struct {
		ctx       context.Context
		permalink string
	}
	tests := []struct {
		name       string
		args       args
		want       *DiskTemplate
		errStr     string
		errResp    *katapult.ResponseError
		respStatus int
		respBody   []byte
	}{
		{
			name: "disk template",
			args: args{
				ctx:       context.Background(),
				permalink: "public",
			},
			want: &DiskTemplate{
				ID:        "dtpl_ytP13XD5DE1RdSL9",
				Name:      "Ubuntu 18.04 Server",
				Permalink: "templates/ubuntu-18-04",
			},
			respStatus: http.StatusOK,
			respBody:   fixture("disk_template_get"),
		},
		{
			name: "non-existent disk template",
			args: args{
				ctx:       context.Background(),
				permalink: "not-here",
			},
			errStr:     fixtureDiskTemplateNotFoundErr,
			errResp:    fixtureDiskTemplateNotFoundResponseError,
			respStatus: http.StatusNotFound,
			respBody:   fixture("disk_template_not_found_error"),
		},
		{
			name: "nil context",
			args: args{
				ctx:       nil,
				permalink: "public",
			},
			errStr: "net/http: nil Context",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rm, mux, _, teardown := prepareTestClient(t)
			defer teardown()
			c := NewDiskTemplatesClient(rm)

			mux.HandleFunc("/core/v1/disk_templates/_",
				func(w http.ResponseWriter, r *http.Request) {
					assert.Equal(t, "GET", r.Method)
					assertEmptyFieldSpec(t, r)
					assertAuthorization(t, r)

					qs := url.Values{
						"disk_template[permalink]": []string{tt.args.permalink},
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
