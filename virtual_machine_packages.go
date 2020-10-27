package katapult

import (
	"context"
	"fmt"
	"net/url"

	"github.com/google/go-querystring/query"
)

type VirtualMachinePackagesService struct {
	*service
	path *url.URL
}

func NewVirtualMachinePackagesService(
	s *service,
) *VirtualMachinePackagesService {
	return &VirtualMachinePackagesService{
		service: s,
		path:    &url.URL{Path: "/core/v1/"},
	}
}

type VirtualMachinePackage struct {
	ID            string `json:"id,omitempty"`
	Name          string `json:"name,omitempty"`
	Permalink     string `json:"permalink,omitempty"`
	CPUCores      int    `json:"cpu_cores,omitempty"`
	IPv4Addresses int    `json:"ipv4_addresses,omitempty"`
	MemoryInGB    int    `json:"memory_in_gb,omitempty"`
	StorageInGB   int    `json:"storage_in_gb,omitempty"`
	Privacy       string `json:"privacy,omitempty"`
	Icon          *Icon  `json:"icon,omitempty"`
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
	u := &url.URL{Path: "virtual_machine_packages"}

	qs, err := query.Values(opts)
	if err != nil {
		return nil, nil, err
	}
	u.RawQuery = qs.Encode()

	body, resp, err := s.doRequest(ctx, "GET", u.String(), nil)
	resp.Pagination = body.Pagination

	return body.VirtualMachinePackages, resp, err
}

func (s *VirtualMachinePackagesService) Get(
	ctx context.Context,
	id string,
) (*VirtualMachinePackage, *Response, error) {
	u := fmt.Sprintf("virtual_machine_packages/%s", id)
	body, resp, err := s.doRequest(ctx, "GET", u, nil)

	return body.VirtualMachinePackage, resp, err
}

func (s *VirtualMachinePackagesService) doRequest(
	ctx context.Context,
	method string,
	urlStr string,
	body interface{},
) (*virtualMachinePackagesResponseBody, *Response, error) {
	u, err := s.path.Parse(urlStr)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequestWithContext(ctx, method, u.String(), body)
	if err != nil {
		return nil, nil, err
	}

	var respBody virtualMachinePackagesResponseBody
	resp, err := s.client.Do(req, &respBody)

	return &respBody, resp, err
}
