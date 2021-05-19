package core

import (
	"context"
	"github.com/krystal/go-katapult"
	"net/url"
)

type DataCenter struct {
	ID        string   `json:"id,omitempty"`
	Name      string   `json:"name,omitempty"`
	Permalink string   `json:"permalink,omitempty"`
	Country   *Country `json:"country,omitempty"`
}

func (dc *DataCenter) Ref() DataCenterRef {
	return DataCenterRef{ID: dc.ID}
}

// DataCenterRef refers to a single data center. Only one field should be set.
type DataCenterRef struct {
	ID        string `json:"id,omitempty"`
	Permalink string `json:"permalink,omitempty"`
}

func (dcr DataCenterRef) queryValues() *url.Values {
	v := &url.Values{}

	switch {
	case dcr.ID != "":
		v.Set("data_center[id]", dcr.ID)
	case dcr.Permalink != "":
		v.Set("data_center[permalink]", dcr.Permalink)
	}

	return v
}

type dataCentersResponseBody struct {
	DataCenter  *DataCenter   `json:"data_center,omitempty"`
	DataCenters []*DataCenter `json:"data_centers,omitempty"`
}

type DataCentersClient struct {
	client   RequestMaker
	basePath *url.URL
}

func NewDataCentersClient(rm RequestMaker) *DataCentersClient {
	return &DataCentersClient{
		client:   rm,
		basePath: &url.URL{Path: "/core/v1/"},
	}
}

func (s *DataCentersClient) List(
	ctx context.Context,
) ([]*DataCenter, *katapult.Response, error) {
	u := &url.URL{Path: "data_centers"}
	body, resp, err := s.doRequest(ctx, "GET", u, nil)

	return body.DataCenters, resp, err
}

func (s *DataCentersClient) Get(
	ctx context.Context,
	ref DataCenterRef,
) (*DataCenter, *katapult.Response, error) {
	u := &url.URL{Path: "data_centers/_", RawQuery: ref.queryValues().Encode()}

	body, resp, err := s.doRequest(ctx, "GET", u, nil)

	return body.DataCenter, resp, err
}

func (s *DataCentersClient) GetByID(
	ctx context.Context,
	id string,
) (*DataCenter, *katapult.Response, error) {
	return s.Get(ctx, DataCenterRef{ID: id})
}

func (s *DataCentersClient) GetByPermalink(
	ctx context.Context,
	permalink string,
) (*DataCenter, *katapult.Response, error) {
	return s.Get(ctx, DataCenterRef{Permalink: permalink})
}

func (s *DataCentersClient) doRequest(
	ctx context.Context,
	method string,
	u *url.URL,
	body interface{},
) (*dataCentersResponseBody, *katapult.Response, error) {
	u = s.basePath.ResolveReference(u)
	respBody := &dataCentersResponseBody{}
	resp := katapult.NewResponse(nil)

	req, err := s.client.NewRequestWithContext(ctx, method, u, body)
	if err == nil {
		resp, err = s.client.Do(req, respBody)
	}

	return respBody, resp, err
}
