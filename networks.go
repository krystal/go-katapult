package katapult

import (
	"context"
	"fmt"
	"net/url"
)

type NetworksService struct {
	client   *apiClient
	basePath *url.URL
}

func newNetworksService(c *apiClient) *NetworksService {
	return &NetworksService{
		client:   c,
		basePath: &url.URL{Path: "/core/v1/"},
	}
}

type Network struct {
	ID         string      `json:"id,omitempty"`
	Name       string      `json:"name,omitempty"`
	DataCenter *DataCenter `json:"data_center,omitempty"`
}

type VirtualNetwork struct {
	ID         string      `json:"id,omitempty"`
	Name       string      `json:"name,omitempty"`
	DataCenter *DataCenter `json:"data_center,omitempty"`
}

type networksResponseBody struct {
	Networks        []*Network        `json:"networks,omitempty"`
	VirtualNetworks []*VirtualNetwork `json:"virtual_networks,omitempty"`
}

func (s *NetworksService) List(
	ctx context.Context,
	orgID string,
) ([]*Network, []*VirtualNetwork, *Response, error) {
	u := &url.URL{
		Path: fmt.Sprintf("organizations/%s/available_networks", orgID),
	}

	body, resp, err := s.doRequest(ctx, "GET", u, nil)

	return body.Networks, body.VirtualNetworks, resp, err
}

func (s *NetworksService) doRequest(
	ctx context.Context,
	method string,
	u *url.URL,
	body interface{},
) (*networksResponseBody, *Response, error) {
	u = s.basePath.ResolveReference(u)
	respBody := &networksResponseBody{}
	resp := &Response{}

	req, err := s.client.NewRequestWithContext(ctx, method, u, body)
	if err == nil {
		resp, err = s.client.Do(req, respBody)
	}

	return respBody, resp, err
}
