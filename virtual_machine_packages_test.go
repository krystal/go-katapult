package katapult

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/google/go-querystring/query"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestVirtualMachinePackagesService_List(t *testing.T) {
	// Correlates to fixtures/virtual_machine_packages_list*.json
	packageList := []*VirtualMachinePackage{
		{
			ID:            "vmpkg_XdNPhGXvyt1dnDts",
			Name:          "X-Small",
			Permalink:     "xsmall",
			CPUCores:      1,
			IPv4Addresses: 1,
			MemoryInGB:    1,
			StorageInGB:   10,
			Privacy:       "private",
			Icon: &Icon{
				URL: "https://my.katapult.io/attachment/" +
					"aa9e51fc-ca56-4a4a-aeba-2f57ffcc9886/cat.jpg",
			},
		},
		{
			ID:            "vmpkg_YlqvfsKqZJODtvjG",
			Name:          "Small",
			Permalink:     "small",
			CPUCores:      2,
			IPv4Addresses: 1,
			MemoryInGB:    2,
			StorageInGB:   10,
			Privacy:       "public",
			Icon: &Icon{
				URL: "https://my.katapult.io/attachment/" +
					"4d014ee8-dae3-4574-a180-e5711fc85f9a/fox.png",
			},
		},
		{
			ID:            "vmpkg_y7NqMMa9TYx0g1Si",
			Name:          "Medium",
			Permalink:     "medium",
			CPUCores:      4,
			IPv4Addresses: 1,
			MemoryInGB:    3,
			StorageInGB:   20,
			Privacy:       "public",
			Icon: &Icon{
				URL: "https://my.katapult.io/attachment/" +
					"23eabfd1-f8a9-4312-80c1-37bc3e563754/lion.png",
			},
		},
	}

	tests := []struct {
		name       string
		opts       *ListOptions
		expected   []*VirtualMachinePackage
		pagination *Pagination
		err        string
		errResp    *ResponseError
		respStatus int
		respBody   []byte
	}{
		{
			name:     "fetch list of virtual machine packages",
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
			name:     "fetch list of virtual machine packages",
			opts:     &ListOptions{Page: 1, PerPage: 2},
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
			name:     "fetch page 2 of virtual machine packages list",
			opts:     &ListOptions{Page: 2, PerPage: 2},
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
			name:       "invalid API token response",
			err:        fixtureInvalidAPITokenErr,
			errResp:    fixtureInvalidAPITokenResponseError,
			respStatus: http.StatusForbidden,
			respBody:   fixture("invalid_api_token_error"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, mux, _, teardown := setup()
			defer teardown()

			mux.HandleFunc("/core/v1/virtual_machine_packages",
				func(w http.ResponseWriter, r *http.Request) {
					assert.Equal(t, "GET", r.Method)
					assert.Equal(t, "", r.Header.Get("X-Field-Spec"))
					if tt.opts != nil {
						qs, err := query.Values(tt.opts)
						require.NoError(t, err)
						assert.Equal(t, qs, r.URL.Query())
					}

					w.WriteHeader(tt.respStatus)
					_, _ = w.Write(tt.respBody)
				},
			)

			got, resp, err := c.VirtualMachinePackages.List(
				context.Background(), tt.opts,
			)

			assert.Equal(t, tt.respStatus, resp.StatusCode)

			if tt.err == "" {
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, err, tt.err)
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

func TestVirtualMachinePackagesService_Get(t *testing.T) {
	// Correlates to fixtures/virtual_machine_package_get.json
	packageGet := &VirtualMachinePackage{
		ID:            "vmpkg_YlqvfsKqZJODtvjG",
		Name:          "Small",
		Permalink:     "small",
		CPUCores:      2,
		IPv4Addresses: 1,
		MemoryInGB:    2,
		StorageInGB:   10,
		Privacy:       "public",
		Icon: &Icon{
			URL: "https://my.katapult.io/attachment/" +
				"4d014ee8-dae3-4574-a180-e5711fc85f9a/fox.png",
			FileName: "fox.png",
			FileType: "image/png",
			FileSize: 4868,
			Digest:   "0f19d773-1166-441b-b146-f25713d20188",
			Token:    "8da34c2a-f444-44b3-b2e5-290daa055a92",
		},
	}

	tests := []struct {
		name       string
		id         string
		expected   *VirtualMachinePackage
		err        string
		errResp    *ResponseError
		respStatus int
		respBody   []byte
	}{
		{
			name:       "get a virtual machine package",
			id:         "vmpkg_YlqvfsKqZJODtvjG",
			expected:   packageGet,
			respStatus: http.StatusOK,
			respBody:   fixture("virtual_machine_package_get"),
		},
		{
			name: "non-existent virtual machine package package",
			id:   "vmpkg_nopethisbegone",
			err: "package_not_found: No package was found matching " +
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
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, mux, _, teardown := setup()
			defer teardown()

			mux.HandleFunc(
				fmt.Sprintf("/core/v1/virtual_machine_packages/%s", tt.id),
				func(w http.ResponseWriter, r *http.Request) {
					assert.Equal(t, "GET", r.Method)
					assert.Equal(t, "", r.Header.Get("X-Field-Spec"))

					w.WriteHeader(tt.respStatus)
					_, _ = w.Write(tt.respBody)
				},
			)

			got, resp, err := c.VirtualMachinePackages.Get(
				context.Background(), tt.id,
			)

			assert.Equal(t, tt.respStatus, resp.StatusCode)

			if tt.err == "" {
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, err, tt.err)
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
