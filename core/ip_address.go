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
	Network NetworkRef
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
	Network      NetworkRef      `json:"network,omitempty"`
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
	reqOpts ...katapult.RequestOption,
) ([]*IPAddress, *katapult.Response, error) {
	qs := queryValues(org, opts)
	u := &url.URL{
		Path:     "organizations/_/ip_addresses",
		RawQuery: qs.Encode(),
	}

	body, resp, err := s.doRequest(ctx, "GET", u, nil, reqOpts...)
	resp.Pagination = body.Pagination

	return body.IPAddresses, resp, err
}

func (s *IPAddressesClient) Get(
	ctx context.Context,
	ref IPAddressRef,
	reqOpts ...katapult.RequestOption,
) (*IPAddress, *katapult.Response, error) {
	u := &url.URL{
		Path:     "ip_addresses/_",
		RawQuery: ref.queryValues().Encode(),
	}

	body, resp, err := s.doRequest(ctx, "GET", u, nil, reqOpts...)

	return body.IPAddress, resp, err
}

func (s *IPAddressesClient) GetByID(
	ctx context.Context,
	id string,
	reqOpts ...katapult.RequestOption,
) (*IPAddress, *katapult.Response, error) {
	return s.Get(ctx, IPAddressRef{ID: id}, reqOpts...)
}

func (s *IPAddressesClient) GetByAddress(
	ctx context.Context,
	address string,
	reqOpts ...katapult.RequestOption,
) (*IPAddress, *katapult.Response, error) {
	return s.Get(ctx, IPAddressRef{Address: address}, reqOpts...)
}

func (s *IPAddressesClient) Create(
	ctx context.Context,
	org OrganizationRef,
	args *IPAddressCreateArguments,
	reqOpts ...katapult.RequestOption,
) (*IPAddress, *katapult.Response, error) {
	u := &url.URL{Path: "organizations/_/ip_addresses"}
	reqBody := &ipAddressCreateRequest{
		Organization: org,
	}

	if args != nil {
		reqBody.Network = args.Network
		reqBody.Version = args.Version
		reqBody.VIP = args.VIP
		reqBody.Label = args.Label
	}

	body, resp, err := s.doRequest(ctx, "POST", u, reqBody, reqOpts...)

	return body.IPAddress, resp, err
}

func (s *IPAddressesClient) Update(
	ctx context.Context,
	ip IPAddressRef,
	args *IPAddressUpdateArguments,
	reqOpts ...katapult.RequestOption,
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

	body, resp, err := s.doRequest(ctx, "PATCH", u, reqBody, reqOpts...)

	return body.IPAddress, resp, err
}

func (s *IPAddressesClient) Delete(
	ctx context.Context,
	ip IPAddressRef,
	reqOpts ...katapult.RequestOption,
) (*katapult.Response, error) {
	qs := queryValues(ip)
	u := &url.URL{Path: "ip_addresses/_", RawQuery: qs.Encode()}

	_, resp, err := s.doRequest(ctx, "DELETE", u, nil, reqOpts...)

	return resp, err
}

func (s *IPAddressesClient) Unallocate(
	ctx context.Context,
	ip IPAddressRef,
	reqOpts ...katapult.RequestOption,
) (*katapult.Response, error) {
	qs := queryValues(ip)
	u := &url.URL{Path: "ip_addresses/_/unallocate", RawQuery: qs.Encode()}

	_, resp, err := s.doRequest(ctx, "POST", u, nil, reqOpts...)

	return resp, err
}

func (s *IPAddressesClient) doRequest(
	ctx context.Context,
	method string,
	u *url.URL,
	body interface{},
	reqOpts ...katapult.RequestOption,
) (*ipAddressesResponseBody, *katapult.Response, error) {
	u = s.basePath.ResolveReference(u)
	respBody := &ipAddressesResponseBody{}

	req := katapult.NewRequest(method, u, body, reqOpts...)
	resp, err := s.client.Do(ctx, req, respBody)
	if resp == nil {
		resp = katapult.NewResponse(nil)
	}

	return respBody, resp, handleResponseError(err)
}
