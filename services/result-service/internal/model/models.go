package model

import (
	"time"
)

// TaskResult represents the result of a task execution
type TaskResult struct {
	TaskID     string            `json:"task_id"`
	Status     string            `json:"status"`
	Result     map[string]string `json:"result,omitempty"`
	CreatedAt  time.Time         `json:"created_at"`
	UpdatedAt  time.Time         `json:"updated_at"`
	FinishedAt time.Time         `json:"finished_at,omitempty"`
}

// SubTaskResult represents the result of a subtask execution
type SubTaskResult struct {
	SubTaskID  string            `json:"subtask_id"`
	TaskID     string            `json:"task_id"`
	WorkerID   string            `json:"worker_id"`
	Status     string            `json:"status"`
	Result     map[string]string `json:"result,omitempty"`
	CreatedAt  time.Time         `json:"created_at"`
	UpdatedAt  time.Time         `json:"updated_at"`
	FinishedAt time.Time         `json:"finished_at,omitempty"`
}
