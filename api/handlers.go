package api

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"
	"time"
	"todoapp/logging"
)

type Status int

const (
	Unknown Status = iota
	NotStarted
	Started
	Completed
)

// Task represents a to-do task
type Task struct {
	ID           int        `json:"id"`
	Title        string     `json:"title"`
	Description  string     `json:"description"`
	StatusID     Status     `json:"status_id"`
	StatusString string     `json:"status"`
	CreatedAt    *time.Time `json:"created_at"`
	UpdatedAt    *time.Time `json:"updated_at"`
	DueDate      *time.Time `json:"due_date"`
	DeletedAt    *time.Time `json:"deleted_at"`
	Deleted      bool       `json:"deleted"`
}

// Handlers struct to hold shared dependencies
type Handlers struct {
	Tasks     *[]Task
	MaxTaskID *int
}

// CreateHandler handles task creation
func (h *Handlers) CreateHandler(w http.ResponseWriter, r *http.Request) {
	now := time.Now()
	slog.Info("Creating a new task")

	var taskToBeCreated Task
	err := json.NewDecoder(r.Body).Decode(&taskToBeCreated)
	if err != nil {
		slog.Error("Failed to decode request body", "error", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	statusID := ConvertStringToStatusID(taskToBeCreated.StatusString)
	if statusID == Unknown {
		slog.Error("Invalid status string", "statusString", taskToBeCreated.StatusString)
		http.Error(w, "Invalid status string", http.StatusBadRequest)
		return
	}

	taskToBeCreated.ID = *h.MaxTaskID + 1
	*h.MaxTaskID++
	taskToBeCreated.CreatedAt = &now
	taskToBeCreated.StatusID = statusID

	*h.Tasks = append(*h.Tasks, taskToBeCreated)

	w.WriteHeader(http.StatusCreated)
	slog.Info("Task created successfully")
}

// GetHandler handles fetching tasks
func (h *Handlers) GetHandler(w http.ResponseWriter, r *http.Request) {
	var currentTasks []Task
	for _, task := range *h.Tasks {
		if !task.Deleted {
			currentTasks = append(currentTasks, task)
		}
	}

	traceId := logging.GetTraceID(r.Context())
	slog.Info("Tasks retrieved successfully",
		"traceId", traceId,
		"taskCount", len(currentTasks),
	)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(currentTasks)
}

// UpdateHandler handles task updates
func (h *Handlers) UpdateHandler(w http.ResponseWriter, r *http.Request) {
	now := time.Now()
	slog.Info("Updating a task")

	var taskToBeUpdated Task
	err := json.NewDecoder(r.Body).Decode(&taskToBeUpdated)
	if err != nil {
		slog.Error("Failed to decode request body", "error", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	statusID := ConvertStringToStatusID(taskToBeUpdated.StatusString)
	if statusID == Unknown {
		slog.Error("Invalid status string", "statusString", taskToBeUpdated.StatusString)
		http.Error(w, "Invalid status string", http.StatusBadRequest)
		return
	}

	var taskFound bool
	for i, task := range *h.Tasks {
		if task.ID == taskToBeUpdated.ID && !task.Deleted {
			taskFound = true

			taskToBeUpdated.StatusID = ConvertStringToStatusID(taskToBeUpdated.StatusString)
			taskToBeUpdated.CreatedAt = task.CreatedAt
			taskToBeUpdated.UpdatedAt = &now
			(*h.Tasks)[i] = taskToBeUpdated
			break
		}
	}

	if !taskFound {
		slog.Error("Task not found", "taskID", taskToBeUpdated.ID)
		http.Error(w, "Task not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	slog.Info("Task updated successfully")
}

// DeleteHandler handles task deletion
func (h *Handlers) DeleteHandler(w http.ResponseWriter, r *http.Request) {
	now := time.Now()
	slog.Info("Deleting a task")

	path := r.URL.Path
	taskID, err := strconv.Atoi(path[len("/delete/"):])
	if err != nil {
		slog.Error("Invalid task ID", "error", err)
		http.Error(w, "Invalid task ID", http.StatusBadRequest)
		return
	}

	var taskFound bool
	for i, task := range *h.Tasks {
		if task.ID == taskID && !task.Deleted {
			taskFound = true
			(*h.Tasks)[i].Deleted = true
			(*h.Tasks)[i].DeletedAt = &now
			break
		}
	}

	if !taskFound {
		slog.Error("Task not found", "taskID", taskID)
		http.Error(w, "Task not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	slog.Info("Task deleted successfully")
}

// ConvertStatus converts status string into status ID
func ConvertStringToStatusID(statusString string) Status {
	switch statusString {
	case "Not Started":
		return NotStarted
	case "Started":
		return Started
	case "Completed":
		return Completed
	default:
		return Unknown
	}
}

// ConvertStatus converts status string into status ID
func ConvertStatusIDToString(statusID Status) string {
	switch statusID {
	case NotStarted:
		return "Not Started"
	case Started:
		return "Started"
	case Completed:
		return "Completed"
	}

	return ""
}
