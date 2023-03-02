package core

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"testing"

	"github.com/jimeh/undent"
	"github.com/krystal/go-katapult"
	"github.com/krystal/go-katapult/internal/testclient"
	"github.com/stretchr/testify/assert"
)

func TestClient_VirtualMachineNetworkInterfaces(t *testing.T) {
	c := New(&testclient.Client{})

	assert.IsType(t,
		&VirtualMachineNetworkInterfacesClient{},
		c.VirtualMachineNetworkInterfaces,
	)
}

//nolint:lll
var (
	fixtureVMNetworkInterfaceNotFoundErr = "katapult: not_found: " +
		"virtual_machine_network_interface_not_found: No network interface " +
		"was found matching any of the criteria provided in the arguments"
	fixtureVMNetworkInterfaceNotFoundResponseError = &katapult.ResponseError{
		Code: "virtual_machine_network_interface_not_found",
		Description: "No network interface was found matching any of the " +
			"criteria provided in the arguments",
		Detail: json.RawMessage(`{}`),
	}

	fixtureNetworkSpeedProfileNotFoundErr = "katapult: not_found: " +
		"network_speed_profile_not_found: No network speed profile was found " +
		"matching any of the criteria provided in the arguments"
	fixtureNetworkSpeedProfileNotFoundResponseError = &katapult.ResponseError{
		Code: "network_speed_profile_not_found",
		Description: "No network speed profile was found matching any of the " +
			"criteria provided in the arguments",
		Detail: json.RawMessage(`{}`),
	}

	fixtureSpeedProfileAlreadyAssignedNotFoundErr = "katapult: " +
		"unprocessable_entity: speed_profile_already_assigned: This network " +
		"speed profile is already assigned to this virtual machine network " +
		"interface."
	fixtureSpeedProfileAlreadyAssignedNotFoundResponseError = &katapult.ResponseError{
		Code: "speed_profile_already_assigned",
		Description: "This network speed profile is already assigned to this " +
			"virtual machine network interface.",
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
		SpeedProfile: &NetworkSpeedProfile{
			ID:        "nsp_H3Mknnus3dtDIbIc",
			Name:      "1 Gbps",
			Permalink: "1gbps",
		},
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

func TestVirtualMachineNetworkInterface_Ref(t *testing.T) {
	tests := []struct {
		name string
		obj  *VirtualMachineNetworkInterface
		want VirtualMachineNetworkInterfaceRef
	}{
		{
			name: "empty",
			obj:  &VirtualMachineNetworkInterface{},
			want: VirtualMachineNetworkInterfaceRef{},
		},
		{
			name: "full",
			obj:  fixtureVirtualMachineNetworkInterfaceFull,
			want: VirtualMachineNetworkInterfaceRef{
				ID: "vmnet_Qlu34yEQgkrIlzql",
			},
		},
		{
			name: "no ID",
			obj:  fixtureVirtualMachineNetworkInterfaceNoID,
			want: VirtualMachineNetworkInterfaceRef{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.obj.Ref()

			assert.Equal(t, tt.want, got)
		})
	}
}

func TestVirtualMachineNetworkInterface_queryValues(t *testing.T) {
	tests := []struct {
		name string
		obj  VirtualMachineNetworkInterfaceRef
	}{
		{
			name: "empty",
			obj:  VirtualMachineNetworkInterfaceRef{},
		},
		{
			name: "full",
			obj: VirtualMachineNetworkInterfaceRef{
				ID: "vmnet_Qlu34yEQgkrIlzql",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testQueryableEncoding(t, tt.obj)
		})
	}
}

func Test_virtualMachineNetworkInterfacesResponseBody_JSONMarshaling(
	t *testing.T,
) {
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
				Pagination: &katapult.Pagination{CurrentPage: 345},
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

func Test_virtualMachineNetworkInterfaceAllocateIPRequest_JSONMarshaling(
	t *testing.T,
) {
	tests := []struct {
		name string
		obj  *virtualMachineNetworkInterfaceAllocateIPRequest
	}{
		{
			name: "empty",
			obj:  &virtualMachineNetworkInterfaceAllocateIPRequest{},
		},
		{
			name: "full",
			obj: &virtualMachineNetworkInterfaceAllocateIPRequest{
				//nolint:lll
				VirtualMachineNetworkInterface: VirtualMachineNetworkInterfaceRef{
					ID: "id1",
				},
				IPAddress: IPAddressRef{ID: "id3"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testJSONMarshaling(t, tt.obj)
		})
	}
}

func Test_virtualMachineNetworkInterfaceAllocateNewIPRequest_JSONMarshaling(
	t *testing.T,
) {
	tests := []struct {
		name string
		obj  *virtualMachineNetworkInterfaceAllocateNewIPRequest
	}{
		{
			name: "empty",
			obj:  &virtualMachineNetworkInterfaceAllocateNewIPRequest{},
		},
		{
			name: "ipv4",
			obj: &virtualMachineNetworkInterfaceAllocateNewIPRequest{
				//nolint:lll
				VirtualMachineNetworkInterface: VirtualMachineNetworkInterfaceRef{
					ID: "id1",
				},
				AddressVersion: IPv4,
			},
		},
		{
			name: "ipv6",
			obj: &virtualMachineNetworkInterfaceAllocateNewIPRequest{
				//nolint:lll
				VirtualMachineNetworkInterface: VirtualMachineNetworkInterfaceRef{
					ID: "id1",
				},
				AddressVersion: IPv6,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testJSONMarshaling(t, tt.obj)
		})
	}
}

//nolint:lll
func Test_virtualMachineNetworkInterfaceUpdateSpeedProfileRequest_JSONMarshaling(
	t *testing.T,
) {
	tests := []struct {
		name string
		obj  *virtualMachineNetworkInterfaceUpdateSpeedProfileRequest
	}{
		{
			name: "empty",
			obj:  &virtualMachineNetworkInterfaceUpdateSpeedProfileRequest{},
		},
		{
			name: "full",
			obj: &virtualMachineNetworkInterfaceUpdateSpeedProfileRequest{
				VirtualMachineNetworkInterface: VirtualMachineNetworkInterfaceRef{
					ID: "id1",
				},
				SpeedProfile: NetworkSpeedProfileRef{
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
		vm   VirtualMachineRef
		opts *ListOptions
	}
	tests := []struct {
		name           string
		args           args
		want           []*VirtualMachineNetworkInterface
		wantQuery      *url.Values
		wantPagination *katapult.Pagination
		errStr         string
		errResp        *katapult.ResponseError
		errIs          error
		respStatus     int
		respBody       []byte
	}{
		{
			name: "by virtual machine ID",
			args: args{
				ctx: context.Background(),
				vm:  VirtualMachineRef{ID: "vm_i5qfOrvEI1CmNrJx"},
			},
			want: virtualMachinesList,
			wantQuery: &url.Values{
				"virtual_machine[id]": []string{"vm_i5qfOrvEI1CmNrJx"},
			},
			wantPagination: &katapult.Pagination{
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
				vm:  VirtualMachineRef{FQDN: "acme"},
			},
			want: virtualMachinesList,
			wantQuery: &url.Values{
				"virtual_machine[fqdn]": []string{"acme"},
			},
			wantPagination: &katapult.Pagination{
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
				vm:   VirtualMachineRef{ID: "vm_i5qfOrvEI1CmNrJx"},
				opts: &ListOptions{Page: 1, PerPage: 2},
			},
			want: virtualMachinesList[0:2],
			wantQuery: &url.Values{
				"virtual_machine[id]": []string{"vm_i5qfOrvEI1CmNrJx"},
				"page":                []string{"1"},
				"per_page":            []string{"2"},
			},
			wantPagination: &katapult.Pagination{
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
				vm:   VirtualMachineRef{ID: "vm_i5qfOrvEI1CmNrJx"},
				opts: &ListOptions{Page: 2, PerPage: 2},
			},
			want: virtualMachinesList[2:],
			wantQuery: &url.Values{
				"virtual_machine[id]": []string{"vm_i5qfOrvEI1CmNrJx"},
				"page":                []string{"2"},
				"per_page":            []string{"2"},
			},
			wantPagination: &katapult.Pagination{
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
				vm:  VirtualMachineRef{ID: "vm_i5qfOrvEI1CmNrJx"},
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
				vm:  VirtualMachineRef{ID: "vm_i5qfOrvEI1CmNrJx"},
			},
			errStr:     fixtureVirtualMachineNotFoundErr,
			errResp:    fixtureVirtualMachineNotFoundResponseError,
			errIs:      ErrVirtualMachineNotFound,
			respStatus: http.StatusNotFound,
			respBody:   fixture("virtual_machine_not_found_error"),
		},
		{
			name: "nil context",
			args: args{
				ctx: nil,
				vm:  VirtualMachineRef{ID: "vm_i5qfOrvEI1CmNrJx"},
			},
			errStr: "net/http: nil Context",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rm, mux, _, teardown := prepareTestClient(t)
			defer teardown()
			c := NewVirtualMachineNetworkInterfacesClient(rm)

			mux.HandleFunc(
				"/core/v1/virtual_machines/_/network_interfaces",
				func(w http.ResponseWriter, r *http.Request) {
					assert.Equal(t, "GET", r.Method)
					assertEmptyFieldSpec(t, r)
					assertAuthorization(t, r)
					assertRequestOptionHeader(t, r)

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

			got, resp, err := c.List(
				tt.args.ctx, tt.args.vm, tt.args.opts, testRequestOption,
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

func TestVirtualMachineNetworkInterfacesClient_GetByID(t *testing.T) {
	type args struct {
		ctx context.Context
		id  string
	}
	tests := []struct {
		name       string
		args       args
		want       *VirtualMachineNetworkInterface
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
			errIs:      ErrVirtualMachineNetworkInterfaceNotFound,
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
			rm, mux, _, teardown := prepareTestClient(t)
			defer teardown()
			c := NewVirtualMachineNetworkInterfacesClient(rm)

			mux.HandleFunc(
				"/core/v1/virtual_machine_network_interfaces/_",
				func(w http.ResponseWriter, r *http.Request) {
					assert.Equal(t, "GET", r.Method)
					assertEmptyFieldSpec(t, r)
					assertAuthorization(t, r)
					assertRequestOptionHeader(t, r)

					assert.Equal(t, url.Values{
						"virtual_machine_network_interface[id]": []string{
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

func TestVirtualMachineNetworkInterfacesClient_Get(t *testing.T) {
	type args struct {
		ctx context.Context
		ref VirtualMachineNetworkInterfaceRef
	}
	tests := []struct {
		name       string
		args       args
		want       *VirtualMachineNetworkInterface
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
				ref: VirtualMachineNetworkInterfaceRef{
					ID: "vmnet_olNAz8ThH0emHvdr",
				},
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
				ref: VirtualMachineNetworkInterfaceRef{
					ID: "vmnet_nopethisbegone",
				},
			},
			errStr:     fixtureVMNetworkInterfaceNotFoundErr,
			errResp:    fixtureVMNetworkInterfaceNotFoundResponseError,
			errIs:      ErrVirtualMachineNetworkInterfaceNotFound,
			respStatus: http.StatusNotFound,
			respBody: fixture(
				"virtual_machine_network_interface_not_found_error",
			),
		},
		{
			name: "nil context",
			args: args{
				ctx: nil,
				ref: VirtualMachineNetworkInterfaceRef{
					ID: "vmnet_nopethisbegone",
				},
			},
			errStr: "net/http: nil Context",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rm, mux, _, teardown := prepareTestClient(t)
			defer teardown()
			c := NewVirtualMachineNetworkInterfacesClient(rm)

			mux.HandleFunc(
				"/core/v1/virtual_machine_network_interfaces/_",
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
		errResp    *katapult.ResponseError
		errIs      error
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
			errIs:      ErrVirtualMachineNetworkInterfaceNotFound,
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
			rm, mux, _, teardown := prepareTestClient(t)
			defer teardown()
			c := NewVirtualMachineNetworkInterfacesClient(rm)

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
					assertRequestOptionHeader(t, r)

					w.WriteHeader(tt.respStatus)
					_, _ = w.Write(tt.respBody)
				},
			)

			got, resp, err := c.AvailableIPs(
				tt.args.ctx, tt.args.vmnet, tt.args.ipVer, testRequestOption,
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

func TestVirtualMachineNetworkInterfacesClient_AllocateIP(t *testing.T) {
	type args struct {
		ctx   context.Context
		vmnet VirtualMachineNetworkInterfaceRef
		ip    IPAddressRef
	}
	tests := []struct {
		name        string
		args        args
		wantReqBody *virtualMachineNetworkInterfaceAllocateIPRequest
		want        *VirtualMachineNetworkInterface
		errStr      string
		errResp     *katapult.ResponseError
		errIs       error
		respStatus  int
		respBody    []byte
	}{
		{
			name: "by IP address ID",
			args: args{
				ctx: context.Background(),
				vmnet: VirtualMachineNetworkInterfaceRef{
					ID: "vmnet_olNAz8ThH0emHvdr",
				},
				ip: IPAddressRef{
					ID: "ip_fAwrdP9NvW0Z25eE",
				},
			},
			wantReqBody: &virtualMachineNetworkInterfaceAllocateIPRequest{
				//nolint:lll
				VirtualMachineNetworkInterface: VirtualMachineNetworkInterfaceRef{
					ID: "vmnet_olNAz8ThH0emHvdr",
				},
				IPAddress: IPAddressRef{
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
				vmnet: VirtualMachineNetworkInterfaceRef{
					ID: "vmnet_olNAz8ThH0emHvdr",
				},
				ip: IPAddressRef{
					Address: "95.135.35.113",
				},
			},
			wantReqBody: &virtualMachineNetworkInterfaceAllocateIPRequest{
				//nolint:lll
				VirtualMachineNetworkInterface: VirtualMachineNetworkInterfaceRef{
					ID: "vmnet_olNAz8ThH0emHvdr",
				},
				IPAddress: IPAddressRef{
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
				vmnet: VirtualMachineNetworkInterfaceRef{
					ID: "vmnet_nopethisbegone",
				},
				ip: IPAddressRef{ID: "ip_fAwrdP9NvW0Z25eE"},
			},
			errStr:     fixtureVMNetworkInterfaceNotFoundErr,
			errResp:    fixtureVMNetworkInterfaceNotFoundResponseError,
			errIs:      ErrVirtualMachineNetworkInterfaceNotFound,
			respStatus: http.StatusNotFound,
			respBody: fixture(
				"virtual_machine_network_interface_not_found_error",
			),
		},
		{
			name: "non-existent ip address",
			args: args{
				ctx: context.Background(),
				vmnet: VirtualMachineNetworkInterfaceRef{
					ID: "vmnet_olNAz8ThH0emHvdr",
				},
				ip: IPAddressRef{ID: "ip_nopethisbegone"},
			},
			errStr:     fixtureIPAddressNotFoundErr,
			errResp:    fixtureIPAddressNotFoundResponseError,
			errIs:      ErrIPAddressNotFound,
			respStatus: http.StatusNotFound,
			respBody:   fixture("ip_address_not_found_error"),
		},
		{
			name: "already allocated ip address",
			args: args{
				ctx: context.Background(),
				vmnet: VirtualMachineNetworkInterfaceRef{
					ID: "vmnet_olNAz8ThH0emHvdr",
				},
				ip: IPAddressRef{ID: "ip_fAwrdP9NvW0Z25eE"},
			},
			errStr:     fixtureIPAlreadyAllocatedErr,
			errResp:    fixtureIPAlreadyAllocatedResponseError,
			errIs:      ErrIPAlreadyAllocated,
			respStatus: http.StatusNotFound,
			respBody:   fixture("ip_already_allocated_error"),
		},
		{
			name: "permission denied",
			args: args{
				ctx: context.Background(),
				vmnet: VirtualMachineNetworkInterfaceRef{
					ID: "vmnet_olNAz8ThH0emHvdr",
				},
				ip: IPAddressRef{ID: "ip_fAwrdP9NvW0Z25eE"},
			},
			errStr:     fixturePermissionDeniedErr,
			errResp:    fixturePermissionDeniedResponseError,
			errIs:      ErrPermissionDenied,
			respStatus: http.StatusForbidden,
			respBody:   fixture("permission_denied_error"),
		},
		{
			name: "nil context",
			args: args{
				ctx: nil,
				vmnet: VirtualMachineNetworkInterfaceRef{
					ID: "vmnet_olNAz8ThH0emHvdr",
				},
				ip: IPAddressRef{
					ID: "ip_fAwrdP9NvW0Z25eE",
				},
			},
			errStr: "net/http: nil Context",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rm, mux, _, teardown := prepareTestClient(t)
			defer teardown()
			c := NewVirtualMachineNetworkInterfacesClient(rm)

			mux.HandleFunc(
				"/core/v1/virtual_machine_network_interfaces/_/allocate_ip",
				func(w http.ResponseWriter, r *http.Request) {
					assert.Equal(t, "POST", r.Method)
					assertEmptyFieldSpec(t, r)
					assertAuthorization(t, r)
					assertRequestOptionHeader(t, r)

					if tt.wantReqBody != nil {
						//nolint:lll
						reqBody := &virtualMachineNetworkInterfaceAllocateIPRequest{}
						err := strictUmarshal(r.Body, reqBody)
						assert.NoError(t, err)
						assert.Equal(t, tt.wantReqBody, reqBody)
					}

					w.WriteHeader(tt.respStatus)
					_, _ = w.Write(tt.respBody)
				},
			)

			got, resp, err := c.AllocateIP(
				tt.args.ctx, tt.args.vmnet, tt.args.ip, testRequestOption,
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

func TestVirtualMachineNetworkInterfacesClient_AllocateNewIP(t *testing.T) {
	type args struct {
		ctx   context.Context
		vmnet VirtualMachineNetworkInterfaceRef
		ipVer IPVersion
	}
	tests := []struct {
		name        string
		args        args
		wantReqBody *virtualMachineNetworkInterfaceAllocateNewIPRequest
		want        *IPAddress
		errStr      string
		errResp     *katapult.ResponseError
		errIs       error
		respStatus  int
		respBody    []byte
	}{
		{
			name: "IPv4",
			args: args{
				ctx: context.Background(),
				vmnet: VirtualMachineNetworkInterfaceRef{
					ID: "vmnet_olNAz8ThH0emHvdr",
				},
				ipVer: IPv4,
			},
			wantReqBody: &virtualMachineNetworkInterfaceAllocateNewIPRequest{
				//nolint:lll
				VirtualMachineNetworkInterface: VirtualMachineNetworkInterfaceRef{
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
				vmnet: VirtualMachineNetworkInterfaceRef{
					ID: "vmnet_olNAz8ThH0emHvdr",
				},
				ipVer: IPv6,
			},
			wantReqBody: &virtualMachineNetworkInterfaceAllocateNewIPRequest{
				//nolint:lll
				VirtualMachineNetworkInterface: VirtualMachineNetworkInterfaceRef{
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
				vmnet: VirtualMachineNetworkInterfaceRef{
					ID: "vmnet_nopethisbegone",
				},
			},
			errStr:     fixtureVMNetworkInterfaceNotFoundErr,
			errResp:    fixtureVMNetworkInterfaceNotFoundResponseError,
			errIs:      ErrVirtualMachineNetworkInterfaceNotFound,
			respStatus: http.StatusNotFound,
			respBody: fixture(
				"virtual_machine_network_interface_not_found_error",
			),
		},
		{
			name: "non-existent ip address",
			args: args{
				ctx: context.Background(),
				vmnet: VirtualMachineNetworkInterfaceRef{
					ID: "vmnet_olNAz8ThH0emHvdr",
				},
			},
			errStr:     fixtureIPAddressNotFoundErr,
			errResp:    fixtureIPAddressNotFoundResponseError,
			errIs:      ErrIPAddressNotFound,
			respStatus: http.StatusNotFound,
			respBody:   fixture("ip_address_not_found_error"),
		},
		{
			name: "no addresses available",
			args: args{
				ctx: context.Background(),
				vmnet: VirtualMachineNetworkInterfaceRef{
					ID: "vmnet_olNAz8ThH0emHvdr",
				},
			},
			errStr:     fixtureNoAvailableAddressesErr,
			errResp:    fixtureNoAvailableAddressesResponseError,
			errIs:      ErrNoAvailableAddresses,
			respStatus: http.StatusNotFound,
			respBody:   fixture("no_available_addresses_error"),
		},
		{
			name: "permission denied",
			args: args{
				ctx: context.Background(),
				vmnet: VirtualMachineNetworkInterfaceRef{
					ID: "vmnet_olNAz8ThH0emHvdr",
				},
			},
			errStr:     fixturePermissionDeniedErr,
			errResp:    fixturePermissionDeniedResponseError,
			errIs:      ErrPermissionDenied,
			respStatus: http.StatusForbidden,
			respBody:   fixture("permission_denied_error"),
		},
		{
			name: "nil context",
			args: args{
				ctx: nil,
				vmnet: VirtualMachineNetworkInterfaceRef{
					ID: "vmnet_olNAz8ThH0emHvdr",
				},
			},
			errStr: "net/http: nil Context",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rm, mux, _, teardown := prepareTestClient(t)
			defer teardown()
			c := NewVirtualMachineNetworkInterfacesClient(rm)

			mux.HandleFunc(
				"/core/v1/virtual_machine_network_interfaces/_/allocate_new_ip",
				func(w http.ResponseWriter, r *http.Request) {
					assert.Equal(t, "POST", r.Method)
					assertEmptyFieldSpec(t, r)
					assertAuthorization(t, r)
					assertRequestOptionHeader(t, r)

					if tt.wantReqBody != nil {
						//nolint:lll
						reqBody := &virtualMachineNetworkInterfaceAllocateNewIPRequest{}
						err := strictUmarshal(r.Body, reqBody)
						assert.NoError(t, err)
						assert.Equal(t, tt.wantReqBody, reqBody)
					}

					w.WriteHeader(tt.respStatus)
					_, _ = w.Write(tt.respBody)
				},
			)

			got, resp, err := c.AllocateNewIP(
				tt.args.ctx, tt.args.vmnet, tt.args.ipVer, testRequestOption,
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

func TestVirtualMachineNetworkInterfacesClient_UpdateSpeedProfile(
	t *testing.T,
) {
	taskResponseBody := undent.Bytes(`
		{
			"task": {
				"id": "task_bdhDJgeDWM0SPskh",
				"name": "Change network interface speed profile",
				"status": "completed"
			}
		}`,
	)

	type args struct {
		ctx          context.Context
		vmnet        VirtualMachineNetworkInterfaceRef
		speedProfile NetworkSpeedProfileRef
	}
	tests := []struct {
		name        string
		args        args
		wantReqBody *virtualMachineNetworkInterfaceUpdateSpeedProfileRequest
		want        *Task
		errStr      string
		errResp     *katapult.ResponseError
		errIs       error
		respStatus  int
		respBody    []byte
	}{
		{
			name: "by NetworkSpeedProfile ID",
			args: args{
				ctx: context.Background(),
				vmnet: VirtualMachineNetworkInterfaceRef{
					ID: "vmnet_olNAz8ThH0emHvdr",
				},
				speedProfile: NetworkSpeedProfileRef{
					ID: "nsp_CReSzkaCt01kWoi7",
				},
			},
			//nolint:lll
			wantReqBody: &virtualMachineNetworkInterfaceUpdateSpeedProfileRequest{
				VirtualMachineNetworkInterface: VirtualMachineNetworkInterfaceRef{
					ID: "vmnet_olNAz8ThH0emHvdr",
				},
				SpeedProfile: NetworkSpeedProfileRef{
					ID: "nsp_CReSzkaCt01kWoi7",
				},
			},
			want: &Task{
				ID:     "task_bdhDJgeDWM0SPskh",
				Name:   "Change network interface speed profile",
				Status: "completed",
			},
			respStatus: http.StatusOK,
			respBody:   taskResponseBody,
		},
		{
			name: "by NetworkSpeedProfile Permalink",
			args: args{
				ctx: context.Background(),
				vmnet: VirtualMachineNetworkInterfaceRef{
					ID: "vmnet_olNAz8ThH0emHvdr",
				},
				speedProfile: NetworkSpeedProfileRef{
					Permalink: "1gbps",
				},
			},
			//nolint:lll
			wantReqBody: &virtualMachineNetworkInterfaceUpdateSpeedProfileRequest{
				VirtualMachineNetworkInterface: VirtualMachineNetworkInterfaceRef{
					ID: "vmnet_olNAz8ThH0emHvdr",
				},
				SpeedProfile: NetworkSpeedProfileRef{
					Permalink: "1gbps",
				},
			},
			want: &Task{
				ID:     "task_bdhDJgeDWM0SPskh",
				Name:   "Change network interface speed profile",
				Status: "completed",
			},
			respStatus: http.StatusOK,
			respBody:   taskResponseBody,
		},
		{
			name: "non-existent virtual machine network interface",
			args: args{
				ctx: context.Background(),
				vmnet: VirtualMachineNetworkInterfaceRef{
					ID: "vmnet_nopethisbegone",
				},
				speedProfile: NetworkSpeedProfileRef{
					ID: "nsp_CReSzkaCt01kWoi7",
				},
			},
			errStr:     fixtureVMNetworkInterfaceNotFoundErr,
			errResp:    fixtureVMNetworkInterfaceNotFoundResponseError,
			errIs:      ErrVirtualMachineNetworkInterfaceNotFound,
			respStatus: http.StatusNotFound,
			respBody: fixture(
				"virtual_machine_network_interface_not_found_error",
			),
		},
		{
			name: "non-existent speed profile",
			args: args{
				ctx: context.Background(),
				vmnet: VirtualMachineNetworkInterfaceRef{
					ID: "vmnet_olNAz8ThH0emHvdr",
				},
				speedProfile: NetworkSpeedProfileRef{
					ID: "nsp_iRIhTnddeHCK9ZBj",
				},
			},
			errStr:     fixtureNetworkSpeedProfileNotFoundErr,
			errResp:    fixtureNetworkSpeedProfileNotFoundResponseError,
			errIs:      ErrNetworkSpeedProfileNotFound,
			respStatus: http.StatusNotFound,
			respBody:   fixture("network_speed_profile_not_found_error"),
		},
		{
			name: "speed profile already assigned",
			args: args{
				ctx: context.Background(),
				vmnet: VirtualMachineNetworkInterfaceRef{
					ID: "vmnet_olNAz8ThH0emHvdr",
				},
				speedProfile: NetworkSpeedProfileRef{
					ID: "nsp_iRIhTnddeHCK9ZBj",
				},
			},
			errStr:     fixtureSpeedProfileAlreadyAssignedNotFoundErr,
			errResp:    fixtureSpeedProfileAlreadyAssignedNotFoundResponseError,
			errIs:      ErrSpeedProfileAlreadyAssigned,
			respStatus: http.StatusNotFound,
			respBody:   fixture("speed_profile_already_assigned_error"),
		},
		{
			name: "task queueing error",
			args: args{
				ctx: context.Background(),
				vmnet: VirtualMachineNetworkInterfaceRef{
					ID: "vmnet_olNAz8ThH0emHvdr",
				},
				speedProfile: NetworkSpeedProfileRef{
					ID: "nsp_iRIhTnddeHCK9ZBj",
				},
			},
			errStr:     fixtureTaskQueueingErrorErr,
			errResp:    fixtureTaskQueueingErrorResponseError,
			errIs:      ErrTaskQueueingError,
			respStatus: http.StatusNotFound,
			respBody:   fixture("task_queueing_error"),
		},
		{
			name: "nil context",
			args: args{
				ctx: nil,
				vmnet: VirtualMachineNetworkInterfaceRef{
					ID: "vmnet_olNAz8ThH0emHvdr",
				},
			},
			errStr: "net/http: nil Context",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rm, mux, _, teardown := prepareTestClient(t)
			defer teardown()
			c := NewVirtualMachineNetworkInterfacesClient(rm)

			mux.HandleFunc(
				"/core/v1/virtual_machine_network_interfaces/"+
					"_/update_speed_profile",
				func(w http.ResponseWriter, r *http.Request) {
					assert.Equal(t, "PATCH", r.Method)
					assertEmptyFieldSpec(t, r)
					assertAuthorization(t, r)
					assertRequestOptionHeader(t, r)

					if tt.wantReqBody != nil {
						//nolint:lll
						reqBody := &virtualMachineNetworkInterfaceUpdateSpeedProfileRequest{}
						err := strictUmarshal(r.Body, reqBody)
						assert.NoError(t, err)
						assert.Equal(t, tt.wantReqBody, reqBody)
					}

					w.WriteHeader(tt.respStatus)
					_, _ = w.Write(tt.respBody)
				},
			)

			got, resp, err := c.UpdateSpeedProfile(
				tt.args.ctx,
				tt.args.vmnet,
				tt.args.speedProfile,
				testRequestOption,
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
