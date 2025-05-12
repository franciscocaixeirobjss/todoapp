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
	manager      = Manager{Tasks: []Task{}, MaxTaskID: 0}
	RequestsChan chan Request
)

func InitChannel(requestChanSize int) {
	RequestsChan = make(chan Request, requestChanSize)
	go processLoop()
}

func SetTasks(tasks []Task, maxTaskID int) {
	manager.Tasks = tasks
	manager.MaxTaskID = maxTaskID
}

func GetManagerTasks() ([]Task, int) {
	return manager.Tasks, manager.MaxTaskID
}

func processLoop() {
	for req := range RequestsChan {
		switch req.Action {
		case CreateRequest:
			err := CreateTask(req.Task)
			req.Response <- Response{Tasks: nil, Error: err}
		case GetRequest:
			tasks := GetTasks()
			req.Response <- Response{Tasks: tasks, Error: nil}
		case UpdateRequest:
			err := UpdateTask(req.Task)
			req.Response <- Response{Tasks: nil, Error: err}
		case DeleteRequest:
			err := DeleteTask(req.TaskID)
			req.Response <- Response{Tasks: nil, Error: err}
		default:
			req.Response <- Response{Tasks: nil, Error: errors.New("unknown action")}
		}
		close(req.Response)
	}
}

// CreateTask adds a new task to the list of tasks
func CreateTask(task Task) error {
	now := time.Now()

	statusID, err := convertStringToStatusID(task.StatusString)
	if err != nil {
		return ErrInvalidStatus
	}

	task.ID = manager.MaxTaskID + 1
	manager.MaxTaskID++
	task.CreatedAt = &now
	task.StatusID = statusID

	manager.Tasks = append(manager.Tasks, task)
	return nil
}

// GetTasks retrieves all non-deleted tasks
func GetTasks() []Task {
	var currentTasks []Task
	for _, task := range manager.Tasks {
		if !task.Deleted {
			currentTasks = append(currentTasks, task)
		}
	}
	return currentTasks
}

// UpdateTask updates an existing task
func UpdateTask(updatedTask Task) error {
	for i, task := range manager.Tasks {
		if task.ID == updatedTask.ID {
			now := time.Now()
			manager.Tasks[i].Title = updatedTask.Title
			manager.Tasks[i].Description = updatedTask.Description
			manager.Tasks[i].StatusID, _ = convertStringToStatusID(updatedTask.StatusString)
			manager.Tasks[i].UpdatedAt = &now
			return nil
		}
	}
	return ErrTaskNotFound
}

// DeleteTask marks a task as deleted
func DeleteTask(taskID int) error {
	now := time.Now()

	for i, task := range manager.Tasks {
		if task.ID == taskID && !task.Deleted {
			manager.Tasks[i].Deleted = true
			manager.Tasks[i].DeletedAt = &now
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
