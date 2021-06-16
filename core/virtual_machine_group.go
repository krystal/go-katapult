package core

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/augurysys/timestamp"
	"github.com/krystal/go-katapult"
)

var (
	nullBytes                  = []byte("null")
	NullVirtualMachineGroupRef = &VirtualMachineGroupRef{null: true}
)

type VirtualMachineGroup struct {
	ID        string               `json:"id,omitempty"`
	Name      string               `json:"name,omitempty"`
	Segregate bool                 `json:"segregate,omitempty"`
	CreatedAt *timestamp.Timestamp `json:"created_at,omitempty"`
}

func (s *VirtualMachineGroup) Ref() VirtualMachineGroupRef {
	return VirtualMachineGroupRef{ID: s.ID}
}

func (s *VirtualMachineGroupRef) UnmarshalJSON(b []byte) error {
	type alias VirtualMachineGroupRef

	if bytes.Equal(b, nullBytes) {
		*s = VirtualMachineGroupRef{null: true}

		return nil
	}

	return json.Unmarshal(b, (*alias)(s))
}

func (s *VirtualMachineGroupRef) MarshalJSON() ([]byte, error) {
	type alias VirtualMachineGroupRef

	if s.null {
		return nullBytes, nil
	}

	return json.Marshal((*alias)(s))
}

type VirtualMachineGroupRef struct {
	ID   string `json:"id,omitempty"`
	null bool
}

func (s VirtualMachineGroupRef) queryValues() *url.Values {
	v := &url.Values{}
	v.Set("virtual_machine_group[id]", s.ID)

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
	Organization OrganizationRef                     `json:"organization"`
	Properties   *VirtualMachineGroupCreateArguments `json:"properties,omitempty"`
}

type virtualMachineGroupUpdateRequest struct {
	VirtualMachineGroup VirtualMachineGroupRef              `json:"virtual_machine_group,omitempty"`
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
	org OrganizationRef,
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
	ref VirtualMachineGroupRef,
) (*VirtualMachineGroup, *katapult.Response, error) {
	return s.GetByID(ctx, ref.ID)
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
	org OrganizationRef,
	args *VirtualMachineGroupCreateArguments,
) (*VirtualMachineGroup, *katapult.Response, error) {
	u := &url.URL{Path: "organizations/_/virtual_machine_groups"}
	reqBody := &virtualMachineGroupCreateRequest{
		Organization: org,
		Properties:   args,
	}

	body, resp, err := s.doRequest(ctx, "POST", u, reqBody)

	return body.VirtualMachineGroup, resp, err
}

func (s *VirtualMachineGroupsClient) Update(
	ctx context.Context,
	ref VirtualMachineGroupRef,
	args *VirtualMachineGroupUpdateArguments,
) (*VirtualMachineGroup, *katapult.Response, error) {
	u := &url.URL{Path: "virtual_machine_groups/_"}
	reqBody := &virtualMachineGroupUpdateRequest{
		VirtualMachineGroup: ref,
		Properties:          args,
	}

	body, resp, err := s.doRequest(ctx, "PATCH", u, reqBody)

	return body.VirtualMachineGroup, resp, err
}

func (s *VirtualMachineGroupsClient) Delete(
	ctx context.Context,
	group VirtualMachineGroupRef,
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

	req := katapult.NewRequest(method, u, body)
	resp, err := s.client.Do(ctx, req, respBody)
	if resp == nil {
		resp = katapult.NewResponse(nil)
	}

	return respBody, resp, handleResponseError(err)
}
