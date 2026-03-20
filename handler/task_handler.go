package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/k07g/g2/domain"
	"github.com/k07g/g2/usecase"
)

type TaskHandler struct {
	uc usecase.TaskUseCase
}

func NewTaskHandler(uc usecase.TaskUseCase) *TaskHandler {
	return &TaskHandler{uc: uc}
}

func (h *TaskHandler) Register(mux *http.ServeMux) {
	mux.HandleFunc("GET /tasks", h.listTasks)
	mux.HandleFunc("POST /tasks", h.createTask)
	mux.HandleFunc("GET /tasks/{id}", h.getTask)
	mux.HandleFunc("PUT /tasks/{id}", h.updateTask)
	mux.HandleFunc("DELETE /tasks/{id}", h.deleteTask)
}

type taskResponse struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Status      string `json:"status"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}

func toResponse(t *domain.Task) taskResponse {
	return taskResponse{
		ID:          t.ID,
		Title:       t.Title,
		Description: t.Description,
		Status:      string(t.Status),
		CreatedAt:   t.CreatedAt.Format(time.RFC3339),
		UpdatedAt:   t.UpdatedAt.Format(time.RFC3339),
	}
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}

func writeError(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, map[string]string{"error": msg})
}

func (h *TaskHandler) listTasks(w http.ResponseWriter, r *http.Request) {
	tasks, err := h.uc.GetAll()
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	resp := make([]taskResponse, len(tasks))
	for i, t := range tasks {
		resp[i] = toResponse(t)
	}
	writeJSON(w, http.StatusOK, resp)
}

func (h *TaskHandler) getTask(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	task, err := h.uc.GetByID(id)
	if err != nil {
		if errors.Is(err, domain.ErrTaskNotFound) {
			writeError(w, http.StatusNotFound, "task not found")
			return
		}
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, toResponse(task))
}

type createTaskRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

func (h *TaskHandler) createTask(w http.ResponseWriter, r *http.Request) {
	var req createTaskRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if req.Title == "" {
		writeError(w, http.StatusBadRequest, "title is required")
		return
	}
	task, err := h.uc.Create(usecase.CreateTaskInput{
		Title:       req.Title,
		Description: req.Description,
	})
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusCreated, toResponse(task))
}

type updateTaskRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Status      string `json:"status"`
}

func (h *TaskHandler) updateTask(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var req updateTaskRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if req.Title == "" {
		writeError(w, http.StatusBadRequest, "title is required")
		return
	}
	status := domain.TaskStatus(req.Status)
	if status == "" {
		status = domain.TaskStatusTodo
	}
	task, err := h.uc.Update(id, usecase.UpdateTaskInput{
		Title:       req.Title,
		Description: req.Description,
		Status:      status,
	})
	if err != nil {
		if errors.Is(err, domain.ErrTaskNotFound) {
			writeError(w, http.StatusNotFound, "task not found")
			return
		}
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, toResponse(task))
}

func (h *TaskHandler) deleteTask(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if err := h.uc.Delete(id); err != nil {
		if errors.Is(err, domain.ErrTaskNotFound) {
			writeError(w, http.StatusNotFound, "task not found")
			return
		}
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
