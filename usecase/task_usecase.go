package usecase

import (
	"time"

	"github.com/google/uuid"
	"github.com/k07g/g2/domain"
)

type CreateTaskInput struct {
	Title       string
	Description string
}

type UpdateTaskInput struct {
	Title       string
	Description string
	Status      domain.TaskStatus
}

type TaskUseCase interface {
	GetAll() ([]*domain.Task, error)
	GetByID(id string) (*domain.Task, error)
	Create(input CreateTaskInput) (*domain.Task, error)
	Update(id string, input UpdateTaskInput) (*domain.Task, error)
	Delete(id string) error
}

type taskUseCase struct {
	repo domain.TaskRepository
}

func NewTaskUseCase(repo domain.TaskRepository) TaskUseCase {
	return &taskUseCase{repo: repo}
}

func (u *taskUseCase) GetAll() ([]*domain.Task, error) {
	return u.repo.FindAll()
}

func (u *taskUseCase) GetByID(id string) (*domain.Task, error) {
	return u.repo.FindByID(id)
}

func (u *taskUseCase) Create(input CreateTaskInput) (*domain.Task, error) {
	now := time.Now()
	task := &domain.Task{
		ID:          uuid.NewString(),
		Title:       input.Title,
		Description: input.Description,
		Status:      domain.TaskStatusTodo,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
	if err := u.repo.Save(task); err != nil {
		return nil, err
	}
	return task, nil
}

func (u *taskUseCase) Update(id string, input UpdateTaskInput) (*domain.Task, error) {
	task, err := u.repo.FindByID(id)
	if err != nil {
		return nil, err
	}
	task.Title = input.Title
	task.Description = input.Description
	task.Status = input.Status
	task.UpdatedAt = time.Now()
	if err := u.repo.Save(task); err != nil {
		return nil, err
	}
	return task, nil
}

func (u *taskUseCase) Delete(id string) error {
	return u.repo.Delete(id)
}
