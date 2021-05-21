package core

import (
	"context"
	"net/url"

	"github.com/krystal/go-katapult"
)

type NetworkSpeedProfile struct {
	ID                  string `json:"id,omitempty"`
	Name                string `json:"name,omitempty"`
	UploadSpeedInMbit   int    `json:"upload_speed_in_mbit,omitempty"`
	DownloadSpeedInMbit int    `json:"download_speed_in_mbit,omitempty"`
	Permalink           string `json:"permalink,omitempty"`
}

func (s *NetworkSpeedProfile) Ref() NetworkSpeedProfileRef {
	return NetworkSpeedProfileRef{ID: s.ID}
}

type NetworkSpeedProfileRef struct {
	ID        string `json:"id,omitempty"`
	Permalink string `json:"permalink,omitempty"`
}

type networkSpeedProfileResponseBody struct {
	Pagination           *katapult.Pagination   `json:"pagination,omitempty"`
	NetworkSpeedProfiles []*NetworkSpeedProfile `json:"network_speed_profiles,omitempty"`
}

type NetworkSpeedProfilesClient struct {
	client   RequestMaker
	basePath *url.URL
}

func NewNetworkSpeedProfilesClient(
	rm RequestMaker,
) *NetworkSpeedProfilesClient {
	return &NetworkSpeedProfilesClient{
		client:   rm,
		basePath: &url.URL{Path: "/core/v1/"},
	}
}

func (s *NetworkSpeedProfilesClient) List(
	ctx context.Context,
	org OrganizationRef,
	opts *ListOptions,
) ([]*NetworkSpeedProfile, *katapult.Response, error) {
	qs := queryValues(org, opts)
	u := &url.URL{
		Path:     "organizations/_/network_speed_profiles",
		RawQuery: qs.Encode(),
	}

	body, resp, err := s.doRequest(ctx, "GET", u, nil)
	resp.Pagination = body.Pagination

	return body.NetworkSpeedProfiles, resp, err
}

func (s *NetworkSpeedProfilesClient) doRequest(
	ctx context.Context,
	method string,
	u *url.URL,
	body interface{},
) (*networkSpeedProfileResponseBody, *katapult.Response, error) {
	u = s.basePath.ResolveReference(u)
	respBody := &networkSpeedProfileResponseBody{}
	resp := katapult.NewResponse(nil)

	req, err := s.client.NewRequestWithContext(ctx, method, u, body)
	if err == nil {
		resp, err = s.client.Do(req, respBody)
	}

	return respBody, resp, err
}
