package handler

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"

	"task-manager/internal/domain"
	"task-manager/internal/service"
)

type TaskHandler struct {
	service *service.TaskService
}

func NewTaskHandler(s *service.TaskService) *TaskHandler {
	return &TaskHandler{service: s}
}

type CreateTaskRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Status      string `json:"status"`
	DueDate     string `json:"due_date"`
}

type UpdateTaskRequest struct {
	Title       *string `json:"title"`
	Description *string `json:"description"`
	Status      *string `json:"status"`
	DueDate     *string `json:"due_date"`
}

type TaskResponse struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Status      string `json:"status"`
	DueDate     string `json:"due_date"`
}

func toTaskResponse(t *domain.Task) TaskResponse {
	return TaskResponse{
		ID:          t.ID,
		Title:       t.Title,
		Description: t.Description,
		Status:      string(t.Status),
		DueDate:     t.DueDate.Format("2006-01-02"),
	}
}

func parseDate(dateStr string) (time.Time, error) {
	return time.Parse("2006-01-02", dateStr)
}

func (h *TaskHandler) CreateTask(c *gin.Context) {
	var req CreateTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	if req.Title == "" || req.DueDate == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "title and due_date are required"})
		return
	}

	dueDate, err := parseDate(req.DueDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid due_date format (YYYY-MM-DD)"})
		return
	}

	var status *domain.Status
	if req.Status != "" {
		s := domain.Status(req.Status)
		status = &s
	}

	task, err := h.service.CreateTask(service.CreateTaskInput{
		Title:       req.Title,
		Description: req.Description,
		Status:      status,
		DueDate:     dueDate,
	})
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusCreated, toTaskResponse(task))
}

func (h *TaskHandler) GetTask(c *gin.Context) {
	id := c.Param("id")

	task, err := h.service.GetTaskByID(id)
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, toTaskResponse(task))
}

func (h *TaskHandler) UpdateTask(c *gin.Context) {
	id := c.Param("id")

	var req UpdateTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	var status *domain.Status
	if req.Status != nil {
		s := domain.Status(*req.Status)
		status = &s
	}

	var dueDate *time.Time
	if req.DueDate != nil {
		parsed, err := parseDate(*req.DueDate)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid due_date format (YYYY-MM-DD)"})
			return
		}
		dueDate = &parsed
	}

	task, err := h.service.UpdateTask(id, service.UpdateTaskInput{
		Title:       req.Title,
		Description: req.Description,
		Status:      status,
		DueDate:     dueDate,
	})
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, toTaskResponse(task))
}

func (h *TaskHandler) DeleteTask(c *gin.Context) {
	id := c.Param("id")

	if err := h.service.DeleteTask(id); err != nil {
		handleError(c, err)
		return
	}

	c.Status(http.StatusNoContent)
}

func (h *TaskHandler) ListTasks(c *gin.Context) {
	statusParam := c.Query("status")
	limitParam := c.Query("limit")
	offsetParam := c.Query("offset")

	var filter service.TaskFilter

	if statusParam != "" {
		s := domain.Status(statusParam)
		if !domain.IsValidStatus(s) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid status"})
			return
		}
		filter.Status = &s
	}

	if limitParam != "" {
		limit, err := strconv.Atoi(limitParam)
		if err != nil || limit < 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid limit"})
			return
		}
		filter.Limit = limit
	}

	if offsetParam != "" {
		offset, err := strconv.Atoi(offsetParam)
		if err != nil || offset < 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid offset"})
			return
		}
		filter.Offset = offset
	}

	tasks, err := h.service.ListTasks(&filter)
	if err != nil {
		handleError(c, err)
		return
	}

	resp := make([]TaskResponse, 0, len(tasks))
	for _, t := range tasks {
		resp = append(resp, toTaskResponse(t))
	}

	c.JSON(http.StatusOK, resp)
}

func handleError(c *gin.Context, err error) {
	switch err {
	case service.ErrTaskNotFound:
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
	case domain.ErrInvalidTitle, domain.ErrInvalidStatus, domain.ErrInvalidDueDate:
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	default:
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
	}
}
