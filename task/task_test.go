package task

import (
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"
)

func TestMain(m *testing.M) {
	// Initialize the TaskManager and RequestsChan before running tests
	InitTaskManager(&Manager{}, 100)
	m.Run()
}

func TestTaskActor_Concurrency(t *testing.T) {
	numGoroutines := 100

	var wg sync.WaitGroup

	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(taskID int) {
			defer wg.Done()

			response := make(chan Response)
			RequestsChan <- Request{
				Action: CreateRequest,
				Task: Task{
					Title:        "Task " + strconv.Itoa(taskID),
					Description:  "Description for Task " + strconv.Itoa(taskID),
					StatusString: "NotStarted",
				},
				Response: response,
			}

			if res := <-response; res.Error != nil {
				t.Errorf("Failed to add task: %v", res.Error)
			}
		}(i)
	}

	wg.Wait()

	response := make(chan Response)
	RequestsChan <- Request{
		Action:   GetRequest,
		Response: response,
	}

	res := <-response
	if len(res.Tasks) != numGoroutines {
		t.Errorf("Expected %d tasks, but got %d", numGoroutines, len(res.Tasks))
	}
}

func TestTaskActor_ConcurrentUpdate(t *testing.T) {
	InitTaskManager(&Manager{}, 100)

	numGoroutines := 50
	var wg sync.WaitGroup

	response := make(chan Response, 1)
	RequestsChan <- Request{
		Action: CreateRequest,
		Task: Task{
			Title:        "Initial Task",
			Description:  "Initial Description",
			StatusString: "NotStarted",
		},
		Response: response,
	}
	if res := <-response; res.Error != nil {
		t.Fatalf("Failed to create task: %v", res.Error)
	}

	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(updateID int) {
			defer wg.Done()

			response := make(chan Response, 1)
			RequestsChan <- Request{
				Action: UpdateRequest,
				Task: Task{
					ID:           1,
					Title:        "Updated Task " + strconv.Itoa(updateID),
					Description:  "Updated Description " + strconv.Itoa(updateID),
					StatusString: "Started",
				},
				Response: response,
			}

			if res := <-response; res.Error != nil {
				t.Errorf("Failed to update task: %v", res.Error)
			}
		}(i)
	}

	wg.Wait()

	response = make(chan Response, 1)
	RequestsChan <- Request{
		Action:   GetRequest,
		Response: response,
	}

	res := <-response
	if len(res.Tasks) != 1 {
		t.Errorf("Expected 1 task, but got %d", len(res.Tasks))
	}

	finalTask := res.Tasks[0]
	if finalTask.ID != 1 {
		t.Errorf("Expected task ID to be 1, but got %d", finalTask.ID)
	}
	if finalTask.StatusID != Started {
		t.Errorf("Expected task status to be Started, but got %d", finalTask.StatusID)
	}
	if !strings.HasPrefix(finalTask.Title, "Updated Task") {
		t.Errorf("Expected task title to start with 'Updated Task', but got '%s'", finalTask.Title)
	}
}

func TestCreateTask(t *testing.T) {
	tm := &Manager{
		Tasks:     []Task{},
		MaxTaskID: 0,
	}

	task := Task{
		Title:        "Test Task",
		Description:  "This is a test task",
		StatusString: "Not Started",
	}

	err := tm.CreateTask(task)
	if err != nil {
		t.Fatalf("AddTask failed: %v", err)
	}

	if len(tm.Tasks) != 1 {
		t.Fatalf("Expected 1 task, got %d", len(tm.Tasks))
	}

	if tm.Tasks[0].Title != "Test Task" {
		t.Errorf("Expected task title to be 'Test Task', got '%s'", tm.Tasks[0].Title)
	}

	if tm.Tasks[0].StatusID != NotStarted {
		t.Errorf("Expected task status to be NotStarted, got %d", tm.Tasks[0].StatusID)
	}
}

func TestGetTasks(t *testing.T) {
	tm := &Manager{
		Tasks: []Task{
			{ID: 1, Title: "Task 1", Deleted: false},
			{ID: 2, Title: "Task 2", Deleted: true},
			{ID: 3, Title: "Task 3", Deleted: false},
		},
	}

	tasks := tm.GetTasks()
	if len(tasks) != 2 {
		t.Fatalf("Expected 2 tasks, got %d", len(tasks))
	}

	if tasks[0].ID != 1 || tasks[1].ID != 3 {
		t.Errorf("Unexpected tasks returned: %+v", tasks)
	}
}

func TestUpdateTask(t *testing.T) {
	now := time.Now()
	tm := &Manager{
		Tasks: []Task{
			{ID: 1, Title: "Task 1", StatusString: "NotStarted", CreatedAt: &now},
		},
	}

	updatedTask := Task{
		ID:           1,
		Title:        "Updated Task 1",
		StatusString: "Completed",
	}

	err := tm.UpdateTask(updatedTask)
	if err != nil {
		t.Fatalf("updateTask failed: %v", err)
	}

	if tm.Tasks[0].Title != "Updated Task 1" {
		t.Errorf("expected task title to be 'Updated Task 1', got '%s'", tm.Tasks[0].Title)
	}

	if tm.Tasks[0].StatusID != Completed {
		t.Errorf("expected task status to be Completed, got %d", tm.Tasks[0].StatusID)
	}

	if tm.Tasks[0].UpdatedAt == nil {
		t.Errorf("expected UpdatedAt to be set, but it was nil")
	}
}

func TestDeleteTask(t *testing.T) {
	now := time.Now()
	tm := &Manager{
		Tasks: []Task{
			{ID: 1, Title: "Task 1", Deleted: false, CreatedAt: &now},
		},
	}

	err := tm.DeleteTask(1)
	if err != nil {
		t.Fatalf("deleteTask failed: %v", err)
	}

	if !tm.Tasks[0].Deleted {
		t.Errorf("expected task to be marked as deleted, but it was not")
	}

	if tm.Tasks[0].DeletedAt == nil {
		t.Errorf("expected DeletedAt to be set, but it was nil")
	}
}

func TestConvertStringToStatusID(t *testing.T) {
	tests := []struct {
		input    string
		expected Status
	}{
		{"NotStarted", NotStarted},
		{" Not Started ", NotStarted},
		{"Started", Started},
		{"Completed", Completed},
		{"Invalid Status", Unknown},
		{"InvalidStatus", Unknown},
	}

	for _, test := range tests {
		t.Run(test.input, func(t *testing.T) {
			result, _ := convertStringToStatusID(test.input)
			if result != test.expected {
				t.Errorf("convertStringToStatusID(%q) = %d; want %d", test.input, result, test.expected)
			}
		})
	}
}
