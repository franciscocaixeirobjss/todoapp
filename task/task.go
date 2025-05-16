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
	manager = Manager{
		Tasks:      make(map[int][]Task),
		MaxTaskIDs: make(map[int]int)}
	RequestsChan chan Request
)

func GetManagerTasks() (map[int][]Task, map[int]int) {
	return manager.Tasks, manager.MaxTaskIDs
}

func processLoop() {
	for req := range RequestsChan {
		switch req.Action {
		case CreateRequest:
			err := CreateTask(req.UserID, req.Task)
			req.Response <- Response{Tasks: nil, Error: err}
		case GetRequest:
			tasks := GetTasks(req.UserID)
			req.Response <- Response{Tasks: tasks, Error: nil}
		case UpdateRequest:
			err := UpdateTask(req.UserID, req.Task)
			req.Response <- Response{Tasks: nil, Error: err}
		case DeleteRequest:
			err := DeleteTask(req.UserID, req.TaskID)
			req.Response <- Response{Tasks: nil, Error: err}
		default:
			req.Response <- Response{Tasks: nil, Error: errors.New("unknown action")}
		}
		close(req.Response)
	}
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

// InitChannel initializes the channel for task requests
// requestChanSize is the size of the channel for task requests
func InitChannel(requestChanSize int) {
	RequestsChan = make(chan Request, requestChanSize)
	go processLoop()
}

// SetTasks sets the tasks and max task IDs for the manager
func SetTasks(tasks map[int][]Task, maxTaskIDs map[int]int) {
	manager.Tasks = tasks
	manager.MaxTaskIDs = maxTaskIDs
}

// CreateTask adds a new task to the list of tasks
func CreateTask(userID int, task Task) error {
	now := time.Now()

	statusID, err := convertStringToStatusID(task.StatusString)
	if err != nil {
		return ErrInvalidStatus
	}

	manager.MaxTaskIDs[userID]++
	task.ID = manager.MaxTaskIDs[userID]
	task.CreatedAt = &now
	task.StatusID = statusID

	manager.Tasks[userID] = append(manager.Tasks[userID], task)
	return nil
}

// GetTasks retrieves all non-deleted tasks
func GetTasks(userID int) []Task {
	var currentTasks []Task
	for _, task := range manager.Tasks[userID] {
		if !task.Deleted {
			currentTasks = append(currentTasks, task)
		}
	}
	return currentTasks
}

// UpdateTask updates an existing task
func UpdateTask(userID int, updatedTask Task) error {
	for i, task := range manager.Tasks[userID] {
		if task.ID == updatedTask.ID {
			statusID, err := convertStringToStatusID(updatedTask.StatusString)
			if err != nil {
				return err
			}

			now := time.Now()
			manager.Tasks[userID][i].Title = updatedTask.Title
			manager.Tasks[userID][i].Description = updatedTask.Description
			manager.Tasks[userID][i].StatusID = statusID
			manager.Tasks[userID][i].StatusString = strings.ReplaceAll(updatedTask.StatusString, " ", "")
			manager.Tasks[userID][i].UpdatedAt = &now
			return nil
		}
	}
	return ErrTaskNotFound
}

// DeleteTask marks a task as deleted
func DeleteTask(userID int, taskID int) error {
	now := time.Now()

	for i, task := range manager.Tasks[userID] {
		if task.ID == taskID && !task.Deleted {
			manager.Tasks[userID][i].Deleted = true
			manager.Tasks[userID][i].DeletedAt = &now
			return nil
		}
	}

	return ErrTaskNotFound
}
