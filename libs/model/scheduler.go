package model

// SchedulerTask represents a task from the scheduler's perspective
type SchedulerTask struct {
	TaskID    string    `json:"task_id"`
	Status    Status    `json:"status"`
	SubTasks  []SubTask `json:"subtasks,omitempty"`
	WorkerIDs []string  `json:"worker_ids,omitempty"`
}
