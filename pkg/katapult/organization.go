package katapult

import (
	"context"
	"fmt"
	"net/url"

	"github.com/augurysys/timestamp"
)

type Organization struct {
	ID                   string               `json:"id,omitempty"`
	Name                 string               `json:"name,omitempty"`
	SubDomain            string               `json:"sub_domain,omitempty"`
	InfrastructureDomain string               `json:"infrastructure_domain,omitempty"`
	Personal             bool                 `json:"personal,omitempty"`
	CreatedAt            *timestamp.Timestamp `json:"created_at,omitempty"`
	Suspended            bool                 `json:"suspended,omitempty"`
	Managed              bool                 `json:"managed,omitempty"`
	BillingName          string               `json:"billing_name,omitempty"`
	Address1             string               `json:"address1,omitempty"`
	Address2             string               `json:"address2,omitempty"`
	Address3             string               `json:"address3,omitempty"`
	Address4             string               `json:"address4,omitempty"`
	Postcode             string               `json:"postcode,omitempty"`
	VatNumber            string               `json:"vat_number,omitempty"`
	Currency             *Currency            `json:"currency,omitempty"`
	Country              *Country             `json:"country,omitempty"`
	CountryState         *CountryState        `json:"country_state,omitempty"`
}

// LookupReference returns a new *Organization stripped down to just ID or
// SubDomain fields, making it suitable for endpoints which require a reference
// to a Organization by ID or SubDomain.
func (s *Organization) LookupReference() *Organization {
	if s == nil {
		return nil
	}

	lr := &Organization{ID: s.ID}
	if lr.ID == "" {
		lr.SubDomain = s.SubDomain
	}

	return lr
}

type organizationsResponseBody struct {
	Organization  *Organization   `json:"organization,omitempty"`
	Organizations []*Organization `json:"organizations,omitempty"`
}

type OrganizationsClient struct {
	client   *apiClient
	basePath *url.URL
}

func newOrganizationsClient(c *apiClient) *OrganizationsClient {
	return &OrganizationsClient{
		client:   c,
		basePath: &url.URL{Path: "/core/v1/"},
	}
}

func (s *OrganizationsClient) List(
	ctx context.Context,
) ([]*Organization, *Response, error) {
	u := &url.URL{Path: "organizations"}
	body, resp, err := s.doRequest(ctx, "GET", u, nil)

	return body.Organizations, resp, err
}

func (s *OrganizationsClient) Get(
	ctx context.Context,
	id string,
) (*Organization, *Response, error) {
	u := &url.URL{Path: fmt.Sprintf("organizations/%s", id)}
	body, resp, err := s.doRequest(ctx, "GET", u, nil)

	return body.Organization, resp, err
}

func (s *OrganizationsClient) GetBySubDomain(
	ctx context.Context,
	subDomain string,
) (*Organization, *Response, error) {
	qs := url.Values{"organization[sub_domain]": []string{subDomain}}
	u := &url.URL{Path: "organizations/_", RawQuery: qs.Encode()}

	body, resp, err := s.doRequest(ctx, "GET", u, nil)

	return body.Organization, resp, err
}

func (s *OrganizationsClient) CreateManaged(
	ctx context.Context,
	parentID string,
	name string,
	subDomain string,
) (*Organization, *Response, error) {
	u := &url.URL{Path: fmt.Sprintf("organizations/%s/managed", parentID)}
	reqBody := &Organization{Name: name, SubDomain: subDomain}
	body, resp, err := s.doRequest(ctx, "POST", u, reqBody)

	return body.Organization, resp, err
}

func (s *OrganizationsClient) doRequest(
	ctx context.Context,
	method string,
	u *url.URL,
	body interface{},
) (*organizationsResponseBody, *Response, error) {
	u = s.basePath.ResolveReference(u)
	respBody := &organizationsResponseBody{}
	resp := &Response{}

	req, err := s.client.NewRequestWithContext(ctx, method, u, body)
	if err == nil {
		resp, err = s.client.Do(req, respBody)
	}

	return respBody, resp, err
}
