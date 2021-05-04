package core

import (
	"context"
	"net/url"
	"strings"

	"github.com/krystal/go-katapult"
)

type DiskTemplate struct {
	ID              string               `json:"id,omitempty"`
	Name            string               `json:"name,omitempty"`
	Description     string               `json:"description,omitempty"`
	Permalink       string               `json:"permalink,omitempty"`
	Universal       bool                 `json:"universal,omitempty"`
	LatestVersion   *DiskTemplateVersion `json:"latest_version,omitempty"`
	OperatingSystem *OperatingSystem     `json:"operating_system,omitempty"`
}

// NewDiskTemplateLookup takes a string that is a DiskTemplate ID or Permalink,
// returning a empty *DiskTemplate struct with either the ID or Permalink field
// populated with the given value. This struct is suitable as input to other
// methods which accept a *DiskTemplate as input.
func NewDiskTemplateLookup(
	idOrPermalink string,
) (lr *DiskTemplate, f FieldName) {
	if strings.HasPrefix(idOrPermalink, "dtpl_") {
		return &DiskTemplate{ID: idOrPermalink}, IDField
	}

	return &DiskTemplate{Permalink: idOrPermalink}, PermalinkField
}

func (s *DiskTemplate) lookupReference() *DiskTemplate {
	if s == nil {
		return nil
	}

	lr := &DiskTemplate{ID: s.ID}
	if lr.ID == "" {
		lr.Permalink = s.Permalink
	}

	return lr
}

func (s *DiskTemplate) queryValues() *url.Values {
	v := &url.Values{}

	if s != nil {
		switch {
		case s.ID != "":
			v.Set("disk_template[id]", s.ID)
		case s.Permalink != "":
			v.Set("disk_template[permalink]", s.Permalink)
		}
	}

	return v
}

type DiskTemplateVersion struct {
	ID       string `json:"id,omitempty"`
	Number   int    `json:"number,omitempty"`
	Stable   bool   `json:"stable,omitempty"`
	SizeInGB int    `json:"size_in_gb,omitempty"`
}

type DiskTemplateOption struct {
	Key   string `json:"key,omitempty"`
	Value string `json:"value,omitempty"`
}

type DiskTemplateListOptions struct {
	IncludeUniversal bool
	Page             int
	PerPage          int
}

func (s *DiskTemplateListOptions) queryValues() *url.Values {
	if s == nil {
		return &url.Values{}
	}

	opts := &ListOptions{
		Page:    s.Page,
		PerPage: s.PerPage,
	}

	values := opts.queryValues()
	if s.IncludeUniversal {
		values.Set("include_universal", "true")
	}

	return values
}

type diskTemplateResponseBody struct {
	Pagination    *katapult.Pagination `json:"pagination,omitempty"`
	DiskTemplate  *DiskTemplate        `json:"disk_template,omitempty"`
	DiskTemplates []*DiskTemplate      `json:"disk_templates,omitempty"`
}

type DiskTemplatesClient struct {
	client   RequestMaker
	basePath *url.URL
}

func NewDiskTemplatesClient(rm RequestMaker) *DiskTemplatesClient {
	return &DiskTemplatesClient{
		client:   rm,
		basePath: &url.URL{Path: "/core/v1/"},
	}
}

func (s *DiskTemplatesClient) List(
	ctx context.Context,
	org *Organization,
	opts *DiskTemplateListOptions,
) ([]*DiskTemplate, *katapult.Response, error) {
	qs := queryValues(org, opts)
	u := &url.URL{
		Path:     "organizations/_/disk_templates",
		RawQuery: qs.Encode(),
	}

	body, resp, err := s.doRequest(ctx, "GET", u, nil)
	resp.Pagination = body.Pagination

	return body.DiskTemplates, resp, err
}

func (s *DiskTemplatesClient) Get(
	ctx context.Context,
	idOrPermalink string,
) (*DiskTemplate, *katapult.Response, error) {
	if _, f := NewDiskTemplateLookup(idOrPermalink); f == IDField {
		return s.GetByID(ctx, idOrPermalink)
	}

	return s.GetByPermalink(ctx, idOrPermalink)
}

func (s *DiskTemplatesClient) GetByID(
	ctx context.Context,
	id string,
) (*DiskTemplate, *katapult.Response, error) {
	return s.get(ctx, &DiskTemplate{ID: id})
}

func (s *DiskTemplatesClient) GetByPermalink(
	ctx context.Context,
	permalink string,
) (*DiskTemplate, *katapult.Response, error) {
	return s.get(ctx, &DiskTemplate{Permalink: permalink})
}

func (s *DiskTemplatesClient) get(
	ctx context.Context,
	dt *DiskTemplate,
) (*DiskTemplate, *katapult.Response, error) {
	u := &url.URL{
		Path:     "disk_templates/_",
		RawQuery: dt.queryValues().Encode(),
	}

	body, resp, err := s.doRequest(ctx, "GET", u, nil)

	return body.DiskTemplate, resp, err
}

func (s *DiskTemplatesClient) doRequest(
	ctx context.Context,
	method string,
	u *url.URL,
	body interface{},
) (*diskTemplateResponseBody, *katapult.Response, error) {
	u = s.basePath.ResolveReference(u)
	respBody := &diskTemplateResponseBody{}
	resp := katapult.NewResponse(nil)

	req, err := s.client.NewRequestWithContext(ctx, method, u, body)
	if err == nil {
		resp, err = s.client.Do(req, respBody)
	}

	return respBody, resp, err
}
