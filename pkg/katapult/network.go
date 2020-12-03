package katapult

import (
	"context"
	"fmt"
	"net/url"
)

type Network struct {
	ID         string      `json:"id,omitempty"`
	Name       string      `json:"name,omitempty"`
	Permalink  string      `json:"permalink,omitempty"`
	DataCenter *DataCenter `json:"data_center,omitempty"`
}

// LookupReference returns a new *Network stripped down to just ID or
// Permalink fields, making it suitable for endpoints which require a reference
// to a Network by ID or Permalink.
func (s *Network) LookupReference() *Network {
	if s == nil {
		return nil
	}

	lr := &Network{ID: s.ID}
	if lr.ID == "" {
		lr.Permalink = s.Permalink
	}

	return lr
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

type NetworksClient struct {
	client   *apiClient
	basePath *url.URL
}

func newNetworksClient(c *apiClient) *NetworksClient {
	return &NetworksClient{
		client:   c,
		basePath: &url.URL{Path: "/core/v1/"},
	}
}

func (s *NetworksClient) List(
	ctx context.Context,
	org *Organization,
) ([]*Network, []*VirtualNetwork, *Response, error) {
	if org == nil {
		org = &Organization{ID: "_"}
	}

	u := &url.URL{
		Path: fmt.Sprintf("organizations/%s/available_networks", org.ID),
	}

	body, resp, err := s.doRequest(ctx, "GET", u, nil)

	return body.Networks, body.VirtualNetworks, resp, err
}

func (s *NetworksClient) doRequest(
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
