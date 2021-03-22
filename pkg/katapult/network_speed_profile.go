package katapult

import (
	"context"
	"net/url"
)

type NetworkSpeedProfile struct {
	ID                  string `json:"id,omitempty"`
	Name                string `json:"name,omitempty"`
	UploadSpeedInMbit   int    `json:"upload_speed_in_mbit,omitempty"`
	DownloadSpeedInMbit int    `json:"download_speed_in_mbit,omitempty"`
	Permalink           string `json:"permalink,omitempty"`
}

func (s *NetworkSpeedProfile) lookupReference() *NetworkSpeedProfile {
	if s == nil {
		return nil
	}

	lr := &NetworkSpeedProfile{ID: s.ID}
	if lr.ID == "" {
		lr.Permalink = s.Permalink
	}

	return lr
}

type networkSpeedProfileResponseBody struct {
	Pagination           *Pagination            `json:"pagination,omitempty"`
	NetworkSpeedProfiles []*NetworkSpeedProfile `json:"network_speed_profiles,omitempty"`
}

type NetworkSpeedProfilesClient struct {
	client   *apiClient
	basePath *url.URL
}

func newNetworkSpeedProfilesClient(c *apiClient) *NetworkSpeedProfilesClient {
	return &NetworkSpeedProfilesClient{
		client:   c,
		basePath: &url.URL{Path: "/core/v1/"},
	}
}

func (s *NetworkSpeedProfilesClient) List(
	ctx context.Context,
	org *Organization,
	opts *ListOptions,
) ([]*NetworkSpeedProfile, *Response, error) {
	qs := queryValues(org, opts)
	u := &url.URL{
		Path:     "organizations/_/network_speed_profiles",
		RawQuery: qs.Encode(),
	}

	body, resp, err := s.doRequest(ctx, "GET", u, nil)

	return body.NetworkSpeedProfiles, resp, err
}

func (s *NetworkSpeedProfilesClient) doRequest(
	ctx context.Context,
	method string,
	u *url.URL,
	body interface{},
) (*networkSpeedProfileResponseBody, *Response, error) {
	u = s.basePath.ResolveReference(u)
	respBody := &networkSpeedProfileResponseBody{}
	resp := newResponse(nil)

	req, err := s.client.NewRequestWithContext(ctx, method, u, body)
	if err == nil {
		resp, err = s.client.Do(req, respBody)
	}

	return respBody, resp, err
}
