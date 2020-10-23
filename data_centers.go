package katapult

import (
	"context"
	"fmt"
	"net/url"
)

type DataCentersService struct {
	*service
	path *url.URL
}

func NewDataCentersService(s *service) *DataCentersService {
	return &DataCentersService{
		service: s,
		path:    &url.URL{Path: "/core/v1/"},
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
	u := "data_centers"
	body, resp, err := s.doRequest(ctx, "GET", u, nil)

	return body.DataCenters, resp, err
}

func (s *DataCentersService) Get(
	ctx context.Context,
	id string,
) (*DataCenter, *Response, error) {
	u := fmt.Sprintf("data_centers/%s", id)
	body, resp, err := s.doRequest(ctx, "GET", u, nil)

	return body.DataCenter, resp, err
}

func (s *DataCentersService) doRequest(
	ctx context.Context,
	method string,
	urlStr string,
	body interface{},
) (*dataCentersResponseBody, *Response, error) {
	u, err := s.path.Parse(urlStr)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequestWithContext(ctx, method, u.String(), body)
	if err != nil {
		return nil, nil, err
	}

	var respBody dataCentersResponseBody
	resp, err := s.client.Do(req, &respBody)

	return &respBody, resp, err
}
