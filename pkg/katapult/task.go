package katapult

import (
	"context"
	"fmt"
	"net/url"

	"github.com/augurysys/timestamp"
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
	client   *apiClient
	basePath *url.URL
}

func newTasksClient(c *apiClient) *TasksClient {
	return &TasksClient{
		client:   c,
		basePath: &url.URL{Path: "/core/v1/"},
	}
}

func (s *TasksClient) Get(
	ctx context.Context,
	id string,
) (*Task, *Response, error) {
	u := &url.URL{
		Path: fmt.Sprintf("tasks/%s", id),
	}

	body, resp, err := s.doRequest(ctx, "GET", u, nil)

	return body.Task, resp, err
}

func (s *TasksClient) doRequest(
	ctx context.Context,
	method string,
	u *url.URL,
	body interface{},
) (*tasksResponseBody, *Response, error) {
	u = s.basePath.ResolveReference(u)
	respBody := &tasksResponseBody{}
	resp := &Response{}

	req, err := s.client.NewRequestWithContext(ctx, method, u, body)
	if err == nil {
		resp, err = s.client.Do(req, respBody)
	}

	return respBody, resp, err
}
