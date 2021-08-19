package core

import (
	"context"
	"net/url"

	"github.com/augurysys/timestamp"
	"github.com/krystal/go-katapult"
)

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

func (s *VirtualMachine) Ref() VirtualMachineRef {
	return VirtualMachineRef{ID: s.ID}
}

type VirtualMachineRef struct {
	ID   string `json:"id,omitempty"`
	FQDN string `json:"fqdn,omitempty"`
}

func (vmr VirtualMachineRef) queryValues() *url.Values {
	v := &url.Values{}

	switch {
	case vmr.ID != "":
		v.Set("virtual_machine[id]", vmr.ID)
	case vmr.FQDN != "":
		v.Set("virtual_machine[fqdn]", vmr.FQDN)
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

type VirtualMachineUpdateArguments struct {
	Name        string                  `json:"name,omitempty"`
	Hostname    string                  `json:"hostname,omitempty"`
	Description string                  `json:"description,omitempty"`
	TagNames    *[]string               `json:"tag_names,omitempty"`
	Group       *VirtualMachineGroupRef `json:"group,omitempty"`
}

type virtualMachinesResponseBody struct {
	Pagination      *katapult.Pagination `json:"pagination,omitempty"`
	Task            *Task                `json:"task,omitempty"`
	TrashObject     *TrashObject         `json:"trash_object,omitempty"`
	VirtualMachine  *VirtualMachine      `json:"virtual_machine,omitempty"`
	VirtualMachines []*VirtualMachine    `json:"virtual_machines,omitempty"`
}

type virtualMachineChangePackageRequest struct {
	VirtualMachine VirtualMachineRef        `json:"virtual_machine,omitempty"`
	Package        VirtualMachinePackageRef `json:"virtual_machine_package,omitempty"`
}

type virtualMachineUpdateRequest struct {
	VirtualMachine VirtualMachineRef              `json:"virtual_machine,omitempty"`
	Properties     *VirtualMachineUpdateArguments `json:"properties,omitempty"`
}

type VirtualMachinesClient struct {
	client   RequestMaker
	basePath *url.URL
}

func NewVirtualMachinesClient(
	rm RequestMaker,
) *VirtualMachinesClient {
	return &VirtualMachinesClient{
		client:   rm,
		basePath: &url.URL{Path: "/core/v1/"},
	}
}

func (s *VirtualMachinesClient) List(
	ctx context.Context,
	org OrganizationRef,
	opts *ListOptions,
	reqOpts ...katapult.RequestOption,
) ([]*VirtualMachine, *katapult.Response, error) {
	qs := queryValues(org, opts)
	u := &url.URL{
		Path:     "organizations/_/virtual_machines",
		RawQuery: qs.Encode(),
	}

	body, resp, err := s.doRequest(ctx, "GET", u, nil, reqOpts...)
	resp.Pagination = body.Pagination

	return body.VirtualMachines, resp, err
}

func (s *VirtualMachinesClient) Get(
	ctx context.Context,
	ref VirtualMachineRef,
	reqOpts ...katapult.RequestOption,
) (*VirtualMachine, *katapult.Response, error) {
	u := &url.URL{
		Path:     "virtual_machines/_",
		RawQuery: ref.queryValues().Encode(),
	}

	body, resp, err := s.doRequest(ctx, "GET", u, nil, reqOpts...)

	return body.VirtualMachine, resp, err
}

func (s *VirtualMachinesClient) GetByID(
	ctx context.Context,
	id string,
	reqOpts ...katapult.RequestOption,
) (*VirtualMachine, *katapult.Response, error) {
	return s.Get(ctx, VirtualMachineRef{ID: id}, reqOpts...)
}

func (s *VirtualMachinesClient) GetByFQDN(
	ctx context.Context,
	fqdn string,
	reqOpts ...katapult.RequestOption,
) (*VirtualMachine, *katapult.Response, error) {
	return s.Get(ctx, VirtualMachineRef{FQDN: fqdn}, reqOpts...)
}

func (s *VirtualMachinesClient) ChangePackage(
	ctx context.Context,
	ref VirtualMachineRef,
	pkg VirtualMachinePackageRef,
	reqOpts ...katapult.RequestOption,
) (*Task, *katapult.Response, error) {
	u := &url.URL{Path: "virtual_machines/_/package"}
	reqBody := &virtualMachineChangePackageRequest{
		VirtualMachine: ref,
		Package:        pkg,
	}
	body, resp, err := s.doRequest(ctx, "PUT", u, reqBody, reqOpts...)

	return body.Task, resp, err
}

func (s *VirtualMachinesClient) Update(
	ctx context.Context,
	ref VirtualMachineRef,
	args *VirtualMachineUpdateArguments,
	reqOpts ...katapult.RequestOption,
) (*VirtualMachine, *katapult.Response, error) {
	u := &url.URL{Path: "virtual_machines/_"}
	reqBody := &virtualMachineUpdateRequest{
		VirtualMachine: ref,
		Properties:     args,
	}
	body, resp, err := s.doRequest(ctx, "PATCH", u, reqBody, reqOpts...)

	return body.VirtualMachine, resp, err
}

func (s *VirtualMachinesClient) Delete(
	ctx context.Context,
	ref VirtualMachineRef,
	reqOpts ...katapult.RequestOption,
) (*TrashObject, *katapult.Response, error) {
	u := &url.URL{
		Path:     "virtual_machines/_",
		RawQuery: ref.queryValues().Encode(),
	}
	body, resp, err := s.doRequest(ctx, "DELETE", u, nil, reqOpts...)

	return body.TrashObject, resp, err
}

func (s *VirtualMachinesClient) Start(
	ctx context.Context,
	ref VirtualMachineRef,
	reqOpts ...katapult.RequestOption,
) (*Task, *katapult.Response, error) {
	u := &url.URL{
		Path:     "virtual_machines/_/start",
		RawQuery: ref.queryValues().Encode(),
	}

	body, resp, err := s.doRequest(ctx, "POST", u, nil, reqOpts...)

	return body.Task, resp, err
}

func (s *VirtualMachinesClient) Stop(
	ctx context.Context,
	ref VirtualMachineRef,
	reqOpts ...katapult.RequestOption,
) (*Task, *katapult.Response, error) {
	u := &url.URL{
		Path:     "virtual_machines/_/stop",
		RawQuery: ref.queryValues().Encode(),
	}

	body, resp, err := s.doRequest(ctx, "POST", u, nil, reqOpts...)

	return body.Task, resp, err
}

func (s *VirtualMachinesClient) Shutdown(
	ctx context.Context,
	ref VirtualMachineRef,
	reqOpts ...katapult.RequestOption,
) (*Task, *katapult.Response, error) {
	u := &url.URL{
		Path:     "virtual_machines/_/shutdown",
		RawQuery: ref.queryValues().Encode(),
	}

	body, resp, err := s.doRequest(ctx, "POST", u, nil, reqOpts...)

	return body.Task, resp, err
}

func (s *VirtualMachinesClient) Reset(
	ctx context.Context,
	ref VirtualMachineRef,
	reqOpts ...katapult.RequestOption,
) (*Task, *katapult.Response, error) {
	u := &url.URL{
		Path:     "virtual_machines/_/reset",
		RawQuery: ref.queryValues().Encode(),
	}

	body, resp, err := s.doRequest(ctx, "POST", u, nil, reqOpts...)

	return body.Task, resp, err
}

func (s *VirtualMachinesClient) doRequest(
	ctx context.Context,
	method string,
	u *url.URL,
	body interface{},
	reqOpts ...katapult.RequestOption,
) (*virtualMachinesResponseBody, *katapult.Response, error) {
	u = s.basePath.ResolveReference(u)
	respBody := &virtualMachinesResponseBody{}

	req := katapult.NewRequest(method, u, body, reqOpts...)
	resp, err := s.client.Do(ctx, req, respBody)
	if resp == nil {
		resp = katapult.NewResponse(nil)
	}

	return respBody, resp, handleResponseError(err)
}
