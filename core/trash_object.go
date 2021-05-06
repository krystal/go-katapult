package core

import (
	"context"
	"net/url"
	"strings"

	"github.com/augurysys/timestamp"
	"github.com/krystal/go-katapult"
)

type TrashObject struct {
	ID         string               `json:"id,omitempty"`
	KeepUntil  *timestamp.Timestamp `json:"keep_until,omitempty"`
	ObjectID   string               `json:"object_id,omitempty"`
	ObjectType string               `json:"object_type,omitempty"`
}

// NewTrashObjectLookup takes a string that is a TrashObject ID or ObjectID
// returning, a empty *TrashObject struct with either the ID or ObjectID field
// populated with the given value. This struct is suitable as input to other
// methods which accept a *TrashObject as input.
func NewTrashObjectLookup(
	idOrObjectID string,
) (lr *TrashObject, f FieldName) {
	if strings.HasPrefix(idOrObjectID, "trsh_") {
		return &TrashObject{ID: idOrObjectID}, IDField
	}

	return &TrashObject{ObjectID: idOrObjectID}, ObjectIDField
}

func (s *TrashObject) lookupReference() *TrashObject {
	if s == nil {
		return nil
	}

	lr := &TrashObject{ID: s.ID}
	if lr.ID == "" {
		lr.ObjectID = s.ObjectID
	}

	return lr
}

func (s *TrashObject) queryValues() *url.Values {
	v := &url.Values{}

	if s != nil {
		switch {
		case s.ID != "":
			v.Set("trash_object[id]", s.ID)
		case s.ObjectID != "":
			v.Set("trash_object[object_id]", s.ObjectID)
		}
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
	org *Organization,
	opts *ListOptions,
) ([]*TrashObject, *katapult.Response, error) {
	qs := queryValues(org, opts)
	u := &url.URL{
		Path:     "organizations/_/trash_objects",
		RawQuery: qs.Encode(),
	}

	body, resp, err := s.doRequest(ctx, "GET", u, nil)
	resp.Pagination = body.Pagination

	return body.TrashObjects, resp, err
}

func (s *TrashObjectsClient) Get(
	ctx context.Context,
	idOrObjectID string,
) (*TrashObject, *katapult.Response, error) {
	if _, f := NewTrashObjectLookup(idOrObjectID); f == IDField {
		return s.GetByID(ctx, idOrObjectID)
	}

	return s.GetByObjectID(ctx, idOrObjectID)
}

func (s *TrashObjectsClient) GetByID(
	ctx context.Context,
	id string,
) (*TrashObject, *katapult.Response, error) {
	return s.get(ctx, &TrashObject{ID: id})
}

func (s *TrashObjectsClient) GetByObjectID(
	ctx context.Context,
	objectID string,
) (*TrashObject, *katapult.Response, error) {
	return s.get(ctx, &TrashObject{ObjectID: objectID})
}

func (s *TrashObjectsClient) get(
	ctx context.Context,
	trsh *TrashObject,
) (*TrashObject, *katapult.Response, error) {
	u := &url.URL{
		Path:     "trash_objects/_",
		RawQuery: trsh.queryValues().Encode(),
	}

	body, resp, err := s.doRequest(ctx, "GET", u, nil)

	return body.TrashObject, resp, err
}

func (s *TrashObjectsClient) Purge(
	ctx context.Context,
	trsh *TrashObject,
) (*Task, *katapult.Response, error) {
	u := &url.URL{
		Path:     "trash_objects/_",
		RawQuery: trsh.queryValues().Encode(),
	}

	body, resp, err := s.doRequest(ctx, "DELETE", u, nil)

	return body.Task, resp, err
}

func (s *TrashObjectsClient) PurgeAll(
	ctx context.Context,
	org *Organization,
) (*Task, *katapult.Response, error) {
	u := &url.URL{
		Path:     "organizations/_/trash_objects/purge_all",
		RawQuery: org.queryValues().Encode(),
	}

	body, resp, err := s.doRequest(ctx, "POST", u, nil)

	return body.Task, resp, err
}

func (s *TrashObjectsClient) Restore(
	ctx context.Context,
	trsh *TrashObject,
) (*TrashObject, *katapult.Response, error) {
	u := &url.URL{
		Path:     "trash_objects/_/restore",
		RawQuery: trsh.queryValues().Encode(),
	}

	body, resp, err := s.doRequest(ctx, "POST", u, nil)

	return body.TrashObject, resp, err
}

func (s *TrashObjectsClient) doRequest(
	ctx context.Context,
	method string,
	u *url.URL,
	body interface{},
) (*trashObjectsResponseBody, *katapult.Response, error) {
	u = s.basePath.ResolveReference(u)
	respBody := &trashObjectsResponseBody{}
	resp := katapult.NewResponse(nil)

	req, err := s.client.NewRequestWithContext(ctx, method, u, body)
	if err == nil {
		resp, err = s.client.Do(req, respBody)
	}

	return respBody, resp, err
}
