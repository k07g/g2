package usecase_test

import (
	"errors"
	"testing"

	"github.com/k07g/g2/domain"
	"github.com/k07g/g2/infrastructure/inmemory"
	"github.com/k07g/g2/usecase"
)

func newUseCase() usecase.TaskUseCase {
	return usecase.NewTaskUseCase(inmemory.NewTaskRepository())
}

func TestCreate(t *testing.T) {
	uc := newUseCase()

	task, err := uc.Create(usecase.CreateTaskInput{Title: "買い物", Description: "牛乳"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if task.ID == "" {
		t.Error("ID should be set")
	}
	if task.Title != "買い物" {
		t.Errorf("Title = %q, want %q", task.Title, "買い物")
	}
	if task.Status != domain.TaskStatusTodo {
		t.Errorf("Status = %q, want %q", task.Status, domain.TaskStatusTodo)
	}
}

func TestCreate_TitleRequired(t *testing.T) {
	uc := newUseCase()

	// title が空でも usecase は弾かない（バリデーションは handler 層の責務）
	task, err := uc.Create(usecase.CreateTaskInput{Title: ""})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if task.Title != "" {
		t.Errorf("Title = %q, want empty", task.Title)
	}
}

func TestGetAll(t *testing.T) {
	uc := newUseCase()

	for _, title := range []string{"タスク1", "タスク2", "タスク3"} {
		if _, err := uc.Create(usecase.CreateTaskInput{Title: title}); err != nil {
			t.Fatalf("Create failed: %v", err)
		}
	}

	tasks, err := uc.GetAll()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(tasks) != 3 {
		t.Errorf("len = %d, want 3", len(tasks))
	}
}

func TestGetByID(t *testing.T) {
	uc := newUseCase()

	created, _ := uc.Create(usecase.CreateTaskInput{Title: "テスト"})

	got, err := uc.GetByID(created.ID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.ID != created.ID {
		t.Errorf("ID = %q, want %q", got.ID, created.ID)
	}
}

func TestGetByID_NotFound(t *testing.T) {
	uc := newUseCase()

	_, err := uc.GetByID("nonexistent")
	if !errors.Is(err, domain.ErrTaskNotFound) {
		t.Errorf("err = %v, want ErrTaskNotFound", err)
	}
}

func TestUpdate(t *testing.T) {
	uc := newUseCase()

	created, _ := uc.Create(usecase.CreateTaskInput{Title: "元タイトル"})

	updated, err := uc.Update(created.ID, usecase.UpdateTaskInput{
		Title:  "新タイトル",
		Status: domain.TaskStatusDone,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if updated.Title != "新タイトル" {
		t.Errorf("Title = %q, want %q", updated.Title, "新タイトル")
	}
	if updated.Status != domain.TaskStatusDone {
		t.Errorf("Status = %q, want %q", updated.Status, domain.TaskStatusDone)
	}
	if updated.UpdatedAt.Before(created.CreatedAt) {
		t.Error("UpdatedAt should not be before CreatedAt")
	}
}

func TestUpdate_NotFound(t *testing.T) {
	uc := newUseCase()

	_, err := uc.Update("nonexistent", usecase.UpdateTaskInput{Title: "x"})
	if !errors.Is(err, domain.ErrTaskNotFound) {
		t.Errorf("err = %v, want ErrTaskNotFound", err)
	}
}

func TestDelete(t *testing.T) {
	uc := newUseCase()

	created, _ := uc.Create(usecase.CreateTaskInput{Title: "削除対象"})

	if err := uc.Delete(created.ID); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	_, err := uc.GetByID(created.ID)
	if !errors.Is(err, domain.ErrTaskNotFound) {
		t.Errorf("err = %v, want ErrTaskNotFound", err)
	}
}

func TestDelete_NotFound(t *testing.T) {
	uc := newUseCase()

	err := uc.Delete("nonexistent")
	if !errors.Is(err, domain.ErrTaskNotFound) {
		t.Errorf("err = %v, want ErrTaskNotFound", err)
	}
}
