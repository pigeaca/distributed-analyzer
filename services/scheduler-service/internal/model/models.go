package model

import (
	"github.com/pigeaca/DistributedMarketplace/libs/model"
)

// SchedulerTask represents a task from the scheduler's perspective
type SchedulerTask struct {
	TaskID    string          `json:"task_id"`
	Status    model.Status    `json:"status"`
	SubTasks  []model.SubTask `json:"subtasks,omitempty"`
	WorkerIDs []string        `json:"worker_ids,omitempty"`
}
