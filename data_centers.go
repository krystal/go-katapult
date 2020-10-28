package katapult

import (
	"context"
	"fmt"
	"net/url"
)

type DataCentersService struct {
	client   *apiClient
	basePath *url.URL
}

func newDataCentersService(c *apiClient) *DataCentersService {
	return &DataCentersService{
		client:   c,
		basePath: &url.URL{Path: "/core/v1/"},
	}
}

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

func (s *DataCentersService) List(
	ctx context.Context,
) ([]*DataCenter, *Response, error) {
	u := &url.URL{Path: "data_centers"}
	body, resp, err := s.doRequest(ctx, "GET", u, nil)

	return body.DataCenters, resp, err
}

func (s *DataCentersService) Get(
	ctx context.Context,
	id string,
) (*DataCenter, *Response, error) {
	u := &url.URL{Path: fmt.Sprintf("data_centers/%s", id)}
	body, resp, err := s.doRequest(ctx, "GET", u, nil)

	return body.DataCenter, resp, err
}

func (s *DataCentersService) doRequest(
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
