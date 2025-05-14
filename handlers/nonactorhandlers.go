package handlers

import (
	"encoding/json"
	"net/http"
	"todoapp/task"
)

func CreateHandlerWithManager(w http.ResponseWriter, r *http.Request, manager *task.NonActorManager) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid HTTP method.", http.StatusMethodNotAllowed)
		return
	}

	var newTask task.Task
	err := json.NewDecoder(r.Body).Decode(&newTask)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	manager.CreateTask(newTask)

	w.WriteHeader(http.StatusCreated)
}

func UpdateHandlerWithManager(w http.ResponseWriter, r *http.Request, manager *task.NonActorManager) {
	if r.Method != http.MethodPut {
		http.Error(w, "Invalid HTTP method.", http.StatusMethodNotAllowed)
		return
	}

	var updatedTask task.Task
	err := json.NewDecoder(r.Body).Decode(&updatedTask)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	manager.UpdateTask(updatedTask)

	w.WriteHeader(http.StatusOK)
}
