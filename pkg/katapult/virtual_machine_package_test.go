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
	fixturePackageNotFoundErr = "package_not_found: No package was found " +
		"matching any of the criteria provided in the arguments"
	fixturePackageNotFoundResponseError = &ResponseError{
		Code: "package_not_found",
		Description: "No package was found matching any of the criteria " +
			"provided in the arguments",
		Detail: json.RawMessage(`{}`),
	}
)

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

func TestVirtualMachinePackage_LookupReference(t *testing.T) {
	tests := []struct {
		name string
		obj  *VirtualMachinePackage
		want *VirtualMachinePackage
	}{
		{
			name: "nil",
			obj:  (*VirtualMachinePackage)(nil),
			want: nil,
		},
		{
			name: "empty",
			obj:  &VirtualMachinePackage{},
			want: &VirtualMachinePackage{},
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
			want: &VirtualMachinePackage{ID: "vmpkg_XdNPhGXvyt1dnDts"},
		},
		{
			name: "no ID",
			obj: &VirtualMachinePackage{
				Name:          "X-Small",
				Permalink:     "xsmall",
				CPUCores:      504684,
				IPv4Addresses: 322134,
				MemoryInGB:    953603,
				StorageInGB:   853121,
				Privacy:       "priv",
				Icon:          &Attachment{URL: "url"},
			},
			want: &VirtualMachinePackage{Permalink: "xsmall"},
		},
		{
			name: "no ID or Permalink",
			obj: &VirtualMachinePackage{
				Name:          "X-Small",
				CPUCores:      504684,
				IPv4Addresses: 322134,
				MemoryInGB:    953603,
				StorageInGB:   853121,
				Privacy:       "priv",
				Icon:          &Attachment{URL: "url"},
			},
			want: &VirtualMachinePackage{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.obj.LookupReference()

			assert.Equal(t, tt.want, got)
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
				Pagination:             &Pagination{CurrentPage: 392},
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
		name       string
		args       args
		expected   []*VirtualMachinePackage
		pagination *Pagination
		errStr     string
		errResp    *ResponseError
		respStatus int
		respBody   []byte
	}{
		{
			name: "virtual machine packages",
			args: args{
				ctx: context.Background(),
			},
			expected: packageList,
			pagination: &Pagination{
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
			name: "page 1 of virtual machine packages",
			args: args{
				ctx:  context.Background(),
				opts: &ListOptions{Page: 1, PerPage: 2},
			},
			expected: packageList[0:2],
			pagination: &Pagination{
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
			name: "page 2 of virtual machine packages",
			args: args{
				ctx:  context.Background(),
				opts: &ListOptions{Page: 2, PerPage: 2},
			},
			expected: packageList[2:],
			pagination: &Pagination{
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
			c, mux, _, teardown := prepareTestClient()
			defer teardown()

			mux.HandleFunc("/core/v1/virtual_machine_packages",
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

			got, resp, err := c.VirtualMachinePackages.List(
				tt.args.ctx, tt.args.opts,
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

func TestVirtualMachinePackagesClient_Get(t *testing.T) {
	type args struct {
		ctx context.Context
		id  string
	}
	tests := []struct {
		name       string
		args       args
		expected   *VirtualMachinePackage
		errStr     string
		errResp    *ResponseError
		respStatus int
		respBody   []byte
	}{
		{
			name: "virtual machine package",
			args: args{
				ctx: context.Background(),
				id:  "vmpkg_YlqvfsKqZJODtvjG",
			},
			expected: &VirtualMachinePackage{
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
			c, mux, _, teardown := prepareTestClient()
			defer teardown()

			mux.HandleFunc(
				fmt.Sprintf("/core/v1/virtual_machine_packages/%s", tt.args.id),
				func(w http.ResponseWriter, r *http.Request) {
					assert.Equal(t, "GET", r.Method)
					assertEmptyFieldSpec(t, r)
					assertAuthorization(t, r)

					w.WriteHeader(tt.respStatus)
					_, _ = w.Write(tt.respBody)
				},
			)

			got, resp, err := c.VirtualMachinePackages.Get(
				tt.args.ctx, tt.args.id,
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

			if tt.errResp != nil {
				assert.Equal(t, tt.errResp, resp.Error)
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
		expected   *VirtualMachinePackage
		errStr     string
		errResp    *ResponseError
		respStatus int
		respBody   []byte
	}{
		{
			name: "virtual machine package",
			args: args{
				ctx:       context.Background(),
				permalink: "small",
			},
			expected: &VirtualMachinePackage{
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
			errStr: "package_not_found: No package was found matching " +
				"any of the criteria provided in the arguments",
			errResp: &ResponseError{
				Code: "package_not_found",
				Description: "No package was found matching any of the " +
					"criteria provided in the arguments",
				Detail: json.RawMessage(`{}`),
			},
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
			c, mux, _, teardown := prepareTestClient()
			defer teardown()

			mux.HandleFunc(
				"/core/v1/virtual_machine_packages/_",
				func(w http.ResponseWriter, r *http.Request) {
					assert.Equal(t, "GET", r.Method)
					assertEmptyFieldSpec(t, r)
					assertAuthorization(t, r)

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

			got, resp, err := c.VirtualMachinePackages.GetByPermalink(
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

			if tt.expected != nil {
				assert.Equal(t, tt.expected, got)
			}

			if tt.errResp != nil {
				assert.Equal(t, tt.errResp, resp.Error)
			}
		})
	}
}
