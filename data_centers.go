package katapult

import (
	"context"
	"fmt"
)

type DataCentersService struct {
	*service
	*pathHelper
}

func NewDataCentersService(s *service) *DataCentersService {
	p, _ := newPathHelper("/core/v1/")

	return &DataCentersService{service: s, pathHelper: p}
}

type DataCenter struct {
	ID        string   `json:"id,omitempty"`
	Name      string   `json:"name,omitempty"`
	Permalink string   `json:"permalink,omitempty"`
	Country   *Country `json:"country,omitempty"`
}

type DataCentersResponseBody struct {
	DataCenter  *DataCenter   `json:"data_center,omitempty"`
	DataCenters []*DataCenter `json:"data_centers,omitempty"`
}

func (s *DataCentersService) List(
	ctx context.Context,
) ([]*DataCenter, *Response, error) {
	u, _ := s.RequestPath("data_centers")

	req, err := s.client.NewRequestWithContext(ctx, "GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	var body *DataCentersResponseBody
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
	u, err := s.RequestPath(fmt.Sprintf("data_centers/%s", id))
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequestWithContext(ctx, "GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	var body *DataCentersResponseBody
	resp, err := s.client.Do(req, &body)
	if err != nil {
		return nil, resp, err
	}

	return body.DataCenter, resp, nil
}
