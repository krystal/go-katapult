package katapult

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
)

type LoadBalancer struct {
	ID                    string       `json:"id,omitempty"`
	Name                  string       `json:"name,omitempty"`
	ResourceType          ResourceType `json:"resource_type,omitempty"`
	ResourceIDs           []string     `json:"-"`
	IPAddress             *IPAddress   `json:"ip_address,omitempty"`
	HTTPSRedirect         bool         `json:"https_redirect,omitempty"`
	BackendCertificate    string       `json:"backend_certificate,omitempty"`
	BackendCertificateKey string       `json:"backend_certificate_key,omitempty"`
	DataCenter            *DataCenter  `json:"-"`
}

func (s *LoadBalancer) lookupReference() *LoadBalancer {
	if s == nil {
		return nil
	}

	return &LoadBalancer{ID: s.ID}
}

func (s *LoadBalancer) queryValues() *url.Values {
	v := &url.Values{}

	if s != nil && s.ID != "" {
		v.Set("load_balancer[id]", s.ID)
	}

	return v
}

func (s *LoadBalancer) MarshalJSON() ([]byte, error) {
	type alias LoadBalancer
	resources := []*loadBalancerResource{}

	for _, id := range s.ResourceIDs {
		resources = append(resources, &loadBalancerResource{
			Type:  s.ResourceType.objectType(),
			Value: &loadBalancerResourceValue{ID: id},
		})
	}

	return json.Marshal(&struct {
		*alias
		Resources []*loadBalancerResource `json:"resources,omitempty"`
	}{
		alias:     (*alias)(s),
		Resources: resources,
	})
}

func (s *LoadBalancer) UnmarshalJSON(b []byte) error {
	type alias LoadBalancer
	aux := &struct {
		*alias
		Resources []*loadBalancerResource `json:"resources,omitempty"`
	}{
		alias: (*alias)(s),
	}

	if err := json.Unmarshal(b, &aux); err != nil {
		return err
	}

	for _, r := range aux.Resources {
		if r.Value != nil {
			s.ResourceIDs = append(s.ResourceIDs, r.Value.ID)
		}
	}

	return nil
}

type LoadBalancerArguments struct {
	Name         string       `json:"name,omitempty"`
	ResourceType ResourceType `json:"resource_type,omitempty"`
	ResourceIDs  []string     `json:"resource_ids"`
	DataCenter   *DataCenter  `json:"data_center,omitempty"`
}

func (
	s *LoadBalancerArguments,
) forRequest() *LoadBalancerArguments {
	if s == nil {
		return nil
	}

	args := *s
	args.DataCenter = s.DataCenter.lookupReference()

	return &args
}

type loadBalancerResource struct {
	Type  string                     `json:"type,omitempty"`
	Value *loadBalancerResourceValue `json:"value,omitempty"`
}

type loadBalancerResourceValue struct {
	ID   string `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
}

type loadBalancerCreateRequest struct {
	Organization *Organization          `json:"organization,omitempty"`
	Properties   *LoadBalancerArguments `json:"properties,omitempty"`
}

type loadBalancerUpdateRequest struct {
	LoadBalancer *LoadBalancer          `json:"load_balancer,omitempty"`
	Properties   *LoadBalancerArguments `json:"properties,omitempty"`
}

type loadBalancersResponseBody struct {
	Pagination    *Pagination     `json:"pagination,omitempty"`
	LoadBalancer  *LoadBalancer   `json:"load_balancer,omitempty"`
	LoadBalancers []*LoadBalancer `json:"load_balancers,omitempty"`
}

type LoadBalancersClient struct {
	client   *apiClient
	basePath *url.URL
}

func newLoadBalancersClient(c *apiClient) *LoadBalancersClient {
	return &LoadBalancersClient{
		client:   c,
		basePath: &url.URL{Path: "/core/v1/"},
	}
}

func (s *LoadBalancersClient) List(
	ctx context.Context,
	org *Organization,
	opts *ListOptions,
) ([]*LoadBalancer, *Response, error) {
	qs := queryValues(org, opts)
	u := &url.URL{
		Path:     "organizations/_/load_balancers",
		RawQuery: qs.Encode(),
	}

	body, resp, err := s.doRequest(ctx, "GET", u, nil)
	resp.Pagination = body.Pagination

	return body.LoadBalancers, resp, err
}

func (s LoadBalancersClient) Get(
	ctx context.Context,
	id string,
) (*LoadBalancer, *Response, error) {
	return s.GetByID(ctx, id)
}

func (s LoadBalancersClient) GetByID(
	ctx context.Context,
	id string,
) (*LoadBalancer, *Response, error) {
	u := &url.URL{Path: fmt.Sprintf("load_balancers/%s", id)}
	body, resp, err := s.doRequest(ctx, "GET", u, nil)

	return body.LoadBalancer, resp, err
}

func (s *LoadBalancersClient) Create(
	ctx context.Context,
	org *Organization,
	args *LoadBalancerArguments,
) (*LoadBalancer, *Response, error) {
	u := &url.URL{Path: "organizations/_/load_balancers"}
	reqBody := &loadBalancerCreateRequest{
		Organization: org.lookupReference(),
		Properties:   args.forRequest(),
	}

	body, resp, err := s.doRequest(ctx, "POST", u, reqBody)

	return body.LoadBalancer, resp, err
}

func (s *LoadBalancersClient) Update(
	ctx context.Context,
	lb *LoadBalancer,
	args *LoadBalancerArguments,
) (*LoadBalancer, *Response, error) {
	u := &url.URL{Path: "load_balancers/_"}
	reqBody := &loadBalancerUpdateRequest{
		LoadBalancer: lb.lookupReference(),
		Properties:   args.forRequest(),
	}

	body, resp, err := s.doRequest(ctx, "PATCH", u, reqBody)

	return body.LoadBalancer, resp, err
}

func (s *LoadBalancersClient) Delete(
	ctx context.Context,
	lb *LoadBalancer,
) (*LoadBalancer, *Response, error) {
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
) (*loadBalancersResponseBody, *Response, error) {
	u = s.basePath.ResolveReference(u)
	respBody := &loadBalancersResponseBody{}
	resp := &Response{}

	req, err := s.client.NewRequestWithContext(ctx, method, u, body)
	if err == nil {
		resp, err = s.client.Do(req, respBody)
	}

	return respBody, resp, err
}
