package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"todoapp/task"
)

type Handlers struct {
	TaskActor *task.TaskActor
}

func (h *Handlers) CreateHandler(w http.ResponseWriter, r *http.Request) {
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

	response := make(chan interface{}, 1)
	request := task.TaskRequest{
		Action:   task.CreateRequest,
		Task:     taskToBeCreated,
		Response: response,
	}

	select {
	case h.TaskActor.RequestsChan <- request:
		if err := <-response; err != nil {
			http.Error(w, err.(error).Error(), http.StatusBadRequest)
			return
		}
		w.WriteHeader(http.StatusCreated)
	default:
		http.Error(w, "Service unavailable. Please try again later.", http.StatusServiceUnavailable)
	}
}

func (h *Handlers) GetHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Invalid HTTP method.", http.StatusMethodNotAllowed)
		return
	}

	response := make(chan interface{}, 1)
	request := task.TaskRequest{
		Action:   task.GetRequest,
		Response: response,
	}

	select {
	case h.TaskActor.RequestsChan <- request:
		tasks := <-response
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(tasks)
	default:
		http.Error(w, "Service unavailable. Please try again later.", http.StatusServiceUnavailable)
	}
}

func (h *Handlers) UpdateHandler(w http.ResponseWriter, r *http.Request) {
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

	response := make(chan interface{}, 1)
	request := task.TaskRequest{
		Action:   task.UpdateRequest,
		Task:     taskToBeUpdated,
		Response: response,
	}

	select {
	case h.TaskActor.RequestsChan <- request:
		if err := <-response; err != nil {
			if err == task.ErrTaskNotFound {
				http.Error(w, err.(error).Error(), http.StatusNotFound)
				return
			}
			http.Error(w, err.(error).Error(), http.StatusBadRequest)
			return
		}
		w.WriteHeader(http.StatusOK)
	default:
		http.Error(w, "Service unavailable. Please try again later.", http.StatusServiceUnavailable)
	}
}

func (h *Handlers) DeleteHandler(w http.ResponseWriter, r *http.Request) {
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

	response := make(chan interface{}, 1)
	request := task.TaskRequest{
		Action:   task.DeleteRequest,
		TaskID:   taskID,
		Response: response,
	}

	select {
	case h.TaskActor.RequestsChan <- request:
		if err := <-response; err != nil {
			if err == task.ErrTaskNotFound {
				http.Error(w, err.(error).Error(), http.StatusNotFound)
				return
			}
			http.Error(w, err.(error).Error(), http.StatusBadRequest)
			return
		}
		w.WriteHeader(http.StatusOK)
	default:
		http.Error(w, "Service unavailable. Please try again later.", http.StatusServiceUnavailable)
	}
}
