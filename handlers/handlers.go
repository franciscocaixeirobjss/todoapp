package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"todoapp/task"
)

// CreateHandler handles task creation
func CreateHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid HTTP method.", http.StatusMethodNotAllowed)
		return
	}

	var taskToBeCreated task.Task
	err := json.NewDecoder(r.Body).Decode(&taskToBeCreated)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	response := make(chan task.Response, 1)
	request := task.Request{
		Action:   task.CreateRequest,
		Task:     taskToBeCreated,
		Response: response,
	}
	select {
	case task.RequestsChan <- request:
		res := <-response
		if res.Error != nil {
			http.Error(w, res.Error.Error(), http.StatusBadRequest)
			return
		}
		w.WriteHeader(http.StatusCreated)
	default:
		http.Error(w, "Service unavailable. Please try again later.", http.StatusServiceUnavailable)
	}
}

// GetHandler handles retrieving tasks
func GetHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Invalid HTTP method.", http.StatusMethodNotAllowed)
		return
	}

	response := make(chan task.Response, 1)
	request := task.Request{
		Action:   task.GetRequest,
		Response: response,
	}

	select {
	case task.RequestsChan <- request:
		res := <-response
		if res.Error != nil {
			http.Error(w, res.Error.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		err := json.NewEncoder(w).Encode(res.Tasks)
		if err != nil {
			http.Error(w, res.Error.Error(), http.StatusInternalServerError)
			return
		}
	default:
		http.Error(w, "Service unavailable. Please try again later.", http.StatusServiceUnavailable)
	}
}

// UpdateHandler handles task updates
func UpdateHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, "Invalid HTTP method.", http.StatusMethodNotAllowed)
		return
	}

	var taskToBeUpdated task.Task
	err := json.NewDecoder(r.Body).Decode(&taskToBeUpdated)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	response := make(chan task.Response, 1)
	request := task.Request{
		Action:   task.UpdateRequest,
		Task:     taskToBeUpdated,
		Response: response,
	}

	select {
	case task.RequestsChan <- request:
		res := <-response
		if res.Error != nil {
			if errors.Is(res.Error, task.ErrTaskNotFound) {
				http.Error(w, res.Error.Error(), http.StatusNotFound)
				return
			}
			http.Error(w, res.Error.Error(), http.StatusBadRequest)
			return
		}
		w.WriteHeader(http.StatusOK)
	default:
		http.Error(w, "Service unavailable. Please try again later.", http.StatusServiceUnavailable)
	}
}

// DeleteHandler handles task deletion
func DeleteHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Invalid HTTP method.", http.StatusMethodNotAllowed)
		return
	}

	path := r.URL.Path
	taskID, err := strconv.Atoi(path[len("/delete/"):])
	if err != nil {
		http.Error(w, "Invalid task ID", http.StatusBadRequest)
		return
	}

	response := make(chan task.Response, 1)
	request := task.Request{
		Action:   task.DeleteRequest,
		TaskID:   taskID,
		Response: response,
	}

	select {
	case task.RequestsChan <- request:
		res := <-response
		if res.Error != nil {
			if errors.Is(res.Error, task.ErrTaskNotFound) {
				http.Error(w, res.Error.Error(), http.StatusNotFound)
				return
			}
			http.Error(w, res.Error.Error(), http.StatusBadRequest)
			return
		}
		w.WriteHeader(http.StatusOK)
	default:
		http.Error(w, "Service unavailable. Please try again later.", http.StatusServiceUnavailable)
	}
}
