package katapult

import (
	"context"
	"fmt"
	"net/url"

	"github.com/augurysys/timestamp"
	"github.com/google/go-querystring/query"
)

type CertificatesService struct {
	*service
	path *url.URL
}

func NewCertificatesService(s *service) *CertificatesService {
	return &CertificatesService{
		service: s,
		path:    &url.URL{Path: "/core/v1/"},
	}
}

type Certificate struct {
	ID                  string               `json:"id,omitempty"`
	Name                string               `json:"name,omitempty"`
	AdditionalNames     []string             `json:"additional_names,omitempty"`
	Issuer              string               `json:"issuer,omitempty"`
	State               string               `json:"state,omitempty"`
	CreatedAt           *timestamp.Timestamp `json:"created_at,omitempty"`
	ExpiresAt           *timestamp.Timestamp `json:"expires_at,omitempty"`
	LastIssuedAt        *timestamp.Timestamp `json:"last_issued_at,omitempty"`
	IssueError          string               `json:"issue_error,omitempty"`
	AuthorizationMethod string               `json:"authorization_method,omitempty"`
	CertificateAPIURL   string               `json:"certificate_api_url,omitempty"`
	Certificate         string               `json:"certificate,omitempty"`
	Chain               string               `json:"chain,omitempty"`
	PrivateKey          string               `json:"private_key,omitempty"`
}

type certificatesResponseBody struct {
	Pagination   *Pagination    `json:"pagination,omitempty"`
	Certificate  *Certificate   `json:"certificate,omitempty"`
	Certificates []*Certificate `json:"certificates,omitempty"`
}

func (s CertificatesService) List(
	ctx context.Context,
	orgID string,
	opts *ListOptions,
) ([]*Certificate, *Response, error) {
	u, err := s.path.Parse(fmt.Sprintf("organizations/%s/certificates", orgID))
	if err != nil {
		return nil, nil, err
	}

	qs, err := query.Values(opts)
	if err != nil {
		return nil, nil, err
	}
	u.RawQuery = qs.Encode()

	req, err := s.client.NewRequestWithContext(ctx, "GET", u.String(), nil)
	if err != nil {
		return nil, nil, err
	}

	var body certificatesResponseBody
	resp, err := s.client.Do(req, &body)
	if err != nil {
		return nil, resp, err
	}
	resp.Pagination = body.Pagination

	return body.Certificates, resp, nil
}

func (s CertificatesService) Get(
	ctx context.Context,
	id string,
) (*Certificate, *Response, error) {
	u, err := s.path.Parse(fmt.Sprintf("certificates/%s", id))
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequestWithContext(ctx, "GET", u.Path, nil)
	if err != nil {
		return nil, nil, err
	}

	var body certificatesResponseBody
	resp, err := s.client.Do(req, &body)
	if err != nil {
		return nil, resp, err
	}

	return body.Certificate, resp, nil
}
