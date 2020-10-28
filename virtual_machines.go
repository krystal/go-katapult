package katapult

import (
	"context"
	"fmt"
	"net/url"

	"github.com/augurysys/timestamp"
)

type VirtualMachinesService struct {
	client   *apiClient
	basePath *url.URL
}

func newVirtualMachinesService(
	c *apiClient,
) *VirtualMachinesService {
	return &VirtualMachinesService{
		client:   c,
		basePath: &url.URL{Path: "/core/v1/"},
	}
}

type VirtualMachine struct {
	ID                  string                 `json:"id,omitempty"`
	Name                string                 `json:"name,omitempty"`
	Hostname            string                 `json:"hostname,omitempty"`
	FQDN                string                 `json:"fqdn,omitempty"`
	CreatedAt           *timestamp.Timestamp   `json:"created_at,omitempty"`
	InitialRootPassword string                 `json:"initial_root_password,omitempty"`
	State               string                 `json:"state,omitempty"`
	Zone                *Zone                  `json:"zone,omitempty"`
	Organization        *Organization          `json:"organization,omitempty"`
	Group               *VirtualMachineGroup   `json:"group,omitempty"`
	Package             *VirtualMachinePackage `json:"package,omitempty"`
	AttachedISO         *ISO                   `json:"attached_iso,omitempty"`
	Tags                []*Tag                 `json:"tags,omitempty"`
	IPAddresses         []*IPAddress           `json:"ip_addresses,omitempty"`
}

type VirtualMachineGroup struct {
	ID        string               `json:"id,omitempty"`
	Name      string               `json:"name,omitempty"`
	Segregate bool                 `json:"segregate,omitempty"`
	CreatedAt *timestamp.Timestamp `json:"created_at,omitempty"`
}

type ISO struct {
	ID              string           `json:"id,omitempty"`
	Name            string           `json:"name,omitempty"`
	OperatingSystem *OperatingSystem `json:"operating_system,omitempty"`
}

type OperatingSystem struct {
	ID    string      `json:"id,omitempty"`
	Name  string      `json:"name,omitempty"`
	Badge *Attachment `json:"badge,omitempty"`
}

type Tag struct {
	ID        string               `json:"id,omitempty"`
	Name      string               `json:"name,omitempty"`
	Color     string               `json:"color,omitempty"`
	CreatedAt *timestamp.Timestamp `json:"created_at,omitempty"`
}

type IPAddress struct {
	ID              string `json:"id,omitempty"`
	Address         string `json:"address,omitempty"`
	ReverseDNS      string `json:"reverse_dns,omitempty"`
	VIP             bool   `json:"vip,omitempty"`
	AddressWithMask string `json:"address_with_mask,omitempty"`
}

type virtualMachinesResponseBody struct {
	Pagination      *Pagination       `json:"pagination,omitempty"`
	VirtualMachine  *VirtualMachine   `json:"virtual_machine,omitempty"`
	VirtualMachines []*VirtualMachine `json:"virtual_machines,omitempty"`
}

func (s VirtualMachinesService) List(
	ctx context.Context,
	orgID string,
	opts *ListOptions,
) ([]*VirtualMachine, *Response, error) {
	u := &url.URL{
		Path:     fmt.Sprintf("organizations/%s/virtual_machines", orgID),
		RawQuery: opts.Values().Encode(),
	}

	body, resp, err := s.doRequest(ctx, "GET", u, nil)
	resp.Pagination = body.Pagination

	return body.VirtualMachines, resp, err
}

func (s VirtualMachinesService) Get(
	ctx context.Context,
	id string,
) (*VirtualMachine, *Response, error) {
	u := &url.URL{Path: fmt.Sprintf("virtual_machines/%s", id)}
	body, resp, err := s.doRequest(ctx, "GET", u, nil)

	return body.VirtualMachine, resp, err
}

func (s *VirtualMachinesService) doRequest(
	ctx context.Context,
	method string,
	u *url.URL,
	body interface{},
) (*virtualMachinesResponseBody, *Response, error) {
	u = s.basePath.ResolveReference(u)
	respBody := &virtualMachinesResponseBody{}
	resp := &Response{}

	req, err := s.client.NewRequestWithContext(ctx, method, u, body)
	if err == nil {
		resp, err = s.client.Do(req, respBody)
	}

	return respBody, resp, err
}
