package katapult

import (
	"context"
	"net/url"
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

// LookupReference returns a new *DiskTemplate stripped down to just ID or
// Permalink fields, making it suitable for endpoints which require a reference
// to a Disk Template by ID or Permalink.
func (s *DiskTemplate) LookupReference() *DiskTemplate {
	if s == nil {
		return nil
	}

	lr := &DiskTemplate{ID: s.ID}
	if lr.ID == "" {
		lr.Permalink = s.Permalink
	}

	return lr
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
	Pagination    *Pagination     `json:"pagination,omitempty"`
	DiskTemplates []*DiskTemplate `json:"disk_templates,omitempty"`
}

type DiskTemplatesClient struct {
	client   *apiClient
	basePath *url.URL
}

func newDiskTemplatesClient(c *apiClient) *DiskTemplatesClient {
	return &DiskTemplatesClient{
		client:   c,
		basePath: &url.URL{Path: "/core/v1/"},
	}
}

func (s *DiskTemplatesClient) List(
	ctx context.Context,
	org *Organization,
	opts *DiskTemplateListOptions,
) ([]*DiskTemplate, *Response, error) {
	qs := queryValues(org, opts)
	u := &url.URL{
		Path:     "organizations/_/disk_templates",
		RawQuery: qs.Encode(),
	}

	body, resp, err := s.doRequest(ctx, "GET", u, nil)
	resp.Pagination = body.Pagination

	return body.DiskTemplates, resp, err
}

func (s *DiskTemplatesClient) doRequest(
	ctx context.Context,
	method string,
	u *url.URL,
	body interface{},
) (*diskTemplateResponseBody, *Response, error) {
	u = s.basePath.ResolveReference(u)
	respBody := &diskTemplateResponseBody{}
	resp := &Response{}

	req, err := s.client.NewRequestWithContext(ctx, method, u, body)
	if err == nil {
		resp, err = s.client.Do(req, respBody)
	}

	return respBody, resp, err
}
