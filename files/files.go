package files

import (
	"encoding/json"
	"log/slog"
	"os"
	"todoapp/task"
)

type dataFormat struct {
	Tasks      map[int][]task.Task `json:"tasks"`
	MaxTaskIDs map[int]int         `json:"maxTaskIDs"`
}

// LoadData initializes the list of tasks and maxTaskID from a JSON file
func LoadData(filePath string, tasks *map[int][]task.Task, maxTaskIDs *map[int]int) error {
	file, err := os.Open(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			*tasks = make(map[int][]task.Task)
			*maxTaskIDs = make(map[int]int)
			slog.Info("No existing data file found. Starting with an empty task list.")
			return nil
		}
		return err
	}

	defer file.Close()

	decoder := json.NewDecoder(file)
	data := dataFormat{}

	if err := decoder.Decode(&data); err != nil {
		return err
	}

	*tasks = data.Tasks
	*maxTaskIDs = data.MaxTaskIDs
	return nil
}

// SaveData saves the tasks to a JSON file
func SaveData(filename string, tasks map[int][]task.Task, maxTaskIDs map[int]int) error {
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
	data := dataFormat{
		Tasks:      tasks,
		MaxTaskIDs: maxTaskIDs,
	}

	if err := encoder.Encode(&data); err != nil {
		return err
	}

	return nil
}
