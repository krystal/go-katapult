package core

import (
	"context"
	"fmt"
	"net/url"

	"github.com/krystal/go-katapult"

	"github.com/augurysys/timestamp"
)

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
	Pagination   *katapult.Pagination `json:"pagination,omitempty"`
	Certificate  *Certificate         `json:"certificate,omitempty"`
	Certificates []*Certificate       `json:"certificates,omitempty"`
}

type CertificatesClient struct {
	client   RequestMaker
	basePath *url.URL
}

func NewCertificatesClient(rm RequestMaker) *CertificatesClient {
	return &CertificatesClient{
		client:   rm,
		basePath: &url.URL{Path: "/core/v1/"},
	}
}

func (s *CertificatesClient) List(
	ctx context.Context,
	org *Organization,
	opts *ListOptions,
) ([]*Certificate, *katapult.Response, error) {
	qs := queryValues(org, opts)
	u := &url.URL{
		Path:     "organizations/_/certificates",
		RawQuery: qs.Encode(),
	}

	body, resp, err := s.doRequest(ctx, "GET", u, nil)
	resp.Pagination = body.Pagination

	return body.Certificates, resp, err
}

func (s *CertificatesClient) Get(
	ctx context.Context,
	id string,
) (*Certificate, *katapult.Response, error) {
	return s.GetByID(ctx, id)
}

func (s *CertificatesClient) GetByID(
	ctx context.Context,
	id string,
) (*Certificate, *katapult.Response, error) {
	u := &url.URL{Path: fmt.Sprintf("certificates/%s", id)}
	body, resp, err := s.doRequest(ctx, "GET", u, nil)

	return body.Certificate, resp, err
}

func (s *CertificatesClient) doRequest(
	ctx context.Context,
	method string,
	u *url.URL,
	body interface{},
) (*certificatesResponseBody, *katapult.Response, error) {
	u = s.basePath.ResolveReference(u)
	respBody := &certificatesResponseBody{}
	resp := katapult.NewResponse(nil)

	req, err := s.client.NewRequestWithContext(ctx, method, u, body)
	if err == nil {
		resp, err = s.client.Do(req, respBody)
	}

	return respBody, resp, err
}
