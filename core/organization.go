package core

import (
	"context"
	"net/url"

	"github.com/augurysys/timestamp"
	"github.com/krystal/go-katapult"
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

func (o *Organization) Ref() OrganizationRef {
	return OrganizationRef{ID: o.ID}
}

type OrganizationRef struct {
	ID        string `json:"id,omitempty"`
	SubDomain string `json:"sub_domain,omitempty"`
}

func (or OrganizationRef) queryValues() *url.Values {
	v := &url.Values{}

	switch {
	case or.ID != "":
		v.Set("organization[id]", or.ID)
	case or.SubDomain != "":
		v.Set("organization[sub_domain]", or.SubDomain)
	}

	return v
}

type OrganizationManagedArguments struct {
	Name      string
	SubDomain string
}

type organizationCreateManagedRequest struct {
	Organization OrganizationRef `json:"organization"`
	Name         string          `json:"name"`
	SubDomain    string          `json:"sub_domain"`
}

type organizationsResponseBody struct {
	Organization  *Organization   `json:"organization,omitempty"`
	Organizations []*Organization `json:"organizations,omitempty"`
}

type OrganizationsClient struct {
	client   RequestMaker
	basePath *url.URL
}

func NewOrganizationsClient(rm RequestMaker) *OrganizationsClient {
	return &OrganizationsClient{
		client:   rm,
		basePath: &url.URL{Path: "/core/v1/"},
	}
}

func (s *OrganizationsClient) List(
	ctx context.Context,
) ([]*Organization, *katapult.Response, error) {
	u := &url.URL{Path: "organizations"}
	body, resp, err := s.doRequest(ctx, "GET", u, nil)

	return body.Organizations, resp, err
}

func (s *OrganizationsClient) Get(
	ctx context.Context,
	ref OrganizationRef,
) (*Organization, *katapult.Response, error) {
	qs := ref.queryValues()
	u := &url.URL{Path: "organizations/_", RawQuery: qs.Encode()}

	body, resp, err := s.doRequest(ctx, "GET", u, nil)

	return body.Organization, resp, err
}

func (s *OrganizationsClient) GetByID(
	ctx context.Context,
	id string,
) (*Organization, *katapult.Response, error) {
	return s.Get(ctx, OrganizationRef{ID: id})
}

func (s *OrganizationsClient) GetBySubDomain(
	ctx context.Context,
	subDomain string,
) (*Organization, *katapult.Response, error) {
	return s.Get(ctx, OrganizationRef{SubDomain: subDomain})
}

func (s *OrganizationsClient) CreateManaged(
	ctx context.Context,
	parent OrganizationRef,
	args *OrganizationManagedArguments,
) (*Organization, *katapult.Response, error) {
	u := &url.URL{Path: "organizations/_/managed"}
	reqBody := &organizationCreateManagedRequest{
		Organization: parent,
	}

	if args != nil {
		reqBody.Name = args.Name
		reqBody.SubDomain = args.SubDomain
	}

	body, resp, err := s.doRequest(ctx, "POST", u, reqBody)

	return body.Organization, resp, err
}

func (s *OrganizationsClient) doRequest(
	ctx context.Context,
	method string,
	u *url.URL,
	body interface{},
) (*organizationsResponseBody, *katapult.Response, error) {
	u = s.basePath.ResolveReference(u)
	respBody := &organizationsResponseBody{}

	req := katapult.NewRequest(method, u, body)
	resp, err := s.client.Do(ctx, req, respBody)
	if resp == nil {
		resp = katapult.NewResponse(nil)
	}

	return respBody, resp, err
}
