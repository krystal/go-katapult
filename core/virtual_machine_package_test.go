package core

import (
	"context"
	"encoding/json"
	"net/http"
	"net/url"
	"testing"

	"github.com/krystal/go-katapult"
	"github.com/krystal/go-katapult/internal/testclient"
	"github.com/stretchr/testify/assert"
)

var (
	fixturePackageNotFoundErr = "katapult: not_found: package_not_found: No " +
		"package was found matching any of the criteria provided in the " +
		"arguments"
	fixturePackageNotFoundResponseError = &katapult.ResponseError{
		Code: "package_not_found",
		Description: "No package was found matching any of the criteria " +
			"provided in the arguments",
		Detail: json.RawMessage(`{}`),
	}
)

func TestClient_VirtualMachinePackages(t *testing.T) {
	c := New(&testclient.Client{})

	assert.IsType(t, &VirtualMachinePackagesClient{}, c.VirtualMachinePackages)
}

func TestVirtualMachinePackage_JSONMarshaling(t *testing.T) {
	tests := []struct {
		name string
		obj  *VirtualMachinePackage
	}{
		{
			name: "empty",
			obj:  &VirtualMachinePackage{},
		},
		{
			name: "full",
			obj: &VirtualMachinePackage{
				ID:            "vmpkg_XdNPhGXvyt1dnDts",
				Name:          "X-Small",
				Permalink:     "xsmall",
				CPUCores:      504684,
				IPv4Addresses: 322134,
				MemoryInGB:    953603,
				StorageInGB:   853121,
				Privacy:       "priv",
				Icon:          &Attachment{URL: "url"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testJSONMarshaling(t, tt.obj)
		})
	}
}

func TestVirtualMachinePackage_Ref(t *testing.T) {
	tests := []struct {
		name string
		obj  *VirtualMachinePackage
		want VirtualMachinePackageRef
	}{
		{
			name: "empty",
			obj:  &VirtualMachinePackage{},
			want: VirtualMachinePackageRef{},
		},
		{
			name: "full",
			obj: &VirtualMachinePackage{
				ID:            "vmpkg_XdNPhGXvyt1dnDts",
				Name:          "X-Small",
				Permalink:     "xsmall",
				CPUCores:      504684,
				IPv4Addresses: 322134,
				MemoryInGB:    953603,
				StorageInGB:   853121,
				Privacy:       "priv",
				Icon:          &Attachment{URL: "url"},
			},
			want: VirtualMachinePackageRef{ID: "vmpkg_XdNPhGXvyt1dnDts"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.obj.Ref()

			assert.Equal(t, tt.want, got)
		})
	}
}

func TestVirtualMachinePackageRef_queryValues(t *testing.T) {
	tests := []struct {
		name string
		obj  VirtualMachinePackageRef
		want *url.Values
	}{
		{
			name: "empty",
			obj:  VirtualMachinePackageRef{},
			want: &url.Values{},
		},
		{
			name: "id and permalink",
			obj: VirtualMachinePackageRef{
				ID:        "vmpkg_XdNPhGXvyt1dnDts",
				Permalink: "xl-rock",
			},
		},
		{
			name: "id",
			obj: VirtualMachinePackageRef{
				ID: "vmpkg_XdNPhGXvyt1dnDts",
			},
		},
		{
			name: "permalink",
			obj: VirtualMachinePackageRef{
				Permalink: "xl-rock",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testQueryableEncoding(t, tt.obj)
		})
	}
}

func Test_virtualMachinePackagesResponseBody_JSONMarshaling(t *testing.T) {
	tests := []struct {
		name string
		obj  *virtualMachinePackagesResponseBody
	}{
		{
			name: "empty",
			obj:  &virtualMachinePackagesResponseBody{},
		},
		{
			name: "full",
			obj: &virtualMachinePackagesResponseBody{
				Pagination:             &katapult.Pagination{CurrentPage: 392},
				VirtualMachinePackage:  &VirtualMachinePackage{ID: "id1"},
				VirtualMachinePackages: []*VirtualMachinePackage{{ID: "id2"}},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testJSONMarshaling(t, tt.obj)
		})
	}
}

func TestVirtualMachinePackagesClient_List(t *testing.T) {
	// Correlates to fixtures/virtual_machine_packages_list*.json
	packageList := []*VirtualMachinePackage{
		{
			ID:        "vmpkg_XdNPhGXvyt1dnDts",
			Name:      "X-Small",
			Permalink: "xsmall",
		},
		{
			ID:        "vmpkg_YlqvfsKqZJODtvjG",
			Name:      "Small",
			Permalink: "small",
		},
		{
			ID:        "vmpkg_y7NqMMa9TYx0g1Si",
			Name:      "Medium",
			Permalink: "medium",
		},
	}

	type args struct {
		ctx  context.Context
		opts *ListOptions
	}
	tests := []struct {
		name           string
		args           args
		want           []*VirtualMachinePackage
		wantQuery      *url.Values
		wantPagination *katapult.Pagination
		errStr         string
		errResp        *katapult.ResponseError
		errIs          error
		respStatus     int
		respBody       []byte
	}{
		{
			name: "without pagination details",
			args: args{
				ctx: context.Background(),
			},
			want:      packageList,
			wantQuery: &url.Values{},
			wantPagination: &katapult.Pagination{
				CurrentPage: 1,
				TotalPages:  1,
				Total:       3,
				PerPage:     30,
				LargeSet:    false,
			},
			respStatus: http.StatusOK,
			respBody:   fixture("virtual_machine_packages_list"),
		},
		{
			name: "page 1",
			args: args{
				ctx:  context.Background(),
				opts: &ListOptions{Page: 1, PerPage: 2},
			},
			want: packageList[0:2],
			wantQuery: &url.Values{
				"page":     []string{"1"},
				"per_page": []string{"2"},
			},
			wantPagination: &katapult.Pagination{
				CurrentPage: 1,
				TotalPages:  2,
				Total:       3,
				PerPage:     2,
				LargeSet:    false,
			},
			respStatus: http.StatusOK,
			respBody:   fixture("virtual_machine_packages_list_page_1"),
		},
		{
			name: "page 2",
			args: args{
				ctx:  context.Background(),
				opts: &ListOptions{Page: 2, PerPage: 2},
			},
			want: packageList[2:],
			wantQuery: &url.Values{
				"page":     []string{"2"},
				"per_page": []string{"2"},
			},
			wantPagination: &katapult.Pagination{
				CurrentPage: 2,
				TotalPages:  2,
				Total:       3,
				PerPage:     2,
				LargeSet:    false,
			},
			respStatus: http.StatusOK,
			respBody:   fixture("virtual_machine_packages_list_page_2"),
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
			c := NewVirtualMachinePackagesClient(rm)

			mux.HandleFunc("/core/v1/virtual_machine_packages",
				func(w http.ResponseWriter, r *http.Request) {
					assert.Equal(t, "GET", r.Method)
					assertEmptyFieldSpec(t, r)
					assertAuthorization(t, r)
					assertRequestOptionHeader(t, r)

					if tt.wantQuery != nil {
						assert.Equal(t, *tt.wantQuery, r.URL.Query())
					} else {
						assert.Equal(t,
							*tt.args.opts.queryValues(), r.URL.Query(),
						)
					}

					w.WriteHeader(tt.respStatus)
					_, _ = w.Write(tt.respBody)
				},
			)

			got, resp, err := c.List(
				tt.args.ctx, tt.args.opts, testRequestOption,
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

func TestVirtualMachinePackagesClient_Get(t *testing.T) {
	type args struct {
		ctx context.Context
		ref VirtualMachinePackageRef
	}
	tests := []struct {
		name       string
		args       args
		reqPath    string
		want       *VirtualMachinePackage
		errStr     string
		errResp    *katapult.ResponseError
		errIs      error
		respStatus int
		respBody   []byte
	}{
		{
			name: "by ID",
			args: args{
				ctx: context.Background(),
				ref: VirtualMachinePackageRef{ID: "vmpkg_YlqvfsKqZJODtvjG"},
			},
			want: &VirtualMachinePackage{
				ID:        "vmpkg_YlqvfsKqZJODtvjG",
				Name:      "Small",
				Permalink: "small",
			},
			respStatus: http.StatusOK,
			respBody:   fixture("virtual_machine_package_get"),
		},
		{
			name: "by Permalink",
			args: args{
				ctx: context.Background(),
				ref: VirtualMachinePackageRef{Permalink: "small"},
			},
			want: &VirtualMachinePackage{
				ID:        "vmpkg_YlqvfsKqZJODtvjG",
				Name:      "Small",
				Permalink: "small",
			},
			respStatus: http.StatusOK,
			respBody:   fixture("virtual_machine_package_get"),
		},
		{
			name: "non-existent virtual machine package",
			args: args{
				ctx: context.Background(),
				ref: VirtualMachinePackageRef{ID: "vmpkg_nopethisbegone"},
			},
			errStr:     fixturePackageNotFoundErr,
			errResp:    fixturePackageNotFoundResponseError,
			errIs:      ErrVirtualMachinePackageNotFound,
			respStatus: http.StatusNotFound,
			respBody:   fixture("package_not_found_error"),
		},
		{
			name: "nil context",
			args: args{
				ctx: nil,
				ref: VirtualMachinePackageRef{ID: "vmpkg_nopethisbegone"},
			},
			errStr: "net/http: nil Context",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rm, mux, _, teardown := prepareTestClient(t)
			defer teardown()
			c := NewVirtualMachinePackagesClient(rm)

			mux.HandleFunc(
				"/core/v1/virtual_machine_packages/_",
				func(w http.ResponseWriter, r *http.Request) {
					assert.Equal(t, "GET", r.Method)
					assertEmptyFieldSpec(t, r)
					assertAuthorization(t, r)
					assertRequestOptionHeader(t, r)

					assert.Equal(t, *tt.args.ref.queryValues(), r.URL.Query())

					w.WriteHeader(tt.respStatus)
					_, _ = w.Write(tt.respBody)
				},
			)

			got, resp, err := c.Get(
				tt.args.ctx, tt.args.ref, testRequestOption,
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

			if tt.errIs != nil {
				assert.ErrorIs(t, err, tt.errIs)
			}
		})
	}
}

func TestVirtualMachinePackagesClient_GetByID(t *testing.T) {
	type args struct {
		ctx context.Context
		id  string
	}
	tests := []struct {
		name       string
		args       args
		want       *VirtualMachinePackage
		errStr     string
		errResp    *katapult.ResponseError
		errIs      error
		respStatus int
		respBody   []byte
	}{
		{
			name: "virtual machine package",
			args: args{
				ctx: context.Background(),
				id:  "vmpkg_YlqvfsKqZJODtvjG",
			},
			want: &VirtualMachinePackage{
				ID:        "vmpkg_YlqvfsKqZJODtvjG",
				Name:      "Small",
				Permalink: "small",
			},
			respStatus: http.StatusOK,
			respBody:   fixture("virtual_machine_package_get"),
		},
		{
			name: "non-existent virtual machine package",
			args: args{
				ctx: context.Background(),
				id:  "vmpkg_nopethisbegone",
			},
			errStr:     fixturePackageNotFoundErr,
			errResp:    fixturePackageNotFoundResponseError,
			errIs:      ErrVirtualMachinePackageNotFound,
			respStatus: http.StatusNotFound,
			respBody:   fixture("package_not_found_error"),
		},
		{
			name: "nil context",
			args: args{
				ctx: nil,
				id:  "vmpkg_nopethisbegone",
			},
			errStr: "net/http: nil Context",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rm, mux, _, teardown := prepareTestClient(t)
			defer teardown()
			c := NewVirtualMachinePackagesClient(rm)

			mux.HandleFunc(
				"/core/v1/virtual_machine_packages/_",
				func(w http.ResponseWriter, r *http.Request) {
					assert.Equal(t, "GET", r.Method)
					assertEmptyFieldSpec(t, r)
					assertAuthorization(t, r)
					assertRequestOptionHeader(t, r)

					assert.Equal(t, url.Values{
						"virtual_machine_package[id]": []string{
							tt.args.id,
						},
					}, r.URL.Query())

					w.WriteHeader(tt.respStatus)
					_, _ = w.Write(tt.respBody)
				},
			)

			got, resp, err := c.GetByID(
				tt.args.ctx, tt.args.id, testRequestOption,
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

			if tt.errIs != nil {
				assert.ErrorIs(t, err, tt.errIs)
			}
		})
	}
}

func TestVirtualMachinePackagesClient_GetByPermalink(t *testing.T) {
	type args struct {
		ctx       context.Context
		permalink string
	}
	tests := []struct {
		name       string
		args       args
		want       *VirtualMachinePackage
		errStr     string
		errResp    *katapult.ResponseError
		errIs      error
		respStatus int
		respBody   []byte
	}{
		{
			name: "virtual machine package",
			args: args{
				ctx:       context.Background(),
				permalink: "small",
			},
			want: &VirtualMachinePackage{
				ID:        "vmpkg_YlqvfsKqZJODtvjG",
				Name:      "Small",
				Permalink: "small",
			},
			respStatus: http.StatusOK,
			respBody:   fixture("virtual_machine_package_get"),
		},
		{
			name: "non-existent virtual machine package",
			args: args{
				ctx:       context.Background(),
				permalink: "nope-not-here",
			},
			errStr:     fixturePackageNotFoundErr,
			errResp:    fixturePackageNotFoundResponseError,
			errIs:      ErrVirtualMachinePackageNotFound,
			respStatus: http.StatusNotFound,
			respBody:   fixture("package_not_found_error"),
		},
		{
			name: "nil context",
			args: args{
				ctx:       nil,
				permalink: "small",
			},
			errStr: "net/http: nil Context",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rm, mux, _, teardown := prepareTestClient(t)
			defer teardown()
			c := NewVirtualMachinePackagesClient(rm)

			mux.HandleFunc(
				"/core/v1/virtual_machine_packages/_",
				func(w http.ResponseWriter, r *http.Request) {
					assert.Equal(t, "GET", r.Method)
					assertEmptyFieldSpec(t, r)
					assertAuthorization(t, r)
					assertRequestOptionHeader(t, r)

					qs := url.Values{
						"virtual_machine_package[permalink]": []string{
							tt.args.permalink,
						},
					}
					assert.Equal(t, qs, r.URL.Query())

					w.WriteHeader(tt.respStatus)
					_, _ = w.Write(tt.respBody)
				},
			)

			got, resp, err := c.GetByPermalink(
				tt.args.ctx, tt.args.permalink, testRequestOption,
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

			if tt.errIs != nil {
				assert.ErrorIs(t, err, tt.errIs)
			}
		})
	}
}
