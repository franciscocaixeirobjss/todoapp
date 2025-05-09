package files

import (
	"encoding/json"
	"io"
	"log/slog"
	"os"
	"todoapp/task"
)

// LoadData initializes the list of tasks and maxTaskID from a JSON file
func LoadData(filePath string, tasks *[]task.Task, maxTaskID *int) error {
	file, err := os.Open(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			*tasks = []task.Task{}
			*maxTaskID = 0
			slog.Info("No existing data file found. Starting with an empty task list.")
			return nil
		}
		return err
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {

		}
	}(file)

	data, err := io.ReadAll(file)
	if err != nil {
		slog.Error("Failed to read file", "error", err)
		return err
	}

	err = json.Unmarshal(data, tasks)
	if err != nil {
		slog.Error("Failed to unmarshal JSON", "error", err)
		return err
	}

	for _, loadedTask := range *tasks {
		if loadedTask.ID > *maxTaskID {
			*maxTaskID = loadedTask.ID
		}
	}

	slog.Info("Loaded tasks", "tasks", *tasks)
	slog.Info("Max task ID", "maxTaskID", *maxTaskID)
	return nil
}

// SaveData saves the tasks to a JSON file
func SaveData(filename string, tasks []task.Task) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {

		}
	}(file)

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	return encoder.Encode(tasks)
}
