package task

import (
	"sync"
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
	Tasks     []Task
	MaxTaskID int
}

type Response struct {
	Tasks []Task
	Error error
}

type Request struct {
	Action   string
	Task     Task
	TaskID   int
	Response chan<- Response
}

// Non-actor implementation using a shared Manager and sync.Mutex
type NonActorManager struct {
	mu        sync.Mutex
	tasks     []Task
	maxTaskID int
}

func (m *NonActorManager) CreateTask(task Task) {
	m.mu.Lock()
	defer m.mu.Unlock()

	task.ID = m.maxTaskID + 1
	m.maxTaskID++
	task.CreatedAt = timePtr(time.Now())
	m.tasks = append(m.tasks, task)
}

func (m *NonActorManager) GetTasks() []Task {
	m.mu.Lock()
	defer m.mu.Unlock()

	var currentTasks []Task
	for _, task := range m.tasks {
		if !task.Deleted {
			currentTasks = append(currentTasks, task)
		}
	}
	return currentTasks
}

// Helper function to create a time pointer
func timePtr(t time.Time) *time.Time {
	return &t
}
