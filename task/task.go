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
// FIXME: Rename to avoid starting with the package name
type TaskManager struct {
	Tasks     []Task
	MaxTaskID int
}

// FIXME: Rename to avoid starting with the package name
type TaskRequest struct {
	Action   string
	Task     Task
	TaskID   int
	Response chan<- interface{}
}

// FIXME: Rename to avoid starting with the package name
type TaskActor struct {
	TaskManager  *TaskManager
	RequestsChan chan TaskRequest
}

var (
	// ErrTaskNotFound is returned when a task is not found
	ErrTaskNotFound = errors.New("task not found")
	// ErrInvalidStatus is returned when a task already exists
	ErrInvalidStatus = errors.New("invalid status string")
)

func NewTaskActor(taskManager *TaskManager, requestChanSize int) *TaskActor {
	actor := &TaskActor{
		TaskManager:  taskManager,
		RequestsChan: make(chan TaskRequest, requestChanSize),
	}
	go actor.processLoop()
	return actor
}

func (ta *TaskActor) processLoop() {
	for req := range ta.RequestsChan {
		switch req.Action {
		case CreateRequest:
			err := ta.TaskManager.CreateTask(req.Task)
			req.Response <- err
		case GetRequest:
			tasks := ta.TaskManager.GetTasks()
			req.Response <- tasks
		case UpdateRequest:
			err := ta.TaskManager.UpdateTask(req.Task)
			req.Response <- err
		case DeleteRequest:
			err := ta.TaskManager.DeleteTask(req.TaskID)
			req.Response <- err
		default:
			req.Response <- errors.New("unknown action")
		}
		close(req.Response)
	}
}

// CreateTask adds a new task to the list of tasks
func (tm *TaskManager) CreateTask(task Task) error {
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

	statusID, err := convertStringToStatusID(updatedTask.StatusString)
	if err != nil {
		return ErrInvalidStatus
	}

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
