package client

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/distributedmarketplace/internal/task/model"
	"io"
	"net/http"
	"strings"
	"time"
)

// TaskClient is a client for the task service
type TaskClient struct {
	baseURL    string
	httpClient *http.Client
}

// NewTaskClient creates a new TaskClient
func NewTaskClient(baseURL string) *TaskClient {
	return &TaskClient{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// CreateTask creates a new task in the system
func (c *TaskClient) CreateTask(ctx context.Context, task *model.Task) (*model.Task, error) {
	url := fmt.Sprintf("%s/tasks", c.baseURL)

	reqBody, err := json.Marshal(task)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal task: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, strings.NewReader(string(reqBody)))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(resp.Body)

	if resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var createdTask model.Task
	if err := json.NewDecoder(resp.Body).Decode(&createdTask); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &createdTask, nil
}

// GetTask retrieves a task by its ID
func (c *TaskClient) GetTask(ctx context.Context, id string) (*model.Task, error) {
	url := fmt.Sprintf("%s/tasks/%s", c.baseURL, id)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(resp.Body)

	if resp.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("task not found")
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var task model.Task
	if err := json.NewDecoder(resp.Body).Decode(&task); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &task, nil
}

// UpdateTask updates an existing task
func (c *TaskClient) UpdateTask(ctx context.Context, task *model.Task) (*model.Task, error) {
	url := fmt.Sprintf("%s/tasks/%s", c.baseURL, task.ID)

	reqBody, err := json.Marshal(task)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal task: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPut, url, strings.NewReader(string(reqBody)))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(resp.Body)

	// Check response status
	if resp.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("task not found")
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	// Parse response
	var updatedTask model.Task
	if err := json.NewDecoder(resp.Body).Decode(&updatedTask); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &updatedTask, nil
}

// DeleteTask removes a task from the system
func (c *TaskClient) DeleteTask(ctx context.Context, id string) error {
	url := fmt.Sprintf("%s/tasks/%s", c.baseURL, id)

	// Create request
	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, url, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	// Send request
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(resp.Body)

	// Check response status
	if resp.StatusCode == http.StatusNotFound {
		return fmt.Errorf("task not found")
	}
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return nil
}

// ListTasks retrieves all tasks with optional filtering
func (c *TaskClient) ListTasks(ctx context.Context) ([]*model.Task, error) {
	url := fmt.Sprintf("%s/tasks", c.baseURL)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(resp.Body)

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var tasks []*model.Task
	if err := json.NewDecoder(resp.Body).Decode(&tasks); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return tasks, nil
}
