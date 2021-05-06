package core

import (
	"context"
	"fmt"
	"net/url"
	"strings"

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

// NewOrganizationLookup takes a string that is a Organization ID or SubDomain,
// returning a empty *Organization struct with either the ID or SubDomain field
// populated with the given value. This struct is suitable as input to other
// methods which accept a *Organization as input.
func NewOrganizationLookup(
	idOrSubDomain string,
) (lr *Organization, f FieldName) {
	if strings.HasPrefix(idOrSubDomain, "org_") {
		return &Organization{ID: idOrSubDomain}, IDField
	}

	return &Organization{SubDomain: idOrSubDomain}, SubDomainField
}

func (s *Organization) lookupReference() *Organization {
	if s == nil {
		return nil
	}

	lr := &Organization{ID: s.ID}
	if lr.ID == "" {
		lr.SubDomain = s.SubDomain
	}

	return lr
}

func (s *Organization) queryValues() *url.Values {
	v := &url.Values{}

	if s != nil {
		switch {
		case s.ID != "":
			v.Set("organization[id]", s.ID)
		case s.SubDomain != "":
			v.Set("organization[sub_domain]", s.SubDomain)
		}
	}

	return v
}

type OrganizationManagedArguments struct {
	Name      string
	SubDomain string
}

type organizationCreateManagedRequest struct {
	Organization *Organization `json:"organization"`
	Name         string        `json:"name"`
	SubDomain    string        `json:"sub_domain"`
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
	idOrSubDomain string,
) (*Organization, *katapult.Response, error) {
	if _, f := NewOrganizationLookup(idOrSubDomain); f == IDField {
		return s.GetByID(ctx, idOrSubDomain)
	}

	return s.GetBySubDomain(ctx, idOrSubDomain)
}

func (s *OrganizationsClient) GetByID(
	ctx context.Context,
	id string,
) (*Organization, *katapult.Response, error) {
	u := &url.URL{Path: fmt.Sprintf("organizations/%s", id)}
	body, resp, err := s.doRequest(ctx, "GET", u, nil)

	return body.Organization, resp, err
}

func (s *OrganizationsClient) GetBySubDomain(
	ctx context.Context,
	subDomain string,
) (*Organization, *katapult.Response, error) {
	qs := url.Values{"organization[sub_domain]": []string{subDomain}}
	u := &url.URL{Path: "organizations/_", RawQuery: qs.Encode()}

	body, resp, err := s.doRequest(ctx, "GET", u, nil)

	return body.Organization, resp, err
}

func (s *OrganizationsClient) CreateManaged(
	ctx context.Context,
	parent *Organization,
	args *OrganizationManagedArguments,
) (*Organization, *katapult.Response, error) {
	u := &url.URL{Path: "organizations/_/managed"}
	reqBody := &organizationCreateManagedRequest{
		Organization: parent.lookupReference(),
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
	resp := katapult.NewResponse(nil)

	req, err := s.client.NewRequestWithContext(ctx, method, u, body)
	if err == nil {
		resp, err = s.client.Do(req, respBody)
	}

	return respBody, resp, err
}
