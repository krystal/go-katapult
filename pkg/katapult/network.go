package katapult

import (
	"context"
	"net/url"
	"strings"
)

const networkIDPrefix = "netw_"

type Network struct {
	ID         string      `json:"id,omitempty"`
	Name       string      `json:"name,omitempty"`
	Permalink  string      `json:"permalink,omitempty"`
	DataCenter *DataCenter `json:"data_center,omitempty"`
}

func (s *Network) lookupReference() *Network {
	if s == nil {
		return nil
	}

	lr := &Network{ID: s.ID}
	if lr.ID == "" {
		lr.Permalink = s.Permalink
	}

	return lr
}

func (s *Network) queryValues() *url.Values {
	v := &url.Values{}

	if s != nil {
		switch {
		case s.ID != "":
			v.Set("network[id]", s.ID)
		case s.Permalink != "":
			v.Set("network[permalink]", s.Permalink)
		}
	}

	return v
}

type VirtualNetwork struct {
	ID         string      `json:"id,omitempty"`
	Name       string      `json:"name,omitempty"`
	DataCenter *DataCenter `json:"data_center,omitempty"`
}

type networksResponseBody struct {
	Network         *Network          `json:"network,omitempty"`
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
	u := &url.URL{
		Path:     "organizations/_/available_networks",
		RawQuery: org.queryValues().Encode(),
	}

	body, resp, err := s.doRequest(ctx, "GET", u, nil)

	return body.Networks, body.VirtualNetworks, resp, err
}

func (s *NetworksClient) Get(
	ctx context.Context,
	idOrPermalink string,
) (*Network, *Response, error) {
	if strings.HasPrefix(idOrPermalink, networkIDPrefix) {
		return s.GetByID(ctx, idOrPermalink)
	}

	return s.GetByPermalink(ctx, idOrPermalink)
}

func (s *NetworksClient) GetByID(
	ctx context.Context,
	id string,
) (*Network, *Response, error) {
	return s.get(ctx, &Network{ID: id})
}

func (s *NetworksClient) GetByPermalink(
	ctx context.Context,
	permalink string,
) (*Network, *Response, error) {
	return s.get(ctx, &Network{Permalink: permalink})
}

func (s *NetworksClient) get(
	ctx context.Context,
	network *Network,
) (*Network, *Response, error) {
	u := &url.URL{
		Path:     "networks/_",
		RawQuery: network.queryValues().Encode(),
	}

	body, resp, err := s.doRequest(ctx, "GET", u, nil)

	return body.Network, resp, err
}

func (s *NetworksClient) doRequest(
	ctx context.Context,
	method string,
	u *url.URL,
	body interface{},
) (*networksResponseBody, *Response, error) {
	u = s.basePath.ResolveReference(u)
	respBody := &networksResponseBody{}
	resp := newResponse(nil)

	req, err := s.client.NewRequestWithContext(ctx, method, u, body)
	if err == nil {
		resp, err = s.client.Do(req, respBody)
	}

	return respBody, resp, err
}
