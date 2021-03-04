package katapult

import (
	"context"
	"fmt"
	"net/url"

	"github.com/augurysys/timestamp"
)

type VirtualMachineGroup struct {
	ID        string               `json:"id,omitempty"`
	Name      string               `json:"name,omitempty"`
	Segregate bool                 `json:"segregate,omitempty"`
	CreatedAt *timestamp.Timestamp `json:"created_at,omitempty"`
}

func (s *VirtualMachineGroup) lookupReference() *VirtualMachineGroup {
	if s == nil {
		return nil
	}

	return &VirtualMachineGroup{ID: s.ID}
}

type VirtualMachineGroupCreateArguments struct {
	Name      string
	Segregate *bool
}

type VirtualMachineGroupUpdateArguments struct {
	Name      string
	Segregate *bool
}

type virtualMachineGroupCreateRequest struct {
	Organization *Organization `json:"organization,omitempty"`
	Name         string        `json:"name,omitempty"`
	Segregate    *bool         `json:"segregate,omitempty"`
}

type virtualMachineGroupUpdateRequest struct {
	VirtualMachineGroup *VirtualMachineGroup `json:"virtual_machine_group,omitempty"`
	Name                string               `json:"name,omitempty"`
	Segregate           *bool                `json:"segregate,omitempty"`
}

type virtualMachineGroupsResponseBody struct {
	VirtualMachineGroups []*VirtualMachineGroup `json:"virtual_machine_groups,omitempty"`
	VirtualMachineGroup  *VirtualMachineGroup   `json:"virtual_machine_group,omitempty"`
}
type VirtualMachineGroupsClient struct {
	client   *apiClient
	basePath *url.URL
}

func newVirtualMachineGroupsClient(
	c *apiClient,
) *VirtualMachineGroupsClient {
	return &VirtualMachineGroupsClient{
		client:   c,
		basePath: &url.URL{Path: "/core/v1/"},
	}
}

func (s *VirtualMachineGroupsClient) List(
	ctx context.Context,
	org *Organization,
) ([]*VirtualMachineGroup, *Response, error) {
	qs := queryValues(org)
	u := &url.URL{
		Path:     "organizations/_/virtual_machine_groups",
		RawQuery: qs.Encode(),
	}

	body, resp, err := s.doRequest(ctx, "GET", u, nil)

	return body.VirtualMachineGroups, resp, err
}

func (s *VirtualMachineGroupsClient) Get(
	ctx context.Context,
	id string,
) (*VirtualMachineGroup, *Response, error) {
	return s.GetByID(ctx, id)
}

func (s *VirtualMachineGroupsClient) GetByID(
	ctx context.Context,
	id string,
) (*VirtualMachineGroup, *Response, error) {
	u := &url.URL{
		Path: fmt.Sprintf("virtual_machine_groups/%s", id),
	}

	body, resp, err := s.doRequest(ctx, "GET", u, nil)

	return body.VirtualMachineGroup, resp, err
}

func (s *VirtualMachineGroupsClient) Create(
	ctx context.Context,
	org *Organization,
	args *VirtualMachineGroupCreateArguments,
) (*VirtualMachineGroup, *Response, error) {
	u := &url.URL{Path: "organizations/_/virtual_machine_groups"}
	reqBody := &virtualMachineGroupCreateRequest{
		Organization: org.lookupReference(),
	}

	if args != nil {
		reqBody.Name = args.Name
		reqBody.Segregate = args.Segregate
	}

	body, resp, err := s.doRequest(ctx, "POST", u, reqBody)

	return body.VirtualMachineGroup, resp, err
}

func (s *VirtualMachineGroupsClient) Update(
	ctx context.Context,
	group *VirtualMachineGroup,
	args *VirtualMachineGroupUpdateArguments,
) (*VirtualMachineGroup, *Response, error) {
	u := &url.URL{Path: "virtual_machine_groups/_"}
	reqBody := &virtualMachineGroupUpdateRequest{
		VirtualMachineGroup: group.lookupReference(),
	}

	if args != nil {
		reqBody.Name = args.Name
		reqBody.Segregate = args.Segregate
	}

	body, resp, err := s.doRequest(ctx, "PATCH", u, reqBody)

	return body.VirtualMachineGroup, resp, err
}

func (s *VirtualMachineGroupsClient) doRequest(
	ctx context.Context,
	method string,
	u *url.URL,
	body interface{},
) (*virtualMachineGroupsResponseBody, *Response, error) {
	u = s.basePath.ResolveReference(u)
	respBody := &virtualMachineGroupsResponseBody{}
	resp := newResponse(nil)

	req, err := s.client.NewRequestWithContext(ctx, method, u, body)
	if err == nil {
		resp, err = s.client.Do(req, respBody)
	}

	return respBody, resp, err
}
