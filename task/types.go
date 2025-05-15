package task

import (
	"time"
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

// Manager struct to manage tasks and their state
type Manager struct {
	Tasks      map[int][]Task
	MaxTaskIDs map[int]int
}

type Response struct {
	Tasks []Task
	Error error
}

type Request struct {
	Action   string
	UserID   int
	Task     Task
	TaskID   int
	Response chan<- Response
}
