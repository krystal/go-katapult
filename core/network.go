package core

import (
	"context"
	"net/url"

	"github.com/krystal/go-katapult"
)

type Network struct {
	ID         string      `json:"id,omitempty"`
	Name       string      `json:"name,omitempty"`
	Permalink  string      `json:"permalink,omitempty"`
	DataCenter *DataCenter `json:"data_center,omitempty"`
}

func (s *Network) Ref() NetworkRef {
	return NetworkRef{ID: s.ID}
}

type NetworkRef struct {
	ID        string `json:"id,omitempty"`
	Permalink string `json:"permalink,omitempty"`
}

func (ref NetworkRef) queryValues() *url.Values {
	v := &url.Values{}

	switch {
	case ref.ID != "":
		v.Set("network[id]", ref.ID)
	case ref.Permalink != "":
		v.Set("network[permalink]", ref.Permalink)
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
	client   RequestMaker
	basePath *url.URL
}

func NewNetworksClient(rm RequestMaker) *NetworksClient {
	return &NetworksClient{
		client:   rm,
		basePath: &url.URL{Path: "/core/v1/"},
	}
}

func (s *NetworksClient) List(
	ctx context.Context,
	org OrganizationRef,
) ([]*Network, []*VirtualNetwork, *katapult.Response, error) {
	u := &url.URL{
		Path:     "organizations/_/available_networks",
		RawQuery: org.queryValues().Encode(),
	}

	body, resp, err := s.doRequest(ctx, "GET", u, nil)

	return body.Networks, body.VirtualNetworks, resp, err
}

func (s *NetworksClient) Get(
	ctx context.Context,
	ref NetworkRef,
) (*Network, *katapult.Response, error) {
	u := &url.URL{
		Path:     "networks/_",
		RawQuery: ref.queryValues().Encode(),
	}

	body, resp, err := s.doRequest(ctx, "GET", u, nil)

	return body.Network, resp, err
}

func (s *NetworksClient) GetByID(
	ctx context.Context,
	id string,
) (*Network, *katapult.Response, error) {
	return s.Get(ctx, NetworkRef{ID: id})
}

func (s *NetworksClient) GetByPermalink(
	ctx context.Context,
	permalink string,
) (*Network, *katapult.Response, error) {
	return s.Get(ctx, NetworkRef{Permalink: permalink})
}

func (s *NetworksClient) doRequest(
	ctx context.Context,
	method string,
	u *url.URL,
	body interface{},
) (*networksResponseBody, *katapult.Response, error) {
	u = s.basePath.ResolveReference(u)
	respBody := &networksResponseBody{}
	resp := katapult.NewResponse(nil)

	req, err := s.client.NewRequestWithContext(ctx, method, u, body)
	if err == nil {
		resp, err = s.client.Do(req, respBody)
	}

	return respBody, resp, err
}