package katapult

import (
	"context"
	"fmt"
	"net/url"
)

type VirtualMachinePackagesService struct {
	client   *apiClient
	basePath *url.URL
}

func newVirtualMachinePackagesService(
	c *apiClient,
) *VirtualMachinePackagesService {
	return &VirtualMachinePackagesService{
		client:   c,
		basePath: &url.URL{Path: "/core/v1/"},
	}
}

type VirtualMachinePackage struct {
	ID            string      `json:"id,omitempty"`
	Name          string      `json:"name,omitempty"`
	Permalink     string      `json:"permalink,omitempty"`
	CPUCores      int         `json:"cpu_cores,omitempty"`
	IPv4Addresses int         `json:"ipv4_addresses,omitempty"`
	MemoryInGB    int         `json:"memory_in_gb,omitempty"`
	StorageInGB   int         `json:"storage_in_gb,omitempty"`
	Privacy       string      `json:"privacy,omitempty"`
	Icon          *Attachment `json:"icon,omitempty"`
}

type virtualMachinePackagesResponseBody struct {
	Pagination             *Pagination              `json:"pagination,omitempty"`
	VirtualMachinePackage  *VirtualMachinePackage   `json:"virtual_machine_package,omitempty"`
	VirtualMachinePackages []*VirtualMachinePackage `json:"virtual_machine_packages,omitempty"`
}

func (s *VirtualMachinePackagesService) List(
	ctx context.Context,
	opts *ListOptions,
) ([]*VirtualMachinePackage, *Response, error) {
	u := &url.URL{
		Path:     "virtual_machine_packages",
		RawQuery: opts.Values().Encode(),
	}

	body, resp, err := s.doRequest(ctx, "GET", u, nil)
	resp.Pagination = body.Pagination

	return body.VirtualMachinePackages, resp, err
}

func (s *VirtualMachinePackagesService) Get(
	ctx context.Context,
	id string,
) (*VirtualMachinePackage, *Response, error) {
	u := &url.URL{Path: fmt.Sprintf("virtual_machine_packages/%s", id)}
	body, resp, err := s.doRequest(ctx, "GET", u, nil)

	return body.VirtualMachinePackage, resp, err
}

func (s *VirtualMachinePackagesService) doRequest(
	ctx context.Context,
	method string,
	u *url.URL,
	body interface{},
) (*virtualMachinePackagesResponseBody, *Response, error) {
	u = s.basePath.ResolveReference(u)
	respBody := &virtualMachinePackagesResponseBody{}
	resp := &Response{}

	req, err := s.client.NewRequestWithContext(ctx, method, u, body)
	if err == nil {
		resp, err = s.client.Do(req, respBody)
	}

	return respBody, resp, err
}