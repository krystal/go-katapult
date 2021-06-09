package core

import (
	"context"
	"net/url"

	"github.com/krystal/go-katapult"
)

type LoadBalancer struct {
	ID                    string       `json:"id,omitempty"`
	Name                  string       `json:"name,omitempty"`
	ResourceType          ResourceType `json:"resource_type,omitempty"`
	ResourceIDs           []string     `json:"resource_ids,omitempty"`
	IPAddress             *IPAddress   `json:"ip_address,omitempty"`
	HTTPSRedirect         bool         `json:"https_redirect,omitempty"`
	BackendCertificate    string       `json:"backend_certificate,omitempty"`
	BackendCertificateKey string       `json:"backend_certificate_key,omitempty"`
	DataCenter            *DataCenter  `json:"-"`
}

func (lb *LoadBalancer) Ref() LoadBalancerRef {
	return LoadBalancerRef{ID: lb.ID}
}

// LoadBalancerRef allows a reference to a load balancer
type LoadBalancerRef struct {
	ID string `json:"id,omitempty"`
}

func (lbr LoadBalancerRef) queryValues() *url.Values {
	v := &url.Values{}
	v.Set("load_balancer[id]", lbr.ID)

	return v
}

type LoadBalancerCreateArguments struct {
	DataCenter    DataCenterRef `json:"data_center"`
	Name          string        `json:"name,omitempty"`
	ResourceType  ResourceType  `json:"resource_type,omitempty"`
	ResourceIDs   *[]string     `json:"resource_ids,omitempty"`
	HTTPSRedirect *bool         `json:"https_redirect,omitempty"`
}

type LoadBalancerUpdateArguments struct {
	Name          string       `json:"name,omitempty"`
	ResourceType  ResourceType `json:"resource_type,omitempty"`
	ResourceIDs   *[]string    `json:"resource_ids,omitempty"`
	HTTPSRedirect *bool        `json:"https_redirect,omitempty"`
}

type loadBalancerCreateRequest struct {
	Organization OrganizationRef              `json:"organization"`
	Properties   *LoadBalancerCreateArguments `json:"properties,omitempty"`
}

type loadBalancerUpdateRequest struct {
	LoadBalancer LoadBalancerRef              `json:"load_balancer"`
	Properties   *LoadBalancerUpdateArguments `json:"properties,omitempty"`
}

type loadBalancersResponseBody struct {
	Pagination    *katapult.Pagination `json:"pagination,omitempty"`
	LoadBalancer  *LoadBalancer        `json:"load_balancer,omitempty"`
	LoadBalancers []*LoadBalancer      `json:"load_balancers,omitempty"`
}

type LoadBalancersClient struct {
	client   RequestMaker
	basePath *url.URL
}

func NewLoadBalancersClient(rm RequestMaker) *LoadBalancersClient {
	return &LoadBalancersClient{
		client:   rm,
		basePath: &url.URL{Path: "/core/v1/"},
	}
}

func (s *LoadBalancersClient) List(
	ctx context.Context,
	org OrganizationRef,
	opts *ListOptions,
) ([]*LoadBalancer, *katapult.Response, error) {
	qs := queryValues(org, opts)
	u := &url.URL{
		Path:     "organizations/_/load_balancers",
		RawQuery: qs.Encode(),
	}

	body, resp, err := s.doRequest(ctx, "GET", u, nil)
	resp.Pagination = body.Pagination

	return body.LoadBalancers, resp, err
}

func (s *LoadBalancersClient) Get(
	ctx context.Context,
	ref LoadBalancerRef,
) (*LoadBalancer, *katapult.Response, error) {
	u := &url.URL{
		Path:     "load_balancers/_",
		RawQuery: ref.queryValues().Encode(),
	}

	body, resp, err := s.doRequest(ctx, "GET", u, nil)

	return body.LoadBalancer, resp, err
}

func (s *LoadBalancersClient) GetByID(
	ctx context.Context,
	id string,
) (*LoadBalancer, *katapult.Response, error) {
	return s.Get(ctx, LoadBalancerRef{ID: id})
}

func (s *LoadBalancersClient) Create(
	ctx context.Context,
	org OrganizationRef,
	args *LoadBalancerCreateArguments,
) (*LoadBalancer, *katapult.Response, error) {
	u := &url.URL{Path: "organizations/_/load_balancers"}
	reqBody := &loadBalancerCreateRequest{
		Organization: org,
		Properties:   args,
	}

	body, resp, err := s.doRequest(ctx, "POST", u, reqBody)

	return body.LoadBalancer, resp, err
}

func (s *LoadBalancersClient) Update(
	ctx context.Context,
	lb LoadBalancerRef,
	args *LoadBalancerUpdateArguments,
) (*LoadBalancer, *katapult.Response, error) {
	u := &url.URL{Path: "load_balancers/_"}
	reqBody := &loadBalancerUpdateRequest{
		LoadBalancer: lb,
		Properties:   args,
	}

	body, resp, err := s.doRequest(ctx, "PATCH", u, reqBody)

	return body.LoadBalancer, resp, err
}

func (s *LoadBalancersClient) Delete(
	ctx context.Context,
	lb LoadBalancerRef,
) (*LoadBalancer, *katapult.Response, error) {
	u := &url.URL{
		Path:     "load_balancers/_",
		RawQuery: lb.queryValues().Encode(),
	}
	body, resp, err := s.doRequest(ctx, "DELETE", u, nil)

	return body.LoadBalancer, resp, err
}

func (s *LoadBalancersClient) doRequest(
	ctx context.Context,
	method string,
	u *url.URL,
	body interface{},
) (*loadBalancersResponseBody, *katapult.Response, error) {
	u = s.basePath.ResolveReference(u)
	respBody := &loadBalancersResponseBody{}

	req := katapult.NewRequest(method, u, body)
	resp, err := s.client.Do(ctx, req, respBody)
	if resp == nil {
		resp = katapult.NewResponse(nil)
	}

	return respBody, resp, err
}
