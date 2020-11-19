package katapult

import (
	"testing"
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
				ID:         "id1",
				Name:       "task name",
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
				ID:     "id1",
				Status: TaskPending,
			},
		},
		{
			name: "running",
			obj: &Task{
				ID:     "id1",
				Status: TaskRunning,
			},
		},
		{
			name: "completed",
			obj: &Task{
				ID:     "id1",
				Status: TaskCompleted,
			},
		},
		{
			name: "failed",
			obj: &Task{
				ID:     "id1",
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
