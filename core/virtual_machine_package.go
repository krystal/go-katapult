package core

import (
	"context"
	"fmt"
	"github.com/krystal/go-katapult"
	"net/url"
	"strings"
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

// NewVirtualMachinePackageLookup takes a string that is a VirtualMachinePackage
// ID or Permalink returning, a empty *VirtualMachinePackage struct with either
// the ID or Permalink field populated with the given value. This struct is
// suitable as input to other methods which accept a *VirtualMachinePackage as
// input.
func NewVirtualMachinePackageLookup(
	idOrPermalink string,
) (lr *VirtualMachinePackage, f FieldName) {
	if strings.HasPrefix(idOrPermalink, "vmpkg_") {
		return &VirtualMachinePackage{ID: idOrPermalink}, IDField
	}

	return &VirtualMachinePackage{Permalink: idOrPermalink}, PermalinkField
}

func (s *VirtualMachinePackage) lookupReference() *VirtualMachinePackage {
	if s == nil {
		return nil
	}

	lr := &VirtualMachinePackage{ID: s.ID}
	if lr.ID == "" {
		lr.Permalink = s.Permalink
	}

	return lr
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
	idOrPermalink string,
) (*VirtualMachinePackage, *katapult.Response, error) {
	if _, f := NewVirtualMachinePackageLookup(idOrPermalink); f == IDField {
		return s.GetByID(ctx, idOrPermalink)
	}

	return s.GetByPermalink(ctx, idOrPermalink)
}

func (s *VirtualMachinePackagesClient) GetByID(
	ctx context.Context,
	id string,
) (*VirtualMachinePackage, *katapult.Response, error) {
	u := &url.URL{Path: fmt.Sprintf("virtual_machine_packages/%s", id)}
	body, resp, err := s.doRequest(ctx, "GET", u, nil)

	return body.VirtualMachinePackage, resp, err
}

func (s *VirtualMachinePackagesClient) GetByPermalink(
	ctx context.Context,
	permalink string,
) (*VirtualMachinePackage, *katapult.Response, error) {
	qs := url.Values{"virtual_machine_package[permalink]": []string{permalink}}
	u := &url.URL{Path: "virtual_machine_packages/_", RawQuery: qs.Encode()}

	body, resp, err := s.doRequest(ctx, "GET", u, nil)

	return body.VirtualMachinePackage, resp, err
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
