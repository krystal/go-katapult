package katapult

import (
	"context"
	"net/url"

	"github.com/augurysys/timestamp"
)

type VirtualMachineGroup struct {
	ID        string               `json:"id,omitempty"`
	Name      string               `json:"name,omitempty"`
	Segregate bool                 `json:"segregate,omitempty"`
	CreatedAt *timestamp.Timestamp `json:"created_at,omitempty"`
}

type VirtualMachineGroupsClient struct {
	client   *apiClient
	basePath *url.URL
}

type virtualMachineGroupsResponseBody struct {
	VirtualMachineGroups []*VirtualMachineGroup `json:"virtual_machine_groups,omitempty"`
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
