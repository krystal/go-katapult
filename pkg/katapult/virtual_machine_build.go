package katapult

import (
	"context"
	"fmt"
	"net/url"

	"github.com/augurysys/timestamp"
)

type VirtualMachineBuild struct {
	ID             string               `json:"id,omitempty"`
	SpecXML        string               `json:"spec_xml,omitempty"`
	State          State                `json:"state,omitempty"`
	VirtualMachine *VirtualMachine      `json:"virtual_machine,omitempty"`
	CreatedAt      *timestamp.Timestamp `json:"created_at,omitempty"`
}

type virtualMachineBuildsResponseBody struct {
	Task                *Task                `json:"task,omitempty"`
	Build               *VirtualMachineBuild `json:"build,omitempty"`
	VirtualMachineBuild *VirtualMachineBuild `json:"virtual_machine_build,omitempty"`
	Hostname            string               `json:"hostname,omitempty"`
}

type VirtualMachineBuildArguments struct {
	Organization        *Organization          `json:"organization,omitempty"`
	Zone                *Zone                  `json:"zone,omitempty"`
	DataCenter          *DataCenter            `json:"data_center,omitempty"`
	Package             *VirtualMachinePackage `json:"package,omitempty"`
	DiskTemplate        *DiskTemplate          `json:"disk_template,omitempty"`
	DiskTemplateOptions []*DiskTemplateOption  `json:"disk_template_options,omitempty"`
	Network             *Network               `json:"network,omitempty"`
	Hostname            string                 `json:"hostname,omitempty"`
}

func (
	s *VirtualMachineBuildArguments,
) forRequest() *VirtualMachineBuildArguments {
	if s == nil {
		return nil
	}

	args := *s
	args.Organization = s.Organization.LookupReference()
	args.Zone = s.Zone.LookupReference()
	args.DataCenter = s.DataCenter.LookupReference()
	args.Package = s.Package.LookupReference()
	args.DiskTemplate = s.DiskTemplate.LookupReference()
	args.Network = s.Network.LookupReference()

	return &args
}

type VirtualMachineBuildsClient struct {
	client   *apiClient
	basePath *url.URL
}

func newVirtualMachineBuildsClient(
	c *apiClient,
) *VirtualMachineBuildsClient {
	return &VirtualMachineBuildsClient{
		client:   c,
		basePath: &url.URL{Path: "/core/v1/"},
	}
}

func (s *VirtualMachineBuildsClient) Get(
	ctx context.Context,
	id string,
) (*VirtualMachineBuild, *Response, error) {
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
	args *VirtualMachineBuildArguments,
) (*VirtualMachineBuild, *Response, error) {
	u := &url.URL{Path: "organizations/_/virtual_machines/build"}
	body, resp, err := s.doRequest(ctx, "POST", u, args.forRequest())

	build := body.VirtualMachineBuild
	if build == nil {
		build = body.Build
	}

	return build, resp, err
}

func (s *VirtualMachineBuildsClient) doRequest(
	ctx context.Context,
	method string,
	u *url.URL,
	body interface{},
) (*virtualMachineBuildsResponseBody, *Response, error) {
	u = s.basePath.ResolveReference(u)
	respBody := &virtualMachineBuildsResponseBody{}
	resp := &Response{}

	req, err := s.client.NewRequestWithContext(ctx, method, u, body)
	if err == nil {
		resp, err = s.client.Do(req, respBody)
	}

	return respBody, resp, err
}
