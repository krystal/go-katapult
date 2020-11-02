package katapult

import (
	"context"
	"errors"
	"fmt"
	"net/url"
)

type DNSZone struct {
	ID                 string `json:"id,omitempty"`
	Name               string `json:"name,omitempty"`
	TTL                int    `json:"ttl,omitempty"`
	Verified           bool   `json:"verified,omitempty"`
	InfrastructureZone bool   `json:"infrastructure_zone,omitempty"`
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

type createDNSZoneRequest struct {
	Details         *DNSZoneDetails `json:"details"`
	SkipVerfication bool            `json:"skip_verification"`
}

type updateDNSZoneTTLRequest struct {
	TTL int `json:"ttl"`
}

type dnsZoneResponseBody struct {
	Pagination          *Pagination                 `json:"pagination,omitempty"`
	DNSZones            []*DNSZone                  `json:"dns_zones"`
	DNSZone             *DNSZone                    `json:"dns_zone"`
	VerificationDetails *DNSZoneVerificationDetails `json:"details,omitempty"`
}

type DNSZonesClient struct {
	client   *apiClient
	basePath *url.URL
}

func newDNSZonesClient(c *apiClient) *DNSZonesClient {
	return &DNSZonesClient{
		client:   c,
		basePath: &url.URL{Path: "/core/v1/"},
	}
}

func (s *DNSZonesClient) List(
	ctx context.Context,
	orgID string,
	opts *ListOptions,
) ([]*DNSZone, *Response, error) {
	u := &url.URL{
		Path:     fmt.Sprintf("organizations/%s/dns/zones", orgID),
		RawQuery: opts.Values().Encode(),
	}

	body, resp, err := s.doRequest(ctx, "GET", u, nil)
	resp.Pagination = body.Pagination

	return body.DNSZones, resp, err
}

func (s *DNSZonesClient) Get(
	ctx context.Context,
	id string,
) (*DNSZone, *Response, error) {
	u := &url.URL{Path: fmt.Sprintf("dns/zones/%s", id)}
	body, resp, err := s.doRequest(ctx, "GET", u, nil)

	return body.DNSZone, resp, err
}

func (s *DNSZonesClient) GetByName(
	ctx context.Context,
	name string,
) (*DNSZone, *Response, error) {
	qs := url.Values{"dns_zone[name]": []string{name}}
	u := &url.URL{Path: "dns/zones/_", RawQuery: qs.Encode()}

	body, resp, err := s.doRequest(ctx, "GET", u, nil)

	return body.DNSZone, resp, err
}

func (s *DNSZonesClient) Create(
	ctx context.Context,
	orgID string,
	zone *DNSZoneArguments,
) (*DNSZone, *Response, error) {
	if zone == nil {
		return nil, nil, errors.New("nil zone arguments")
	}

	u := &url.URL{Path: fmt.Sprintf("organizations/%s/dns/zones", orgID)}
	reqBody := &createDNSZoneRequest{
		Details: &DNSZoneDetails{
			Name: zone.Name,
			TTL:  zone.TTL,
		},
		SkipVerfication: zone.SkipVerfication,
	}

	body, resp, err := s.doRequest(ctx, "POST", u, reqBody)

	return body.DNSZone, resp, err
}

func (s *DNSZonesClient) Delete(
	ctx context.Context,
	id string,
) (*DNSZone, *Response, error) {
	u := &url.URL{Path: fmt.Sprintf("dns/zones/%s", id)}
	body, resp, err := s.doRequest(ctx, "DELETE", u, nil)

	return body.DNSZone, resp, err
}

func (s *DNSZonesClient) VerificationDetails(
	ctx context.Context,
	id string,
) (*DNSZoneVerificationDetails, *Response, error) {
	u := &url.URL{Path: fmt.Sprintf("dns/zones/%s/verification_details", id)}
	body, resp, err := s.doRequest(ctx, "GET", u, nil)

	return body.VerificationDetails, resp, err
}

func (s *DNSZonesClient) Verify(
	ctx context.Context,
	id string,
) (*DNSZone, *Response, error) {
	u := &url.URL{Path: fmt.Sprintf("dns/zones/%s/verify", id)}
	body, resp, err := s.doRequest(ctx, "POST", u, nil)

	return body.DNSZone, resp, err
}

func (s *DNSZonesClient) UpdateTTL(
	ctx context.Context,
	id string,
	ttl int,
) (*DNSZone, *Response, error) {
	u := &url.URL{Path: fmt.Sprintf("dns/zones/%s/update_ttl", id)}
	reqBody := &updateDNSZoneTTLRequest{TTL: ttl}
	body, resp, err := s.doRequest(ctx, "POST", u, reqBody)

	return body.DNSZone, resp, err
}

func (s *DNSZonesClient) doRequest(
	ctx context.Context,
	method string,
	u *url.URL,
	body interface{},
) (*dnsZoneResponseBody, *Response, error) {
	u = s.basePath.ResolveReference(u)
	respBody := &dnsZoneResponseBody{}
	resp := &Response{}

	req, err := s.client.NewRequestWithContext(ctx, method, u, body)
	if err == nil {
		resp, err = s.client.Do(req, respBody)
	}

	return respBody, resp, err
}
