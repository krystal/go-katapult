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

//nolint:lll
var (
	fixtureVMNetworkInterfaceNotFoundErr = "virtual_machine_network_interface_not_found: " +
		"No network interface was found matching any of the criteria " +
		"provided in the arguments"
	fixtureVMNetworkInterfaceNotFoundResponseError = &ResponseError{
		Code: "virtual_machine_network_interface_not_found",
		Description: "No network interface was found matching any of the " +
			"criteria provided in the arguments",
		Detail: json.RawMessage(`{}`),
	}

	fixtureVirtualMachineNetworkInterfaceFull = &VirtualMachineNetworkInterface{
		ID:             "vmnet_Qlu34yEQgkrIlzql",
		VirtualMachine: &VirtualMachine{ID: "vm_i5qfOrvEI1CmNrJx"},
		Name:           "Public Network on foo-bar",
		Network:        &Network{ID: "net_HjwzDggBv9gsHZ1T"},
		MACAddress:     "ab:cd:ef:12:34:56",
		State:          "attached",
		IPAddresses:    []*IPAddress{{ID: "ip_7S6uoasz4jM5mkMs"}},
	}
	fixtureVirtualMachineNetworkInterfaceNoID = &VirtualMachineNetworkInterface{
		VirtualMachine: fixtureVirtualMachineNetworkInterfaceFull.VirtualMachine,
		Name:           fixtureVirtualMachineNetworkInterfaceFull.Name,
		Network:        fixtureVirtualMachineNetworkInterfaceFull.Network,
		MACAddress:     fixtureVirtualMachineNetworkInterfaceFull.MACAddress,
		State:          fixtureVirtualMachineNetworkInterfaceFull.State,
		IPAddresses:    fixtureVirtualMachineNetworkInterfaceFull.IPAddresses,
	}
)

func TestVirtualMachineNetworkInterface_JSONMarshaling(t *testing.T) {
	tests := []struct {
		name string
		obj  *VirtualMachineNetworkInterface
	}{
		{
			name: "empty",
			obj:  &VirtualMachineNetworkInterface{},
		},
		{
			name: "full",
			obj:  fixtureVirtualMachineNetworkInterfaceFull,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testJSONMarshaling(t, tt.obj)
		})
	}
}

func TestVirtualMachineNetworkInterface_lookupReference(t *testing.T) {
	tests := []struct {
		name string
		obj  *VirtualMachineNetworkInterface
		want *VirtualMachineNetworkInterface
	}{
		{
			name: "nil",
			obj:  nil,
			want: nil,
		},
		{
			name: "empty",
			obj:  &VirtualMachineNetworkInterface{},
			want: &VirtualMachineNetworkInterface{},
		},
		{
			name: "full",
			obj:  fixtureVirtualMachineNetworkInterfaceFull,
			want: &VirtualMachineNetworkInterface{ID: "vmnet_Qlu34yEQgkrIlzql"},
		},
		{
			name: "no ID",
			obj:  fixtureVirtualMachineNetworkInterfaceNoID,
			want: &VirtualMachineNetworkInterface{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.obj.lookupReference()

			assert.Equal(t, tt.want, got)
		})
	}
}

func TestVirtualMachineNetworkInterface_queryValues(t *testing.T) {
	tests := []struct {
		name string
		obj  *VirtualMachineNetworkInterface
	}{
		{
			name: "nil",
			obj:  nil,
		},
		{
			name: "empty",
			obj:  &VirtualMachineNetworkInterface{},
		},
		{
			name: "full",
			obj:  fixtureVirtualMachineNetworkInterfaceFull,
		},
		{
			name: "no ID",
			obj:  fixtureVirtualMachineNetworkInterfaceNoID,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testQueryableEncoding(t, tt.obj)
		})
	}
}

func Test_virtualMachineNetworkInterfacesResponseBody_JSONMarshaling(t *testing.T) { //nolint:lll
	tests := []struct {
		name string
		obj  *virtualMachineNetworkInterfacesResponseBody
	}{
		{
			name: "empty",
			obj:  &virtualMachineNetworkInterfacesResponseBody{},
		},
		{
			name: "full",
			obj: &virtualMachineNetworkInterfacesResponseBody{
				Pagination: &Pagination{CurrentPage: 345},
				VirtualMachineNetworkInterface: &VirtualMachineNetworkInterface{
					ID: "id1",
				},
				//nolint:lll
				VirtualMachineNetworkInterfaces: []*VirtualMachineNetworkInterface{
					{ID: "id2"},
				},
				IPAddress:   &IPAddress{ID: "id3"},
				IPAddresses: []*IPAddress{{ID: "id4"}},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testJSONMarshaling(t, tt.obj)
		})
	}
}

func TestVirtualMachineNetworkInterfacesClient_List(t *testing.T) {
	// Correlates to fixtures/virtual_machine_network_interfaces_list*.json
	virtualMachinesList := []*VirtualMachineNetworkInterface{
		{
			ID:   "vmnet_olNAz8ThH0emHvdr",
			Name: "Public Network on bitter-beautiful-mango",
		},
		{
			ID:   "vmnet_KxKhb8M7jpN8hTBL",
			Name: "Private Network on bitter-beautiful-mango",
		},
		{
			ID:   "vmnet_19JWZO4oHJ51J79y",
			Name: "Internal Network on bitter-beautiful-mango",
		},
	}

	type args struct {
		ctx  context.Context
		vm   *VirtualMachine
		opts *ListOptions
	}
	tests := []struct {
		name           string
		args           args
		want           []*VirtualMachineNetworkInterface
		wantQuery      *url.Values
		wantPagination *Pagination
		errStr         string
		errResp        *ResponseError
		respStatus     int
		respBody       []byte
	}{
		{
			name: "by virtual machine ID",
			args: args{
				ctx: context.Background(),
				vm:  &VirtualMachine{ID: "vm_i5qfOrvEI1CmNrJx"},
			},
			want: virtualMachinesList,
			wantQuery: &url.Values{
				"virtual_machine[id]": []string{"vm_i5qfOrvEI1CmNrJx"},
			},
			wantPagination: &Pagination{
				CurrentPage: 1,
				TotalPages:  1,
				Total:       3,
				PerPage:     30,
				LargeSet:    false,
			},
			respStatus: http.StatusOK,
			respBody:   fixture("virtual_machine_network_interfaces_list"),
		},
		{
			name: "by virtual machine FQDN",
			args: args{
				ctx: context.Background(),
				vm:  &VirtualMachine{FQDN: "acme"},
			},
			want: virtualMachinesList,
			wantQuery: &url.Values{
				"virtual_machine[fqdn]": []string{"acme"},
			},
			wantPagination: &Pagination{
				CurrentPage: 1,
				TotalPages:  1,
				Total:       3,
				PerPage:     30,
				LargeSet:    false,
			},
			respStatus: http.StatusOK,
			respBody:   fixture("virtual_machine_network_interfaces_list"),
		},
		{
			name: "page 1",
			args: args{
				ctx:  context.Background(),
				vm:   &VirtualMachine{ID: "vm_i5qfOrvEI1CmNrJx"},
				opts: &ListOptions{Page: 1, PerPage: 2},
			},
			want: virtualMachinesList[0:2],
			wantQuery: &url.Values{
				"virtual_machine[id]": []string{"vm_i5qfOrvEI1CmNrJx"},
				"page":                []string{"1"},
				"per_page":            []string{"2"},
			},
			wantPagination: &Pagination{
				CurrentPage: 1,
				TotalPages:  2,
				Total:       3,
				PerPage:     2,
				LargeSet:    false,
			},
			respStatus: http.StatusOK,
			respBody: fixture(
				"virtual_machine_network_interfaces_list_page_1",
			),
		},
		{
			name: "page 2",
			args: args{
				ctx:  context.Background(),
				vm:   &VirtualMachine{ID: "vm_i5qfOrvEI1CmNrJx"},
				opts: &ListOptions{Page: 2, PerPage: 2},
			},
			want: virtualMachinesList[2:],
			wantQuery: &url.Values{
				"virtual_machine[id]": []string{"vm_i5qfOrvEI1CmNrJx"},
				"page":                []string{"2"},
				"per_page":            []string{"2"},
			},
			wantPagination: &Pagination{
				CurrentPage: 2,
				TotalPages:  2,
				Total:       3,
				PerPage:     2,
				LargeSet:    false,
			},
			respStatus: http.StatusOK,
			respBody: fixture(
				"virtual_machine_network_interfaces_list_page_2",
			),
		},
		{
			name: "invalid API token response",
			args: args{
				ctx: context.Background(),
				vm:  &VirtualMachine{ID: "vm_i5qfOrvEI1CmNrJx"},
			},
			errStr:     fixtureInvalidAPITokenErr,
			errResp:    fixtureInvalidAPITokenResponseError,
			respStatus: http.StatusForbidden,
			respBody:   fixture("invalid_api_token_error"),
		},
		{
			name: "non-existent virtual machine",
			args: args{
				ctx: context.Background(),
				vm:  &VirtualMachine{ID: "vm_i5qfOrvEI1CmNrJx"},
			},
			errStr:     fixtureVirtualMachineNotFoundErr,
			errResp:    fixtureVirtualMachineNotFoundResponseError,
			respStatus: http.StatusNotFound,
			respBody:   fixture("virtual_machine_not_found_error"),
		},
		{
			name: "nil virtual_machine",
			args: args{
				ctx: context.Background(),
				vm:  nil,
			},
			errStr:     fixtureVirtualMachineNotFoundErr,
			errResp:    fixtureVirtualMachineNotFoundResponseError,
			respStatus: http.StatusNotFound,
			respBody:   fixture("virtual_machine_not_found_error"),
		},
		{
			name: "nil context",
			args: args{
				ctx: nil,
				vm:  &VirtualMachine{ID: "vm_i5qfOrvEI1CmNrJx"},
			},
			errStr: "net/http: nil Context",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, mux, _, teardown := prepareTestClient()
			defer teardown()

			mux.HandleFunc(
				"/core/v1/virtual_machines/_/network_interfaces",
				func(w http.ResponseWriter, r *http.Request) {
					assert.Equal(t, "GET", r.Method)
					assertEmptyFieldSpec(t, r)
					assertAuthorization(t, r)

					if tt.wantQuery != nil {
						assert.Equal(t, *tt.wantQuery, r.URL.Query())
					} else {
						qs := queryValues(tt.args.vm, tt.args.opts)
						assert.Equal(t, *qs, r.URL.Query())
					}

					w.WriteHeader(tt.respStatus)
					_, _ = w.Write(tt.respBody)
				},
			)

			got, resp, err := c.VirtualMachineNetworkInterfaces.List(
				tt.args.ctx, tt.args.vm, tt.args.opts,
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

func TestVirtualMachineNetworkInterfacesClient_Get(t *testing.T) {
	type args struct {
		ctx context.Context
		id  string
	}
	tests := []struct {
		name       string
		args       args
		want       *VirtualMachineNetworkInterface
		errStr     string
		errResp    *ResponseError
		respStatus int
		respBody   []byte
	}{
		{
			name: "by ID",
			args: args{
				ctx: context.Background(),
				id:  "vmnet_olNAz8ThH0emHvdr",
			},
			want: &VirtualMachineNetworkInterface{
				ID:      "vmnet_olNAz8ThH0emHvdr",
				Name:    "Public Network on bitter-beautiful-mango",
				Network: &Network{ID: "net_4s5J6gMQXhcwqIqs"},
			},
			respStatus: http.StatusOK,
			respBody:   fixture("virtual_machine_network_interface_get"),
		},
		{
			name: "non-existent virtual machine network interface",
			args: args{
				ctx: context.Background(),
				id:  "vmnet_nopethisbegone",
			},
			errStr:     fixtureVMNetworkInterfaceNotFoundErr,
			errResp:    fixtureVMNetworkInterfaceNotFoundResponseError,
			respStatus: http.StatusNotFound,
			respBody: fixture(
				"virtual_machine_network_interface_not_found_error",
			),
		},
		{
			name: "nil context",
			args: args{
				ctx: nil,
				id:  "vmnet_olNAz8ThH0emHvdr",
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
					"/core/v1/virtual_machine_network_interfaces/%s",
					tt.args.id,
				),
				func(w http.ResponseWriter, r *http.Request) {
					assert.Equal(t, "GET", r.Method)
					assertEmptyFieldSpec(t, r)
					assertAuthorization(t, r)

					w.WriteHeader(tt.respStatus)
					_, _ = w.Write(tt.respBody)
				},
			)

			got, resp, err := c.VirtualMachineNetworkInterfaces.Get(
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

			if tt.want != nil {
				assert.Equal(t, tt.want, got)
			}

			if tt.errResp != nil {
				assert.Equal(t, tt.errResp, resp.Error)
			}
		})
	}
}

func TestVirtualMachineNetworkInterfacesClient_AvailableIPs(t *testing.T) {
	type args struct {
		ctx   context.Context
		vmnet *VirtualMachineNetworkInterface
		ipVer IPVersion
	}
	tests := []struct {
		name       string
		args       args
		want       []*IPAddress
		errStr     string
		errResp    *ResponseError
		respStatus int
		respBody   []byte
	}{
		{
			name: "IPv4",
			args: args{
				ctx: context.Background(),
				vmnet: &VirtualMachineNetworkInterface{
					ID: "vmnet_olNAz8ThH0emHvdr",
				},
				ipVer: IPv4,
			},
			want: []*IPAddress{
				{
					ID:         "ip_dZLqwQifQFtboHXW",
					Address:    "169.37.118.179",
					ReverseDNS: "bitter-beautiful-mango.acme.katapult.cloud",
				},
				{
					ID:         "ip_fAwrdP9NvW0Z25eE",
					Address:    "95.135.35.113",
					ReverseDNS: "popular-shapely-tank.acme.katapult.cloud",
				},
			},
			respStatus: http.StatusOK,
			respBody: fixture(
				"virtual_machine_network_interface_available_ips_ipv4",
			),
		},
		{
			name: "IPv6",
			args: args{
				ctx: context.Background(),
				vmnet: &VirtualMachineNetworkInterface{
					ID: "vmnet_olNAz8ThH0emHvdr",
				},
				ipVer: IPv6,
			},
			want: []*IPAddress{
				{
					ID:         "ip_iFFZTQqaMjf78SYk",
					Address:    "f90c:8b3a:547b:7674:2ea7:ceeb:bbc1:5490",
					ReverseDNS: "bitter-beautiful-mango.acme.katapult.cloud",
				},
				{
					ID:         "ip_ewXuxoUZCik6a44X",
					Address:    "1247:0426:be1c:9335:0751:dbeb:2564:a9e3",
					ReverseDNS: "popular-shapely-tank.acme.katapult.cloud",
				},
			},
			respStatus: http.StatusOK,
			respBody: fixture(
				"virtual_machine_network_interface_available_ips_ipv6",
			),
		},
		{
			name: "non-existent virtual machine network interface",
			args: args{
				ctx: context.Background(),
				vmnet: &VirtualMachineNetworkInterface{
					ID: "vmnet_nopethisbegone",
				},
			},
			errStr:     fixtureVMNetworkInterfaceNotFoundErr,
			errResp:    fixtureVMNetworkInterfaceNotFoundResponseError,
			respStatus: http.StatusNotFound,
			respBody: fixture(
				"virtual_machine_network_interface_not_found_error",
			),
		},
		{
			name: "nil context",
			args: args{
				ctx: nil,
				vmnet: &VirtualMachineNetworkInterface{
					ID: "vmnet_olNAz8ThH0emHvdr",
				},
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
					"/core/v1/virtual_machine_network_interfaces"+
						"/%s/available_ips/%s",
					tt.args.vmnet.ID, tt.args.ipVer,
				),
				func(w http.ResponseWriter, r *http.Request) {
					assert.Equal(t, "GET", r.Method)
					assertEmptyFieldSpec(t, r)
					assertAuthorization(t, r)

					w.WriteHeader(tt.respStatus)
					_, _ = w.Write(tt.respBody)
				},
			)

			got, resp, err := c.VirtualMachineNetworkInterfaces.AvailableIPs(
				tt.args.ctx, tt.args.vmnet, tt.args.ipVer,
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

func TestVirtualMachineNetworkInterfacesClient_AllocateIP(t *testing.T) {
	type args struct {
		ctx   context.Context
		vmnet *VirtualMachineNetworkInterface
		ip    *IPAddress
	}
	tests := []struct {
		name        string
		args        args
		wantReqBody *virtualMachineNetworkInterfaceAllocateIPRequest
		want        *VirtualMachineNetworkInterface
		errStr      string
		errResp     *ResponseError
		respStatus  int
		respBody    []byte
	}{
		{
			name: "by IP address ID",
			args: args{
				ctx: context.Background(),
				vmnet: &VirtualMachineNetworkInterface{
					ID:   "vmnet_olNAz8ThH0emHvdr",
					Name: "Public Network on bitter-beautiful-mango",
				},
				ip: &IPAddress{
					ID:         "ip_fAwrdP9NvW0Z25eE",
					Address:    "95.135.35.113",
					ReverseDNS: "bitter-beautiful-mango.acme.katapult.cloud",
				},
			},
			wantReqBody: &virtualMachineNetworkInterfaceAllocateIPRequest{
				VirtualMachineNetworkInterface: &VirtualMachineNetworkInterface{
					ID: "vmnet_olNAz8ThH0emHvdr",
				},
				IPAddress: &IPAddress{
					ID: "ip_fAwrdP9NvW0Z25eE",
				},
			},
			want: &VirtualMachineNetworkInterface{
				ID:      "vmnet_olNAz8ThH0emHvdr",
				Name:    "Public Network on bitter-beautiful-mango",
				Network: &Network{ID: "net_4s5J6gMQXhcwqIqs"},
			},
			respStatus: http.StatusOK,
			respBody: fixture(
				"virtual_machine_network_interface_allocate_ip",
			),
		},
		{
			name: "by IP address Address",
			args: args{
				ctx: context.Background(),
				vmnet: &VirtualMachineNetworkInterface{
					ID:   "vmnet_olNAz8ThH0emHvdr",
					Name: "Public Network on bitter-beautiful-mango",
				},
				ip: &IPAddress{
					Address:    "95.135.35.113",
					ReverseDNS: "bitter-beautiful-mango.acme.katapult.cloud",
				},
			},
			wantReqBody: &virtualMachineNetworkInterfaceAllocateIPRequest{
				VirtualMachineNetworkInterface: &VirtualMachineNetworkInterface{
					ID: "vmnet_olNAz8ThH0emHvdr",
				},
				IPAddress: &IPAddress{
					Address: "95.135.35.113",
				},
			},
			want: &VirtualMachineNetworkInterface{
				ID:      "vmnet_olNAz8ThH0emHvdr",
				Name:    "Public Network on bitter-beautiful-mango",
				Network: &Network{ID: "net_4s5J6gMQXhcwqIqs"},
			},
			respStatus: http.StatusOK,
			respBody: fixture(
				"virtual_machine_network_interface_allocate_ip",
			),
		},
		{
			name: "non-existent virtual machine network interface",
			args: args{
				ctx: context.Background(),
				vmnet: &VirtualMachineNetworkInterface{
					ID: "vmnet_nopethisbegone",
				},
				ip: &IPAddress{ID: "ip_fAwrdP9NvW0Z25eE"},
			},
			errStr:     fixtureVMNetworkInterfaceNotFoundErr,
			errResp:    fixtureVMNetworkInterfaceNotFoundResponseError,
			respStatus: http.StatusNotFound,
			respBody: fixture(
				"virtual_machine_network_interface_not_found_error",
			),
		},
		{
			name: "non-existent ip address",
			args: args{
				ctx: context.Background(),
				vmnet: &VirtualMachineNetworkInterface{
					ID: "vmnet_olNAz8ThH0emHvdr",
				},
				ip: &IPAddress{ID: "ip_nopethisbegone"},
			},
			errStr:     fixtureIPAddressNotFoundErr,
			errResp:    fixtureIPAddressNotFoundResponseError,
			respStatus: http.StatusNotFound,
			respBody:   fixture("ip_address_not_found_error"),
		},
		{
			name: "already allocated ip address",
			args: args{
				ctx: context.Background(),
				vmnet: &VirtualMachineNetworkInterface{
					ID: "vmnet_olNAz8ThH0emHvdr",
				},
				ip: &IPAddress{ID: "ip_fAwrdP9NvW0Z25eE"},
			},
			errStr:     fixtureIPAlreadyAllocatedErr,
			errResp:    fixtureIPAlreadyAllocatedResponseError,
			respStatus: http.StatusNotFound,
			respBody:   fixture("ip_already_allocated_error"),
		},
		{
			name: "permission denied",
			args: args{
				ctx: context.Background(),
				vmnet: &VirtualMachineNetworkInterface{
					ID: "vmnet_olNAz8ThH0emHvdr",
				},
				ip: &IPAddress{ID: "ip_fAwrdP9NvW0Z25eE"},
			},
			errStr:     fixturePermissionDeniedErr,
			errResp:    fixturePermissionDeniedResponseError,
			respStatus: http.StatusForbidden,
			respBody:   fixture("permission_denied_error"),
		},
		{
			name: "nil virtual machine network interface",
			args: args{
				ctx:   context.Background(),
				vmnet: nil,
				ip:    &IPAddress{ID: "ip_fAwrdP9NvW0Z25eE"},
			},
			errStr:     fixtureValidationErrorErr,
			errResp:    fixtureValidationErrorResponseError,
			respStatus: http.StatusUnprocessableEntity,
			respBody:   fixture("validation_error"),
		},
		{
			name: "nil ip address",
			args: args{
				ctx: context.Background(),
				vmnet: &VirtualMachineNetworkInterface{
					ID: "vmnet_olNAz8ThH0emHvdr",
				},
				ip: nil,
			},
			errStr:     fixtureValidationErrorErr,
			errResp:    fixtureValidationErrorResponseError,
			respStatus: http.StatusUnprocessableEntity,
			respBody:   fixture("validation_error"),
		},
		{
			name: "nil context",
			args: args{
				ctx: nil,
				vmnet: &VirtualMachineNetworkInterface{
					ID: "vmnet_olNAz8ThH0emHvdr",
				},
				ip: &IPAddress{
					ID: "ip_fAwrdP9NvW0Z25eE",
				},
			},
			errStr: "net/http: nil Context",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, mux, _, teardown := prepareTestClient()
			defer teardown()

			mux.HandleFunc(
				"/core/v1/virtual_machine_network_interfaces/_/allocate_ip",
				func(w http.ResponseWriter, r *http.Request) {
					assert.Equal(t, "POST", r.Method)
					assertEmptyFieldSpec(t, r)
					assertAuthorization(t, r)

					if tt.wantReqBody != nil {
						reqBody :=
							&virtualMachineNetworkInterfaceAllocateIPRequest{}
						err := strictUmarshal(r.Body, reqBody)
						assert.NoError(t, err)
						assert.Equal(t, tt.wantReqBody, reqBody)
					}

					w.WriteHeader(tt.respStatus)
					_, _ = w.Write(tt.respBody)
				},
			)

			got, resp, err := c.VirtualMachineNetworkInterfaces.AllocateIP(
				tt.args.ctx, tt.args.vmnet, tt.args.ip,
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

func TestVirtualMachineNetworkInterfacesClient_AllocateNewIP(t *testing.T) {
	type args struct {
		ctx   context.Context
		vmnet *VirtualMachineNetworkInterface
		ipVer IPVersion
	}
	tests := []struct {
		name        string
		args        args
		wantReqBody *virtualMachineNetworkInterfaceAllocateNewIPRequest
		want        *IPAddress
		errStr      string
		errResp     *ResponseError
		respStatus  int
		respBody    []byte
	}{
		{
			name: "IPv4",
			args: args{
				ctx: context.Background(),
				vmnet: &VirtualMachineNetworkInterface{
					ID: "vmnet_olNAz8ThH0emHvdr",
				},
				ipVer: IPv4,
			},
			wantReqBody: &virtualMachineNetworkInterfaceAllocateNewIPRequest{
				VirtualMachineNetworkInterface: &VirtualMachineNetworkInterface{
					ID: "vmnet_olNAz8ThH0emHvdr",
				},
				AddressVersion: IPv4,
			},
			want: &IPAddress{
				ID:         "ip_dZLqwQifQFtboHXW",
				Address:    "169.37.118.179",
				ReverseDNS: "bitter-beautiful-mango.acme.katapult.cloud",
			},
			respStatus: http.StatusOK,
			respBody: fixture(
				"virtual_machine_network_interface_allocate_new_ipv4",
			),
		},
		{
			name: "IPv6",
			args: args{
				ctx: context.Background(),
				vmnet: &VirtualMachineNetworkInterface{
					ID: "vmnet_olNAz8ThH0emHvdr",
				},
				ipVer: IPv6,
			},
			wantReqBody: &virtualMachineNetworkInterfaceAllocateNewIPRequest{
				VirtualMachineNetworkInterface: &VirtualMachineNetworkInterface{
					ID: "vmnet_olNAz8ThH0emHvdr",
				},
				AddressVersion: IPv6,
			},
			want: &IPAddress{
				ID:         "ip_iFFZTQqaMjf78SYk",
				Address:    "f90c:8b3a:547b:7674:2ea7:ceeb:bbc1:5490",
				ReverseDNS: "bitter-beautiful-mango.acme.katapult.cloud",
			},
			respStatus: http.StatusOK,
			respBody: fixture(
				"virtual_machine_network_interface_allocate_new_ipv6",
			),
		},
		{
			name: "non-existent virtual machine network interface",
			args: args{
				ctx: context.Background(),
				vmnet: &VirtualMachineNetworkInterface{
					ID: "vmnet_nopethisbegone",
				},
			},
			errStr:     fixtureVMNetworkInterfaceNotFoundErr,
			errResp:    fixtureVMNetworkInterfaceNotFoundResponseError,
			respStatus: http.StatusNotFound,
			respBody: fixture(
				"virtual_machine_network_interface_not_found_error",
			),
		},
		{
			name: "non-existent ip address",
			args: args{
				ctx: context.Background(),
				vmnet: &VirtualMachineNetworkInterface{
					ID: "vmnet_olNAz8ThH0emHvdr",
				},
			},
			errStr:     fixtureIPAddressNotFoundErr,
			errResp:    fixtureIPAddressNotFoundResponseError,
			respStatus: http.StatusNotFound,
			respBody:   fixture("ip_address_not_found_error"),
		},
		{
			name: "no addresses available",
			args: args{
				ctx: context.Background(),
				vmnet: &VirtualMachineNetworkInterface{
					ID: "vmnet_olNAz8ThH0emHvdr",
				},
			},
			errStr:     fixtureNoAvailableAddressesErr,
			errResp:    fixtureNoAvailableAddressesResponseError,
			respStatus: http.StatusNotFound,
			respBody:   fixture("no_available_addresses_error"),
		},
		{
			name: "permission denied",
			args: args{
				ctx: context.Background(),
				vmnet: &VirtualMachineNetworkInterface{
					ID: "vmnet_olNAz8ThH0emHvdr",
				},
			},
			errStr:     fixturePermissionDeniedErr,
			errResp:    fixturePermissionDeniedResponseError,
			respStatus: http.StatusForbidden,
			respBody:   fixture("permission_denied_error"),
		},
		{
			name: "nil virtual machine network interface",
			args: args{
				ctx:   context.Background(),
				vmnet: nil,
			},
			errStr:     fixtureValidationErrorErr,
			errResp:    fixtureValidationErrorResponseError,
			respStatus: http.StatusUnprocessableEntity,
			respBody:   fixture("validation_error"),
		},
		{
			name: "nil context",
			args: args{
				ctx: nil,
				vmnet: &VirtualMachineNetworkInterface{
					ID: "vmnet_olNAz8ThH0emHvdr",
				},
			},
			errStr: "net/http: nil Context",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, mux, _, teardown := prepareTestClient()
			defer teardown()

			mux.HandleFunc(
				"/core/v1/virtual_machine_network_interfaces/_/allocate_new_ip",
				func(w http.ResponseWriter, r *http.Request) {
					assert.Equal(t, "POST", r.Method)
					assertEmptyFieldSpec(t, r)
					assertAuthorization(t, r)

					if tt.wantReqBody != nil {
						//nolint:lll
						reqBody :=
							&virtualMachineNetworkInterfaceAllocateNewIPRequest{}
						err := strictUmarshal(r.Body, reqBody)
						assert.NoError(t, err)
						assert.Equal(t, tt.wantReqBody, reqBody)
					}

					w.WriteHeader(tt.respStatus)
					_, _ = w.Write(tt.respBody)
				},
			)

			got, resp, err := c.VirtualMachineNetworkInterfaces.AllocateNewIP(
				tt.args.ctx, tt.args.vmnet, tt.args.ipVer,
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
