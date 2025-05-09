package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
	"todoapp/task"
)

func TestMain(m *testing.M) {
	m.Run()
}

func TestCreateHandler_ServiceUnavailable(t *testing.T) {
	mockRequestsChan := make(chan task.Request, 1)
	task.RequestsChan = mockRequestsChan

	response := make(chan task.Response, 1)
	mockRequestsChan <- task.Request{
		Action:   task.CreateRequest,
		Task:     task.Task{Title: "Mock Task"},
		Response: response,
	}

	req := httptest.NewRequest(http.MethodPost, "/create", strings.NewReader(`{"title": "New Task", "description": "Task", "status": "NotStarted"}`))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	CreateHandler(w, req)

	resp := w.Result()
	if resp.StatusCode != http.StatusServiceUnavailable {
		t.Errorf("Expected status %d, got %d", http.StatusServiceUnavailable, resp.StatusCode)
	}
}

func TestCreateHandler_Parallel(t *testing.T) {
	task.InitTaskManager(&task.Manager{}, 100)

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
			w := httptest.NewRecorder()

			CreateHandler(w, req)

			resp := w.Result()
			if resp.StatusCode != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d for task %d", tt.expectedStatus, resp.StatusCode, tt.taskID)
			}
		})
	}
}

func TestCreateHandler_Goroutine_Parallel(t *testing.T) {
	numRequests := 100
	task.InitTaskManager(&task.Manager{}, numRequests)

	var wg sync.WaitGroup

	for i := 0; i < numRequests; i++ {
		wg.Add(1)
		go func(taskID int) {
			defer wg.Done()

			req := httptest.NewRequest(http.MethodPost, "/create", strings.NewReader(fmt.Sprintf(`{"title": "Task %d", "description": "Description %d", "status": "NotStarted"}`, taskID, taskID)))
			req.Header.Set("Content-Type", "application/json")
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
		Action:   task.GetRequest,
		Response: response,
	}

	res := <-response
	if len(res.Tasks) != numRequests {
		t.Errorf("Expected %d tasks, but got %d", numRequests, len(res.Tasks))
	}
	// tasks := (*task.TaskManager).GetTasks()
	// if len(tasks) != numRequests {
	// 	t.Errorf("Expected %d tasks, but got %d", numRequests, len(tasks))
	// }
}

func TestCreateHandler(t *testing.T) {
	task.InitTaskManager(&task.Manager{}, 10)

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
			rec := httptest.NewRecorder()

			CreateHandler(rec, req)

			if rec.Code != test.expectedStatus {
				t.Errorf("expected status code %d, got %d", test.expectedStatus, rec.Code)
			}
		})
	}
}

func TestGetHandler(t *testing.T) {
	task.InitTaskManager(&task.Manager{
		Tasks: []task.Task{
			{ID: 1, Title: "Task 1", Deleted: false},
			{ID: 2, Title: "Task 2", Deleted: true},
		},
	}, 10)

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

func TestUpdateHandler(t *testing.T) {
	task.InitTaskManager(&task.Manager{
		Tasks: []task.Task{
			{ID: 1, Title: "Task 1", StatusString: "NotStarted"},
		},
	}, 10)

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
			rec := httptest.NewRecorder()

			UpdateHandler(rec, req)

			if rec.Code != test.expectedStatus {
				t.Errorf("Expected status code %d, got %d", test.expectedStatus, rec.Code)
			}
		})
	}
}

func TestDeleteHandler(t *testing.T) {
	task.InitTaskManager(&task.Manager{
		Tasks: []task.Task{
			{ID: 1, Title: "Task 1", Deleted: false},
			{ID: 2, Title: "Task 2", Deleted: true},
		},
	}, 10)

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
			rec := httptest.NewRecorder()

			DeleteHandler(rec, req)

			if rec.Code != test.expectedStatus {
				t.Errorf("Expected status code %d, got %d", test.expectedStatus, rec.Code)
			}
		})
	}
}
