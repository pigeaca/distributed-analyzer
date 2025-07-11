package service

import (
	"context"
	"errors"
	"github.com/pigeaca/DistributedMarketplace/libs/model"
	"sync"
	"time"
)

var ErrTaskNotFound = errors.New("task not found")

// TaskServiceImpl implements the TaskService interface
type TaskServiceImpl struct {
	tasks  map[string]*model.Task
	taskMu sync.RWMutex
}

// NewTaskServiceImpl creates a new instance of TaskServiceImpl
func NewTaskServiceImpl() *TaskServiceImpl {
	return &TaskServiceImpl{
		tasks: make(map[string]*model.Task),
	}
}

// CreateTask creates a new task in the system
func (t *TaskServiceImpl) CreateTask(ctx context.Context, task *model.Task) (*model.Task, error) {
	t.taskMu.Lock()
	defer t.taskMu.Unlock()

	// In a real implementation, we would generate a proper UUID
	if task.ID == "" {
		task.ID = time.Now().Format("20060102150405")
	}

	task.Status = model.StatusPending
	task.CreatedAt = time.Now()
	task.UpdatedAt = time.Now()

	t.tasks[task.ID] = task
	return task, nil
}

// GetTask retrieves a task by its ID
func (t *TaskServiceImpl) GetTask(ctx context.Context, id string) (*model.Task, error) {
	t.taskMu.RLock()
	defer t.taskMu.RUnlock()

	task, exists := t.tasks[id]
	if !exists {
		return nil, ErrTaskNotFound
	}

	return task, nil
}

// UpdateTask updates an existing task
func (t *TaskServiceImpl) UpdateTask(ctx context.Context, task *model.Task) (*model.Task, error) {
	t.taskMu.Lock()
	defer t.taskMu.Unlock()

	existingTask, exists := t.tasks[task.ID]
	if !exists {
		return nil, ErrTaskNotFound
	}

	// Update fields
	existingTask.Name = task.Name
	existingTask.Description = task.Description
	existingTask.Status = task.Status
	existingTask.UpdatedAt = time.Now()

	if task.Status == model.StatusCompleted && existingTask.CompletedAt.IsZero() {
		existingTask.CompletedAt = time.Now()
	}

	return existingTask, nil
}

// DeleteTask removes a task from the system
func (t *TaskServiceImpl) DeleteTask(ctx context.Context, id string) error {
	t.taskMu.Lock()
	defer t.taskMu.Unlock()

	if _, exists := t.tasks[id]; !exists {
		return ErrTaskNotFound
	}

	delete(t.tasks, id)
	return nil
}

// ListTasks retrieves all tasks with optional filtering
func (t *TaskServiceImpl) ListTasks(ctx context.Context) ([]*model.Task, error) {
	t.taskMu.RLock()
	defer t.taskMu.RUnlock()

	tasks := make([]*model.Task, 0, len(t.tasks))
	for _, task := range t.tasks {
		tasks = append(tasks, task)
	}

	return tasks, nil
}
