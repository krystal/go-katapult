package core

import (
	"context"
	"net/url"

	"github.com/augurysys/timestamp"
	"github.com/krystal/go-katapult"
)

type Tag struct {
	ID        string               `json:"id,omitempty"`
	Name      string               `json:"name,omitempty"`
	Color     string               `json:"color,omitempty"`
	CreatedAt *timestamp.Timestamp `json:"created_at,omitempty"`
}

type TagsClient struct {
	client   RequestMaker
	basePath *url.URL
}

func NewTagsClient(rm RequestMaker) *TagsClient {
	return &TagsClient{
		client:   rm,
		basePath: &url.URL{Path: "/core/v1/"},
	}
}

type tagsResponseBody struct {
	Pagination *katapult.Pagination `json:"pagination,omitempty"`
	Tags       []*Tag               `json:"tags,omitempty"`
	Tag        *Tag                 `json:"tag,omitempty"`
}

func (s *TagsClient) List(
	ctx context.Context,
	ref OrganizationRef,
	opts *ListOptions,
	reqOpts ...katapult.RequestOption,
) ([]*Tag, *katapult.Response, error) {
	qs := queryValues(opts, ref)
	u := &url.URL{Path: "organizations/_/tags", RawQuery: qs.Encode()}

	body, resp, err := s.doRequest(ctx, "GET", u, nil, reqOpts...)
	resp.Pagination = body.Pagination

	return body.Tags, resp, err
}

type TagRef struct {
	ID string `json:"id"`
}

func (tr TagRef) queryValues() *url.Values {
	return &url.Values{"tag[id]": []string{tr.ID}}
}

func (s *TagsClient) Get(
	ctx context.Context,
	ref TagRef,
	reqOpts ...katapult.RequestOption,
) (*Tag, *katapult.Response, error) {
	qs := ref.queryValues()
	u := &url.URL{Path: "tags/_", RawQuery: qs.Encode()}

	body, resp, err := s.doRequest(ctx, "GET", u, nil, reqOpts...)

	return body.Tag, resp, err
}

type TagArguments struct {
	Name  string `json:"name,omitempty"`
	Color string `json:"color,omitempty"`
}

func (s *TagsClient) Create(
	ctx context.Context,
	ref OrganizationRef,
	args TagArguments,
	reqOpts ...katapult.RequestOption,
) (*Tag, *katapult.Response, error) {
	qs := ref.queryValues()
	u := &url.URL{Path: "organizations/_/tags", RawQuery: qs.Encode()}

	body, resp, err := s.doRequest(ctx, "POST", u, args, reqOpts...)

	return body.Tag, resp, err
}

func (s *TagsClient) Update(
	ctx context.Context,
	ref TagRef,
	args TagArguments,
	reqOpts ...katapult.RequestOption,
) (*Tag, *katapult.Response, error) {
	qs := ref.queryValues()
	u := &url.URL{Path: "tags/_", RawQuery: qs.Encode()}

	body, resp, err := s.doRequest(ctx, "PATCH", u, args, reqOpts...)

	return body.Tag, resp, err
}

func (s *TagsClient) Delete(
	ctx context.Context,
	ref TagRef,
	reqOpts ...katapult.RequestOption,
) (*Tag, *katapult.Response, error) {
	qs := ref.queryValues()
	u := &url.URL{Path: "tags/_", RawQuery: qs.Encode()}

	body, resp, err := s.doRequest(
		ctx, "DELETE", u, nil, reqOpts...,
	)

	return body.Tag, resp, err
}

func (s *TagsClient) doRequest(
	ctx context.Context,
	method string,
	u *url.URL,
	body interface{},
	reqOpts ...katapult.RequestOption,
) (*tagsResponseBody, *katapult.Response, error) {
	u = s.basePath.ResolveReference(u)
	respBody := &tagsResponseBody{}

	req := katapult.NewRequest(method, u, body, reqOpts...)
	resp, err := s.client.Do(ctx, req, respBody)
	if resp == nil {
		resp = katapult.NewResponse(nil)
	}

	return respBody, resp, handleResponseError(err)
}
