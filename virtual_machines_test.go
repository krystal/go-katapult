package katapult

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestVirtualMachine_JSONMarshaling(t *testing.T) {
	tests := []struct {
		name string
		obj  *VirtualMachine
	}{
		{
			name: "empty",
			obj:  &VirtualMachine{},
		},
		{
			name: "full",
			obj: &VirtualMachine{
				ID:                  "id",
				Name:                "name",
				Hostname:            "hostname",
				FQDN:                "Fqdn",
				CreatedAt:           timestampPtr(934834834),
				InitialRootPassword: "initial_root_password",
				State:               "state",
				Zone:                &Zone{ID: "id0"},
				Organization:        &Organization{ID: "id1"},
				Group:               &VirtualMachineGroup{ID: "id2"},
				Package:             &VirtualMachinePackage{ID: "id3"},
				AttachedISO:         &ISO{ID: "id4"},
				Tags:                []*Tag{{ID: "id5"}},
				IPAddresses:         []*IPAddress{{ID: "id6"}},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testJSONMarshaling(t, tt.obj)
		})
	}
}

func TestVirtualMachineGroup_JSONMarshaling(t *testing.T) {
	tests := []struct {
		name string
		obj  *VirtualMachineGroup
	}{
		{
			name: "empty",
			obj:  &VirtualMachineGroup{},
		},
		{
			name: "full",
			obj: &VirtualMachineGroup{
				ID:        "id",
				Name:      "name",
				Segregate: true,
				CreatedAt: timestampPtr(934834834),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testJSONMarshaling(t, tt.obj)
		})
	}
}

func TestISO_JSONMarshaling(t *testing.T) {
	tests := []struct {
		name string
		obj  *ISO
	}{
		{
			name: "empty",
			obj:  &ISO{},
		},
		{
			name: "full",
			obj: &ISO{
				ID:              "id1",
				Name:            "name",
				OperatingSystem: &OperatingSystem{ID: "id2"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testJSONMarshaling(t, tt.obj)
		})
	}
}

func TestOperatingSystem_JSONMarshaling(t *testing.T) {
	tests := []struct {
		name string
		obj  *OperatingSystem
	}{
		{
			name: "empty",
			obj:  &OperatingSystem{},
		},
		{
			name: "full",
			obj: &OperatingSystem{
				ID:    "id1",
				Name:  "name",
				Badge: &Attachment{URL: "url2"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testJSONMarshaling(t, tt.obj)
		})
	}
}

func TestTag_JSONMarshaling(t *testing.T) {
	tests := []struct {
		name string
		obj  *Tag
	}{
		{
			name: "empty",
			obj:  &Tag{},
		},
		{
			name: "full",
			obj: &Tag{
				ID:        "id1",
				Name:      "name",
				Color:     "color",
				CreatedAt: timestampPtr(3043009),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testJSONMarshaling(t, tt.obj)
		})
	}
}

func TestIPAddress_JSONMarshaling(t *testing.T) {
	tests := []struct {
		name string
		obj  *IPAddress
	}{
		{
			name: "empty",
			obj:  &IPAddress{},
		},
		{
			name: "full",
			obj: &IPAddress{
				ID:              "id1",
				Address:         "address",
				ReverseDNS:      "reverse_dns",
				VIP:             true,
				AddressWithMask: "address_with_mask",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testJSONMarshaling(t, tt.obj)
		})
	}
}

func Test_virtualMachinesResponseBody_JSONMarshaling(t *testing.T) {
	tests := []struct {
		name string
		obj  *virtualMachinesResponseBody
	}{
		{
			name: "empty",
			obj:  &virtualMachinesResponseBody{},
		},
		{
			name: "full",
			obj: &virtualMachinesResponseBody{
				Pagination:      &Pagination{CurrentPage: 345},
				VirtualMachine:  &VirtualMachine{ID: "id1"},
				VirtualMachines: []*VirtualMachine{{ID: "id2"}},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testJSONMarshaling(t, tt.obj)
		})
	}
}

func TestVirtualMachinesService_List(t *testing.T) {
	// Correlates to fixtures/virtual_machines_list*.json
	zone := &Zone{
		ID:        "zone_FZ7ZpVppByTepLqf",
		Name:      "Main Zone",
		Permalink: "main",
		DataCenter: &DataCenter{
			ID:        "loc_25d48761871e4bf",
			Name:      "Shirebury",
			Permalink: "shirebury",
		},
	}
	virtualMachinesList := []*VirtualMachine{
		{
			ID:          "vm_t8yomYsG4bccKw5D",
			Name:        "bitter-beautiful-mango",
			Hostname:    "bitter-beautiful-mango",
			FQDN:        "bitter-beautiful-mango.test.kpult.com",
			CreatedAt:   timestampPtr(1579883611),
			Zone:        zone,
			Package:     &VirtualMachinePackage{Name: "Small"},
			IPAddresses: []*IPAddress{{Address: "31.228.217.44"}},
		},
		{
			ID:          "vm_h7bzdXXHa0GvJYMc",
			Name:        "popular-shapely-tank",
			Hostname:    "popular-shapely-tank",
			FQDN:        "popular-shapely-tank.test.kpult.com",
			CreatedAt:   timestampPtr(1568651782),
			Zone:        zone,
			Package:     &VirtualMachinePackage{Name: "Medium"},
			IPAddresses: []*IPAddress{{Address: "188.56.184.51"}},
		},
		{
			ID:          "vm_1kpkjQeMEI43tztr",
			Name:        "popular-blue-kumquat",
			Hostname:    "popular-blue-kumquat",
			FQDN:        "popular-blue-kumquat.test.kpult.com",
			CreatedAt:   timestampPtr(1563381453),
			Zone:        zone,
			Package:     &VirtualMachinePackage{Name: "Small"},
			IPAddresses: []*IPAddress{{Address: "106.127.29.51"}},
		},
	}

	type args struct {
		ctx   context.Context
		orgID string
		opts  *ListOptions
	}
	tests := []struct {
		name       string
		args       args
		expected   []*VirtualMachine
		pagination *Pagination
		errStr     string
		errResp    *ResponseError
		respStatus int
		respBody   []byte
	}{
		{
			name: "virtual machines",
			args: args{
				ctx:   context.Background(),
				orgID: "org_O648YDMEYeLmqdmn",
			},
			expected: virtualMachinesList,
			pagination: &Pagination{
				CurrentPage: 1,
				TotalPages:  1,
				Total:       3,
				PerPage:     30,
				LargeSet:    false,
			},
			respStatus: http.StatusOK,
			respBody:   fixture("virtual_machines_list"),
		},
		{
			name: "page 1 of virtual machines",
			args: args{
				ctx:   context.Background(),
				orgID: "org_O648YDMEYeLmqdmn",
				opts:  &ListOptions{Page: 1, PerPage: 2},
			},
			expected: virtualMachinesList[0:2],
			pagination: &Pagination{
				CurrentPage: 1,
				TotalPages:  2,
				Total:       3,
				PerPage:     2,
				LargeSet:    false,
			},
			respStatus: http.StatusOK,
			respBody:   fixture("virtual_machines_list_page_1"),
		},
		{
			name: "page 2 of virtual machines",
			args: args{
				ctx:   context.Background(),
				orgID: "org_O648YDMEYeLmqdmn",
				opts:  &ListOptions{Page: 2, PerPage: 2},
			},
			expected: virtualMachinesList[2:],
			pagination: &Pagination{
				CurrentPage: 2,
				TotalPages:  2,
				Total:       3,
				PerPage:     2,
				LargeSet:    false,
			},
			respStatus: http.StatusOK,
			respBody:   fixture("virtual_machines_list_page_2"),
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
					"/core/v1/organizations/%s/virtual_machines", tt.args.orgID,
				),
				func(w http.ResponseWriter, r *http.Request) {
					assert.Equal(t, "GET", r.Method)
					assert.Equal(t, "", r.Header.Get("X-Field-Spec"))
					if tt.args.opts != nil {
						assert.Equal(t, *tt.args.opts.Values(), r.URL.Query())
					}

					w.WriteHeader(tt.respStatus)
					_, _ = w.Write(tt.respBody)
				},
			)

			got, resp, err := c.VirtualMachines.List(
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

func TestVirtualMachinesService_Get(t *testing.T) {
	// Correlates to fixtures/virtual_machine_get.json
	zone := &Zone{
		ID:        "zone_FZ7ZpVppByTepLqf",
		Name:      "Main Zone",
		Permalink: "main",
		DataCenter: &DataCenter{
			ID:        "loc_25d48761871e4bf",
			Name:      "Shirebury",
			Permalink: "shirebury",
			Country: &Country{
				ID:       "ctry_2f2dc89a5956437",
				Name:     "United Kingdom",
				ISOCode2: "GB",
				ISOCode3: "GBR",
				TimeZone: "Europe/London",
				EU:       true,
			},
		},
	}
	virtualMachine := &VirtualMachine{
		ID:                  "vm_t8yomYsG4bccKw5D",
		Name:                "bitter-beautiful-mango",
		Hostname:            "bitter-beautiful-mango",
		FQDN:                "bitter-beautiful-mango.test.kpult.com",
		CreatedAt:           timestampPtr(1579883611),
		InitialRootPassword: "w7OrYGqImpwoNZOM",
		State:               "started",
		Zone:                zone,
		Organization: &Organization{
			ID:        "org_BVxE1c44NazNmoOb",
			Name:      "Sherlock & Co.",
			SubDomain: "sherlock",
		},
		Group: &VirtualMachineGroup{
			ID:        "vmgrp_CQZ3NlZq37aMKaDU",
			Name:      "misc",
			Segregate: false,
			CreatedAt: timestampPtr(1558199798),
		},
		Package: &VirtualMachinePackage{
			ID:   "vmpkg_YlqvfsKqZJODtvjG",
			Name: "Small",
		},
		AttachedISO: &ISO{
			ID:   "iso_AkqXHtTlQXwSwQJ9",
			Name: "Ubuntu 20.04",
			OperatingSystem: &OperatingSystem{
				ID:   "os_D1Z90eOlbCSqGCU4",
				Name: "Ubuntu",
				Badge: &Attachment{
					URL: "https://my.katapult.io/attachment/" +
						"87ef444a-43e0-440e-9205-4285e7cccc3b/ubuntu.png",
					FileName: "ubuntu.png",
					FileType: "image/png",
					FileSize: 4924,
					Digest:   "0964db5f-bf6f-430e-bd77-85078d16b85a",
					Token:    "f0b3faf5-313d-43f6-87e0-7f142c757466",
				},
			},
		},
		Tags: []*Tag{
			{
				ID:        "tag_dpkrkYssrXDryOhn",
				Name:      "web",
				Color:     "light_brown",
				CreatedAt: timestampPtr(1596993883),
			},
		},
		IPAddresses: []*IPAddress{
			{
				ID:              "ip_K4t3ya3Werh6zPCd",
				Address:         "210.103.10.205",
				ReverseDNS:      "210-103-10-205.infra.katapult.dev",
				VIP:             false,
				AddressWithMask: "210.103.10.205/26",
			},
		},
	}

	type args struct {
		ctx context.Context
		id  string
	}
	tests := []struct {
		name       string
		args       args
		expected   *VirtualMachine
		errStr     string
		errResp    *ResponseError
		respStatus int
		respBody   []byte
	}{
		{
			name: "virtual machine",
			args: args{
				ctx: context.Background(),
				id:  "vm_t8yomYsG4bccKw5D",
			},
			expected:   virtualMachine,
			respStatus: http.StatusOK,
			respBody:   fixture("virtual_machine_get"),
		},
		{
			name: "non-existent virtual machine",
			args: args{
				ctx: context.Background(),
				id:  "vm_nopethisbegone",
			},
			errStr: "virtual_machine_not_found: No virtual machine was found " +
				"matching any of the criteria provided in the arguments",
			errResp: &ResponseError{
				Code: "virtual_machine_not_found",
				Description: "No virtual machine was found matching any of " +
					"the criteria provided in the arguments",
				Detail: json.RawMessage(`{}`),
			},
			respStatus: http.StatusNotFound,
			respBody:   fixture("virtual_machine_not_found_error"),
		},
		{
			name: "virtual machine is in trash",
			args: args{
				ctx: context.Background(),
				id:  "vm_t8yomYsG4bccKw5D",
			},
			errStr:     fixtureObjectInTrashErr,
			errResp:    fixtureObjectInTrashResponseError,
			respStatus: http.StatusNotAcceptable,
			respBody:   fixture("object_in_trash_error"),
		},
		{
			name: "nil context",
			args: args{
				ctx: nil,
				id:  "vm_t8yomYsG4bccKw5D",
			},
			errStr: "net/http: nil Context",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, mux, _, teardown := prepareTestClient()
			defer teardown()

			mux.HandleFunc(
				fmt.Sprintf("/core/v1/virtual_machines/%s", tt.args.id),
				func(w http.ResponseWriter, r *http.Request) {
					assert.Equal(t, "GET", r.Method)
					assert.Equal(t, "", r.Header.Get("X-Field-Spec"))

					w.WriteHeader(tt.respStatus)
					_, _ = w.Write(tt.respBody)
				},
			)

			got, resp, err := c.VirtualMachines.Get(tt.args.ctx, tt.args.id)

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
