package handler_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/k07g/g2/handler"
	"github.com/k07g/g2/infrastructure/inmemory"
	"github.com/k07g/g2/usecase"
)

func newMux() *http.ServeMux {
	uc := usecase.NewTaskUseCase(inmemory.NewTaskRepository())
	mux := http.NewServeMux()
	handler.NewTaskHandler(uc).Register(mux)
	return mux
}

func do(mux *http.ServeMux, method, path, body string) *httptest.ResponseRecorder {
	var b *strings.Reader
	if body != "" {
		b = strings.NewReader(body)
	} else {
		b = strings.NewReader("")
	}
	req := httptest.NewRequest(method, path, b)
	if body != "" {
		req.Header.Set("Content-Type", "application/json")
	}
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)
	return w
}

// POST /tasks でタスクを作成してIDを返す
func createTask(t *testing.T, mux *http.ServeMux, title string) string {
	t.Helper()
	w := do(mux, http.MethodPost, "/tasks", `{"title":"`+title+`"}`)
	if w.Code != http.StatusCreated {
		t.Fatalf("POST /tasks status = %d, want %d", w.Code, http.StatusCreated)
	}
	var resp map[string]any
	json.NewDecoder(w.Body).Decode(&resp)
	return resp["id"].(string)
}

func TestListTasks_Empty(t *testing.T) {
	w := do(newMux(), http.MethodGet, "/tasks", "")

	if w.Code != http.StatusOK {
		t.Errorf("status = %d, want %d", w.Code, http.StatusOK)
	}
	var tasks []any
	json.NewDecoder(w.Body).Decode(&tasks)
	if len(tasks) != 0 {
		t.Errorf("len = %d, want 0", len(tasks))
	}
}

func TestListTasks(t *testing.T) {
	mux := newMux()
	createTask(t, mux, "タスク1")
	createTask(t, mux, "タスク2")

	w := do(mux, http.MethodGet, "/tasks", "")

	if w.Code != http.StatusOK {
		t.Errorf("status = %d, want %d", w.Code, http.StatusOK)
	}
	var tasks []any
	json.NewDecoder(w.Body).Decode(&tasks)
	if len(tasks) != 2 {
		t.Errorf("len = %d, want 2", len(tasks))
	}
}

func TestCreateTask(t *testing.T) {
	w := do(newMux(), http.MethodPost, "/tasks", `{"title":"買い物","description":"牛乳"}`)

	if w.Code != http.StatusCreated {
		t.Errorf("status = %d, want %d", w.Code, http.StatusCreated)
	}
	var resp map[string]any
	json.NewDecoder(w.Body).Decode(&resp)
	if resp["id"] == "" {
		t.Error("id should not be empty")
	}
	if resp["title"] != "買い物" {
		t.Errorf("title = %v, want 買い物", resp["title"])
	}
	if resp["status"] != "todo" {
		t.Errorf("status = %v, want todo", resp["status"])
	}
}

func TestCreateTask_NoTitle(t *testing.T) {
	w := do(newMux(), http.MethodPost, "/tasks", `{"description":"説明のみ"}`)

	if w.Code != http.StatusBadRequest {
		t.Errorf("status = %d, want %d", w.Code, http.StatusBadRequest)
	}
}

func TestCreateTask_InvalidJSON(t *testing.T) {
	w := do(newMux(), http.MethodPost, "/tasks", `not-json`)

	if w.Code != http.StatusBadRequest {
		t.Errorf("status = %d, want %d", w.Code, http.StatusBadRequest)
	}
}

func TestGetTask(t *testing.T) {
	mux := newMux()
	id := createTask(t, mux, "テスト")

	w := do(mux, http.MethodGet, "/tasks/"+id, "")

	if w.Code != http.StatusOK {
		t.Errorf("status = %d, want %d", w.Code, http.StatusOK)
	}
	var resp map[string]any
	json.NewDecoder(w.Body).Decode(&resp)
	if resp["id"] != id {
		t.Errorf("id = %v, want %v", resp["id"], id)
	}
}

func TestGetTask_NotFound(t *testing.T) {
	w := do(newMux(), http.MethodGet, "/tasks/nonexistent", "")

	if w.Code != http.StatusNotFound {
		t.Errorf("status = %d, want %d", w.Code, http.StatusNotFound)
	}
}

func TestUpdateTask(t *testing.T) {
	mux := newMux()
	id := createTask(t, mux, "元タイトル")

	w := do(mux, http.MethodPut, "/tasks/"+id, `{"title":"新タイトル","status":"done"}`)

	if w.Code != http.StatusOK {
		t.Errorf("status = %d, want %d", w.Code, http.StatusOK)
	}
	var resp map[string]any
	json.NewDecoder(w.Body).Decode(&resp)
	if resp["title"] != "新タイトル" {
		t.Errorf("title = %v, want 新タイトル", resp["title"])
	}
	if resp["status"] != "done" {
		t.Errorf("status = %v, want done", resp["status"])
	}
}

func TestUpdateTask_NotFound(t *testing.T) {
	w := do(newMux(), http.MethodPut, "/tasks/nonexistent", `{"title":"x"}`)

	if w.Code != http.StatusNotFound {
		t.Errorf("status = %d, want %d", w.Code, http.StatusNotFound)
	}
}

func TestUpdateTask_NoTitle(t *testing.T) {
	mux := newMux()
	id := createTask(t, mux, "タスク")

	w := do(mux, http.MethodPut, "/tasks/"+id, `{"status":"done"}`)

	if w.Code != http.StatusBadRequest {
		t.Errorf("status = %d, want %d", w.Code, http.StatusBadRequest)
	}
}

func TestDeleteTask(t *testing.T) {
	mux := newMux()
	id := createTask(t, mux, "削除対象")

	w := do(mux, http.MethodDelete, "/tasks/"+id, "")

	if w.Code != http.StatusNoContent {
		t.Errorf("status = %d, want %d", w.Code, http.StatusNoContent)
	}

	w2 := do(mux, http.MethodGet, "/tasks/"+id, "")
	if w2.Code != http.StatusNotFound {
		t.Errorf("after delete: status = %d, want %d", w2.Code, http.StatusNotFound)
	}
}

func TestDeleteTask_NotFound(t *testing.T) {
	w := do(newMux(), http.MethodDelete, "/tasks/nonexistent", "")

	if w.Code != http.StatusNotFound {
		t.Errorf("status = %d, want %d", w.Code, http.StatusNotFound)
	}
}
