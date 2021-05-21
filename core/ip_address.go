package core

import (
	"context"
	"net/url"
	"strings"

	"github.com/krystal/go-katapult"
)

type IPAddress struct {
	ID              string   `json:"id,omitempty"`
	Address         string   `json:"address,omitempty"`
	ReverseDNS      string   `json:"reverse_dns,omitempty"`
	VIP             bool     `json:"vip,omitempty"`
	Label           string   `json:"label,omitempty"`
	AddressWithMask string   `json:"address_with_mask,omitempty"`
	Network         *Network `json:"network,omitempty"`
	AllocationID    string   `json:"allocation_id,omitempty"`
	AllocationType  string   `json:"allocation_type,omitempty"`
}

func (s *IPAddress) Ref() IPAddressRef {
	return IPAddressRef{ID: s.ID}
}

func (s *IPAddress) Version() IPVersion {
	if strings.Count(s.Address, ":") < 2 {
		return IPv4
	}

	return IPv6
}

type IPAddressRef struct {
	ID      string `json:"id,omitempty"`
	Address string `json:"address,omitempty"`
}

func (ref IPAddressRef) queryValues() *url.Values {
	v := &url.Values{}

	switch {
	case ref.ID != "":
		v.Set("ip_address[id]", ref.ID)
	case ref.Address != "":
		v.Set("ip_address[address]", ref.Address)
	}

	return v
}

type IPAddressCreateArguments struct {
	Network *Network
	Version IPVersion
	VIP     *bool
	Label   string
}

type IPAddressUpdateArguments struct {
	VIP        *bool
	Label      string
	ReverseDNS string
}

type ipAddressCreateRequest struct {
	Organization OrganizationRef `json:"organization,omitempty"`
	Network      *Network        `json:"network,omitempty"`
	Version      IPVersion       `json:"version,omitempty"`
	VIP          *bool           `json:"vip,omitempty"`
	Label        string          `json:"label,omitempty"`
}

type ipAddressUpdateRequest struct {
	IPAddress  IPAddressRef `json:"ip_address"`
	VIP        *bool        `json:"vip,omitempty"`
	Label      string       `json:"label,omitempty"`
	ReverseDNS string       `json:"reverse_dns,omitempty"`
}

type ipAddressesResponseBody struct {
	Pagination  *katapult.Pagination `json:"pagination,omitempty"`
	IPAddress   *IPAddress           `json:"ip_address,omitempty"`
	IPAddresses []*IPAddress         `json:"ip_addresses,omitempty"`
}

type IPAddressesClient struct {
	client   RequestMaker
	basePath *url.URL
}

func NewIPAddressesClient(rm RequestMaker) *IPAddressesClient {
	return &IPAddressesClient{
		client:   rm,
		basePath: &url.URL{Path: "/core/v1/"},
	}
}

func (s *IPAddressesClient) List(
	ctx context.Context,
	org OrganizationRef,
	opts *ListOptions,
) ([]*IPAddress, *katapult.Response, error) {
	qs := queryValues(org, opts)
	u := &url.URL{
		Path:     "organizations/_/ip_addresses",
		RawQuery: qs.Encode(),
	}

	body, resp, err := s.doRequest(ctx, "GET", u, nil)
	resp.Pagination = body.Pagination

	return body.IPAddresses, resp, err
}

func (s *IPAddressesClient) Get(
	ctx context.Context,
	ref IPAddressRef,
) (*IPAddress, *katapult.Response, error) {
	u := &url.URL{
		Path:     "ip_addresses/_",
		RawQuery: ref.queryValues().Encode(),
	}

	body, resp, err := s.doRequest(ctx, "GET", u, nil)

	return body.IPAddress, resp, err
}

func (s *IPAddressesClient) GetByID(
	ctx context.Context,
	id string,
) (*IPAddress, *katapult.Response, error) {
	return s.Get(ctx, IPAddressRef{ID: id})
}

func (s *IPAddressesClient) GetByAddress(
	ctx context.Context,
	address string,
) (*IPAddress, *katapult.Response, error) {
	return s.Get(ctx, IPAddressRef{Address: address})
}

func (s *IPAddressesClient) Create(
	ctx context.Context,
	org OrganizationRef,
	args *IPAddressCreateArguments,
) (*IPAddress, *katapult.Response, error) {
	u := &url.URL{Path: "organizations/_/ip_addresses"}
	reqBody := &ipAddressCreateRequest{
		Organization: org,
	}

	if args != nil {
		reqBody.Network = args.Network.lookupReference()
		reqBody.Version = args.Version
		reqBody.VIP = args.VIP
		reqBody.Label = args.Label
	}

	body, resp, err := s.doRequest(ctx, "POST", u, reqBody)

	return body.IPAddress, resp, err
}

func (s *IPAddressesClient) Update(
	ctx context.Context,
	ip IPAddressRef,
	args *IPAddressUpdateArguments,
) (*IPAddress, *katapult.Response, error) {
	u := &url.URL{Path: "ip_addresses/_"}
	reqBody := &ipAddressUpdateRequest{
		IPAddress: ip,
	}

	if args != nil {
		reqBody.VIP = args.VIP
		reqBody.Label = args.Label
		reqBody.ReverseDNS = args.ReverseDNS
	}

	body, resp, err := s.doRequest(ctx, "PATCH", u, reqBody)

	return body.IPAddress, resp, err
}

func (s *IPAddressesClient) Delete(
	ctx context.Context,
	ip IPAddressRef,
) (*katapult.Response, error) {
	qs := queryValues(ip)
	u := &url.URL{Path: "ip_addresses/_", RawQuery: qs.Encode()}

	_, resp, err := s.doRequest(ctx, "DELETE", u, nil)

	return resp, err
}

func (s *IPAddressesClient) Unallocate(
	ctx context.Context,
	ip IPAddressRef,
) (*katapult.Response, error) {
	qs := queryValues(ip)
	u := &url.URL{Path: "ip_addresses/_/unallocate", RawQuery: qs.Encode()}

	_, resp, err := s.doRequest(ctx, "POST", u, nil)

	return resp, err
}

func (s *IPAddressesClient) doRequest(
	ctx context.Context,
	method string,
	u *url.URL,
	body interface{},
) (*ipAddressesResponseBody, *katapult.Response, error) {
	u = s.basePath.ResolveReference(u)
	respBody := &ipAddressesResponseBody{}
	resp := katapult.NewResponse(nil)

	req, err := s.client.NewRequestWithContext(ctx, method, u, body)
	if err == nil {
		resp, err = s.client.Do(req, respBody)
	}

	return respBody, resp, err
}
