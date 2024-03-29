package core

import (
	"context"
	"net/url"

	"github.com/augurysys/timestamp"
	"github.com/krystal/go-katapult"
)

type TrashObject struct {
	ID         string               `json:"id,omitempty"`
	KeepUntil  *timestamp.Timestamp `json:"keep_until,omitempty"`
	ObjectID   string               `json:"object_id,omitempty"`
	ObjectType string               `json:"object_type,omitempty"`
}

func (s *TrashObject) Ref() TrashObjectRef {
	return TrashObjectRef{ID: s.ID}
}

type TrashObjectRef struct {
	ID       string `json:"id,omitempty"`
	ObjectID string `json:"object_id,omitempty"`
}

func (s TrashObjectRef) queryValues() *url.Values {
	v := &url.Values{}
	switch {
	case s.ID != "":
		v.Set("trash_object[id]", s.ID)
	case s.ObjectID != "":
		v.Set("trash_object[object_id]", s.ObjectID)
	}

	return v
}

type trashObjectsResponseBody struct {
	Pagination   *katapult.Pagination `json:"pagination,omitempty"`
	TrashObject  *TrashObject         `json:"trash_object,omitempty"`
	TrashObjects []*TrashObject       `json:"trash_objects,omitempty"`
	Task         *Task                `json:"task,omitempty"`
}

type TrashObjectsClient struct {
	client   RequestMaker
	basePath *url.URL
}

func NewTrashObjectsClient(rm RequestMaker) *TrashObjectsClient {
	return &TrashObjectsClient{
		client:   rm,
		basePath: &url.URL{Path: "/core/v1/"},
	}
}

func (s *TrashObjectsClient) List(
	ctx context.Context,
	org OrganizationRef,
	opts *ListOptions,
	reqOpts ...katapult.RequestOption,
) ([]*TrashObject, *katapult.Response, error) {
	qs := queryValues(org, opts)
	u := &url.URL{
		Path:     "organizations/_/trash_objects",
		RawQuery: qs.Encode(),
	}

	body, resp, err := s.doRequest(ctx, "GET", u, nil, reqOpts...)
	resp.Pagination = body.Pagination

	return body.TrashObjects, resp, err
}

func (s *TrashObjectsClient) Get(
	ctx context.Context,
	ref TrashObjectRef,
	reqOpts ...katapult.RequestOption,
) (*TrashObject, *katapult.Response, error) {
	u := &url.URL{
		Path:     "trash_objects/_",
		RawQuery: ref.queryValues().Encode(),
	}
	body, resp, err := s.doRequest(ctx, "GET", u, nil, reqOpts...)

	return body.TrashObject, resp, err
}

func (s *TrashObjectsClient) GetByID(
	ctx context.Context,
	id string,
	reqOpts ...katapult.RequestOption,
) (*TrashObject, *katapult.Response, error) {
	return s.Get(ctx, TrashObjectRef{ID: id}, reqOpts...)
}

func (s *TrashObjectsClient) GetByObjectID(
	ctx context.Context,
	objectID string,
	reqOpts ...katapult.RequestOption,
) (*TrashObject, *katapult.Response, error) {
	return s.Get(ctx, TrashObjectRef{ObjectID: objectID}, reqOpts...)
}

func (s *TrashObjectsClient) Purge(
	ctx context.Context,
	ref TrashObjectRef,
	reqOpts ...katapult.RequestOption,
) (*Task, *katapult.Response, error) {
	u := &url.URL{
		Path:     "trash_objects/_",
		RawQuery: ref.queryValues().Encode(),
	}

	body, resp, err := s.doRequest(ctx, "DELETE", u, nil, reqOpts...)

	return body.Task, resp, err
}

func (s *TrashObjectsClient) PurgeAll(
	ctx context.Context,
	org OrganizationRef,
	reqOpts ...katapult.RequestOption,
) (*Task, *katapult.Response, error) {
	u := &url.URL{
		Path:     "organizations/_/trash_objects/purge_all",
		RawQuery: org.queryValues().Encode(),
	}

	body, resp, err := s.doRequest(ctx, "POST", u, nil, reqOpts...)

	return body.Task, resp, err
}

func (s *TrashObjectsClient) Restore(
	ctx context.Context,
	ref TrashObjectRef,
	reqOpts ...katapult.RequestOption,
) (*TrashObject, *katapult.Response, error) {
	u := &url.URL{
		Path:     "trash_objects/_/restore",
		RawQuery: ref.queryValues().Encode(),
	}

	body, resp, err := s.doRequest(ctx, "POST", u, nil, reqOpts...)

	return body.TrashObject, resp, err
}

func (s *TrashObjectsClient) doRequest(
	ctx context.Context,
	method string,
	u *url.URL,
	body interface{}, //nolint:unparam
	reqOpts ...katapult.RequestOption,
) (*trashObjectsResponseBody, *katapult.Response, error) {
	u = s.basePath.ResolveReference(u)
	respBody := &trashObjectsResponseBody{}

	req := katapult.NewRequest(method, u, body, reqOpts...)
	resp, err := s.client.Do(ctx, req, respBody)
	if resp == nil {
		resp = katapult.NewResponse(nil)
	}

	return respBody, resp, handleResponseError(err)
}
