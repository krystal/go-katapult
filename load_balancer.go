package katapult

import (
	"context"
	"encoding/json"
	"errors"
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

func (s *LoadBalancer) MarshalJSON() ([]byte, error) {
	type alias LoadBalancer
	resources := []*loadBalancerResource{}

	for _, id := range s.ResourceIDs {
		resources = append(resources, &loadBalancerResource{
			Type:   s.ResourceType.objectType(),
			Object: &loadBalancerResourceObject{ID: id},
		})
	}

	return json.Marshal(&struct {
		*alias
		Resources []*loadBalancerResource `json:"resource,omitempty"`
	}{
		alias:     (*alias)(s),
		Resources: resources,
	})
}

func (s *LoadBalancer) UnmarshalJSON(b []byte) error {
	type alias LoadBalancer
	aux := &struct {
		*alias
		Resources []*loadBalancerResource `json:"resource,omitempty"`
	}{
		alias: (*alias)(s),
	}

	if err := json.Unmarshal(b, &aux); err != nil {
		return err
	}

	for _, r := range aux.Resources {
		s.ResourceIDs = append(s.ResourceIDs, r.Object.ID)
	}

	return nil
}

type loadBalancerResource struct {
	Type   string                      `json:"type,omitempty"`
	Object *loadBalancerResourceObject `json:"object,omitempty"`
}

type loadBalancerResourceObject struct {
	ID   string `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
}

type loadBalancersResponseBody struct {
	Pagination    *Pagination     `json:"pagination,omitempty"`
	LoadBalancer  *LoadBalancer   `json:"load_balancer,omitempty"`
	LoadBalancers []*LoadBalancer `json:"load_balancers,omitempty"`
}

type LoadBalancerArguments struct {
	Name         string       `json:"name,omitempty"`
	ResourceType ResourceType `json:"resource_type,omitempty"`
	ResourceIDs  []string     `json:"resource_ids"`
	DataCenter   *DataCenter  `json:"data_center,omitempty"`
}

func (s *LoadBalancerArguments) MarshalJSON() ([]byte, error) {
	// do not perform destructive changes directly against s
	type alias LoadBalancerArguments
	args := alias(*s)

	if args.DataCenter != nil {
		args.DataCenter = &DataCenter{ID: s.DataCenter.ID}
	}

	return json.Marshal(&args)
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
	orgID string,
	opts *ListOptions,
) ([]*LoadBalancer, *Response, error) {
	u := &url.URL{
		Path:     fmt.Sprintf("organizations/%s/load_balancers", orgID),
		RawQuery: opts.Values().Encode(),
	}

	body, resp, err := s.doRequest(ctx, "GET", u, nil)
	resp.Pagination = body.Pagination

	return body.LoadBalancers, resp, err
}

func (s LoadBalancersClient) Get(
	ctx context.Context,
	id string,
) (*LoadBalancer, *Response, error) {
	u := &url.URL{Path: fmt.Sprintf("load_balancers/%s", id)}
	body, resp, err := s.doRequest(ctx, "GET", u, nil)

	return body.LoadBalancer, resp, err
}

func (s *LoadBalancersClient) Create(
	ctx context.Context,
	orgID string,
	args *LoadBalancerArguments,
) (*LoadBalancer, *Response, error) {
	if args == nil {
		return nil, nil, errors.New("nil load balancer arguments")
	}

	u := &url.URL{Path: fmt.Sprintf("organizations/%s/load_balancers", orgID)}
	reqBody := &struct {
		Properties *LoadBalancerArguments `json:"properties"`
	}{
		Properties: args,
	}

	body, resp, err := s.doRequest(ctx, "POST", u, reqBody)

	return body.LoadBalancer, resp, err
}

func (s *LoadBalancersClient) Update(
	ctx context.Context,
	lb *LoadBalancer,
) (*LoadBalancer, *Response, error) {
	if lb == nil {
		return nil, nil, errors.New("nil load balancer arguments")
	}

	if lb.ID == "" {
		return nil, nil, errors.New("ID value is empty")
	}

	u := &url.URL{Path: fmt.Sprintf("load_balancers/%s", lb.ID)}
	reqBody := &struct {
		Properties *LoadBalancerArguments `json:"properties"`
	}{
		Properties: &LoadBalancerArguments{
			Name:         lb.Name,
			ResourceType: lb.ResourceType,
			ResourceIDs:  lb.ResourceIDs,
			DataCenter:   lb.DataCenter,
		},
	}

	body, resp, err := s.doRequest(ctx, "PATCH", u, reqBody)

	return body.LoadBalancer, resp, err
}

func (s *LoadBalancersClient) Delete(
	ctx context.Context,
	id string,
) (*LoadBalancer, *Response, error) {
	u := &url.URL{Path: fmt.Sprintf("load_balancers/%s", id)}
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
