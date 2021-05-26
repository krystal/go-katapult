package core

import (
	"context"
	"net/url"

	"github.com/krystal/go-katapult"
)

type DNSZone struct {
	ID                 string `json:"id,omitempty"`
	Name               string `json:"name,omitempty"`
	TTL                int    `json:"ttl,omitempty"`
	Verified           bool   `json:"verified,omitempty"`
	InfrastructureZone bool   `json:"infrastructure_zone,omitempty"`
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

type DNSZoneVerificationDetails struct {
	Nameservers []string `json:"nameservers,omitempty"`
	TXTRecord   string   `json:"txt_record,omitempty"`
}

type DNSZoneArguments struct {
	Name            string
	TTL             int
	SkipVerfication bool
}

type DNSZoneDetails struct {
	Name string `json:"name"`
	TTL  int    `json:"ttl,omitempty"`
}

type dnsZoneCreateRequest struct {
	Organization    OrganizationRef `json:"organization"`
	Details         *DNSZoneDetails `json:"details"`
	SkipVerfication bool            `json:"skip_verification"`
}

type dnsZoneVerifyRequest struct {
	DNSZone DNSZoneRef `json:"dns_zone"`
}

type dnsZoneUpdateTTLRequest struct {
	DNSZone DNSZoneRef `json:"dns_zone"`
	TTL     int        `json:"ttl"`
}

type dnsZoneResponseBody struct {
	Pagination          *katapult.Pagination        `json:"pagination,omitempty"`
	DNSZones            []*DNSZone                  `json:"dns_zones"`
	DNSZone             *DNSZone                    `json:"dns_zone"`
	VerificationDetails *DNSZoneVerificationDetails `json:"details,omitempty"`
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
) ([]*DNSZone, *katapult.Response, error) {
	qs := queryValues(org, opts)

	u := &url.URL{
		Path:     "organizations/_/dns/zones",
		RawQuery: qs.Encode(),
	}

	body, resp, err := s.doRequest(ctx, "GET", u, nil)
	resp.Pagination = body.Pagination

	return body.DNSZones, resp, err
}

func (s *DNSZonesClient) Get(
	ctx context.Context,
	ref DNSZoneRef,
) (*DNSZone, *katapult.Response, error) {
	u := &url.URL{Path: "dns/zones/_", RawQuery: ref.queryValues().Encode()}

	body, resp, err := s.doRequest(ctx, "GET", u, nil)

	return body.DNSZone, resp, err
}

func (s *DNSZonesClient) GetByID(
	ctx context.Context,
	id string,
) (*DNSZone, *katapult.Response, error) {
	return s.Get(ctx, DNSZoneRef{ID: id})
}

func (s *DNSZonesClient) GetByName(
	ctx context.Context,
	name string,
) (*DNSZone, *katapult.Response, error) {
	return s.Get(ctx, DNSZoneRef{Name: name})
}

func (s *DNSZonesClient) Create(
	ctx context.Context,
	org OrganizationRef,
	args *DNSZoneArguments,
) (*DNSZone, *katapult.Response, error) {
	u := &url.URL{Path: "organizations/_/dns/zones"}
	reqBody := &dnsZoneCreateRequest{
		Organization: org,
	}

	if args != nil {
		reqBody.Details = &DNSZoneDetails{
			Name: args.Name,
			TTL:  args.TTL,
		}
		reqBody.SkipVerfication = args.SkipVerfication
	}

	body, resp, err := s.doRequest(ctx, "POST", u, reqBody)

	return body.DNSZone, resp, err
}

func (s *DNSZonesClient) Delete(
	ctx context.Context,
	zone DNSZoneRef,
) (*DNSZone, *katapult.Response, error) {
	u := &url.URL{
		Path:     "dns/zones/_",
		RawQuery: zone.queryValues().Encode(),
	}
	body, resp, err := s.doRequest(ctx, "DELETE", u, nil)

	return body.DNSZone, resp, err
}

func (s *DNSZonesClient) VerificationDetails(
	ctx context.Context,
	zone DNSZoneRef,
) (*DNSZoneVerificationDetails, *katapult.Response, error) {
	u := &url.URL{
		Path:     "dns/zones/_/verification_details",
		RawQuery: zone.queryValues().Encode(),
	}
	body, resp, err := s.doRequest(ctx, "GET", u, nil)

	return body.VerificationDetails, resp, err
}

func (s *DNSZonesClient) Verify(
	ctx context.Context,
	ref DNSZoneRef,
) (*DNSZone, *katapult.Response, error) {
	u := &url.URL{Path: "dns/zones/_/verify"}
	reqBody := &dnsZoneVerifyRequest{
		DNSZone: ref,
	}
	body, resp, err := s.doRequest(ctx, "POST", u, reqBody)

	return body.DNSZone, resp, err
}

func (s *DNSZonesClient) UpdateTTL(
	ctx context.Context,
	ref DNSZoneRef,
	ttl int,
) (*DNSZone, *katapult.Response, error) {
	u := &url.URL{Path: "dns/zones/_/update_ttl"}
	reqBody := &dnsZoneUpdateTTLRequest{
		DNSZone: ref,
		TTL:     ttl,
	}
	body, resp, err := s.doRequest(ctx, "POST", u, reqBody)

	return body.DNSZone, resp, err
}

func (s *DNSZonesClient) doRequest(
	ctx context.Context,
	method string,
	u *url.URL,
	body interface{},
) (*dnsZoneResponseBody, *katapult.Response, error) {
	u = s.basePath.ResolveReference(u)
	respBody := &dnsZoneResponseBody{}
	resp := katapult.NewResponse(nil)

	req, err := s.client.NewRequestWithContext(ctx, method, u, body)
	if err == nil {
		resp, err = s.client.Do(req, respBody)
	}

	return respBody, resp, err
}
