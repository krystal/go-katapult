package core

import (
	"context"
	"net/url"

	"github.com/krystal/go-katapult"
)

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

func (s *VirtualMachinePackage) Ref() VirtualMachinePackageRef {
	return VirtualMachinePackageRef{ID: s.ID}
}

type VirtualMachinePackageRef struct {
	ID        string `json:"id,omitempty"`
	Permalink string `json:"permalink,omitempty"`
}

func (vmpf VirtualMachinePackageRef) queryValues() *url.Values {
	v := &url.Values{}
	switch {
	case vmpf.ID != "":
		v.Set("virtual_machine_package[id]", vmpf.ID)
	case vmpf.Permalink != "":
		v.Set("virtual_machine_package[permalink]", vmpf.Permalink)
	}

	return v
}

type virtualMachinePackagesResponseBody struct {
	Pagination             *katapult.Pagination     `json:"pagination,omitempty"`
	VirtualMachinePackage  *VirtualMachinePackage   `json:"virtual_machine_package,omitempty"`
	VirtualMachinePackages []*VirtualMachinePackage `json:"virtual_machine_packages,omitempty"`
}

type VirtualMachinePackagesClient struct {
	client   RequestMaker
	basePath *url.URL
}

func NewVirtualMachinePackagesClient(
	rm RequestMaker,
) *VirtualMachinePackagesClient {
	return &VirtualMachinePackagesClient{
		client:   rm,
		basePath: &url.URL{Path: "/core/v1/"},
	}
}

func (s *VirtualMachinePackagesClient) List(
	ctx context.Context,
	opts *ListOptions,
) ([]*VirtualMachinePackage, *katapult.Response, error) {
	u := &url.URL{
		Path:     "virtual_machine_packages",
		RawQuery: opts.queryValues().Encode(),
	}

	body, resp, err := s.doRequest(ctx, "GET", u, nil)
	resp.Pagination = body.Pagination

	return body.VirtualMachinePackages, resp, err
}

func (s *VirtualMachinePackagesClient) Get(
	ctx context.Context,
	ref VirtualMachinePackageRef,
) (*VirtualMachinePackage, *katapult.Response, error) {
	u := &url.URL{
		Path:     "virtual_machine_packages/_",
		RawQuery: ref.queryValues().Encode(),
	}
	body, resp, err := s.doRequest(ctx, "GET", u, nil)

	return body.VirtualMachinePackage, resp, err
}

func (s *VirtualMachinePackagesClient) GetByID(
	ctx context.Context,
	id string,
) (*VirtualMachinePackage, *katapult.Response, error) {
	return s.Get(ctx, VirtualMachinePackageRef{ID: id})
}

func (s *VirtualMachinePackagesClient) GetByPermalink(
	ctx context.Context,
	permalink string,
) (*VirtualMachinePackage, *katapult.Response, error) {
	return s.Get(ctx, VirtualMachinePackageRef{Permalink: permalink})
}

func (s *VirtualMachinePackagesClient) doRequest(
	ctx context.Context,
	method string,
	u *url.URL,
	body interface{},
) (*virtualMachinePackagesResponseBody, *katapult.Response, error) {
	u = s.basePath.ResolveReference(u)
	respBody := &virtualMachinePackagesResponseBody{}
	resp := katapult.NewResponse(nil)

	req, err := s.client.NewRequestWithContext(ctx, method, u, body)
	if err == nil {
		resp, err = s.client.Do(req, respBody)
	}

	return respBody, resp, err
}
