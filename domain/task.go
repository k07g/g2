package domain

import (
	"errors"
	"time"
)

var ErrTaskNotFound = errors.New("task not found")

type TaskStatus string

const (
	TaskStatusTodo       TaskStatus = "todo"
	TaskStatusInProgress TaskStatus = "in_progress"
	TaskStatusDone       TaskStatus = "done"
)

type Task struct {
	ID          string
	Title       string
	Description string
	Status      TaskStatus
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type TaskRepository interface {
	FindAll() ([]*Task, error)
	FindByID(id string) (*Task, error)
	Save(task *Task) error
	Delete(id string) error
}
