package core

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/krystal/go-katapult"

	"github.com/stretchr/testify/assert"
)

var (
	fixtureTaskQueueingErrorErr = "task_queueing_error: This error means " +
		"that a background task that was needed to complete your request " +
		"could not be queued"
	fixtureTaskQueueingErrorResponseError = &katapult.ResponseError{
		Code: "task_queueing_error",
		Description: "This error means that a background task that was " +
			"needed to complete your request could not be queued",
		Detail: json.RawMessage(`{}`),
	}
	fixtureTaskNotFoundErr = "task_not_found: No task was found matching any " +
		"of the criteria provided in the arguments"
	fixtureTaskNotFoundResponseError = &katapult.ResponseError{
		Code: "task_not_found",
		Description: "No task was found matching any of the criteria " +
			"provided in the arguments",
		Detail: json.RawMessage(`{}`),
	}
)

func TestTask_JSONMarshaling(t *testing.T) {
	tests := []struct {
		name string
		obj  *Task
	}{
		{
			name: "empty",
			obj:  &Task{},
		},
		{
			name: "full",
			obj: &Task{
				ID:         "task_wZgsjyVjrYEtw0Wl",
				Name:       "Purge items from trash",
				Status:     TaskPending,
				CreatedAt:  timestampPtr(1599412748),
				StartedAt:  timestampPtr(1591636763),
				FinishedAt: timestampPtr(1598203165),
				Progress:   42,
			},
		},
		{
			name: "pending",
			obj: &Task{
				ID:     "task_wZgsjyVjrYEtw0Wl",
				Status: TaskPending,
			},
		},
		{
			name: "running",
			obj: &Task{
				ID:     "task_wZgsjyVjrYEtw0Wl",
				Status: TaskRunning,
			},
		},
		{
			name: "completed",
			obj: &Task{
				ID:     "task_wZgsjyVjrYEtw0Wl",
				Status: TaskCompleted,
			},
		},
		{
			name: "failed",
			obj: &Task{
				ID:     "task_wZgsjyVjrYEtw0Wl",
				Status: TaskFailed,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testJSONMarshaling(t, tt.obj)
		})
	}
}

func TestTaskStatuses(t *testing.T) {
	tests := []struct {
		name  string
		enum  TaskStatus
		value string
	}{
		{name: "TaskPending", enum: TaskPending, value: "pending"},
		{name: "TaskRunning", enum: TaskRunning, value: "running"},
		{name: "TaskCompleted", enum: TaskCompleted, value: "completed"},
		{name: "TaskFailed", enum: TaskFailed, value: "failed"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.value, string(tt.enum))
		})
	}
}

func Test_tasksResponseBody_JSONMarshaling(t *testing.T) {
	tests := []struct {
		name string
		obj  *tasksResponseBody
	}{
		{
			name: "empty",
			obj:  &tasksResponseBody{},
		},
		{
			name: "full",
			obj: &tasksResponseBody{
				Task: &Task{ID: "id3"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testJSONMarshaling(t, tt.obj)
		})
	}
}

func TestTasksClient_Get(t *testing.T) {
	type args struct {
		ctx context.Context
		id  string
	}
	tests := []struct {
		name       string
		args       args
		want       *Task
		errStr     string
		errResp    *katapult.ResponseError
		respStatus int
		respBody   []byte
	}{
		{
			name: "by ID",
			args: args{
				ctx: context.Background(),
				id:  "task_EvA0bkUBGZh6ATca",
			},
			want: &Task{
				ID:     "task_EvA0bkUBGZh6ATca",
				Name:   "Purge items from trash",
				Status: "completed",
			},
			respStatus: http.StatusOK,
			respBody:   fixture("task_get"),
		},
		{
			name: "non-existent task",
			args: args{
				ctx: context.Background(),
				id:  "task_nopethisbegone",
			},
			errStr:     fixtureTaskNotFoundErr,
			errResp:    fixtureTaskNotFoundResponseError,
			respStatus: http.StatusNotFound,
			respBody:   fixture("task_not_found_error"),
		},
		{
			name: "empty id",
			args: args{
				ctx: context.Background(),
				id:  "",
			},
			errStr:     fixtureTaskNotFoundErr,
			errResp:    fixtureTaskNotFoundResponseError,
			respStatus: http.StatusNotFound,
			respBody:   fixture("task_not_found_error"),
		},
		{
			name: "nil context",
			args: args{
				ctx: nil,
				id:  "task_EvA0bkUBGZh6ATca",
			},
			errStr: "net/http: nil Context",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rm, mux, _, teardown := katapult.MockClient(t)
			defer teardown()
			c := NewTasksClient(rm)

			mux.HandleFunc(
				fmt.Sprintf("/core/v1/tasks/%s", tt.args.id),
				func(w http.ResponseWriter, r *http.Request) {
					assert.Equal(t, "GET", r.Method)
					assertEmptyFieldSpec(t, r)
					assertAuthorization(t, r)

					w.WriteHeader(tt.respStatus)
					_, _ = w.Write(tt.respBody)
				},
			)

			got, resp, err := c.Get(
				tt.args.ctx, tt.args.id,
			)

			if tt.respStatus != 0 {
				assert.Equal(t, tt.respStatus, resp.StatusCode)
			}

			if tt.errStr == "" {
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, err, tt.errStr)
			}

			if tt.want != nil {
				assert.Equal(t, tt.want, got)
			}

			if tt.errResp != nil {
				assert.Equal(t, tt.errResp, resp.Error)
			}
		})
	}
}
