package handlers

import (
	"context"
	"net/http"
	"time"

	"distributed-analyzer/libs/model"
	"distributed-analyzer/services/api-gateway/internal/service"
	"github.com/gin-gonic/gin"
)

type TaskHandler struct {
	taskServiceClient service.TaskServiceClient
}

func NewTaskHandler(taskService service.TaskServiceClient) *TaskHandler {
	return &TaskHandler{taskServiceClient: taskService}
}

type TaskRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
}

type UpdateTaskRequest struct {
	ID          string            `json:"id" binding:"required"`
	Name        string            `json:"name"`
	Description string            `json:"description"`
	Status      string            `json:"status"`
	Input       map[string]string `json:"input,omitempty"`
	Output      map[string]string `json:"output,omitempty"`
}

type TaskResponse struct {
	ID          string            `json:"id"`
	Name        string            `json:"name"`
	Description string            `json:"description"`
	Status      string            `json:"status"`
	Input       map[string]string `json:"input,omitempty"`
	Output      map[string]string `json:"output,omitempty"`
	CreatedAt   time.Time         `json:"created_at"`
	UpdatedAt   time.Time         `json:"updated_at"`
}

func (h *TaskHandler) Register(rg *gin.RouterGroup) {
	rg.POST("/submit", h.SubmitTask)
	rg.GET("/status/:id", h.GetTaskStatus)
	rg.PUT("/update", h.UpdateTask)
	rg.DELETE("/delete/:id", h.DeleteTask)
	rg.GET("/list", h.ListTasks)
}

// SubmitTask Submit a new task
// @Summary Submit a new task
// @Description Creates a new task in the system
// @Tags tasks
// @Accept json
// @Produce json
// @Param task body TaskRequest true "Task information"
// @Success 201 {object} TaskResponse "Task created successfully"
// @Failure 400 {object} map[string]string "Invalid request"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /api/task/submit [post]
func (h *TaskHandler) SubmitTask(c *gin.Context) {
	var req TaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request: " + err.Error()})
		return
	}

	// Create a new task
	task := &model.Task{
		Name:        req.Name,
		Description: req.Description,
		Status:      model.StatusPending,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	// Call the task service to create the task
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	createdTask, err := h.taskServiceClient.CreateTask(ctx, task)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create task: " + err.Error()})
		return
	}

	// Return the created task
	c.JSON(http.StatusCreated, TaskResponse{
		ID:          createdTask.ID,
		Name:        createdTask.Name,
		Description: createdTask.Description,
		Status:      string(createdTask.Status),
		Input:       createdTask.Input,
		Output:      createdTask.Output,
		CreatedAt:   createdTask.CreatedAt,
		UpdatedAt:   createdTask.UpdatedAt,
	})
}

// GetTaskStatus Get task status
// @Summary Get task status
// @Description Retrieves the status of a task by its ID
// @Tags tasks
// @Produce json
// @Param id path string true "Task ID"
// @Success 200 {object} TaskResponse "Task details"
// @Failure 400 {object} map[string]string "Invalid request"
// @Failure 404 {object} map[string]string "Task not found"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /api/task/status/{id} [get]
func (h *TaskHandler) GetTaskStatus(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Task ID is required"})
		return
	}

	// Call the task service to get the task
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	task, err := h.taskServiceClient.GetTask(ctx, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get task: " + err.Error()})
		return
	}

	if task == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
		return
	}

	// Return the task
	c.JSON(http.StatusOK, TaskResponse{
		ID:          task.ID,
		Name:        task.Name,
		Description: task.Description,
		Status:      string(task.Status),
		Input:       task.Input,
		Output:      task.Output,
		CreatedAt:   task.CreatedAt,
		UpdatedAt:   task.UpdatedAt,
	})
}

// UpdateTask Update a task
// @Summary Update a task
// @Description Updates an existing task in the system
// @Tags tasks
// @Accept json
// @Produce json
// @Param task body UpdateTaskRequest true "Updated task information"
// @Success 200 {object} TaskResponse "Task updated successfully"
// @Failure 400 {object} map[string]string "Invalid request"
// @Failure 404 {object} map[string]string "Task not found"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /api/task/update [put]
func (h *TaskHandler) UpdateTask(c *gin.Context) {
	var req UpdateTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request: " + err.Error()})
		return
	}

	// Get the existing task first
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	existingTask, err := h.taskServiceClient.GetTask(ctx, req.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get task: " + err.Error()})
		return
	}

	if existingTask == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
		return
	}

	// Update the task fields if provided
	if req.Name != "" {
		existingTask.Name = req.Name
	}
	if req.Description != "" {
		existingTask.Description = req.Description
	}
	if req.Status != "" {
		existingTask.Status = model.Status(req.Status)
	}
	if req.Input != nil {
		existingTask.Input = req.Input
	}
	if req.Output != nil {
		existingTask.Output = req.Output
	}
	existingTask.UpdatedAt = time.Now()

	// Call the task service to update the task
	updatedTask, err := h.taskServiceClient.UpdateTask(ctx, existingTask)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update task: " + err.Error()})
		return
	}

	// Return the updated task
	c.JSON(http.StatusOK, TaskResponse{
		ID:          updatedTask.ID,
		Name:        updatedTask.Name,
		Description: updatedTask.Description,
		Status:      string(updatedTask.Status),
		Input:       updatedTask.Input,
		Output:      updatedTask.Output,
		CreatedAt:   updatedTask.CreatedAt,
		UpdatedAt:   updatedTask.UpdatedAt,
	})
}

// DeleteTask Delete a task
// @Summary Delete a task
// @Description Deletes a task from the system
// @Tags tasks
// @Produce json
// @Param id path string true "Task ID"
// @Success 200 {object} map[string]string "Task deleted successfully"
// @Failure 400 {object} map[string]string "Invalid request"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /api/task/delete/{id} [delete]
func (h *TaskHandler) DeleteTask(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Task ID is required"})
		return
	}

	// Call the task service to delete the task
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	err := h.taskServiceClient.DeleteTask(ctx, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete task: " + err.Error()})
		return
	}

	// Return a success message
	c.JSON(http.StatusOK, gin.H{"message": "Task deleted successfully"})
}

// ListTasks List all tasks
// @Summary List all tasks
// @Description Retrieves a list of all tasks in the system
// @Tags tasks
// @Produce json
// @Success 200 {array} TaskResponse "List of tasks"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /api/task/list [get]
func (h *TaskHandler) ListTasks(c *gin.Context) {
	// Call the task service to list all tasks
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	tasks, err := h.taskServiceClient.ListTasks(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to list tasks: " + err.Error()})
		return
	}

	// Convert tasks to a response format
	var response []TaskResponse
	for _, task := range tasks {
		response = append(response, TaskResponse{
			ID:          task.ID,
			Name:        task.Name,
			Description: task.Description,
			Status:      string(task.Status),
			Input:       task.Input,
			Output:      task.Output,
			CreatedAt:   task.CreatedAt,
			UpdatedAt:   task.UpdatedAt,
		})
	}

	// Return the list of tasks
	c.JSON(http.StatusOK, response)
}
