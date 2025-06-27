package handlers

import (
	"errors"
	"github.com/distributedmarketplace/internal/gateway/service"
	"github.com/distributedmarketplace/internal/task/model"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

var ErrTaskNotFound = errors.New("task not found")

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

func (h *TaskHandler) Register(rg *gin.RouterGroup) {
	rg.POST("/tasks", h.SubmitTask)
}

// SubmitTask handles the creation of a new task
func (h *TaskHandler) SubmitTask(c *gin.Context) {
	var req TaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request: " + err.Error(),
		})
		return
	}

	task := &model.Task{
		Name:        req.Name,
		Description: req.Description,
		Status:      model.StatusPending,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	createdTask, err := h.taskServiceClient.CreateTask(c.Request.Context(), task)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to create task: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, createdTask)
}

func (h *TaskHandler) GetTaskStatus(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Task ID is required",
		})
		return
	}

	task, err := h.taskServiceClient.GetTask(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, ErrTaskNotFound) {
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
