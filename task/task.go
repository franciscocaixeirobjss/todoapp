package task

import (
	"errors"
	"strings"
	"time"
)

// Status represents the status of a task
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

// TaskManager struct to manage tasks and their state
type TaskManager struct {
	Tasks     []Task
	MaxTaskID int
}

var (
	// ErrTaskNotFound is returned when a task is not found
	ErrTaskNotFound = errors.New("task not found")
	// ErrTaskAlreadyExists is returned when a task already exists
	ErrInvalidStatus = errors.New("invalid status string")
)

// AddTask adds a new task to the list of tasks
func (tm *TaskManager) AddTask(task Task) error {
	now := time.Now()

	statusID, err := convertStringToStatusID(task.StatusString)
	if err != nil {
		return ErrInvalidStatus
	}

	task.ID = tm.MaxTaskID + 1
	tm.MaxTaskID++
	task.CreatedAt = &now
	task.StatusID = statusID

	tm.Tasks = append(tm.Tasks, task)
	return nil
}

// GetTasks retrieves all non-deleted tasks
func (tm *TaskManager) GetTasks() []Task {
	var currentTasks []Task
	for _, task := range tm.Tasks {
		if !task.Deleted {
			currentTasks = append(currentTasks, task)
		}
	}
	return currentTasks
}

// UpdateTask updates an existing task
func (tm *TaskManager) UpdateTask(updatedTask Task) error {
	now := time.Now()

	// Validate the task's status
	statusID, err := convertStringToStatusID(updatedTask.StatusString)
	if err != nil {
		return ErrInvalidStatus
	}

	// Find and update the task
	for i, task := range tm.Tasks {
		if task.ID == updatedTask.ID && !task.Deleted {
			updatedTask.StatusID = statusID
			updatedTask.CreatedAt = task.CreatedAt
			updatedTask.UpdatedAt = &now
			tm.Tasks[i] = updatedTask
			return nil
		}
	}

	return ErrTaskNotFound
}

// DeleteTask marks a task as deleted
func (tm *TaskManager) DeleteTask(taskID int) error {
	now := time.Now()

	for i, task := range tm.Tasks {
		if task.ID == taskID && !task.Deleted {
			tm.Tasks[i].Deleted = true
			tm.Tasks[i].DeletedAt = &now
			return nil
		}
	}

	return ErrTaskNotFound
}

func convertStringToStatusID(status string) (Status, error) {
	switch strings.ReplaceAll(status, " ", "") {
	case "NotStarted":
		return NotStarted, nil
	case "Started":
		return Started, nil
	case "Completed":
		return Completed, nil
	default:
		return Unknown, ErrInvalidStatus
	}
}
