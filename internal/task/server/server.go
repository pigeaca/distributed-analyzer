package server

import (
	"github.com/distributedmarketplace/internal/task/model"
	"github.com/distributedmarketplace/internal/task/service"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

// TaskServer is a server for the task service
type TaskServer struct {
	taskService service.TaskService
	router      *gin.Engine
}

// NewTaskServer creates a new TaskServer
func NewTaskServer(taskService service.TaskService) *TaskServer {
	server := &TaskServer{
		taskService: taskService,
		router:      gin.Default(),
	}
	server.setupRoutes()
	return server
}

// setupRoutes sets up the HTTP routes for the task service
func (s *TaskServer) setupRoutes() {
	s.router.POST("/tasks", s.createTask)
	s.router.GET("/tasks/:id", s.getTask)
	s.router.PUT("/tasks/:id", s.updateTask)
	s.router.DELETE("/tasks/:id", s.deleteTask)
	s.router.GET("/tasks", s.listTasks)
}

// Run starts the HTTP server
func (s *TaskServer) Run(addr string) error {
	return s.router.Run(addr)
}

// createTask handles the creation of a new task
func (s *TaskServer) createTask(c *gin.Context) {
	var taskRequest struct {
		Name        string            `json:"name" binding:"required"`
		Description string            `json:"description"`
		Input       map[string]string `json:"input,omitempty"`
		Resources   []model.Resource  `json:"resources,omitempty"`
	}

	if err := c.ShouldBindJSON(&taskRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request: " + err.Error(),
		})
		return
	}

	task := &model.Task{
		Name:        taskRequest.Name,
		Description: taskRequest.Description,
		Input:       taskRequest.Input,
		Resources:   taskRequest.Resources,
		Status:      model.StatusPending,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	createdTask, err := s.taskService.CreateTask(c.Request.Context(), task)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to create task: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, createdTask)
}

// getTask retrieves a task by its ID
func (s *TaskServer) getTask(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Task ID is required",
		})
		return
	}

	task, err := s.taskService.GetTask(c.Request.Context(), id)
	if err != nil {
		if err.Error() == "task not found" {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Task not found",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to get task: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, task)
}

// updateTask updates an existing task
func (s *TaskServer) updateTask(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Task ID is required",
		})
		return
	}

	var taskRequest model.Task
	if err := c.ShouldBindJSON(&taskRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request: " + err.Error(),
		})
		return
	}

	// Ensure the ID in the path matches the ID in the request body
	taskRequest.ID = id

	updatedTask, err := s.taskService.UpdateTask(c.Request.Context(), &taskRequest)
	if err != nil {
		if err.Error() == "task not found" {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Task not found",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to update task: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, updatedTask)
}

// deleteTask removes a task from the system
func (s *TaskServer) deleteTask(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Task ID is required",
		})
		return
	}

	err := s.taskService.DeleteTask(c.Request.Context(), id)
	if err != nil {
		if err.Error() == "task not found" {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Task not found",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to delete task: " + err.Error(),
		})
		return
	}

	c.Status(http.StatusNoContent)
}

// listTasks retrieves all tasks with optional filtering
func (s *TaskServer) listTasks(c *gin.Context) {
	tasks, err := s.taskService.ListTasks(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to list tasks: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, tasks)
}
