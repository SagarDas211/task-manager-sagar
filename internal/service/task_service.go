package service

import (
	"errors"
	"sort"
	"time"

	"github.com/google/uuid"

	"task-manager/internal/domain"
)

//
// -------- Errors (Service-Level) --------
//

var (
	ErrTaskNotFound = errors.New("task not found")
)

//
// -------- Input Models --------
//

type CreateTaskInput struct {
	Title       string
	Description string
	Status      *domain.Status
	DueDate     time.Time
}

type UpdateTaskInput struct {
	Title       *string
	Description *string
	Status      *domain.Status
	DueDate     *time.Time
}

type TaskFilter struct {
	Status *domain.Status
	Limit  int
	Offset int
}

//
// -------- Service Struct --------
//

type TaskService struct {
	repo domain.TaskRepository
}

//
// -------- Constructor --------
//

func NewTaskService(repo domain.TaskRepository) *TaskService {
	return &TaskService{repo: repo}
}

//
// -------- Create Task --------
//

func (s *TaskService) CreateTask(input CreateTaskInput) (*domain.Task, error) {

	// Business validation
	if input.DueDate.IsZero() {
		return nil, domain.ErrInvalidDueDate
	}

	// Bonus rule: due date must be future
	if input.DueDate.Before(time.Now()) {
		return nil, domain.ErrInvalidDueDate
	}

	// Default status handling
	var status domain.Status
	if input.Status != nil {
		status = *input.Status
	}

	// Generate ID
	id := uuid.NewString()

	// Create entity
	task, err := domain.NewTask(
		id,
		input.Title,
		input.Description,
		status,
		input.DueDate,
	)
	if err != nil {
		return nil, err
	}

	// Persist
	if err := s.repo.Create(task); err != nil {
		return nil, err
	}

	return task, nil
}

//
// -------- Get Task --------
//

func (s *TaskService) GetTaskByID(id string) (*domain.Task, error) {
	task, err := s.repo.GetByID(id)
	if err != nil {
		return nil, ErrTaskNotFound
	}
	return task, nil
}

//
// -------- Update Task --------
//

func (s *TaskService) UpdateTask(id string, input UpdateTaskInput) (*domain.Task, error) {

	// Fetch existing
	task, err := s.repo.GetByID(id)
	if err != nil {
		return nil, ErrTaskNotFound
	}

	// Business validation (due date future check if provided)
	if input.DueDate != nil {
		if input.DueDate.IsZero() {
			return nil, domain.ErrInvalidDueDate
		}
		if input.DueDate.Before(time.Now()) {
			return nil, domain.ErrInvalidDueDate
		}
	}

	// Apply update via domain
	if err := task.Update(
		input.Title,
		input.Description,
		input.Status,
		input.DueDate,
	); err != nil {
		return nil, err
	}

	// Persist updated task
	if err := s.repo.Update(task); err != nil {
		return nil, err
	}

	return task, nil
}

//
// -------- Delete Task --------
//

func (s *TaskService) DeleteTask(id string) error {

	// Optional: check existence first
	_, err := s.repo.GetByID(id)
	if err != nil {
		return ErrTaskNotFound
	}

	if err := s.repo.Delete(id); err != nil {
		return err
	}

	return nil
}

// -------- List Tasks --------
func (s *TaskService) ListTasks(filter *TaskFilter) ([]*domain.Task, error) {

	tasks, err := s.repo.List()
	if err != nil {
		return nil, err
	}

	// -------- Filtering --------
	if filter != nil && filter.Status != nil {
		filtered := make([]*domain.Task, 0)
		for _, t := range tasks {
			if t.Status == *filter.Status {
				filtered = append(filtered, t)
			}
		}
		tasks = filtered
	}

	// -------- Sorting (by due date) --------
	sort.Slice(tasks, func(i, j int) bool {
		return tasks[i].DueDate.Before(tasks[j].DueDate)
	})

	// -------- Pagination --------
	if filter != nil {
		start := filter.Offset
		if start > len(tasks) {
			return []*domain.Task{}, nil
		}

		end := start + filter.Limit
		if filter.Limit == 0 || end > len(tasks) {
			end = len(tasks)
		}

		tasks = tasks[start:end]
	}

	return tasks, nil
}
