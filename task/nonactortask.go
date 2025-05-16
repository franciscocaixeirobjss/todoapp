// Non-actor implementation of a task manager
// This code is a non-actor implementation of a task manager using a shared Manager and sync.Mutex.
// It provides methods to create, update, delete, and retrieve tasks.
package task

import (
	"sync"
	"time"
)

// Non-actor implementation using a shared Manager and sync.Mutex
type NonActorManager struct {
	mu        sync.Mutex
	Tasks     []Task
	MaxTaskID int
}

// CreateTasks creates a new task and assigns it a unique ID
func (m *NonActorManager) CreateTask(task Task) {
	m.mu.Lock()
	defer m.mu.Unlock()

	task.ID = m.MaxTaskID + 1
	m.MaxTaskID++
	task.CreatedAt = timePtr(time.Now())
	m.Tasks = append(m.Tasks, task)
}

// GetTasks retrieves all non-deleted tasks
func (m *NonActorManager) GetTasks() []Task {
	m.mu.Lock()
	defer m.mu.Unlock()

	var currentTasks []Task
	for _, task := range m.Tasks {
		if !task.Deleted {
			currentTasks = append(currentTasks, task)
		}
	}
	return currentTasks
}

// UpdateTask updates an existing task
func (m *NonActorManager) UpdateTask(updatedTask Task) {
	m.mu.Lock()
	defer m.mu.Unlock()

	for i, task := range m.Tasks {
		if task.ID == updatedTask.ID {
			now := time.Now()
			m.Tasks[i].Title = updatedTask.Title
			m.Tasks[i].Description = updatedTask.Description
			m.Tasks[i].StatusID, _ = convertStringToStatusID(updatedTask.StatusString)
			m.Tasks[i].UpdatedAt = &now
			return
		}
	}
}

// DeleteTask marks a task as deleted
func (m *NonActorManager) DeleteTask(taskID int) {
	m.mu.Lock()
	defer m.mu.Unlock()

	for i, task := range m.Tasks {
		if task.ID == taskID && !task.Deleted {
			now := time.Now()
			m.Tasks[i].Deleted = true
			m.Tasks[i].DeletedAt = &now
			return
		}
	}
}

// Helper function to create a time pointer
func timePtr(t time.Time) *time.Time {
	return &t
}
