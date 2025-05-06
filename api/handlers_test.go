package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"todoapp/task"
)

func TestCreateHandler(t *testing.T) {
	mockTaskManager := &task.TaskManager{}
	handlers := &Handlers{TaskManager: mockTaskManager}

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
			name:           "status bad request - innvalid json format",
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

			handlers.CreateHandler(rec, req)

			if rec.Code != test.expectedStatus {
				t.Errorf("expected status code %d, got %d", test.expectedStatus, rec.Code)
			}
		})
	}
}

func TestGetHandler(t *testing.T) {
	mockTaskManager := &task.TaskManager{
		Tasks: []task.Task{
			{ID: 1, Title: "Task 1", Deleted: false},
			{ID: 2, Title: "Task 2", Deleted: true},
		},
	}
	handlers := &Handlers{TaskManager: mockTaskManager}

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

			handlers.GetHandler(rec, req)

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
	mockTaskManager := &task.TaskManager{
		Tasks: []task.Task{
			{ID: 1, Title: "Task 1", StatusString: "NotStarted"},
		},
	}
	handlers := &Handlers{TaskManager: mockTaskManager}

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

			handlers.UpdateHandler(rec, req)

			if rec.Code != test.expectedStatus {
				t.Errorf("Expected status code %d, got %d", test.expectedStatus, rec.Code)
			}
		})
	}
}

func TestDeleteHandler(t *testing.T) {
	mockTaskManager := &task.TaskManager{
		Tasks: []task.Task{
			{ID: 1, Title: "Task 1", Deleted: false},
			{ID: 2, Title: "Task 2", Deleted: true},
		},
	}
	handlers := &Handlers{TaskManager: mockTaskManager}

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

			handlers.DeleteHandler(rec, req)

			if rec.Code != test.expectedStatus {
				t.Errorf("Expected status code %d, got %d", test.expectedStatus, rec.Code)
			}
		})
	}
}
