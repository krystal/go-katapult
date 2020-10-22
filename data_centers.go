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
	u, err := s.path.Parse("data_centers")
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequestWithContext(ctx, "GET", u.Path, nil)
	if err != nil {
		return nil, nil, err
	}

	var body dataCentersResponseBody
	resp, err := s.client.Do(req, &body)
	if err != nil {
		return nil, resp, err
	}

	return body.DataCenters, resp, nil
}

func (s *DataCentersService) Get(
	ctx context.Context,
	id string,
) (*DataCenter, *Response, error) {
	u, err := s.path.Parse(fmt.Sprintf("data_centers/%s", id))
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequestWithContext(ctx, "GET", u.Path, nil)
	if err != nil {
		return nil, nil, err
	}

	var body dataCentersResponseBody
	resp, err := s.client.Do(req, &body)
	if err != nil {
		return nil, resp, err
	}

	return body.DataCenter, resp, nil
}
