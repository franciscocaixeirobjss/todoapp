package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"sync"
	"testing"
	"todoapp/middleware"
	"todoapp/task"
)

var channelSizeBenchmark = 100

func TestMain(m *testing.M) {
	m.Run()
}

func BenchmarkUpdateHandlerNonActor(b *testing.B) {
	manager := &task.NonActorManager{}

	for i := 0; i < 1; i++ {
		manager.CreateTask(task.Task{
			Title:        "Task " + strconv.Itoa(i),
			Description:  "Description " + strconv.Itoa(i),
			StatusString: "NotStarted",
		})
	}

	updatedTask := task.Task{
		ID:           1,
		Title:        "Updated Task",
		Description:  "Updated Description",
		StatusString: "Started",
	}
	updatedTaskJSON, _ := json.Marshal(updatedTask)

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			req := httptest.NewRequest(http.MethodPut, "/update", bytes.NewReader(updatedTaskJSON))
			req = addUserIDToContext(req, 1)
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			UpdateHandlerWithManager(w, req, manager)

			resp := w.Result()
			if resp.StatusCode != http.StatusOK {
				b.Errorf("Expected status %d, got %d", http.StatusOK, resp.StatusCode)
			}
		}
	})
}

func BenchmarkUpdateHandlerActor_NoParallel(b *testing.B) {
	task.InitChannel(channelSizeBenchmark)

	response := make(chan task.Response, 1)
	task.RequestsChan <- task.Request{
		UserID: 1,
		Action: task.CreateRequest,
		Task: task.Task{
			Title:        "Initial Task",
			Description:  "Initial Description",
			StatusString: "NotStarted",
		},
		Response: response,
	}
	<-response

	updatedTask := task.Task{
		ID:           1,
		Title:        "Updated Task",
		Description:  "Updated Description",
		StatusString: "Started",
	}
	updatedTaskJSON, _ := json.Marshal(updatedTask)

	for i := 0; i < b.N; i++ {
		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodPut, "/update", bytes.NewReader(updatedTaskJSON))
		r = addUserIDToContext(r, 1)

		r.Header.Set("Content-Type", "application/json")
		UpdateHandler(w, r)

		if w.Code != http.StatusOK {
			b.Errorf("wanted status %d, got %d", http.StatusOK, w.Code)
		}

	}
}

func BenchmarkUpdateHandlerActor(b *testing.B) {
	task.InitChannel(channelSizeBenchmark)

	response := make(chan task.Response, 1)
	task.RequestsChan <- task.Request{
		UserID: 1,
		Action: task.CreateRequest,
		Task: task.Task{
			Title:        "Initial Task",
			Description:  "Initial Description",
			StatusString: "NotStarted",
		},
		Response: response,
	}
	<-response

	updatedTask := task.Task{
		ID:           1,
		Title:        "Updated Task",
		Description:  "Updated Description",
		StatusString: "Started",
	}
	updatedTaskJSON, _ := json.Marshal(updatedTask)

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			req := httptest.NewRequest(http.MethodPut, "/update", bytes.NewReader(updatedTaskJSON))
			req = addUserIDToContext(req, 1)
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			UpdateHandler(w, req)

			resp := w.Result()
			if resp.StatusCode != http.StatusOK {
				b.Errorf("Expected status %d, got %d", http.StatusOK, resp.StatusCode)
			}
		}
	})
}

func BenchmarkCreateHandlerNonActor(b *testing.B) {
	manager := &task.NonActorManager{}

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			req := httptest.NewRequest(http.MethodPost, "/create", strings.NewReader(`{"title": "New Task", "description": "Task Description", "status": "NotStarted"}`))
			req = addUserIDToContext(req, 1)
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()

			CreateHandlerWithManager(w, req, manager)

			resp := w.Result()
			if resp.StatusCode != http.StatusCreated {
				b.Errorf("Expected status %d, got %d", http.StatusCreated, resp.StatusCode)
			}
		}
	})
}

func BenchmarkCreateHandlerActor_NoParallel(b *testing.B) {
	task.InitChannel(channelSizeBenchmark)

	for i := 0; i < b.N; i++ {
		w := httptest.NewRecorder()

		r := httptest.NewRequest(http.MethodPost, "/create", strings.NewReader(`{"title": "New Task", "description": "Task Description", "status": "NotStarted"}`))
		r = addUserIDToContext(r, 1)
		r.Header.Set("Content-Type", "application/json")

		CreateHandler(w, r)

		if w.Code != http.StatusCreated {
			b.Errorf("wanted status %d, got %d", http.StatusOK, w.Code)
		}

	}
}

func BenchmarkCreateHandlerActor(b *testing.B) {
	task.InitChannel(channelSizeBenchmark)

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			req := httptest.NewRequest(http.MethodPost, "/create", strings.NewReader(`{"title": "New Task", "description": "Task Description", "status": "NotStarted"}`))
			req = addUserIDToContext(req, 1)
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()

			CreateHandler(w, req)

			resp := w.Result()
			if resp.StatusCode != http.StatusCreated {
				b.Errorf("Expected status %d, got %d", http.StatusCreated, resp.StatusCode)
			}
		}
	})
}

func TestCreateHandler_ServiceUnavailable(t *testing.T) {
	mockRequestsChan := make(chan task.Request, 1)
	task.RequestsChan = mockRequestsChan

	response := make(chan task.Response, 1)
	mockRequestsChan <- task.Request{
		UserID:   1,
		Action:   task.CreateRequest,
		Task:     task.Task{Title: "Mock Task"},
		Response: response,
	}

	req := httptest.NewRequest(http.MethodPost, "/create", strings.NewReader(`{"title": "New Task", "description": "Task", "status": "NotStarted"}`))
	req = addUserIDToContext(req, 1)

	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	CreateHandler(w, req)

	resp := w.Result()
	if resp.StatusCode != http.StatusServiceUnavailable {
		t.Errorf("Expected status %d, got %d", http.StatusServiceUnavailable, resp.StatusCode)
	}
}

func TestCreateHandler_Parallel(t *testing.T) {
	task.InitChannel(100)

	task.SetTasks(createEmptyTaskMap(), createEmptyTaskCountMap())

	// Define test cases
	var tests []struct {
		name           string
		taskID         int
		expectedStatus int
	}

	for i := 0; i < 100; i++ {
		tests = append(tests, struct {
			name           string
			taskID         int
			expectedStatus int
		}{
			name:           fmt.Sprintf("Task %d", i),
			taskID:         i,
			expectedStatus: http.StatusCreated,
		})
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			req := httptest.NewRequest(http.MethodPost, "/create", strings.NewReader(fmt.Sprintf(`{"title": "Task %d", "description": "Description %d", "status": "NotStarted"}`, tt.taskID, tt.taskID)))
			req.Header.Set("Content-Type", "application/json")
			req = addUserIDToContext(req, 1) // Add UserID to context
			w := httptest.NewRecorder()

			CreateHandler(w, req)

			resp := w.Result()
			if resp.StatusCode != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d for task %d", tt.expectedStatus, resp.StatusCode, tt.taskID)
			}
		})
	}
}

func addUserIDToContext(req *http.Request, userID int) *http.Request {
	ctx := req.Context()
	ctx = context.WithValue(ctx, middleware.UserIDKey, userID)
	return req.WithContext(ctx)
}

func TestCreateHandler_Goroutine_Parallel(t *testing.T) {
	numRequests := 100
	task.InitChannel(numRequests)

	task.SetTasks(createEmptyTaskMap(), createEmptyTaskCountMap())

	var wg sync.WaitGroup

	for i := 0; i < numRequests; i++ {
		wg.Add(1)
		go func(taskID int) {
			defer wg.Done()

			req := httptest.NewRequest(http.MethodPost, "/create", strings.NewReader(fmt.Sprintf(`{"title": "Task %d", "description": "Description %d", "status": "NotStarted"}`, taskID, taskID)))
			req.Header.Set("Content-Type", "application/json")
			req = addUserIDToContext(req, 1) // Add UserID to context
			w := httptest.NewRecorder()

			CreateHandler(w, req)

			resp := w.Result()
			if resp.StatusCode != http.StatusCreated {
				t.Errorf("Expected status %d, got %d for task %d", http.StatusCreated, resp.StatusCode, taskID)
			}
		}(i)
	}

	wg.Wait()

	response := make(chan task.Response, 1)
	task.RequestsChan <- task.Request{
		UserID:   1,
		Action:   task.GetRequest,
		Response: response,
	}

	res := <-response
	if len(res.Tasks) != numRequests {
		t.Errorf("Expected %d tasks, but got %d", numRequests, len(res.Tasks))
	}
}

func TestCreateHandler(t *testing.T) {
	task.InitChannel(10)

	tests := []struct {
		name             string
		method           string
		data             string
		expectedStatus   int
		expectedResponse string
	}{
		{
			name:           "valid create request",
			method:         http.MethodPost,
			data:           `{"title": "Test Task", "description": "A valid task", "status": "NotStarted"}`,
			expectedStatus: http.StatusCreated,
		},
		{
			name:           "status bad request - invalid json format",
			method:         http.MethodPost,
			data:           `{"title": "Invalid Task", "status":`,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "method not allowed - get request instead of post",
			method:         http.MethodGet,
			expectedStatus: http.StatusMethodNotAllowed,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			req, _ := http.NewRequest(test.method, "/create", strings.NewReader(test.data))
			req = addUserIDToContext(req, 1)
			rec := httptest.NewRecorder()

			CreateHandler(rec, req)

			if rec.Code != test.expectedStatus {
				t.Errorf("expected status code %d, got %d", test.expectedStatus, rec.Code)
			}
		})
	}
}

func TestGetHandler(t *testing.T) {
	task.InitChannel(10)

	task.SetTasks(map[int][]task.Task{
		1: {
			{ID: 1, Title: "Task 1", Deleted: false},
			{ID: 2, Title: "Task 2", Deleted: true},
		},
	}, map[int]int{
		1: 2,
	})

	tests := []struct {
		name             string
		method           string
		expectedStatus   int
		expectedResponse []task.Task
	}{
		{
			name:           "valid get request",
			method:         http.MethodGet,
			expectedStatus: http.StatusOK,
			expectedResponse: []task.Task{
				{ID: 1, Title: "Task 1", Deleted: false},
			},
		},
		{
			name:           "method not allowed - post request instead of get",
			method:         http.MethodPost,
			expectedStatus: http.StatusMethodNotAllowed,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			req, _ := http.NewRequest(test.method, "/get", nil)
			req = addUserIDToContext(req, 1)

			rec := httptest.NewRecorder()

			GetHandler(rec, req)

			if rec.Code != test.expectedStatus {
				t.Errorf("expected status code %d, got %d", test.expectedStatus, rec.Code)
			}

			if test.expectedStatus == http.StatusOK {
				var tasks []task.Task
				err := json.Unmarshal(rec.Body.Bytes(), &tasks)
				if err != nil {
					t.Fatalf("failed to parse response body: %v", err)
				}

				if len(tasks) != len(test.expectedResponse) {
					t.Errorf("expected %d tasks, got %d", len(test.expectedResponse), len(tasks))
				}
			}
		})
	}
}

func TestTaskActor_Concurrency(t *testing.T) {
	task.InitChannel(100)
	numGoroutines := 100
	var wg sync.WaitGroup

	task.SetTasks(createEmptyTaskMap(), createEmptyTaskCountMap())

	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(taskID int) {
			defer wg.Done()

			response := make(chan task.Response)
			task.RequestsChan <- task.Request{
				UserID: 1,
				Action: task.CreateRequest,
				Task: task.Task{
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

	response := make(chan task.Response)
	task.RequestsChan <- task.Request{
		UserID:   1,
		Action:   task.GetRequest,
		Response: response,
	}

	res := <-response
	if len(res.Tasks) != numGoroutines {
		t.Errorf("Expected %d tasks, but got %d", numGoroutines, len(res.Tasks))
	}
}

func TestTaskActor_ConcurrentUpdate(t *testing.T) {
	task.InitChannel(100)

	numGoroutines := 50
	var wg sync.WaitGroup

	task.SetTasks(createEmptyTaskMap(), createEmptyTaskCountMap())

	response := make(chan task.Response, 1)
	task.RequestsChan <- task.Request{
		UserID: 1,
		Action: task.CreateRequest,
		Task: task.Task{
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

			response := make(chan task.Response, 1)
			task.RequestsChan <- task.Request{
				UserID: 1,
				Action: task.UpdateRequest,
				Task: task.Task{
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

	response = make(chan task.Response, 1)
	task.RequestsChan <- task.Request{
		UserID:   1,
		Action:   task.GetRequest,
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
	if finalTask.StatusID != task.Started {
		t.Errorf("Expected task status to be Started, but got %d", finalTask.StatusID)
	}
	if !strings.HasPrefix(finalTask.Title, "Updated Task") {
		t.Errorf("Expected task title to start with 'Updated Task', but got '%s'", finalTask.Title)
	}
}

func TestUpdateHandler(t *testing.T) {
	task.InitChannel(10)

	task.SetTasks(map[int][]task.Task{
		1: {
			{ID: 1, Title: "Task 1", StatusString: "NotStarted"},
		},
	}, map[int]int{
		1: 1,
	})

	tests := []struct {
		name           string
		method         string
		data           string
		expectedStatus int
	}{
		{
			name:           "valid update request",
			method:         http.MethodPut,
			data:           `{"id": 1, "title": "Updated Task", "status": "Completed"}`,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "not found - non-existing id",
			method:         http.MethodPut,
			data:           `{"id": 999, "title": "Non-existent Task", "status": "Started"}`,
			expectedStatus: http.StatusNotFound,
		},
		{
			name:           "method not allowed - post request instead of put",
			method:         http.MethodPost,
			expectedStatus: http.StatusMethodNotAllowed,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			req, _ := http.NewRequest(test.method, "/update", strings.NewReader(test.data))
			req = addUserIDToContext(req, 1)

			rec := httptest.NewRecorder()

			UpdateHandler(rec, req)

			if rec.Code != test.expectedStatus {
				t.Errorf("Expected status code %d, got %d", test.expectedStatus, rec.Code)
			}
		})
	}
}

func TestDeleteHandler(t *testing.T) {
	task.InitChannel(10)

	task.SetTasks(map[int][]task.Task{
		1: {
			{ID: 1, Title: "Task 1", Deleted: false},
			{ID: 2, Title: "Task 2", Deleted: true},
		},
	}, map[int]int{
		1: 2,
	})

	tests := []struct {
		name           string
		method         string
		url            string
		expectedStatus int
	}{
		{
			name:           "valid delete request",
			method:         http.MethodDelete,
			url:            "/delete/1",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "not found - non-existing id",
			method:         http.MethodDelete,
			url:            "/delete/999",
			expectedStatus: http.StatusNotFound,
		},
		{
			name:           "not found - already deleted task",
			method:         http.MethodDelete,
			url:            "/delete/2",
			expectedStatus: http.StatusNotFound,
		},
		{
			name:           "bad request - invalid task ID",
			method:         http.MethodDelete,
			url:            "/delete/abc",
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "method not allowed - post request instead of delete",
			method:         http.MethodPost,
			url:            "/delete/1",
			expectedStatus: http.StatusMethodNotAllowed,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			req, _ := http.NewRequest(test.method, test.url, nil)
			req = addUserIDToContext(req, 1)

			rec := httptest.NewRecorder()

			DeleteHandler(rec, req)

			if rec.Code != test.expectedStatus {
				t.Errorf("Expected status code %d, got %d", test.expectedStatus, rec.Code)
			}
		})
	}
}

func createEmptyTaskMap() map[int][]task.Task {
	return make(map[int][]task.Task)
}

func createEmptyTaskCountMap() map[int]int {
	return make(map[int]int)
}
