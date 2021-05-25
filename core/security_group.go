package core

import (
	"context"
	"fmt"
	"net/url"

	"github.com/krystal/go-katapult"
)

type SecurityGroup struct {
	ID               string   `json:"id,omitempty"`
	Name             string   `json:"name,omitempty"`
	AllowAllInbound  bool     `json:"allow_all_inbound,omitempty"`
	AllowAllOutbound bool     `json:"allow_all_outbound,omitempty"`
	Associations     []string `json:"associations,omitempty"`
}

func (s *SecurityGroup) lookupReference() *SecurityGroup {
	if s == nil {
		return nil
	}

	return &SecurityGroup{ID: s.ID}
}

func (s *SecurityGroup) queryValues() *url.Values {
	v := &url.Values{}

	if s != nil && s.ID != "" {
		v.Set("security_group[id]", s.ID)
	}

	return v
}

type SecurityGroupCreateArguments struct {
	Name             string    `json:"name,omitempty"`
	Associations     *[]string `json:"associations,omitempty"`
	AllowAllInbound  *bool     `json:"allow_all_inbound,omitempty"`
	AllowAllOutbound *bool     `json:"allow_all_outbound,omitempty"`
}

func (
	s *SecurityGroupCreateArguments,
) forRequest() *SecurityGroupCreateArguments {
	if s == nil {
		return nil
	}

	args := *s

	return &args
}

type SecurityGroupUpdateArguments struct {
	Name             string    `json:"name,omitempty"`
	Associations     *[]string `json:"associations,omitempty"`
	AllowAllInbound  *bool     `json:"allow_all_inbound,omitempty"`
	AllowAllOutbound *bool     `json:"allow_all_outbound,omitempty"`
}

type SecurityGroupCreateRequest struct {
	Organization *Organization                 `json:"organization,omitempty"`
	Properties   *SecurityGroupCreateArguments `json:"properties,omitempty"`
}

type SecurityGroupUpdateRequest struct {
	SecurityGroup *SecurityGroup                `json:"security_group,omitempty"`
	Properties    *SecurityGroupUpdateArguments `json:"properties,omitempty"`
}

type SecurityGroupsResponseBody struct {
	Pagination     *katapult.Pagination `json:"pagination,omitempty"`
	SecurityGroup  *SecurityGroup       `json:"security_group,omitempty"`
	SecurityGroups []*SecurityGroup     `json:"security_groups,omitempty"`
}

type SecurityGroupsClient struct {
	client   RequestMaker
	basePath *url.URL
}

func NewSecurityGroupsClient(rm RequestMaker) *SecurityGroupsClient {
	return &SecurityGroupsClient{
		client:   rm,
		basePath: &url.URL{Path: "/core/v1/"},
	}
}

func (s *SecurityGroupsClient) List(
	ctx context.Context,
	org *Organization,
	opts *ListOptions,
) ([]*SecurityGroup, *katapult.Response, error) {
	qs := queryValues(org, opts)
	u := &url.URL{
		Path:     "organizations/_/security_groups",
		RawQuery: qs.Encode(),
	}

	body, resp, err := s.doRequest(ctx, "GET", u, nil)
	resp.Pagination = body.Pagination

	return body.SecurityGroups, resp, err
}

func (s *SecurityGroupsClient) Get(
	ctx context.Context,
	id string,
) (*SecurityGroup, *katapult.Response, error) {
	return s.GetByID(ctx, id)
}

func (s *SecurityGroupsClient) GetByID(
	ctx context.Context,
	id string,
) (*SecurityGroup, *katapult.Response, error) {
	u := &url.URL{Path: fmt.Sprintf("security_groups/%s", id)}
	body, resp, err := s.doRequest(ctx, "GET", u, nil)

	return body.SecurityGroup, resp, err
}

func (s *SecurityGroupsClient) Create(
	ctx context.Context,
	org *Organization,
	args *SecurityGroupCreateArguments,
) (*SecurityGroup, *katapult.Response, error) {
	u := &url.URL{Path: "organizations/_/security_groups"}
	reqBody := &SecurityGroupCreateRequest{
		Organization: org.lookupReference(),
		Properties:   args.forRequest(),
	}

	body, resp, err := s.doRequest(ctx, "POST", u, reqBody)

	return body.SecurityGroup, resp, err
}

func (s *SecurityGroupsClient) Update(
	ctx context.Context,
	lb *SecurityGroup,
	args *SecurityGroupUpdateArguments,
) (*SecurityGroup, *katapult.Response, error) {
	u := &url.URL{Path: "security_groups/_"}
	reqBody := &SecurityGroupUpdateRequest{
		SecurityGroup: lb.lookupReference(),
		Properties:    args,
	}

	body, resp, err := s.doRequest(ctx, "PATCH", u, reqBody)

	return body.SecurityGroup, resp, err
}

func (s *SecurityGroupsClient) Delete(
	ctx context.Context,
	lb *SecurityGroup,
) (*SecurityGroup, *katapult.Response, error) {
	u := &url.URL{
		Path:     "security_groups/_",
		RawQuery: lb.queryValues().Encode(),
	}
	body, resp, err := s.doRequest(ctx, "DELETE", u, nil)

	return body.SecurityGroup, resp, err
}

func (s *SecurityGroupsClient) doRequest(
	ctx context.Context,
	method string,
	u *url.URL,
	body interface{},
) (*SecurityGroupsResponseBody, *katapult.Response, error) {
	u = s.basePath.ResolveReference(u)
	respBody := &SecurityGroupsResponseBody{}
	resp := katapult.NewResponse(nil)

	req, err := s.client.NewRequestWithContext(ctx, method, u, body)
	if err == nil {
		resp, err = s.client.Do(req, respBody)
	}

	return respBody, resp, err
}
