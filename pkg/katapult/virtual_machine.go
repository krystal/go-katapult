package katapult

import (
	"context"
	"net/url"
	"strings"

	"github.com/augurysys/timestamp"
)

const virtualMachineIDPrefix = "vm_"

type VirtualMachine struct {
	ID                  string                 `json:"id,omitempty"`
	Name                string                 `json:"name,omitempty"`
	Hostname            string                 `json:"hostname,omitempty"`
	FQDN                string                 `json:"fqdn,omitempty"`
	Description         string                 `json:"description,omitempty"`
	CreatedAt           *timestamp.Timestamp   `json:"created_at,omitempty"`
	InitialRootPassword string                 `json:"initial_root_password,omitempty"`
	State               VirtualMachineState    `json:"state,omitempty"`
	Zone                *Zone                  `json:"zone,omitempty"`
	Organization        *Organization          `json:"organization,omitempty"`
	Group               *VirtualMachineGroup   `json:"group,omitempty"`
	Package             *VirtualMachinePackage `json:"package,omitempty"`
	AttachedISO         *ISO                   `json:"attached_iso,omitempty"`
	Tags                []*Tag                 `json:"tags,omitempty"`
	TagNames            []string               `json:"tag_names,omitempty"`
	IPAddresses         []*IPAddress           `json:"ip_addresses,omitempty"`
}

func (s *VirtualMachine) lookupReference() *VirtualMachine {
	if s == nil {
		return nil
	}

	lr := &VirtualMachine{ID: s.ID}
	if lr.ID == "" {
		lr.FQDN = s.FQDN
	}

	return lr
}

func (s *VirtualMachine) queryValues() *url.Values {
	v := &url.Values{}

	if s != nil {
		switch {
		case s.ID != "":
			v.Set("virtual_machine[id]", s.ID)
		case s.FQDN != "":
			v.Set("virtual_machine[fqdn]", s.FQDN)
		}
	}

	return v
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

type VirtualMachineUpdateArguments struct {
	Name        string    `json:"name,omitempty"`
	Hostname    string    `json:"hostname,omitempty"`
	Description string    `json:"description,omitempty"`
	TagNames    *[]string `json:"tag_names,omitempty"`
}

type virtualMachinesResponseBody struct {
	Pagination      *Pagination       `json:"pagination,omitempty"`
	Task            *Task             `json:"task,omitempty"`
	TrashObject     *TrashObject      `json:"trash_object,omitempty"`
	VirtualMachine  *VirtualMachine   `json:"virtual_machine,omitempty"`
	VirtualMachines []*VirtualMachine `json:"virtual_machines,omitempty"`
}

type virtualMachineChangePackageRequest struct {
	VirtualMachine *VirtualMachine        `json:"virtual_machine,omitempty"`
	Package        *VirtualMachinePackage `json:"virtual_machine_package,omitempty"`
}

type virtualMachineUpdateRequest struct {
	VirtualMachine *VirtualMachine                `json:"virtual_machine,omitempty"`
	Properties     *VirtualMachineUpdateArguments `json:"properties,omitempty"`
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

func (s *VirtualMachinesClient) List(
	ctx context.Context,
	org *Organization,
	opts *ListOptions,
) ([]*VirtualMachine, *Response, error) {
	qs := queryValues(org, opts)
	u := &url.URL{
		Path:     "organizations/_/virtual_machines",
		RawQuery: qs.Encode(),
	}

	body, resp, err := s.doRequest(ctx, "GET", u, nil)
	resp.Pagination = body.Pagination

	return body.VirtualMachines, resp, err
}

func (s *VirtualMachinesClient) Get(
	ctx context.Context,
	idOrFQDN string,
) (*VirtualMachine, *Response, error) {
	if strings.HasPrefix(idOrFQDN, virtualMachineIDPrefix) {
		return s.GetByID(ctx, idOrFQDN)
	}

	return s.GetByFQDN(ctx, idOrFQDN)
}

func (s *VirtualMachinesClient) GetByID(
	ctx context.Context,
	id string,
) (*VirtualMachine, *Response, error) {
	return s.get(ctx, &VirtualMachine{ID: id})
}

func (s *VirtualMachinesClient) GetByFQDN(
	ctx context.Context,
	fqdn string,
) (*VirtualMachine, *Response, error) {
	return s.get(ctx, &VirtualMachine{FQDN: fqdn})
}

func (s *VirtualMachinesClient) get(
	ctx context.Context,
	vm *VirtualMachine,
) (*VirtualMachine, *Response, error) {
	u := &url.URL{
		Path:     "virtual_machines/_",
		RawQuery: vm.queryValues().Encode(),
	}

	body, resp, err := s.doRequest(ctx, "GET", u, nil)

	return body.VirtualMachine, resp, err
}

func (s *VirtualMachinesClient) ChangePackage(
	ctx context.Context,
	vm *VirtualMachine,
	pkg *VirtualMachinePackage,
) (*Task, *Response, error) {
	u := &url.URL{Path: "virtual_machines/_/package"}
	reqBody := &virtualMachineChangePackageRequest{
		VirtualMachine: vm.lookupReference(),
		Package:        pkg.lookupReference(),
	}
	body, resp, err := s.doRequest(ctx, "PUT", u, reqBody)

	return body.Task, resp, err
}

func (s *VirtualMachinesClient) Update(
	ctx context.Context,
	vm *VirtualMachine,
	args *VirtualMachineUpdateArguments,
) (*VirtualMachine, *Response, error) {
	u := &url.URL{Path: "virtual_machines/_"}
	reqBody := &virtualMachineUpdateRequest{
		VirtualMachine: vm.lookupReference(),
		Properties:     args,
	}
	body, resp, err := s.doRequest(ctx, "PATCH", u, reqBody)

	return body.VirtualMachine, resp, err
}

func (s *VirtualMachinesClient) Delete(
	ctx context.Context,
	vm *VirtualMachine,
) (*TrashObject, *Response, error) {
	u := &url.URL{
		Path:     "virtual_machines/_",
		RawQuery: vm.queryValues().Encode(),
	}
	body, resp, err := s.doRequest(ctx, "DELETE", u, nil)

	return body.TrashObject, resp, err
}

func (s *VirtualMachinesClient) Start(
	ctx context.Context,
	vm *VirtualMachine,
) (*Task, *Response, error) {
	u := &url.URL{
		Path:     "virtual_machines/_/start",
		RawQuery: vm.queryValues().Encode(),
	}

	body, resp, err := s.doRequest(ctx, "POST", u, nil)

	return body.Task, resp, err
}

func (s *VirtualMachinesClient) Stop(
	ctx context.Context,
	vm *VirtualMachine,
) (*Task, *Response, error) {
	u := &url.URL{
		Path:     "virtual_machines/_/stop",
		RawQuery: vm.queryValues().Encode(),
	}

	body, resp, err := s.doRequest(ctx, "POST", u, nil)

	return body.Task, resp, err
}

func (s *VirtualMachinesClient) Shutdown(
	ctx context.Context,
	vm *VirtualMachine,
) (*Task, *Response, error) {
	u := &url.URL{
		Path:     "virtual_machines/_/shutdown",
		RawQuery: vm.queryValues().Encode(),
	}

	body, resp, err := s.doRequest(ctx, "POST", u, nil)

	return body.Task, resp, err
}

func (s *VirtualMachinesClient) Reset(
	ctx context.Context,
	vm *VirtualMachine,
) (*Task, *Response, error) {
	u := &url.URL{
		Path:     "virtual_machines/_/reset",
		RawQuery: vm.queryValues().Encode(),
	}

	body, resp, err := s.doRequest(ctx, "POST", u, nil)

	return body.Task, resp, err
}

func (s *VirtualMachinesClient) doRequest(
	ctx context.Context,
	method string,
	u *url.URL,
	body interface{},
) (*virtualMachinesResponseBody, *Response, error) {
	u = s.basePath.ResolveReference(u)
	respBody := &virtualMachinesResponseBody{}
	resp := newResponse(nil)

	req, err := s.client.NewRequestWithContext(ctx, method, u, body)
	if err == nil {
		resp, err = s.client.Do(req, respBody)
	}

	return respBody, resp, err
}
