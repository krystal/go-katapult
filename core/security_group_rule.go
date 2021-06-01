package core

import (
	"context"
	"fmt"
	"net/url"

	"github.com/krystal/go-katapult"
)

type SecurityGroupRule struct {
	ID        string    `json:"id,omitempty"`
	Direction string    `json:"direction,omitempty"`
	Protocol  string    `json:"protocol,omitempty"`
	Ports     *string   `json:"ports,omitempty"`
	Targets   *[]string `json:"targets,omitempty"`
	Notes     *string   `json:"notes,omitempty"`
}

func (sgr *SecurityGroupRule) Ref() SecurityGroupRuleRef {
	return SecurityGroupRuleRef{ID: sgr.ID}
}

type SecurityGroupRuleRef struct {
	ID string `json:"id,omitempty"`
}

func (sgr SecurityGroupRuleRef) queryValues() *url.Values {
	return &url.Values{"security_group_rule[id]": []string{sgr.ID}}
}

type SecurityGroupRuleArguments struct {
	Direction string    `json:"direction,omitempty"`
	Protocol  string    `json:"protocol,omitempty"`
	Ports     *string   `json:"ports,omitempty"`
	Targets   *[]string `json:"targets,omitempty"`
	Notes     *string   `json:"notes,omitempty"`
}

type SecurityGroupRulesResponseBody struct {
	Pagination         *katapult.Pagination `json:"pagination,omitempty"`
	SecurityGroupRule  *SecurityGroupRule   `json:"security_group_rule,omitempty"`
	SecurityGroupRules []SecurityGroupRule  `json:"security_group_rules,omitempty"`
}

type SecurityGroupRulesClient struct {
	client   RequestMaker
	basePath *url.URL
}

// NewSecurityGroupRulesClient returns a new SecurityGroupRulesClient for
// interacting with SecurityGroup Rules.
func NewSecurityGroupRulesClient(rm RequestMaker) *SecurityGroupRulesClient {
	return &SecurityGroupRulesClient{
		client:   rm,
		basePath: &url.URL{Path: "/core/v1/"},
	}
}

// List returns SecurityGroup Rules for the specified SecurityGroup.
func (s *SecurityGroupRulesClient) List(
	ctx context.Context,
	sg SecurityGroupRef,
	opts *ListOptions,
) ([]SecurityGroupRule, *katapult.Response, error) {
	qs := queryValues(opts, sg)
	u := &url.URL{
		Path:     "security_groups/_/rules",
		RawQuery: qs.Encode(),
	}

	body, resp, err := s.doRequest(ctx, "GET", u, nil)
	if err != nil {
		return nil, resp, err
	}

	resp.Pagination = body.Pagination

	return body.SecurityGroupRules, resp, err
}

func (s *SecurityGroupRulesClient) Get(
	ctx context.Context,
	ref SecurityGroupRuleRef,
) (*SecurityGroupRule, *katapult.Response, error) {
	u := &url.URL{
		Path:     "security_groups/rules/_",
		RawQuery: ref.queryValues().Encode(),
	}
	body, resp, err := s.doRequest(ctx, "GET", u, nil)
	if err != nil {
		return nil, resp, err
	}

	return body.SecurityGroupRule, resp, err
}

func (s *SecurityGroupRulesClient) GetByID(
	ctx context.Context,
	id string,
) (*SecurityGroupRule, *katapult.Response, error) {
	return s.Get(ctx, SecurityGroupRuleRef{ID: id})
}

type securityGroupRuleCreateRequest struct {
	Properties *SecurityGroupRuleArguments `json:"properties,omitempty"`
}

func (s *SecurityGroupRulesClient) Create(
	ctx context.Context,
	sg SecurityGroupRef,
	args *SecurityGroupRuleArguments,
) (*SecurityGroupRule, *katapult.Response, error) {
	u := &url.URL{Path: fmt.Sprintf("security_groups/%s/rules", sg.ID)}
	reqBody := &securityGroupRuleCreateRequest{
		Properties: args,
	}

	body, resp, err := s.doRequest(ctx, "POST", u, reqBody)
	if err != nil {
		return nil, resp, err
	}

	return body.SecurityGroupRule, resp, nil
}

type securityGroupRuleUpdateRequest struct {
	Properties *SecurityGroupRuleArguments `json:"properties,omitempty"`
}

func (s *SecurityGroupRulesClient) Update(
	ctx context.Context,
	ref SecurityGroupRuleRef,
	args *SecurityGroupRuleArguments,
) (*SecurityGroupRule, *katapult.Response, error) {
	u := &url.URL{
		Path:     "security_groups/rules/_",
		RawQuery: ref.queryValues().Encode(),
	}
	reqBody := &securityGroupRuleUpdateRequest{
		Properties: args,
	}

	body, resp, err := s.doRequest(ctx, "PATCH", u, reqBody)
	if err != nil {
		return nil, resp, err
	}

	return body.SecurityGroupRule, resp, nil
}

func (s *SecurityGroupRulesClient) Delete(
	ctx context.Context,
	ref SecurityGroupRuleRef,
) (*SecurityGroupRule, *katapult.Response, error) {
	u := &url.URL{
		Path:     "security_groups/rules/_",
		RawQuery: ref.queryValues().Encode(),
	}
	body, resp, err := s.doRequest(ctx, "DELETE", u, nil)
	if err != nil {
		return nil, resp, err
	}

	return body.SecurityGroupRule, resp, nil
}

func (s *SecurityGroupRulesClient) doRequest(
	ctx context.Context,
	method string,
	u *url.URL,
	body interface{},
) (*SecurityGroupRulesResponseBody, *katapult.Response, error) {
	u = s.basePath.ResolveReference(u)

	req, err := s.client.NewRequestWithContext(ctx, method, u, body)
	if err != nil {
		return nil, nil, err
	}

	respBody := &SecurityGroupRulesResponseBody{}
	resp, err := s.client.Do(req, respBody)
	if err != nil {
		return nil, resp, err
	}

	return respBody, resp, nil
}
