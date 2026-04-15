package memory

import (
	"errors"
	"sync"

	"task-manager/internal/domain"
)

//
// -------- Errors --------
//

var ErrTaskNotFound = errors.New("task not found")

//
// -------- Repository Struct --------
//

type TaskRepository struct {
	mu    sync.RWMutex
	tasks map[string]*domain.Task
}

//
// -------- Constructor --------
//

func NewTaskRepository() *TaskRepository {
	return &TaskRepository{
		tasks: make(map[string]*domain.Task),
	}
}

//
// -------- Create --------
//

func (r *TaskRepository) Create(task *domain.Task) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	// Defensive check (should not happen ideally)
	if _, exists := r.tasks[task.ID]; exists {
		return errors.New("task with given ID already exists")
	}

	r.tasks[task.ID] = task
	return nil
}

//
// -------- GetByID --------
//

func (r *TaskRepository) GetByID(id string) (*domain.Task, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	task, exists := r.tasks[id]
	if !exists {
		return nil, ErrTaskNotFound
	}

	return task, nil
}

//
// -------- Update --------
//

func (r *TaskRepository) Update(task *domain.Task) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.tasks[task.ID]; !exists {
		return ErrTaskNotFound
	}

	r.tasks[task.ID] = task
	return nil
}

//
// -------- Delete --------
//

func (r *TaskRepository) Delete(id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.tasks[id]; !exists {
		return ErrTaskNotFound
	}

	delete(r.tasks, id)
	return nil
}

//
// -------- List --------
//

func (r *TaskRepository) List() ([]*domain.Task, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	tasks := make([]*domain.Task, 0, len(r.tasks))
	for _, task := range r.tasks {
		copy := *task
		tasks = append(tasks, &copy)
	}

	return tasks, nil
}
