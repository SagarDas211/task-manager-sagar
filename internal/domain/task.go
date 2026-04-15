package domain

import (
	"errors"
	"strings"
	"time"
)

//
// -------- Status Enum --------
//

type Status string

const (
	StatusPending    Status = "PENDING"
	StatusInProgress Status = "IN_PROGRESS"
	StatusDone       Status = "DONE"
)

func IsValidStatus(s Status) bool {
	switch s {
	case StatusPending, StatusInProgress, StatusDone:
		return true
	default:
		return false
	}
}

//
// -------- Errors --------
//

var (
	ErrInvalidTitle   = errors.New("title cannot be empty")
	ErrInvalidStatus  = errors.New("invalid status")
	ErrInvalidDueDate = errors.New("invalid due date")
)

//
// -------- Task Entity --------
//

type Task struct {
	ID          string
	Title       string
	Description string
	Status      Status
	DueDate     time.Time
}

//
// -------- Constructor --------
//

func NewTask(
	id string,
	title string,
	description string,
	status Status,
	dueDate time.Time,
) (*Task, error) {

	title = strings.TrimSpace(title)
	if title == "" {
		return nil, ErrInvalidTitle
	}

	if dueDate.IsZero() {
		return nil, ErrInvalidDueDate
	}

	// Default status
	if status == "" {
		status = StatusPending
	}

	if !IsValidStatus(status) {
		return nil, ErrInvalidStatus
	}

	task := &Task{
		ID:          id,
		Title:       title,
		Description: description,
		Status:      status,
		DueDate:     dueDate,
	}

	return task, nil
}

//
// -------- Validation Method --------
//

func (t *Task) Validate() error {
	if strings.TrimSpace(t.Title) == "" {
		return ErrInvalidTitle
	}

	if !IsValidStatus(t.Status) {
		return ErrInvalidStatus
	}

	if t.DueDate.IsZero() {
		return ErrInvalidDueDate
	}

	return nil
}

//
// -------- Update Method (Partial Update) --------
//

func (t *Task) Update(
	title *string,
	description *string,
	status *Status,
	dueDate *time.Time,
) error {

	if title != nil {
		trimmed := strings.TrimSpace(*title)
		if trimmed == "" {
			return ErrInvalidTitle
		}
		t.Title = trimmed
	}

	if description != nil {
		t.Description = *description
	}

	if status != nil {
		if !IsValidStatus(*status) {
			return ErrInvalidStatus
		}
		t.Status = *status
	}

	if dueDate != nil {
		if dueDate.IsZero() {
			return ErrInvalidDueDate
		}
		t.DueDate = *dueDate
	}

	return nil
}
