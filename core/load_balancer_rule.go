package core

import (
	"context"
	"fmt"
	"net/url"

	"github.com/krystal/go-katapult"
)

type LoadBalancerRuleAlgorithm string

//nolint:lll
const (
	RoundRobinRuleAlgorithm       LoadBalancerRuleAlgorithm = "round_robin"
	LeastConnectionsRuleAlgorithm LoadBalancerRuleAlgorithm = "least_connections"
	StickyRuleAlgorithm           LoadBalancerRuleAlgorithm = "sticky"
)

type Protocol string

const (
	HTTPSProtocol Protocol = "HTTPS"
	TCPProtocol   Protocol = "TCP"
	HTTPProtocol  Protocol = "HTTP"
)

type LoadBalancerRule struct {
	ID              string                    `json:"id,omitempty"`
	Algorithm       LoadBalancerRuleAlgorithm `json:"algorithm,omitempty"`
	DestinationPort int                       `json:"destination_port,omitempty"`
	ListenPort      int                       `json:"listen_port,omitempty"`
	Protocol        Protocol                  `json:"protocol,omitempty"`
	ProxyProtocol   bool                      `json:"proxy_protocol,omitempty"`
	Certificates    []Certificate             `json:"certificates,omitempty"`
	BackendSSL      bool                      `json:"backend_ssl,omitempty"`
	PassthroughSSL  bool                      `json:"passthrough_ssl,omitempty"`
	CheckEnabled    bool                      `json:"check_enabled,omitempty"`
	CheckFall       int                       `json:"check_fall,omitempty"`
	CheckInterval   int                       `json:"check_interval,omitempty"`
	CheckPath       string                    `json:"check_path,omitempty"`
	CheckProtocol   Protocol                  `json:"check_protocol,omitempty"`
	CheckRise       int                       `json:"check_rise,omitempty"`
	CheckTimeout    int                       `json:"check_timeout,omitempty"`
}

func (lbr *LoadBalancerRule) Ref() LoadBalancerRuleRef {
	return LoadBalancerRuleRef{ID: lbr.ID}
}

type LoadBalancerRuleRef struct {
	ID string `json:"id,omitempty"`
}

func (lbr LoadBalancerRuleRef) queryValues() *url.Values {
	v := &url.Values{}
	v.Set("load_balancer_rule[id]", lbr.ID)

	return v
}

type LoadBalancerRuleArguments struct {
	Algorithm       LoadBalancerRuleAlgorithm `json:"algorithm,omitempty"`
	DestinationPort int                       `json:"destination_port,omitempty"`
	ListenPort      int                       `json:"listen_port,omitempty"`
	Protocol        Protocol                  `json:"protocol,omitempty"`
	ProxyProtocol   *bool                     `json:"proxy_protocol,omitempty"`
	Certificates    *[]CertificateRef         `json:"certificates,omitempty"`
	CheckEnabled    *bool                     `json:"check_enabled,omitempty"`
	CheckFall       int                       `json:"check_fall,omitempty"`
	CheckInterval   int                       `json:"check_interval,omitempty"`
	CheckPath       string                    `json:"check_path,omitempty"`
	CheckProtocol   Protocol                  `json:"check_protocol,omitempty"`
	CheckRise       int                       `json:"check_rise,omitempty"`
	CheckTimeout    int                       `json:"check_timeout,omitempty"`
}

type loadBalancerRulesResponseBody struct {
	Pagination        *katapult.Pagination `json:"pagination,omitempty"`
	LoadBalancerRule  *LoadBalancerRule    `json:"load_balancer_rule,omitempty"`
	LoadBalancerRules []LoadBalancerRule   `json:"load_balancer_rules,omitempty"`
}

type LoadBalancerRulesClient struct {
	client   RequestMaker
	basePath *url.URL
}

// NewLoadBalancerRulesClient returns a new LoadBalancerRulesClient for
// interacting with LoadBalancer Rules.
func NewLoadBalancerRulesClient(rm RequestMaker) *LoadBalancerRulesClient {
	return &LoadBalancerRulesClient{
		client:   rm,
		basePath: &url.URL{Path: "/core/v1/"},
	}
}

// List returns LoadBalancer Rules for the specified LoadBalancer.
func (s *LoadBalancerRulesClient) List(
	ctx context.Context,
	lb LoadBalancerRef,
	opts *ListOptions,
) ([]LoadBalancerRule, *katapult.Response, error) {
	qs := queryValues(opts, lb)
	u := &url.URL{
		Path:     "load_balancers/_/rules",
		RawQuery: qs.Encode(),
	}

	body, resp, err := s.doRequest(ctx, "GET", u, nil)
	if err != nil {
		return nil, resp, err
	}

	resp.Pagination = body.Pagination

	return body.LoadBalancerRules, resp, err
}

func (s *LoadBalancerRulesClient) Get(
	ctx context.Context,
	ref LoadBalancerRuleRef,
) (*LoadBalancerRule, *katapult.Response, error) {
	u := &url.URL{
		Path:     "load_balancers/rules/_",
		RawQuery: ref.queryValues().Encode(),
	}
	body, resp, err := s.doRequest(ctx, "GET", u, nil)
	if err != nil {
		return nil, resp, err
	}

	return body.LoadBalancerRule, resp, err
}

func (s *LoadBalancerRulesClient) GetByID(
	ctx context.Context,
	id string,
) (*LoadBalancerRule, *katapult.Response, error) {
	return s.Get(ctx, LoadBalancerRuleRef{ID: id})
}

type loadBalancerRuleCreateRequest struct {
	Properties *LoadBalancerRuleArguments `json:"properties,omitempty"`
}

func (s *LoadBalancerRulesClient) Create(
	ctx context.Context,
	lb LoadBalancerRef,
	args *LoadBalancerRuleArguments,
) (*LoadBalancerRule, *katapult.Response, error) {
	u := &url.URL{Path: fmt.Sprintf("load_balancers/%s/rules", lb.ID)}
	reqBody := &loadBalancerRuleCreateRequest{
		Properties: args,
	}

	body, resp, err := s.doRequest(ctx, "POST", u, reqBody)
	if err != nil {
		return nil, resp, err
	}

	return body.LoadBalancerRule, resp, nil
}

type loadBalancerRuleUpdateRequest struct {
	Properties *LoadBalancerRuleArguments `json:"properties,omitempty"`
}

func (s *LoadBalancerRulesClient) Update(
	ctx context.Context,
	ref LoadBalancerRuleRef,
	args *LoadBalancerRuleArguments,
) (*LoadBalancerRule, *katapult.Response, error) {
	u := &url.URL{
		Path:     "load_balancers/rules/_",
		RawQuery: ref.queryValues().Encode(),
	}
	reqBody := &loadBalancerRuleUpdateRequest{
		Properties: args,
	}

	body, resp, err := s.doRequest(ctx, "PATCH", u, reqBody)
	if err != nil {
		return nil, resp, err
	}

	return body.LoadBalancerRule, resp, nil
}

func (s *LoadBalancerRulesClient) Delete(
	ctx context.Context,
	ref LoadBalancerRuleRef,
) (*LoadBalancerRule, *katapult.Response, error) {
	u := &url.URL{
		Path:     "load_balancers/rules/_",
		RawQuery: ref.queryValues().Encode(),
	}
	body, resp, err := s.doRequest(ctx, "DELETE", u, nil)
	if err != nil {
		return nil, resp, err
	}

	return body.LoadBalancerRule, resp, nil
}

func (s *LoadBalancerRulesClient) doRequest(
	ctx context.Context,
	method string,
	u *url.URL,
	body interface{},
) (*loadBalancerRulesResponseBody, *katapult.Response, error) {
	u = s.basePath.ResolveReference(u)
	respBody := &loadBalancerRulesResponseBody{}

	req := katapult.NewRequest(method, u, body)
	resp, err := s.client.Do(ctx, req, respBody)
	if resp == nil {
		resp = katapult.NewResponse(nil)
	}

	return respBody, resp, handleResponseError(err)
}
