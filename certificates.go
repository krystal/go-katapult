package katapult

import (
	"context"
	"fmt"
	"net/url"

	"github.com/augurysys/timestamp"
)

type CertificatesService struct {
	client   *apiClient
	basePath *url.URL
}

func newCertificatesService(c *apiClient) *CertificatesService {
	return &CertificatesService{
		client:   c,
		basePath: &url.URL{Path: "/core/v1/"},
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
	u := &url.URL{
		Path:     fmt.Sprintf("organizations/%s/certificates", orgID),
		RawQuery: opts.Values().Encode(),
	}

	body, resp, err := s.doRequest(ctx, "GET", u, nil)
	resp.Pagination = body.Pagination

	return body.Certificates, resp, err
}

func (s CertificatesService) Get(
	ctx context.Context,
	id string,
) (*Certificate, *Response, error) {
	u := &url.URL{Path: fmt.Sprintf("certificates/%s", id)}
	body, resp, err := s.doRequest(ctx, "GET", u, nil)

	return body.Certificate, resp, err
}

func (s *CertificatesService) doRequest(
	ctx context.Context,
	method string,
	u *url.URL,
	body interface{},
) (*certificatesResponseBody, *Response, error) {
	u = s.basePath.ResolveReference(u)
	respBody := &certificatesResponseBody{}
	resp := &Response{}

	req, err := s.client.NewRequestWithContext(ctx, method, u, body)
	if err == nil {
		resp, err = s.client.Do(req, respBody)
	}

	return respBody, resp, err
}
