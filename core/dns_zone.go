package core

import (
	"context"
	"net/url"

	"github.com/krystal/go-katapult"
)

type DNSZone struct {
	ID         string `json:"id,omitempty"`
	Name       string `json:"name,omitempty"`
	DefaultTTL int    `json:"default_ttl,omitempty"`
	Verified   bool   `json:"verified,omitempty"`
}

func (s *DNSZone) Ref() DNSZoneRef {
	return DNSZoneRef{ID: s.ID}
}

type DNSZoneRef struct {
	ID   string `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
}

func (s DNSZoneRef) queryValues() *url.Values {
	v := &url.Values{}

	switch {
	case s.ID != "":
		v.Set("dns_zone[id]", s.ID)
	case s.Name != "":
		v.Set("dns_zone[name]", s.Name)
	}

	return v
}

type DNSZoneCreateArguments struct {
	Name       string `json:"name"`
	DefaultTTL int    `json:"default_ttl,omitempty"`
}

type DNSZoneUpdateArguments struct {
	Name       string `json:"name"`
	DefaultTTL int    `json:"default_ttl,omitempty"`
}

type dnsZoneCreateRequest struct {
	Organization OrganizationRef         `json:"organization"`
	Properties   *DNSZoneCreateArguments `json:"properties"`
}

type dnsZoneUpdateRequest struct {
	DNSZone    DNSZoneRef              `json:"dns_zone"`
	Properties *DNSZoneUpdateArguments `json:"properties"`
}

type dnsZoneResponseBody struct {
	Pagination  *katapult.Pagination `json:"pagination,omitempty"`
	DNSZones    []*DNSZone           `json:"dns_zones"`
	DNSZone     *DNSZone             `json:"dns_zone"`
	Deleted     *bool                `json:"deleted,omitempty"`
	Nameservers []string             `json:"nameservers,omitempty"`
}

type DNSZonesClient struct {
	client   RequestMaker
	basePath *url.URL
}

func NewDNSZonesClient(rm RequestMaker) *DNSZonesClient {
	return &DNSZonesClient{
		client:   rm,
		basePath: &url.URL{Path: "/core/v1/"},
	}
}

func (s *DNSZonesClient) List(
	ctx context.Context,
	org OrganizationRef,
	opts *ListOptions,
	reqOpts ...katapult.RequestOption,
) ([]*DNSZone, *katapult.Response, error) {
	qs := queryValues(org, opts)

	u := &url.URL{
		Path:     "organizations/_/dns_zones",
		RawQuery: qs.Encode(),
	}

	body, resp, err := s.doRequest(ctx, "GET", u, nil, reqOpts...)
	resp.Pagination = body.Pagination

	return body.DNSZones, resp, err
}

func (s *DNSZonesClient) Nameservers(
	ctx context.Context,
	org OrganizationRef,
	reqOpts ...katapult.RequestOption,
) ([]string, *katapult.Response, error) {
	u := &url.URL{
		Path:     "organizations/_/dns_zones/nameservers",
		RawQuery: org.queryValues().Encode(),
	}

	body, resp, err := s.doRequest(ctx, "GET", u, nil, reqOpts...)

	return body.Nameservers, resp, err
}

func (s *DNSZonesClient) Get(
	ctx context.Context,
	ref DNSZoneRef,
	reqOpts ...katapult.RequestOption,
) (*DNSZone, *katapult.Response, error) {
	u := &url.URL{Path: "dns_zones/_", RawQuery: ref.queryValues().Encode()}

	body, resp, err := s.doRequest(ctx, "GET", u, nil, reqOpts...)

	return body.DNSZone, resp, err
}

func (s *DNSZonesClient) GetByID(
	ctx context.Context,
	id string,
	reqOpts ...katapult.RequestOption,
) (*DNSZone, *katapult.Response, error) {
	return s.Get(ctx, DNSZoneRef{ID: id}, reqOpts...)
}

func (s *DNSZonesClient) GetByName(
	ctx context.Context,
	name string,
	reqOpts ...katapult.RequestOption,
) (*DNSZone, *katapult.Response, error) {
	return s.Get(ctx, DNSZoneRef{Name: name}, reqOpts...)
}

func (s *DNSZonesClient) Create(
	ctx context.Context,
	org OrganizationRef,
	args *DNSZoneCreateArguments,
	reqOpts ...katapult.RequestOption,
) (*DNSZone, *katapult.Response, error) {
	u := &url.URL{
		Path:     "organizations/_/dns_zones",
		RawQuery: org.queryValues().Encode(),
	}
	reqBody := &dnsZoneCreateRequest{Properties: args}

	body, resp, err := s.doRequest(ctx, "POST", u, reqBody, reqOpts...)

	return body.DNSZone, resp, err
}

func (s *DNSZonesClient) Update(
	ctx context.Context,
	zone DNSZoneRef,
	args *DNSZoneUpdateArguments,
	reqOpts ...katapult.RequestOption,
) (*DNSZone, *katapult.Response, error) {
	u := &url.URL{
		Path:     "dns_zones/_",
		RawQuery: zone.queryValues().Encode(),
	}
	reqBody := &dnsZoneUpdateRequest{Properties: args}

	body, resp, err := s.doRequest(ctx, "PATCH", u, reqBody, reqOpts...)

	return body.DNSZone, resp, err
}

func (s *DNSZonesClient) Delete(
	ctx context.Context,
	zone DNSZoneRef,
	reqOpts ...katapult.RequestOption,
) (*bool, *katapult.Response, error) {
	u := &url.URL{
		Path:     "dns_zones/_",
		RawQuery: zone.queryValues().Encode(),
	}
	body, resp, err := s.doRequest(ctx, "DELETE", u, nil, reqOpts...)

	return body.Deleted, resp, err
}

func (s *DNSZonesClient) Verify(
	ctx context.Context,
	ref DNSZoneRef,
	reqOpts ...katapult.RequestOption,
) (*DNSZone, *katapult.Response, error) {
	u := &url.URL{
		Path:     "dns_zones/_/verify",
		RawQuery: ref.queryValues().Encode(),
	}

	body, resp, err := s.doRequest(ctx, "POST", u, nil, reqOpts...)

	return body.DNSZone, resp, err
}

func (s *DNSZonesClient) doRequest(
	ctx context.Context,
	method string,
	u *url.URL,
	body interface{},
	reqOpts ...katapult.RequestOption,
) (*dnsZoneResponseBody, *katapult.Response, error) {
	u = s.basePath.ResolveReference(u)
	respBody := &dnsZoneResponseBody{}

	req := katapult.NewRequest(method, u, body, reqOpts...)
	resp, err := s.client.Do(ctx, req, respBody)
	if resp == nil {
		resp = katapult.NewResponse(nil)
	}

	return respBody, resp, handleResponseError(err)
}
