package core

import (
	"context"
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

func (sg *SecurityGroup) Ref() SecurityGroupRef {
	return SecurityGroupRef{ID: sg.ID}
}

// allows a reference to a security group
type SecurityGroupRef struct {
	ID string `json:"id,omitempty"`
}

func (sg SecurityGroupRef) queryValues() *url.Values {
	return &url.Values{
		"security_group[id]": []string{sg.ID},
	}
}

type SecurityGroupCreateArguments struct {
	Name             string    `json:"name,omitempty"`
	Associations     *[]string `json:"associations,omitempty"`
	AllowAllInbound  *bool     `json:"allow_all_inbound,omitempty"`
	AllowAllOutbound *bool     `json:"allow_all_outbound,omitempty"`
}

type SecurityGroupUpdateArguments struct {
	Name             string    `json:"name,omitempty"`
	Associations     *[]string `json:"associations,omitempty"`
	AllowAllInbound  *bool     `json:"allow_all_inbound,omitempty"`
	AllowAllOutbound *bool     `json:"allow_all_outbound,omitempty"`
}

type securityGroupCreateRequest struct {
	Organization OrganizationRef               `json:"organization"`
	Properties   *SecurityGroupCreateArguments `json:"properties,omitempty"`
}

type securityGroupUpdateRequest struct {
	SecurityGroup SecurityGroupRef              `json:"security_group"`
	Properties    *SecurityGroupUpdateArguments `json:"properties,omitempty"`
}

type securityGroupsResponseBody struct {
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

func (sgc *SecurityGroupsClient) List(
	ctx context.Context,
	org OrganizationRef,
	opts *ListOptions,
) ([]*SecurityGroup, *katapult.Response, error) {
	qs := queryValues(org, opts)
	u := &url.URL{
		Path:     "organizations/_/security_groups",
		RawQuery: qs.Encode(),
	}

	body, resp, err := sgc.doRequest(ctx, "GET", u, nil)
	resp.Pagination = body.Pagination

	return body.SecurityGroups, resp, err
}

func (sgc *SecurityGroupsClient) Get(
	ctx context.Context,
	ref SecurityGroupRef,
) (*SecurityGroup, *katapult.Response, error) {
	u := &url.URL{
		Path:     "security_groups/_",
		RawQuery: ref.queryValues().Encode(),
	}
	body, resp, err := sgc.doRequest(ctx, "GET", u, nil)

	return body.SecurityGroup, resp, err
}

func (sgc *SecurityGroupsClient) GetByID(
	ctx context.Context,
	id string,
) (*SecurityGroup, *katapult.Response, error) {
	return sgc.Get(ctx, SecurityGroupRef{ID: id})
}

func (sgc *SecurityGroupsClient) Create(
	ctx context.Context,
	org OrganizationRef,
	args *SecurityGroupCreateArguments,
) (*SecurityGroup, *katapult.Response, error) {
	u := &url.URL{Path: "organizations/_/security_groups"}
	reqBody := &securityGroupCreateRequest{
		Organization: org,
		Properties:   args,
	}

	body, resp, err := sgc.doRequest(ctx, "POST", u, reqBody)

	return body.SecurityGroup, resp, err
}

func (sgc *SecurityGroupsClient) Update(
	ctx context.Context,
	sg SecurityGroupRef,
	args *SecurityGroupUpdateArguments,
) (*SecurityGroup, *katapult.Response, error) {
	u := &url.URL{Path: "security_groups/_"}
	reqBody := &securityGroupUpdateRequest{
		SecurityGroup: sg,
		Properties:    args,
	}

	body, resp, err := sgc.doRequest(ctx, "PATCH", u, reqBody)

	return body.SecurityGroup, resp, err
}

func (sgc *SecurityGroupsClient) Delete(
	ctx context.Context,
	sg SecurityGroupRef,
) (*SecurityGroup, *katapult.Response, error) {
	u := &url.URL{
		Path:     "security_groups/_",
		RawQuery: sg.queryValues().Encode(),
	}
	body, resp, err := sgc.doRequest(ctx, "DELETE", u, nil)

	return body.SecurityGroup, resp, err
}

func (sgc *SecurityGroupsClient) doRequest(
	ctx context.Context,
	method string,
	u *url.URL,
	body interface{},
) (*securityGroupsResponseBody, *katapult.Response, error) {
	u = sgc.basePath.ResolveReference(u)
	respBody := &securityGroupsResponseBody{}

	req := katapult.NewRequest(method, u, body)
	resp, err := sgc.client.Do(ctx, req, respBody)
	if resp == nil {
		resp = katapult.NewResponse(nil)
	}

	return respBody, resp, err
}
