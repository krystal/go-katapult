package core

import (
	"context"
	"github.com/krystal/go-katapult"
	"net/url"
	"strings"
)

type Network struct {
	ID         string      `json:"id,omitempty"`
	Name       string      `json:"name,omitempty"`
	Permalink  string      `json:"permalink,omitempty"`
	DataCenter *DataCenter `json:"data_center,omitempty"`
}

// NewNetworkLookup takes a string that is a Network ID or Permalink, returning
// a empty *Network struct with either the ID or Permalink field populated with
// the given value. This struct is suitable as input to other methods which
// accept a *Network as input.
func NewNetworkLookup(
	idOrPermalink string,
) (lr *Network, f FieldName) {
	if strings.HasPrefix(idOrPermalink, "netw_") {
		return &Network{ID: idOrPermalink}, IDField
	}

	return &Network{Permalink: idOrPermalink}, PermalinkField
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
	org *Organization,
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
	idOrPermalink string,
) (*Network, *katapult.Response, error) {
	if _, f := NewNetworkLookup(idOrPermalink); f == IDField {
		return s.GetByID(ctx, idOrPermalink)
	}

	return s.GetByPermalink(ctx, idOrPermalink)
}

func (s *NetworksClient) GetByID(
	ctx context.Context,
	id string,
) (*Network, *katapult.Response, error) {
	return s.get(ctx, &Network{ID: id})
}

func (s *NetworksClient) GetByPermalink(
	ctx context.Context,
	permalink string,
) (*Network, *katapult.Response, error) {
	return s.get(ctx, &Network{Permalink: permalink})
}

func (s *NetworksClient) get(
	ctx context.Context,
	network *Network,
) (*Network, *katapult.Response, error) {
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
