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
	reqOpts ...katapult.RequestOption,
) ([]*Network, []*VirtualNetwork, *katapult.Response, error) {
	u := &url.URL{
		Path:     "organizations/_/available_networks",
		RawQuery: org.queryValues().Encode(),
	}

	body, resp, err := s.doRequest(ctx, "GET", u, nil, reqOpts...)

	return body.Networks, body.VirtualNetworks, resp, err
}

func (s *NetworksClient) Get(
	ctx context.Context,
	ref NetworkRef,
	reqOpts ...katapult.RequestOption,
) (*Network, *katapult.Response, error) {
	u := &url.URL{
		Path:     "networks/_",
		RawQuery: ref.queryValues().Encode(),
	}

	body, resp, err := s.doRequest(ctx, "GET", u, nil, reqOpts...)

	return body.Network, resp, err
}

func (s *NetworksClient) GetByID(
	ctx context.Context,
	id string,
	reqOpts ...katapult.RequestOption,
) (*Network, *katapult.Response, error) {
	return s.Get(ctx, NetworkRef{ID: id}, reqOpts...)
}

func (s *NetworksClient) GetByPermalink(
	ctx context.Context,
	permalink string,
	reqOpts ...katapult.RequestOption,
) (*Network, *katapult.Response, error) {
	return s.Get(ctx, NetworkRef{Permalink: permalink}, reqOpts...)
}

func (s *NetworksClient) doRequest(
	ctx context.Context,
	method string,
	u *url.URL,
	body interface{},
	reqOpts ...katapult.RequestOption,
) (*networksResponseBody, *katapult.Response, error) {
	u = s.basePath.ResolveReference(u)
	respBody := &networksResponseBody{}

	req := katapult.NewRequest(method, u, body, reqOpts...)
	resp, err := s.client.Do(ctx, req, respBody)
	if resp == nil {
		resp = katapult.NewResponse(nil)
	}

	return respBody, resp, handleResponseError(err)
}
