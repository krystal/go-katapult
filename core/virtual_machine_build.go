package core

import (
	"context"
	"fmt"
	"net/url"

	"github.com/krystal/go-katapult"
	"github.com/krystal/go-katapult/buildspec"

	"github.com/augurysys/timestamp"
)

type VirtualMachineBuild struct {
	ID             string                   `json:"id,omitempty"`
	SpecXML        string                   `json:"spec_xml,omitempty"`
	State          VirtualMachineBuildState `json:"state,omitempty"`
	VirtualMachine *VirtualMachine          `json:"virtual_machine,omitempty"`
	CreatedAt      *timestamp.Timestamp     `json:"created_at,omitempty"`
}

func (vmb *VirtualMachineBuild) Ref() VirtualMachineBuildRef {
	return VirtualMachineBuildRef{ID: vmb.ID}
}

type VirtualMachineBuildRef struct {
	ID string `json:"id,omitempty"`
}

type VirtualMachineBuildState string

const (
	VirtualMachineBuildDraft    VirtualMachineBuildState = "draft"
	VirtualMachineBuildFailed   VirtualMachineBuildState = "failed"
	VirtualMachineBuildPending  VirtualMachineBuildState = "pending"
	VirtualMachineBuildComplete VirtualMachineBuildState = "complete"
	VirtualMachineBuildBuilding VirtualMachineBuildState = "building"
)

type VirtualMachineBuildArguments struct {
	Zone                ZoneRef
	DataCenter          DataCenterRef
	Package             VirtualMachinePackageRef
	DiskTemplate        DiskTemplateRef
	DiskTemplateOptions []*DiskTemplateOption
	Network             NetworkRef
	Hostname            string
}

type virtualMachineBuildCreateRequest struct {
	Hostname            string                   `json:"hostname,omitempty"`
	Organization        OrganizationRef          `json:"organization"`
	Zone                ZoneRef                  `json:"zone,omitempty"`
	DataCenter          DataCenterRef            `json:"data_center"`
	Package             VirtualMachinePackageRef `json:"package,omitempty"`
	DiskTemplate        DiskTemplateRef          `json:"disk_template,omitempty"`
	DiskTemplateOptions []*DiskTemplateOption    `json:"disk_template_options,omitempty"`
	Network             NetworkRef               `json:"network"`
}

type virtualMachineBuildCreateFromSpecRequest struct {
	Organization OrganizationRef `json:"organization"`
	XML          string          `json:"xml,omitempty"`
}

type virtualMachineBuildsResponseBody struct {
	Task                *Task                `json:"task,omitempty"`
	Build               *VirtualMachineBuild `json:"build,omitempty"`
	VirtualMachineBuild *VirtualMachineBuild `json:"virtual_machine_build,omitempty"`
	Hostname            string               `json:"hostname,omitempty"`
}

type VirtualMachineBuildsClient struct {
	client   RequestMaker
	basePath *url.URL
}

func NewVirtualMachineBuildsClient(
	rm RequestMaker,
) *VirtualMachineBuildsClient {
	return &VirtualMachineBuildsClient{
		client:   rm,
		basePath: &url.URL{Path: "/core/v1/"},
	}
}

func (s *VirtualMachineBuildsClient) Get(
	ctx context.Context,
	ref VirtualMachineBuildRef,
) (*VirtualMachineBuild, *katapult.Response, error) {
	return s.GetByID(ctx, ref.ID)
}

func (s *VirtualMachineBuildsClient) GetByID(
	ctx context.Context,
	id string,
) (*VirtualMachineBuild, *katapult.Response, error) {
	u := &url.URL{Path: fmt.Sprintf("virtual_machines/builds/%s", id)}
	body, resp, err := s.doRequest(ctx, "GET", u, nil)

	build := body.VirtualMachineBuild
	if build == nil {
		build = body.Build
	}

	return build, resp, err
}

func (s *VirtualMachineBuildsClient) Create(
	ctx context.Context,
	org OrganizationRef,
	args *VirtualMachineBuildArguments,
) (*VirtualMachineBuild, *katapult.Response, error) {
	u := &url.URL{Path: "organizations/_/virtual_machines/build"}
	reqBody := &virtualMachineBuildCreateRequest{
		Hostname:            args.Hostname,
		Organization:        org,
		Zone:                args.Zone,
		DataCenter:          args.DataCenter,
		Package:             args.Package,
		DiskTemplate:        args.DiskTemplate,
		DiskTemplateOptions: args.DiskTemplateOptions,
		Network:             args.Network,
	}

	body, resp, err := s.doRequest(ctx, "POST", u, reqBody)

	return body.VirtualMachineBuild, resp, err
}

func (s *VirtualMachineBuildsClient) CreateFromSpec(
	ctx context.Context,
	org OrganizationRef,
	spec *buildspec.VirtualMachineSpec,
) (*VirtualMachineBuild, *katapult.Response, error) {
	specXML, _ := spec.XML()

	u := &url.URL{Path: "organizations/_/virtual_machines/build_from_spec"}
	reqBody := &virtualMachineBuildCreateFromSpecRequest{
		Organization: org,
		XML:          string(specXML),
	}

	body, resp, err := s.doRequest(ctx, "POST", u, reqBody)

	return body.VirtualMachineBuild, resp, err
}

func (s *VirtualMachineBuildsClient) CreateFromSpecXML(
	ctx context.Context,
	org OrganizationRef,
	specXML string,
) (*VirtualMachineBuild, *katapult.Response, error) {
	u := &url.URL{Path: "organizations/_/virtual_machines/build_from_spec"}
	reqBody := &virtualMachineBuildCreateFromSpecRequest{
		Organization: org,
		XML:          specXML,
	}

	body, resp, err := s.doRequest(ctx, "POST", u, reqBody)

	return body.VirtualMachineBuild, resp, err
}

func (s *VirtualMachineBuildsClient) doRequest(
	ctx context.Context,
	method string,
	u *url.URL,
	body interface{},
) (*virtualMachineBuildsResponseBody, *katapult.Response, error) {
	u = s.basePath.ResolveReference(u)
	respBody := &virtualMachineBuildsResponseBody{}
	resp := katapult.NewResponse(nil)

	req, err := s.client.NewRequestWithContext(ctx, method, u, body)
	if err == nil {
		resp, err = s.client.Do(req, respBody)
	}
	if respBody.VirtualMachineBuild == nil {
		respBody.VirtualMachineBuild = respBody.Build
	}

	return respBody, resp, err
}
