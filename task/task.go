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

const (
	GetRequest    = "get"
	CreateRequest = "create"
	UpdateRequest = "update"
	DeleteRequest = "delete"
)

var (
	// ErrTaskNotFound is returned when a task is not found
	ErrTaskNotFound = errors.New("task not found")
	// ErrInvalidStatus is returned when a task already exists
	ErrInvalidStatus = errors.New("invalid status string")
)

var (
	taskManager  *Manager
	RequestsChan chan Request
)

// InitTaskManager initializes the TaskManager and RequestsChan
func InitTaskManager(tm *Manager, requestChanSize int) {
	taskManager = tm
	RequestsChan = make(chan Request, requestChanSize)
	go processLoop()
}

func processLoop() {
	for req := range RequestsChan {
		switch req.Action {
		case CreateRequest:
			err := taskManager.CreateTask(req.Task)
			req.Response <- Response{Tasks: nil, Error: err}
		case GetRequest:
			tasks := taskManager.GetTasks()
			req.Response <- Response{Tasks: tasks, Error: nil}
		case UpdateRequest:
			err := taskManager.UpdateTask(req.Task)
			req.Response <- Response{Tasks: nil, Error: err}
		case DeleteRequest:
			err := taskManager.DeleteTask(req.TaskID)
			req.Response <- Response{Tasks: nil, Error: err}
		default:
			req.Response <- Response{Tasks: nil, Error: errors.New("unknown action")}
		}
		close(req.Response)
	}
}

// CreateTask adds a new task to the list of tasks
func (tm *Manager) CreateTask(task Task) error {
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
func (tm *Manager) GetTasks() []Task {
	var currentTasks []Task
	for _, task := range tm.Tasks {
		if !task.Deleted {
			currentTasks = append(currentTasks, task)
		}
	}
	return currentTasks
}

// UpdateTask updates an existing task
func (tm *Manager) UpdateTask(updatedTask Task) error {
	for i, task := range tm.Tasks {
		if task.ID == updatedTask.ID {
			now := time.Now()
			tm.Tasks[i].Title = updatedTask.Title
			tm.Tasks[i].Description = updatedTask.Description
			tm.Tasks[i].StatusID, _ = convertStringToStatusID(updatedTask.StatusString)
			tm.Tasks[i].UpdatedAt = &now
			return nil
		}
	}
	return ErrTaskNotFound
}

// DeleteTask marks a task as deleted
func (tm *Manager) DeleteTask(taskID int) error {
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
