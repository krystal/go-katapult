package katapult

import (
	"context"
	"fmt"
	"net/url"
)

type NetworksService struct {
	*service
	path *url.URL
}

func NewNetworksService(s *service) *NetworksService {
	return &NetworksService{
		service: s,
		path:    &url.URL{Path: "/core/v1/"},
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
	u, err := s.path.Parse(
		fmt.Sprintf("organizations/%s/available_networks", orgID),
	)
	if err != nil {
		return nil, nil, nil, err
	}

	req, err := s.client.NewRequestWithContext(ctx, "GET", u.Path, nil)
	if err != nil {
		return nil, nil, nil, err
	}

	var body networksResponseBody
	resp, err := s.client.Do(req, &body)
	if err != nil {
		return nil, nil, resp, err
	}

	return body.Networks, body.VirtualNetworks, resp, nil
}
