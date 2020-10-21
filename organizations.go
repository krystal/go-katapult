package katapult

import (
	"context"
	"fmt"
	"net/url"

	"github.com/augurysys/timestamp"
)

type OrganizationsService struct {
	*service
	path *url.URL
}

func NewOrganizationsService(s *service) *OrganizationsService {
	return &OrganizationsService{
		service: s,
		path:    &url.URL{Path: "/core/v1/"},
	}
}

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

type OrganizationsResponseBody struct {
	Organization  *Organization   `json:"organization,omitempty"`
	Organizations []*Organization `json:"organizations,omitempty"`
}

func (s *OrganizationsService) List(
	ctx context.Context,
) ([]*Organization, *Response, error) {
	u, err := s.path.Parse("organizations")
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequestWithContext(ctx, "GET", u.Path, nil)
	if err != nil {
		return nil, nil, err
	}

	var body *OrganizationsResponseBody
	resp, err := s.client.Do(req, &body)
	if err != nil {
		return nil, resp, err
	}

	return body.Organizations, resp, nil
}

func (s *OrganizationsService) Get(
	ctx context.Context,
	subDomainOrID string,
) (*Organization, *Response, error) {
	u, err := s.path.Parse(fmt.Sprintf("organizations/%s", subDomainOrID))
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequestWithContext(ctx, "GET", u.Path, nil)
	if err != nil {
		return nil, nil, err
	}

	var body *OrganizationsResponseBody
	resp, err := s.client.Do(req, &body)
	if err != nil {
		return nil, resp, err
	}

	return body.Organization, resp, nil
}

func (s *OrganizationsService) CreateManaged(
	ctx context.Context,
	parent *Organization,
	name string,
	subDomain string,
) (*Organization, *Response, error) {
	u, err := s.path.Parse(fmt.Sprintf("organizations/%s/managed", parent.ID))
	if err != nil {
		return nil, nil, err
	}

	reqBody := &Organization{
		Name:      name,
		SubDomain: subDomain,
	}

	req, err := s.client.NewRequestWithContext(ctx, "POST", u.Path, reqBody)
	if err != nil {
		return nil, nil, err
	}

	var body *OrganizationsResponseBody
	resp, err := s.client.Do(req, &body)
	if err != nil {
		return nil, resp, err
	}

	return body.Organization, resp, nil
}
