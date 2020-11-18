package katapult

import (
	"context"
	"fmt"
	"net/url"
)

type DataCenter struct {
	ID        string   `json:"id,omitempty"`
	Name      string   `json:"name,omitempty"`
	Permalink string   `json:"permalink,omitempty"`
	Country   *Country `json:"country,omitempty"`
}

type dataCentersResponseBody struct {
	DataCenter  *DataCenter   `json:"data_center,omitempty"`
	DataCenters []*DataCenter `json:"data_centers,omitempty"`
}

type DataCentersClient struct {
	client   *apiClient
	basePath *url.URL
}

func newDataCentersClient(c *apiClient) *DataCentersClient {
	return &DataCentersClient{
		client:   c,
		basePath: &url.URL{Path: "/core/v1/"},
	}
}

func (s *DataCentersClient) List(
	ctx context.Context,
) ([]*DataCenter, *Response, error) {
	u := &url.URL{Path: "data_centers"}
	body, resp, err := s.doRequest(ctx, "GET", u, nil)

	return body.DataCenters, resp, err
}

func (s *DataCentersClient) Get(
	ctx context.Context,
	id string,
) (*DataCenter, *Response, error) {
	u := &url.URL{Path: fmt.Sprintf("data_centers/%s", id)}
	body, resp, err := s.doRequest(ctx, "GET", u, nil)

	return body.DataCenter, resp, err
}

func (s *DataCentersClient) GetByPermalink(
	ctx context.Context,
	permalink string,
) (*DataCenter, *Response, error) {
	qs := url.Values{"data_center[permalink]": []string{permalink}}
	u := &url.URL{Path: "data_centers/_", RawQuery: qs.Encode()}

	body, resp, err := s.doRequest(ctx, "GET", u, nil)

	return body.DataCenter, resp, err
}

func (s *DataCentersClient) doRequest(
	ctx context.Context,
	method string,
	u *url.URL,
	body interface{},
) (*dataCentersResponseBody, *Response, error) {
	u = s.basePath.ResolveReference(u)
	respBody := &dataCentersResponseBody{}
	resp := &Response{}

	req, err := s.client.NewRequestWithContext(ctx, method, u, body)
	if err == nil {
		resp, err = s.client.Do(req, respBody)
	}

	return respBody, resp, err
}