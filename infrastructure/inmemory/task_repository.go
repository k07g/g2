package inmemory

import (
	"sync"

	"github.com/k07g/g2/domain"
)

type taskRepository struct {
	mu    sync.RWMutex
	tasks map[string]*domain.Task
}

func NewTaskRepository() domain.TaskRepository {
	return &taskRepository{
		tasks: make(map[string]*domain.Task),
	}
}

func (r *taskRepository) FindAll() ([]*domain.Task, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	result := make([]*domain.Task, 0, len(r.tasks))
	for _, t := range r.tasks {
		result = append(result, t)
	}
	return result, nil
}

func (r *taskRepository) FindByID(id string) (*domain.Task, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	t, ok := r.tasks[id]
	if !ok {
		return nil, domain.ErrTaskNotFound
	}
	return t, nil
}

func (r *taskRepository) Save(task *domain.Task) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.tasks[task.ID] = task
	return nil
}

func (r *taskRepository) Delete(id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.tasks[id]; !ok {
		return domain.ErrTaskNotFound
	}
	delete(r.tasks, id)
	return nil
}
