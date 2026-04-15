package service

import (
	"errors"
	"sync"
	"testing"
	"time"

	"task-manager/internal/domain"
)

type mockTaskRepository struct {
	mu    sync.Mutex
	tasks map[string]*domain.Task
}

func newMockTaskRepository() *mockTaskRepository {
	return &mockTaskRepository{tasks: make(map[string]*domain.Task)}
}

func (m *mockTaskRepository) Create(task *domain.Task) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.tasks[task.ID] = task
	return nil
}

func (m *mockTaskRepository) GetByID(id string) (*domain.Task, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	task, ok := m.tasks[id]
	if !ok {
		return nil, errors.New("not found")
	}
	return task, nil
}

func (m *mockTaskRepository) Update(task *domain.Task) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if _, ok := m.tasks[task.ID]; !ok {
		return errors.New("not found")
	}
	m.tasks[task.ID] = task
	return nil
}

func (m *mockTaskRepository) Delete(id string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if _, ok := m.tasks[id]; !ok {
		return errors.New("not found")
	}
	delete(m.tasks, id)
	return nil
}

func (m *mockTaskRepository) List() ([]*domain.Task, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	out := make([]*domain.Task, 0, len(m.tasks))
	for _, t := range m.tasks {
		out = append(out, t)
	}
	return out, nil
}

func TestCreateTask_Success(t *testing.T) {
	repo := newMockTaskRepository()
	svc := NewTaskService(repo)

	input := CreateTaskInput{
		Title:   "Test Task",
		DueDate: time.Now().Add(24 * time.Hour),
	}

	task, err := svc.CreateTask(input)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if task.Title != "Test Task" {
		t.Fatalf("expected title to be set")
	}

	if task.Status != domain.StatusPending {
		t.Fatalf("expected default status %s, got %s", domain.StatusPending, task.Status)
	}
}

func TestCreateTask_InvalidDueDate(t *testing.T) {
	repo := newMockTaskRepository()
	svc := NewTaskService(repo)

	input := CreateTaskInput{
		Title:   "Test",
		DueDate: time.Time{},
	}

	_, err := svc.CreateTask(input)
	if err == nil {
		t.Fatalf("expected error for invalid due date")
	}
}

func TestGetTask_Success(t *testing.T) {
	repo := newMockTaskRepository()
	svc := NewTaskService(repo)

	task, err := domain.NewTask("1", "Test", "", domain.StatusPending, time.Now().Add(time.Hour))
	if err != nil {
		t.Fatalf("failed to create task: %v", err)
	}
	repo.tasks["1"] = task

	result, err := svc.GetTaskByID("1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.ID != "1" {
		t.Fatalf("expected ID 1")
	}
}

func TestGetTask_NotFound(t *testing.T) {
	repo := newMockTaskRepository()
	svc := NewTaskService(repo)

	_, err := svc.GetTaskByID("999")
	if !errors.Is(err, ErrTaskNotFound) {
		t.Fatalf("expected ErrTaskNotFound, got %v", err)
	}
}

func TestUpdateTask_Success(t *testing.T) {
	repo := newMockTaskRepository()
	svc := NewTaskService(repo)

	task, err := domain.NewTask("1", "Old", "", domain.StatusPending, time.Now().Add(time.Hour))
	if err != nil {
		t.Fatalf("failed to create task: %v", err)
	}
	repo.tasks["1"] = task

	newTitle := "Updated"
	input := UpdateTaskInput{Title: &newTitle}

	updated, err := svc.UpdateTask("1", input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if updated.Title != "Updated" {
		t.Fatalf("expected updated title")
	}
}

func TestUpdateTask_NotFound(t *testing.T) {
	repo := newMockTaskRepository()
	svc := NewTaskService(repo)

	title := "New"
	input := UpdateTaskInput{Title: &title}

	_, err := svc.UpdateTask("1", input)
	if !errors.Is(err, ErrTaskNotFound) {
		t.Fatalf("expected ErrTaskNotFound, got %v", err)
	}
}

func TestDeleteTask_Success(t *testing.T) {
	repo := newMockTaskRepository()
	svc := NewTaskService(repo)

	task, err := domain.NewTask("1", "Test", "", domain.StatusPending, time.Now().Add(time.Hour))
	if err != nil {
		t.Fatalf("failed to create task: %v", err)
	}
	repo.tasks["1"] = task

	err = svc.DeleteTask("1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestDeleteTask_NotFound(t *testing.T) {
	repo := newMockTaskRepository()
	svc := NewTaskService(repo)

	err := svc.DeleteTask("1")
	if !errors.Is(err, ErrTaskNotFound) {
		t.Fatalf("expected ErrTaskNotFound, got %v", err)
	}
}
