package katapult

import (
	"context"
	"fmt"
	"net/url"

	"github.com/augurysys/timestamp"
)

type VirtualMachine struct {
	ID                  string                 `json:"id,omitempty"`
	Name                string                 `json:"name,omitempty"`
	Hostname            string                 `json:"hostname,omitempty"`
	FQDN                string                 `json:"fqdn,omitempty"`
	CreatedAt           *timestamp.Timestamp   `json:"created_at,omitempty"`
	InitialRootPassword string                 `json:"initial_root_password,omitempty"`
	State               VirtualMachineState    `json:"state,omitempty"`
	Zone                *Zone                  `json:"zone,omitempty"`
	Organization        *Organization          `json:"organization,omitempty"`
	Group               *VirtualMachineGroup   `json:"group,omitempty"`
	Package             *VirtualMachinePackage `json:"package,omitempty"`
	AttachedISO         *ISO                   `json:"attached_iso,omitempty"`
	Tags                []*Tag                 `json:"tags,omitempty"`
	IPAddresses         []*IPAddress           `json:"ip_addresses,omitempty"`
}

// LookupReference returns a new *VirtualMachine stripped down to just ID or
// FQDN fields, making it suitable for endpoints which require a reference to a
// Virtual Machine by ID or FQDN.
func (s *VirtualMachine) LookupReference() *VirtualMachine {
	if s == nil {
		return nil
	}

	lr := &VirtualMachine{ID: s.ID}
	if lr.ID == "" {
		lr.FQDN = s.FQDN
	}

	return lr
}

type VirtualMachineState string

const (
	VirtualMachineStopped      VirtualMachineState = "stopped"
	VirtualMachineFailed       VirtualMachineState = "failed"
	VirtualMachineStarted      VirtualMachineState = "started"
	VirtualMachineStarting     VirtualMachineState = "starting"
	VirtualMachineResetting    VirtualMachineState = "resetting"
	VirtualMachineMigrating    VirtualMachineState = "migrating"
	VirtualMachineStopping     VirtualMachineState = "stopping"
	VirtualMachineShuttingDown VirtualMachineState = "shutting_down"
	VirtualMachineOrphaned     VirtualMachineState = "orphaned"
)

type VirtualMachineGroup struct {
	ID        string               `json:"id,omitempty"`
	Name      string               `json:"name,omitempty"`
	Segregate bool                 `json:"segregate,omitempty"`
	CreatedAt *timestamp.Timestamp `json:"created_at,omitempty"`
}

type virtualMachinesResponseBody struct {
	Pagination      *Pagination       `json:"pagination,omitempty"`
	Task            *Task             `json:"task,omitempty"`
	TrashObject     *TrashObject      `json:"trash_object,omitempty"`
	VirtualMachine  *VirtualMachine   `json:"virtual_machine,omitempty"`
	VirtualMachines []*VirtualMachine `json:"virtual_machines,omitempty"`
}

type virtualMachineChangePackageRequestBody struct {
	VirtualMachine *VirtualMachine        `json:"virtual_machine,omitempty"`
	Package        *VirtualMachinePackage `json:"virtual_machine_package,omitempty"`
}

type VirtualMachinesClient struct {
	client   *apiClient
	basePath *url.URL
}

func newVirtualMachinesClient(
	c *apiClient,
) *VirtualMachinesClient {
	return &VirtualMachinesClient{
		client:   c,
		basePath: &url.URL{Path: "/core/v1/"},
	}
}

func (s VirtualMachinesClient) List(
	ctx context.Context,
	org *Organization,
	opts *ListOptions,
) ([]*VirtualMachine, *Response, error) {
	if org == nil {
		org = &Organization{ID: "_"}
	}

	u := &url.URL{
		Path:     fmt.Sprintf("organizations/%s/virtual_machines", org.ID),
		RawQuery: opts.Values().Encode(),
	}

	body, resp, err := s.doRequest(ctx, "GET", u, nil)
	resp.Pagination = body.Pagination

	return body.VirtualMachines, resp, err
}

func (s VirtualMachinesClient) Get(
	ctx context.Context,
	id string,
) (*VirtualMachine, *Response, error) {
	u := &url.URL{Path: fmt.Sprintf("virtual_machines/%s", id)}
	body, resp, err := s.doRequest(ctx, "GET", u, nil)

	return body.VirtualMachine, resp, err
}

func (s VirtualMachinesClient) GetByFQDN(
	ctx context.Context,
	fqdn string,
) (*VirtualMachine, *Response, error) {
	qs := url.Values{"virtual_machine[fqdn]": []string{fqdn}}
	u := &url.URL{Path: "virtual_machines/_", RawQuery: qs.Encode()}

	body, resp, err := s.doRequest(ctx, "GET", u, nil)

	return body.VirtualMachine, resp, err
}

func (s *VirtualMachinesClient) ChangePackage(
	ctx context.Context,
	vm *VirtualMachine,
	pkg *VirtualMachinePackage,
) (*Task, *Response, error) {
	u := &url.URL{Path: "virtual_machines/_/package"}
	reqBody := &virtualMachineChangePackageRequestBody{
		VirtualMachine: vm.LookupReference(),
		Package:        pkg.LookupReference(),
	}
	body, resp, err := s.doRequest(ctx, "PUT", u, reqBody)

	return body.Task, resp, err
}

func (s *VirtualMachinesClient) Delete(
	ctx context.Context,
	vm *VirtualMachine,
) (*TrashObject, *Response, error) {
	if vm == nil {
		vm = &VirtualMachine{ID: "_"}
	}

	u := &url.URL{Path: fmt.Sprintf("virtual_machines/%s", vm.ID)}
	body, resp, err := s.doRequest(ctx, "DELETE", u, nil)

	return body.TrashObject, resp, err
}

func (s *VirtualMachinesClient) doRequest(
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
