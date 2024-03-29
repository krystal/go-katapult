package core

import (
	"context"
	"fmt"
	"net/url"

	"github.com/krystal/go-katapult"
)

type VirtualMachineNetworkInterface struct {
	ID             string               `json:"id,omitempty"`
	VirtualMachine *VirtualMachine      `json:"virtual_machine,omitempty"`
	Name           string               `json:"name,omitempty"`
	Network        *Network             `json:"network,omitempty"`
	MACAddress     string               `json:"mac_address,omitempty"`
	State          string               `json:"state,omitempty"`
	IPAddresses    []*IPAddress         `json:"ip_addresses,omitempty"`
	SpeedProfile   *NetworkSpeedProfile `json:"speed_profile,omitempty"`
}

//nolint:lll
func (s *VirtualMachineNetworkInterface) Ref() VirtualMachineNetworkInterfaceRef {
	return VirtualMachineNetworkInterfaceRef{ID: s.ID}
}

type VirtualMachineNetworkInterfaceRef struct {
	ID string `json:"id,omitempty"`
}

func (s VirtualMachineNetworkInterfaceRef) queryValues() *url.Values {
	v := &url.Values{}
	v.Set("virtual_machine_network_interface[id]", s.ID)

	return v
}

type virtualMachineNetworkInterfacesResponseBody struct {
	Pagination                      *katapult.Pagination              `json:"pagination,omitempty"`
	Task                            *Task                             `json:"task,omitempty"`
	VirtualMachineNetworkInterface  *VirtualMachineNetworkInterface   `json:"virtual_machine_network_interface,omitempty"`
	VirtualMachineNetworkInterfaces []*VirtualMachineNetworkInterface `json:"virtual_machine_network_interfaces,omitempty"`
	IPAddress                       *IPAddress                        `json:"ip_address,omitempty"`
	IPAddresses                     []*IPAddress                      `json:"ip_addresses,omitempty"`
}

type virtualMachineNetworkInterfaceAllocateIPRequest struct {
	VirtualMachineNetworkInterface VirtualMachineNetworkInterfaceRef `json:"virtual_machine_network_interface"`
	IPAddress                      IPAddressRef                      `json:"ip_address"`
}

type virtualMachineNetworkInterfaceAllocateNewIPRequest struct {
	VirtualMachineNetworkInterface VirtualMachineNetworkInterfaceRef `json:"virtual_machine_network_interface"`
	AddressVersion                 IPVersion                         `json:"address_version"`
}

type virtualMachineNetworkInterfaceUpdateSpeedProfileRequest struct {
	VirtualMachineNetworkInterface VirtualMachineNetworkInterfaceRef `json:"virtual_machine_network_interface"`
	SpeedProfile                   NetworkSpeedProfileRef            `json:"speed_profile"`
}

type VirtualMachineNetworkInterfacesClient struct {
	client   RequestMaker
	basePath *url.URL
}

func NewVirtualMachineNetworkInterfacesClient(
	rm RequestMaker,
) *VirtualMachineNetworkInterfacesClient {
	return &VirtualMachineNetworkInterfacesClient{
		client:   rm,
		basePath: &url.URL{Path: "/core/v1/"},
	}
}

func (s *VirtualMachineNetworkInterfacesClient) List(
	ctx context.Context,
	vm VirtualMachineRef,
	opts *ListOptions,
	reqOpts ...katapult.RequestOption,
) ([]*VirtualMachineNetworkInterface, *katapult.Response, error) {
	qs := queryValues(vm, opts)
	u := &url.URL{
		Path:     "virtual_machines/_/network_interfaces",
		RawQuery: qs.Encode(),
	}

	body, resp, err := s.doRequest(ctx, "GET", u, nil, reqOpts...)
	resp.Pagination = body.Pagination

	return body.VirtualMachineNetworkInterfaces, resp, err
}

func (s *VirtualMachineNetworkInterfacesClient) Get(
	ctx context.Context,
	ref VirtualMachineNetworkInterfaceRef,
	reqOpts ...katapult.RequestOption,
) (*VirtualMachineNetworkInterface, *katapult.Response, error) {
	u := &url.URL{
		Path:     "virtual_machine_network_interfaces/_",
		RawQuery: ref.queryValues().Encode(),
	}

	body, resp, err := s.doRequest(ctx, "GET", u, nil, reqOpts...)

	return body.VirtualMachineNetworkInterface, resp, err
}

func (s *VirtualMachineNetworkInterfacesClient) GetByID(
	ctx context.Context,
	id string,
	reqOpts ...katapult.RequestOption,
) (*VirtualMachineNetworkInterface, *katapult.Response, error) {
	return s.Get(ctx, VirtualMachineNetworkInterfaceRef{ID: id}, reqOpts...)
}

func (s *VirtualMachineNetworkInterfacesClient) AvailableIPs(
	ctx context.Context,
	vmnet *VirtualMachineNetworkInterface,
	ipVer IPVersion,
	reqOpts ...katapult.RequestOption,
) ([]*IPAddress, *katapult.Response, error) {
	u := &url.URL{
		Path: fmt.Sprintf(
			"virtual_machine_network_interfaces/%s/available_ips/%s",
			vmnet.ID, ipVer,
		),
	}

	body, resp, err := s.doRequest(ctx, "GET", u, nil, reqOpts...)

	return body.IPAddresses, resp, err
}

func (s *VirtualMachineNetworkInterfacesClient) AllocateIP(
	ctx context.Context,
	vmnet VirtualMachineNetworkInterfaceRef,
	ip IPAddressRef,
	reqOpts ...katapult.RequestOption,
) (*VirtualMachineNetworkInterface, *katapult.Response, error) {
	u := &url.URL{Path: "virtual_machine_network_interfaces/_/allocate_ip"}
	reqBody := &virtualMachineNetworkInterfaceAllocateIPRequest{
		VirtualMachineNetworkInterface: vmnet,
		IPAddress:                      ip,
	}

	body, resp, err := s.doRequest(ctx, "POST", u, reqBody, reqOpts...)

	return body.VirtualMachineNetworkInterface, resp, err
}

func (s *VirtualMachineNetworkInterfacesClient) AllocateNewIP(
	ctx context.Context,
	vmnet VirtualMachineNetworkInterfaceRef,
	ipVer IPVersion,
	reqOpts ...katapult.RequestOption,
) (*IPAddress, *katapult.Response, error) {
	u := &url.URL{Path: "virtual_machine_network_interfaces/_/allocate_new_ip"}
	reqBody := &virtualMachineNetworkInterfaceAllocateNewIPRequest{
		VirtualMachineNetworkInterface: vmnet,
		AddressVersion:                 ipVer,
	}

	body, resp, err := s.doRequest(ctx, "POST", u, reqBody, reqOpts...)

	return body.IPAddress, resp, err
}

func (s *VirtualMachineNetworkInterfacesClient) UpdateSpeedProfile(
	ctx context.Context,
	vmnet VirtualMachineNetworkInterfaceRef,
	speedProfile NetworkSpeedProfileRef,
	reqOpts ...katapult.RequestOption,
) (*Task, *katapult.Response, error) {
	u := &url.URL{
		Path: "virtual_machine_network_interfaces/_/update_speed_profile",
	}
	reqBody := &virtualMachineNetworkInterfaceUpdateSpeedProfileRequest{
		VirtualMachineNetworkInterface: vmnet,
		SpeedProfile:                   speedProfile,
	}

	body, resp, err := s.doRequest(ctx, "PATCH", u, reqBody, reqOpts...)

	return body.Task, resp, err
}

func (s *VirtualMachineNetworkInterfacesClient) doRequest(
	ctx context.Context,
	method string,
	u *url.URL,
	body interface{},
	reqOpts ...katapult.RequestOption,
) (*virtualMachineNetworkInterfacesResponseBody, *katapult.Response, error) {
	u = s.basePath.ResolveReference(u)
	respBody := &virtualMachineNetworkInterfacesResponseBody{}

	req := katapult.NewRequest(method, u, body, reqOpts...)
	resp, err := s.client.Do(ctx, req, respBody)
	if resp == nil {
		resp = katapult.NewResponse(nil)
	}

	return respBody, resp, handleResponseError(err)
}
