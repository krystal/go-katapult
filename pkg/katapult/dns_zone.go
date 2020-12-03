package katapult

import (
	"context"
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

// LookupReference returns a new *DNSZone stripped down to just ID or Name
// fields, making it suitable for endpoints which require a reference to a
// DNSZone by ID or Name.
func (s *DNSZone) LookupReference() *DNSZone {
	if s == nil {
		return nil
	}

	lr := &DNSZone{ID: s.ID}
	if lr.ID == "" {
		lr.Name = s.Name
	}

	return lr
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
	Organization    *Organization   `json:"organization"`
	Details         *DNSZoneDetails `json:"details"`
	SkipVerfication bool            `json:"skip_verification"`
}

type dnsZoneVerifyRequest struct {
	DNSZone *DNSZone `json:"dns_zone"`
}

type dnsZoneUpdateTTLRequest struct {
	DNSZone *DNSZone `json:"dns_zone"`
	TTL     int      `json:"ttl"`
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
	org *Organization,
	opts *ListOptions,
) ([]*DNSZone, *Response, error) {
	if org == nil {
		org = &Organization{ID: "_"}
	}

	u := &url.URL{
		Path:     fmt.Sprintf("organizations/%s/dns/zones", org.ID),
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
	org *Organization,
	args *DNSZoneArguments,
) (*DNSZone, *Response, error) {
	u := &url.URL{Path: "organizations/_/dns/zones"}
	reqBody := &dnsZoneCreateRequest{
		Organization: org.LookupReference(),
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
	zone *DNSZone,
) (*DNSZone, *Response, error) {
	if zone == nil {
		zone = &DNSZone{ID: "_"}
	}

	u := &url.URL{Path: fmt.Sprintf("dns/zones/%s", zone.ID)}
	body, resp, err := s.doRequest(ctx, "DELETE", u, nil)

	return body.DNSZone, resp, err
}

func (s *DNSZonesClient) VerificationDetails(
	ctx context.Context,
	zone *DNSZone,
) (*DNSZoneVerificationDetails, *Response, error) {
	if zone == nil {
		zone = &DNSZone{ID: "_"}
	}

	u := &url.URL{
		Path: fmt.Sprintf("dns/zones/%s/verification_details", zone.ID),
	}
	body, resp, err := s.doRequest(ctx, "GET", u, nil)

	return body.VerificationDetails, resp, err
}

func (s *DNSZonesClient) Verify(
	ctx context.Context,
	zone *DNSZone,
) (*DNSZone, *Response, error) {
	u := &url.URL{Path: "dns/zones/_/verify"}
	reqBody := &dnsZoneVerifyRequest{
		DNSZone: zone.LookupReference(),
	}
	body, resp, err := s.doRequest(ctx, "POST", u, reqBody)

	return body.DNSZone, resp, err
}

func (s *DNSZonesClient) UpdateTTL(
	ctx context.Context,
	zone *DNSZone,
	ttl int,
) (*DNSZone, *Response, error) {
	u := &url.URL{Path: "dns/zones/_/update_ttl"}
	reqBody := &dnsZoneUpdateTTLRequest{
		DNSZone: zone.LookupReference(),
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
