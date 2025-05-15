package task

import (
	"testing"
	"time"
)

func TestMain(m *testing.M) {
	m.Run()
}

func BenchmarkUpdateActorPattern(b *testing.B) {
	InitChannel(1000)

	response := make(chan Response)
	RequestsChan <- Request{
		Action: CreateRequest,
		Task: Task{
			Title:        "Initial Task",
			Description:  "Initial Description",
			StatusString: "NotStarted",
		},
		Response: response,
	}
	<-response

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			response := make(chan Response)
			RequestsChan <- Request{
				Action: UpdateRequest,
				Task: Task{
					ID:           1,
					Title:        "Updated Task",
					Description:  "Updated Description",
					StatusString: "Started",
				},
				Response: response,
			}
			<-response
		}
	})
}

func BenchmarkUpdateNonActorPattern(b *testing.B) {
	manager := &NonActorManager{}

	manager.CreateTask(Task{
		Title:        "Initial Task",
		Description:  "Initial Description",
		StatusString: "NotStarted",
	})

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			manager.UpdateTask(Task{
				ID:           1,
				Title:        "Updated Task",
				Description:  "Updated Description",
				StatusString: "Started",
			})
		}
	})
}

func BenchmarkCreateActorPattern(b *testing.B) {
	InitChannel(1000)

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			response := make(chan Response)
			RequestsChan <- Request{
				Action: CreateRequest,
				Task: Task{
					Title:        "Benchmark Task",
					Description:  "This is a benchmark task",
					StatusString: "NotStarted",
				},
				Response: response,
			}
			<-response
		}
	})
}

// Benchmark for the non-actor pattern
func BenchmarkCreateNonActorPattern(b *testing.B) {
	manager := &NonActorManager{}

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			manager.CreateTask(Task{
				Title:        "Benchmark Task",
				Description:  "This is a benchmark task",
				StatusString: "NotStarted",
			})
		}
	})
}

func TestCreateTask(t *testing.T) {
	SetTasks(map[int][]Task{}, map[int]int{})

	task := Task{
		Title:        "Test Task",
		Description:  "This is a test task",
		StatusString: "Not Started",
	}

	err := CreateTask(1, task)
	if err != nil {
		t.Fatalf("CreateTask failed: %v", err)
	}

	tasks, _ := GetManagerTasks()

	if len(tasks[1]) != 1 {
		t.Fatalf("Expected 1 task, got %d", len(tasks[1]))
	}

	if tasks[1][0].Title != "Test Task" {
		t.Errorf("Expected task title to be 'Test Task', got '%s'", tasks[1][0].Title)
	}

	if tasks[1][0].StatusID != NotStarted {
		t.Errorf("Expected task status to be NotStarted, got %d", tasks[1][0].StatusID)
	}
}

func TestGetTasks(t *testing.T) {
	taskToSave := map[int][]Task{
		1: {
			{ID: 1, Title: "Task 1", Deleted: false},
			{ID: 2, Title: "Task 2", Deleted: true},
			{ID: 3, Title: "Task 3", Deleted: false},
		},
	}
	SetTasks(taskToSave, map[int]int{})

	tasks := GetTasks(1)
	if len(tasks) != 2 {
		t.Fatalf("Expected 2 tasks, got %d", len(tasks))
	}

	if tasks[0].ID != 1 || tasks[1].ID != 3 {
		t.Errorf("Unexpected tasks returned: %+v", tasks)
	}
}

func TestUpdateTask(t *testing.T) {
	now := time.Now()
	tasksToSave := map[int][]Task{
		1: {
			{ID: 1, Title: "Task 1", StatusString: "NotStarted", CreatedAt: &now},
		},
	}
	SetTasks(tasksToSave, map[int]int{1: 1})

	updatedTask := Task{
		ID:           1,
		Title:        "Updated Task 1",
		StatusString: "Completed",
	}

	err := UpdateTask(1, updatedTask)
	if err != nil {
		t.Fatalf("UpdateTask failed: %v", err)
	}

	tasks, _ := GetManagerTasks()

	if tasks[1][0].Title != "Updated Task 1" {
		t.Errorf("Expected task title to be 'Updated Task 1', got '%s'", tasks[1][0].Title)
	}

	if tasks[1][0].StatusID != Completed {
		t.Errorf("Expected task status to be Completed, got %d", tasks[1][0].StatusID)
	}

	if tasks[1][0].UpdatedAt == nil {
		t.Errorf("Expected UpdatedAt to be set, but it was nil")
	}
}

func TestDeleteTask(t *testing.T) {
	now := time.Now()
	taskToDelete := map[int][]Task{
		1: {
			{ID: 1, Title: "Task 1", Deleted: false, CreatedAt: &now},
		},
	}
	SetTasks(taskToDelete, map[int]int{1: 1})

	err := DeleteTask(1, 1)
	if err != nil {
		t.Fatalf("DeleteTask failed: %v", err)
	}

	tasks, _ := GetManagerTasks()

	if !tasks[1][0].Deleted {
		t.Errorf("Expected task to be marked as deleted, but it was not")
	}

	if tasks[1][0].DeletedAt == nil {
		t.Errorf("Expected DeletedAt to be set, but it was nil")
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
