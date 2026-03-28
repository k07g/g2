package inmemory_test

import (
	"errors"
	"testing"

	"github.com/k07g/g2/domain"
	"github.com/k07g/g2/infrastructure/inmemory"
)

func seed(t *testing.T, repo domain.TaskRepository, titles ...string) []*domain.Task {
	t.Helper()
	tasks := make([]*domain.Task, len(titles))
	for i, title := range titles {
		task := &domain.Task{ID: title, Title: title, Status: domain.TaskStatusTodo}
		if err := repo.Save(task); err != nil {
			t.Fatalf("Save failed: %v", err)
		}
		tasks[i] = task
	}
	return tasks
}

func TestFindAll_Empty(t *testing.T) {
	repo := inmemory.NewTaskRepository()

	tasks, err := repo.FindAll()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(tasks) != 0 {
		t.Errorf("len = %d, want 0", len(tasks))
	}
}

func TestFindAll(t *testing.T) {
	repo := inmemory.NewTaskRepository()
	seed(t, repo, "a", "b", "c")

	tasks, err := repo.FindAll()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(tasks) != 3 {
		t.Errorf("len = %d, want 3", len(tasks))
	}
}

func TestFindByID(t *testing.T) {
	repo := inmemory.NewTaskRepository()
	seed(t, repo, "task-1")

	got, err := repo.FindByID("task-1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.ID != "task-1" {
		t.Errorf("ID = %q, want %q", got.ID, "task-1")
	}
}

func TestFindByID_NotFound(t *testing.T) {
	repo := inmemory.NewTaskRepository()

	_, err := repo.FindByID("nonexistent")
	if !errors.Is(err, domain.ErrTaskNotFound) {
		t.Errorf("err = %v, want ErrTaskNotFound", err)
	}
}

func TestSave_Update(t *testing.T) {
	repo := inmemory.NewTaskRepository()
	seed(t, repo, "task-1")

	updated := &domain.Task{ID: "task-1", Title: "updated", Status: domain.TaskStatusDone}
	if err := repo.Save(updated); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	got, _ := repo.FindByID("task-1")
	if got.Title != "updated" {
		t.Errorf("Title = %q, want %q", got.Title, "updated")
	}
	if got.Status != domain.TaskStatusDone {
		t.Errorf("Status = %q, want %q", got.Status, domain.TaskStatusDone)
	}

	all, _ := repo.FindAll()
	if len(all) != 1 {
		t.Errorf("len = %d, want 1 (Save should not duplicate)", len(all))
	}
}

func TestDelete(t *testing.T) {
	repo := inmemory.NewTaskRepository()
	seed(t, repo, "task-1", "task-2")

	if err := repo.Delete("task-1"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	_, err := repo.FindByID("task-1")
	if !errors.Is(err, domain.ErrTaskNotFound) {
		t.Errorf("err = %v, want ErrTaskNotFound", err)
	}

	all, _ := repo.FindAll()
	if len(all) != 1 {
		t.Errorf("len = %d, want 1", len(all))
	}
}

func TestDelete_NotFound(t *testing.T) {
	repo := inmemory.NewTaskRepository()

	err := repo.Delete("nonexistent")
	if !errors.Is(err, domain.ErrTaskNotFound) {
		t.Errorf("err = %v, want ErrTaskNotFound", err)
	}
}
