package core

import (
	"context"
	"fmt"
	"net/url"
	"strings"

	"github.com/krystal/go-katapult"
)

type DNSZone struct {
	ID                 string `json:"id,omitempty"`
	Name               string `json:"name,omitempty"`
	TTL                int    `json:"ttl,omitempty"`
	Verified           bool   `json:"verified,omitempty"`
	InfrastructureZone bool   `json:"infrastructure_zone,omitempty"`
}

// NewDNSZoneLookup takes a string that is a DNSZone ID or Name, returning a
// empty *DNSZone struct with either the ID or Name field populated with the
// given value. This struct is suitable as input to other methods which accept a
// *DNSZone as input.
func NewDNSZoneLookup(
	idOrName string,
) (lr *DNSZone, f FieldName) {
	if strings.HasPrefix(idOrName, "dnszone_") {
		return &DNSZone{ID: idOrName}, IDField
	}

	return &DNSZone{Name: idOrName}, NameField
}

func (s *DNSZone) lookupReference() *DNSZone {
	if s == nil {
		return nil
	}

	lr := &DNSZone{ID: s.ID}
	if lr.ID == "" {
		lr.Name = s.Name
	}

	return lr
}

func (s *DNSZone) queryValues() *url.Values {
	v := &url.Values{}

	if s != nil {
		switch {
		case s.ID != "":
			v.Set("dns_zone[id]", s.ID)
		case s.Name != "":
			v.Set("dns_zone[name]", s.Name)
		}
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
	org *Organization,
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
	idOrName string,
) (*DNSZone, *katapult.Response, error) {
	if _, f := NewDNSZoneLookup(idOrName); f == IDField {
		return s.GetByID(ctx, idOrName)
	}

	return s.GetByName(ctx, idOrName)
}

func (s *DNSZonesClient) GetByID(
	ctx context.Context,
	id string,
) (*DNSZone, *katapult.Response, error) {
	u := &url.URL{Path: fmt.Sprintf("dns/zones/%s", id)}
	body, resp, err := s.doRequest(ctx, "GET", u, nil)

	return body.DNSZone, resp, err
}

func (s *DNSZonesClient) GetByName(
	ctx context.Context,
	name string,
) (*DNSZone, *katapult.Response, error) {
	qs := url.Values{"dns_zone[name]": []string{name}}
	u := &url.URL{Path: "dns/zones/_", RawQuery: qs.Encode()}

	body, resp, err := s.doRequest(ctx, "GET", u, nil)

	return body.DNSZone, resp, err
}

func (s *DNSZonesClient) Create(
	ctx context.Context,
	org *Organization,
	args *DNSZoneArguments,
) (*DNSZone, *katapult.Response, error) {
	u := &url.URL{Path: "organizations/_/dns/zones"}
	reqBody := &dnsZoneCreateRequest{
		Organization: org.lookupReference(),
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
	zone *DNSZone,
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
	zone *DNSZone,
) (*DNSZone, *katapult.Response, error) {
	u := &url.URL{Path: "dns/zones/_/verify"}
	reqBody := &dnsZoneVerifyRequest{
		DNSZone: zone.lookupReference(),
	}
	body, resp, err := s.doRequest(ctx, "POST", u, reqBody)

	return body.DNSZone, resp, err
}

func (s *DNSZonesClient) UpdateTTL(
	ctx context.Context,
	zone *DNSZone,
	ttl int,
) (*DNSZone, *katapult.Response, error) {
	u := &url.URL{Path: "dns/zones/_/update_ttl"}
	reqBody := &dnsZoneUpdateTTLRequest{
		DNSZone: zone.lookupReference(),
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
