package core

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/krystal/go-katapult"
	"net/url"

	"github.com/augurysys/timestamp"
)

var (
	nullBytes               = []byte("null")
	NullVirtualMachineGroup = &VirtualMachineGroup{null: true}
)

type VirtualMachineGroup struct {
	ID        string               `json:"id,omitempty"`
	Name      string               `json:"name,omitempty"`
	Segregate bool                 `json:"segregate,omitempty"`
	CreatedAt *timestamp.Timestamp `json:"created_at,omitempty"`
	null      bool
}

func (s *VirtualMachineGroup) UnmarshalJSON(b []byte) error {
	type alias VirtualMachineGroup

	if bytes.Equal(b, nullBytes) {
		*s = VirtualMachineGroup{null: true}

		return nil
	}

	return json.Unmarshal(b, (*alias)(s))
}

func (s *VirtualMachineGroup) MarshalJSON() ([]byte, error) {
	type alias VirtualMachineGroup

	if s.null {
		return nullBytes, nil
	}

	return json.Marshal((*alias)(s))
}

func (s *VirtualMachineGroup) lookupReference() *VirtualMachineGroup {
	if s == nil {
		return nil
	}

	return &VirtualMachineGroup{ID: s.ID}
}

func (s *VirtualMachineGroup) queryValues() *url.Values {
	v := &url.Values{}

	if s != nil {
		v.Set("virtual_machine_group[id]", s.ID)
	}

	return v
}

type VirtualMachineGroupCreateArguments struct {
	Name      string `json:"name,omitempty"`
	Segregate *bool  `json:"segregate,omitempty"`
}

type VirtualMachineGroupUpdateArguments struct {
	Name      string `json:"name,omitempty"`
	Segregate *bool  `json:"segregate,omitempty"`
}

type virtualMachineGroupCreateRequest struct {
	Organization *Organization                       `json:"organization,omitempty"`
	Properties   *VirtualMachineGroupCreateArguments `json:"properties,omitempty"`
}

type virtualMachineGroupUpdateRequest struct {
	VirtualMachineGroup *VirtualMachineGroup                `json:"virtual_machine_group,omitempty"`
	Properties          *VirtualMachineGroupUpdateArguments `json:"properties,omitempty"`
}

type virtualMachineGroupsResponseBody struct {
	VirtualMachineGroups []*VirtualMachineGroup `json:"virtual_machine_groups,omitempty"`
	VirtualMachineGroup  *VirtualMachineGroup   `json:"virtual_machine_group,omitempty"`
}

type VirtualMachineGroupsClient struct {
	client   RequestMaker
	basePath *url.URL
}

func NewVirtualMachineGroupsClient(
	rm RequestMaker,
) *VirtualMachineGroupsClient {
	return &VirtualMachineGroupsClient{
		client:   rm,
		basePath: &url.URL{Path: "/core/v1/"},
	}
}

func (s *VirtualMachineGroupsClient) List(
	ctx context.Context,
	org *Organization,
) ([]*VirtualMachineGroup, *katapult.Response, error) {
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
) (*VirtualMachineGroup, *katapult.Response, error) {
	return s.GetByID(ctx, id)
}

func (s *VirtualMachineGroupsClient) GetByID(
	ctx context.Context,
	id string,
) (*VirtualMachineGroup, *katapult.Response, error) {
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
) (*VirtualMachineGroup, *katapult.Response, error) {
	u := &url.URL{Path: "organizations/_/virtual_machine_groups"}
	reqBody := &virtualMachineGroupCreateRequest{
		Organization: org.lookupReference(),
		Properties:   args,
	}

	body, resp, err := s.doRequest(ctx, "POST", u, reqBody)

	return body.VirtualMachineGroup, resp, err
}

func (s *VirtualMachineGroupsClient) Update(
	ctx context.Context,
	group *VirtualMachineGroup,
	args *VirtualMachineGroupUpdateArguments,
) (*VirtualMachineGroup, *katapult.Response, error) {
	u := &url.URL{Path: "virtual_machine_groups/_"}
	reqBody := &virtualMachineGroupUpdateRequest{
		VirtualMachineGroup: group.lookupReference(),
		Properties:          args,
	}

	body, resp, err := s.doRequest(ctx, "PATCH", u, reqBody)

	return body.VirtualMachineGroup, resp, err
}

func (s *VirtualMachineGroupsClient) Delete(
	ctx context.Context,
	group *VirtualMachineGroup,
) (*katapult.Response, error) {
	qs := queryValues(group)
	u := &url.URL{Path: "virtual_machine_groups/_", RawQuery: qs.Encode()}

	_, resp, err := s.doRequest(ctx, "DELETE", u, nil)

	return resp, err
}

func (s *VirtualMachineGroupsClient) doRequest(
	ctx context.Context,
	method string,
	u *url.URL,
	body interface{},
) (*virtualMachineGroupsResponseBody, *katapult.Response, error) {
	u = s.basePath.ResolveReference(u)
	respBody := &virtualMachineGroupsResponseBody{}
	resp := katapult.NewResponse(nil)

	req, err := s.client.NewRequestWithContext(ctx, method, u, body)
	if err == nil {
		resp, err = s.client.Do(req, respBody)
	}

	return respBody, resp, err
}
