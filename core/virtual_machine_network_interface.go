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
func (s *VirtualMachineNetworkInterface) lookupReference() *VirtualMachineNetworkInterface {
	if s == nil {
		return nil
	}

	return &VirtualMachineNetworkInterface{ID: s.ID}
}

func (s *VirtualMachineNetworkInterface) queryValues() *url.Values {
	v := &url.Values{}

	if s != nil {
		v.Set("virtual_machine_network_interface[id]", s.ID)
	}

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
	VirtualMachineNetworkInterface *VirtualMachineNetworkInterface `json:"virtual_machine_network_interface,omitempty"`
	IPAddress                      *IPAddress                      `json:"ip_address,omitempty"`
}

type virtualMachineNetworkInterfaceAllocateNewIPRequest struct {
	VirtualMachineNetworkInterface *VirtualMachineNetworkInterface `json:"virtual_machine_network_interface,omitempty"`
	AddressVersion                 IPVersion                       `json:"address_version,omitempty"`
}

type virtualMachineNetworkInterfaceUpdateSpeedProfileRequest struct {
	VirtualMachineNetworkInterface *VirtualMachineNetworkInterface `json:"virtual_machine_network_interface,omitempty"`
	SpeedProfile                   *NetworkSpeedProfile            `json:"speed_profile,omitempty"`
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
	vm *VirtualMachine,
	opts *ListOptions,
) ([]*VirtualMachineNetworkInterface, *katapult.Response, error) {
	qs := queryValues(vm, opts)
	u := &url.URL{
		Path:     "virtual_machines/_/network_interfaces",
		RawQuery: qs.Encode(),
	}

	body, resp, err := s.doRequest(ctx, "GET", u, nil)
	resp.Pagination = body.Pagination

	return body.VirtualMachineNetworkInterfaces, resp, err
}

func (s *VirtualMachineNetworkInterfacesClient) Get(
	ctx context.Context,
	id string,
) (*VirtualMachineNetworkInterface, *katapult.Response, error) {
	u := &url.URL{
		Path: fmt.Sprintf("virtual_machine_network_interfaces/%s", id),
	}

	body, resp, err := s.doRequest(ctx, "GET", u, nil)

	return body.VirtualMachineNetworkInterface, resp, err
}

func (s *VirtualMachineNetworkInterfacesClient) AvailableIPs(
	ctx context.Context,
	vmnet *VirtualMachineNetworkInterface,
	ipVer IPVersion,
) ([]*IPAddress, *katapult.Response, error) {
	u := &url.URL{
		Path: fmt.Sprintf(
			"virtual_machine_network_interfaces/%s/available_ips/%s",
			vmnet.ID, ipVer,
		),
	}

	body, resp, err := s.doRequest(ctx, "GET", u, nil)

	return body.IPAddresses, resp, err
}

func (s *VirtualMachineNetworkInterfacesClient) AllocateIP(
	ctx context.Context,
	vmnet *VirtualMachineNetworkInterface,
	ip *IPAddress,
) (*VirtualMachineNetworkInterface, *katapult.Response, error) {
	u := &url.URL{Path: "virtual_machine_network_interfaces/_/allocate_ip"}
	reqBody := &virtualMachineNetworkInterfaceAllocateIPRequest{
		VirtualMachineNetworkInterface: vmnet.lookupReference(),
		IPAddress:                      ip.lookupReference(),
	}

	body, resp, err := s.doRequest(ctx, "POST", u, reqBody)

	return body.VirtualMachineNetworkInterface, resp, err
}

func (s *VirtualMachineNetworkInterfacesClient) AllocateNewIP(
	ctx context.Context,
	vmnet *VirtualMachineNetworkInterface,
	ipVer IPVersion,
) (*IPAddress, *katapult.Response, error) {
	u := &url.URL{Path: "virtual_machine_network_interfaces/_/allocate_new_ip"}
	reqBody := &virtualMachineNetworkInterfaceAllocateNewIPRequest{
		VirtualMachineNetworkInterface: vmnet.lookupReference(),
		AddressVersion:                 ipVer,
	}

	body, resp, err := s.doRequest(ctx, "POST", u, reqBody)

	return body.IPAddress, resp, err
}

func (s *VirtualMachineNetworkInterfacesClient) UpdateSpeedProfile(
	ctx context.Context,
	vmnet *VirtualMachineNetworkInterface,
	speedProfile *NetworkSpeedProfile,
) (*Task, *katapult.Response, error) {
	u := &url.URL{
		Path: "virtual_machine_network_interfaces/_/update_speed_profile",
	}
	reqBody := &virtualMachineNetworkInterfaceUpdateSpeedProfileRequest{
		VirtualMachineNetworkInterface: vmnet.lookupReference(),
		SpeedProfile:                   speedProfile.lookupReference(),
	}

	body, resp, err := s.doRequest(ctx, "PATCH", u, reqBody)

	return body.Task, resp, err
}

func (s *VirtualMachineNetworkInterfacesClient) doRequest(
	ctx context.Context,
	method string,
	u *url.URL,
	body interface{},
) (*virtualMachineNetworkInterfacesResponseBody, *katapult.Response, error) {
	u = s.basePath.ResolveReference(u)
	respBody := &virtualMachineNetworkInterfacesResponseBody{}
	resp := katapult.NewResponse(nil)

	req, err := s.client.NewRequestWithContext(ctx, method, u, body)
	if err == nil {
		resp, err = s.client.Do(req, respBody)
	}

	return respBody, resp, err
}
