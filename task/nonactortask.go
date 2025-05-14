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

func (m *NonActorManager) CreateTask(task Task) {
	m.mu.Lock()
	defer m.mu.Unlock()

	task.ID = m.MaxTaskID + 1
	m.MaxTaskID++
	task.CreatedAt = timePtr(time.Now())
	m.Tasks = append(m.Tasks, task)
}

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
