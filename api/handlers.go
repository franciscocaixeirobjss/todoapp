package api

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"
	"todoapp/task"
)

type Handlers struct {
	TaskManager *task.TaskManager
}

func (h *Handlers) CreateHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		slog.Error("Invalid HTTP method", "method", r.Method)
		http.Error(w, "Invalid HTTP method.", http.StatusMethodNotAllowed)
		return
	}

	var taskToBeCreated task.Task
	err := json.NewDecoder(r.Body).Decode(&taskToBeCreated)
	if err != nil {
		slog.Error("Failed to decode request body", "error", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	err = h.TaskManager.AddTask(taskToBeCreated)
	if err != nil {
		slog.Error("Failed to create task", "error", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusCreated)
	slog.Info("Task created successfully")
}

func (h *Handlers) GetHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		slog.Error("Invalid HTTP method", "method", r.Method)
		http.Error(w, "Invalid HTTP method.", http.StatusMethodNotAllowed)
		return
	}

	currentTasks := h.TaskManager.GetTasks()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(currentTasks)
}

func (h *Handlers) UpdateHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		slog.Error("Invalid HTTP method", "method", r.Method)
		http.Error(w, "Invalid HTTP method.", http.StatusMethodNotAllowed)
		return
	}

	var taskToBeUpdated task.Task
	err := json.NewDecoder(r.Body).Decode(&taskToBeUpdated)
	if err != nil {
		slog.Error("Failed to decode request body", "error", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	err = h.TaskManager.UpdateTask(taskToBeUpdated)
	if err != nil {
		if err == task.ErrTaskNotFound {
			slog.Error("Task not found", "taskID", taskToBeUpdated.ID)
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}

		slog.Error("Failed to update task", "error", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	slog.Info("Task updated successfully", "taskID", taskToBeUpdated.ID)
}

func (h *Handlers) DeleteHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		slog.Error("Invalid HTTP method", "method", r.Method)
		http.Error(w, "Invalid HTTP method.", http.StatusMethodNotAllowed)
		return
	}

	// Extract the task ID from the URL path
	path := r.URL.Path
	taskID, err := strconv.Atoi(path[len("/delete/"):])
	if err != nil {
		slog.Error("Invalid task ID", "error", err)
		http.Error(w, "Invalid task ID", http.StatusBadRequest)
		return
	}

	err = h.TaskManager.DeleteTask(taskID)
	if err != nil {
		slog.Error("Failed to delete task", "error", err)
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	slog.Info("Task deleted successfully", "taskID", taskID)
}
