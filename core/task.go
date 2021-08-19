package core

import (
	"context"
	"fmt"
	"net/url"

	"github.com/augurysys/timestamp"
	"github.com/krystal/go-katapult"
)

type TaskStatus string

const (
	TaskPending   TaskStatus = "pending"
	TaskRunning   TaskStatus = "running"
	TaskCompleted TaskStatus = "completed"
	TaskFailed    TaskStatus = "failed"
)

type Task struct {
	ID         string               `json:"id,omitempty"`
	Name       string               `json:"name,omitempty"`
	Status     TaskStatus           `json:"status,omitempty"`
	CreatedAt  *timestamp.Timestamp `json:"created_at,omitempty"`
	StartedAt  *timestamp.Timestamp `json:"started_at,omitempty"`
	FinishedAt *timestamp.Timestamp `json:"finished_at,omitempty"`
	Progress   int                  `json:"progress,omitempty"`
}

type tasksResponseBody struct {
	Task *Task `json:"task,omitempty"`
}

type TasksClient struct {
	client   RequestMaker
	basePath *url.URL
}

func NewTasksClient(rm RequestMaker) *TasksClient {
	return &TasksClient{
		client:   rm,
		basePath: &url.URL{Path: "/core/v1/"},
	}
}

func (s *TasksClient) Get(
	ctx context.Context,
	id string,
	reqOpts ...katapult.RequestOption,
) (*Task, *katapult.Response, error) {
	u := &url.URL{
		Path: fmt.Sprintf("tasks/%s", id),
	}

	body, resp, err := s.doRequest(ctx, "GET", u, nil, reqOpts...)

	return body.Task, resp, err
}

func (s *TasksClient) doRequest(
	ctx context.Context,
	method string,
	u *url.URL,
	body interface{},
	reqOpts ...katapult.RequestOption,
) (*tasksResponseBody, *katapult.Response, error) {
	u = s.basePath.ResolveReference(u)
	respBody := &tasksResponseBody{}

	req := katapult.NewRequest(method, u, body, reqOpts...)
	resp, err := s.client.Do(ctx, req, respBody)
	if resp == nil {
		resp = katapult.NewResponse(nil)
	}

	return respBody, resp, handleResponseError(err)
}
