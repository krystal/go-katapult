package katapult

import (
	"context"
	"fmt"
	"net/url"

	"github.com/google/go-querystring/query"
)

type DNSZonesService struct {
	*service
	path *url.URL
}

func NewDNSZonesService(s *service) *DNSZonesService {
	return &DNSZonesService{
		service: s,
		path:    &url.URL{Path: "/core/v1/"},
	}
}

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

func (s *DNSZonesService) List(
	ctx context.Context,
	orgID string,
	opts *ListOptions,
) ([]*DNSZone, *Response, error) {
	u := &url.URL{Path: fmt.Sprintf("organizations/%s/dns/zones", orgID)}

	qs, err := query.Values(opts)
	if err != nil {
		return nil, nil, err
	}
	u.RawQuery = qs.Encode()

	body, resp, err := s.doRequest(ctx, "GET", u.String(), nil)
	resp.Pagination = body.Pagination

	return body.DNSZones, resp, err
}

func (s *DNSZonesService) Get(
	ctx context.Context,
	id string,
) (*DNSZone, *Response, error) {
	u := fmt.Sprintf("dns/zones/%s", id)
	body, resp, err := s.doRequest(ctx, "GET", u, nil)

	return body.DNSZone, resp, err
}

func (s *DNSZonesService) Create(
	ctx context.Context,
	orgID string,
	zone *DNSZoneArguments,
) (*DNSZone, *Response, error) {
	u := fmt.Sprintf("organizations/%s/dns/zones", orgID)
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

func (s *DNSZonesService) Delete(
	ctx context.Context,
	id string,
) (*DNSZone, *Response, error) {
	u := fmt.Sprintf("dns/zones/%s", id)
	body, resp, err := s.doRequest(ctx, "DELETE", u, nil)

	return body.DNSZone, resp, err
}

func (s *DNSZonesService) VerificationDetails(
	ctx context.Context,
	id string,
) (*DNSZoneVerificationDetails, *Response, error) {
	u := fmt.Sprintf("dns/zones/%s/verification_details", id)
	body, resp, err := s.doRequest(ctx, "GET", u, nil)

	return body.VerificationDetails, resp, err
}

func (s *DNSZonesService) Verify(
	ctx context.Context,
	id string,
) (*DNSZone, *Response, error) {
	u := fmt.Sprintf("dns/zones/%s/verify", id)
	body, resp, err := s.doRequest(ctx, "POST", u, nil)

	return body.DNSZone, resp, err
}

func (s *DNSZonesService) UpdateTTL(
	ctx context.Context,
	id string,
	ttl int,
) (*DNSZone, *Response, error) {
	u := fmt.Sprintf("dns/zones/%s/update_ttl", id)
	reqBody := &updateDNSZoneTTLRequest{TTL: ttl}
	body, resp, err := s.doRequest(ctx, "POST", u, reqBody)

	return body.DNSZone, resp, err
}

func (s *DNSZonesService) doRequest(
	ctx context.Context,
	method string,
	urlStr string,
	body interface{},
) (*dnsZoneResponseBody, *Response, error) {
	u, err := s.path.Parse(urlStr)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequestWithContext(ctx, method, u.String(), body)
	if err != nil {
		return nil, nil, err
	}

	var respBody dnsZoneResponseBody
	resp, err := s.client.Do(req, &respBody)

	return &respBody, resp, err
}
