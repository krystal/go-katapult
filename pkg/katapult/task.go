package katapult

import "github.com/augurysys/timestamp"

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
