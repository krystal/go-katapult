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

// NewIPAddressLookup takes a string that is a IPAddress ID or Address,
// returning a empty *IPAddress struct with either the ID or Address field
// populated with the given value. This struct is suitable as input to other
// methods which accept a *IPAddress as input.
func NewIPAddressLookup(
	idOrAddress string,
) (lr *IPAddress, f FieldName) {
	if strings.HasPrefix(idOrAddress, "ip_") {
		return &IPAddress{ID: idOrAddress}, IDField
	}

	return &IPAddress{Address: idOrAddress}, AddressField
}

func (s *IPAddress) Version() IPVersion {
	if strings.Count(s.Address, ":") < 2 {
		return IPv4
	}

	return IPv6
}

func (s *IPAddress) lookupReference() *IPAddress {
	if s == nil {
		return nil
	}

	lr := &IPAddress{ID: s.ID}
	if lr.ID == "" {
		lr.Address = s.Address
	}

	return lr
}

func (s *IPAddress) queryValues() *url.Values {
	v := &url.Values{}

	if s != nil {
		switch {
		case s.ID != "":
			v.Set("ip_address[id]", s.ID)
		case s.Address != "":
			v.Set("ip_address[address]", s.Address)
		}
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
	Organization *Organization `json:"organization,omitempty"`
	Network      *Network      `json:"network,omitempty"`
	Version      IPVersion     `json:"version,omitempty"`
	VIP          *bool         `json:"vip,omitempty"`
	Label        string        `json:"label,omitempty"`
}

type ipAddressUpdateRequest struct {
	IPAddress  *IPAddress `json:"ip_address,omitempty"`
	VIP        *bool      `json:"vip,omitempty"`
	Label      string     `json:"label,omitempty"`
	ReverseDNS string     `json:"reverse_dns,omitempty"`
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
	org *Organization,
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
	idOrAddress string,
) (*IPAddress, *katapult.Response, error) {
	if _, f := NewIPAddressLookup(idOrAddress); f == IDField {
		return s.GetByID(ctx, idOrAddress)
	}

	return s.GetByAddress(ctx, idOrAddress)
}

func (s *IPAddressesClient) GetByID(
	ctx context.Context,
	id string,
) (*IPAddress, *katapult.Response, error) {
	return s.get(ctx, &IPAddress{ID: id})
}

func (s *IPAddressesClient) GetByAddress(
	ctx context.Context,
	address string,
) (*IPAddress, *katapult.Response, error) {
	return s.get(ctx, &IPAddress{Address: address})
}

func (s *IPAddressesClient) get(
	ctx context.Context,
	ip *IPAddress,
) (*IPAddress, *katapult.Response, error) {
	u := &url.URL{
		Path:     "ip_addresses/_",
		RawQuery: ip.queryValues().Encode(),
	}

	body, resp, err := s.doRequest(ctx, "GET", u, nil)

	return body.IPAddress, resp, err
}

func (s *IPAddressesClient) Create(
	ctx context.Context,
	org *Organization,
	args *IPAddressCreateArguments,
) (*IPAddress, *katapult.Response, error) {
	u := &url.URL{Path: "organizations/_/ip_addresses"}
	reqBody := &ipAddressCreateRequest{
		Organization: org.lookupReference(),
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
	ip *IPAddress,
	args *IPAddressUpdateArguments,
) (*IPAddress, *katapult.Response, error) {
	u := &url.URL{Path: "ip_addresses/_"}
	reqBody := &ipAddressUpdateRequest{
		IPAddress: ip.lookupReference(),
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
	ip *IPAddress,
) (*katapult.Response, error) {
	qs := queryValues(ip)
	u := &url.URL{Path: "ip_addresses/_", RawQuery: qs.Encode()}

	_, resp, err := s.doRequest(ctx, "DELETE", u, nil)

	return resp, err
}

func (s *IPAddressesClient) Unallocate(
	ctx context.Context,
	ip *IPAddress,
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