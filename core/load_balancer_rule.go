package core

import (
	"context"
	"fmt"
	"github.com/krystal/go-katapult"
	"net/url"
)

type LoadBalancerRule struct {
	ID              string     `json:"id,omitempty"`
	Algorithm       string     `json:"algorithm,omitempty"` // TODO: replace with constrained type?
	DestinationPort int        `json:"destination_port,omitempty"`
	ListenPort      int        `json:"listen_port,omitempty"`
	Protocol        string     `json:"protocol,omitempty"`     // TODO: replace with type?
	Certificates    []struct{} `json:"certificates,omitempty"` // TODO: is this the same certificate type as certificate.go
	BackendSSL      bool       `json:"backend_ssl,omitempty"`
	PassthroughSSL  bool       `json:"passthrough_ssl,omitempty"`
	CheckEnabled    bool       `json:"check_enabled,omitempty"`
	CheckFall       int        `json:"check_fall,omitempty"`
	CheckInterval   int        `json:"check_interval,omitempty"`
	CheckPath       string     `json:"check_path,omitempty"`
	CheckProtocol   string     `json:"check_protocol,omitempty"` // TODO: replace with type?
	CheckRise       int        `json:"check_rise,omitempty"`
	CheckTimeout    int        `json:"check_timeout,omitempty"`
}

type LoadBalancerRuleCreateArguments struct {
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
	loadBalancerID string,
	opts *ListOptions,
) ([]LoadBalancerRule, *katapult.Response, error) {
	qs := queryValues(opts)
	u := &url.URL{
		Path:     fmt.Sprintf("load_balancers/%s/rules", loadBalancerID),
		RawQuery: qs.Encode(),
	}

	body, resp, err := s.doRequest(ctx, "GET", u, nil)
	resp.Pagination = body.Pagination

	return body.LoadBalancerRules, resp, err
}

func (s *LoadBalancerRulesClient) doRequest(
	ctx context.Context,
	method string,
	u *url.URL,
	body interface{},
) (*loadBalancerRulesResponseBody, *katapult.Response, error) {
	u = s.basePath.ResolveReference(u)
	respBody := &loadBalancerRulesResponseBody{}
	resp := katapult.NewResponse(nil)

	req, err := s.client.NewRequestWithContext(ctx, method, u, body)
	if err == nil {
		resp, err = s.client.Do(req, respBody)
	}

	return respBody, resp, err
}
