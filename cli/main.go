package cli

import (
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"todoapp/files"
	"todoapp/task"
)

func main() {
	opts := &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}
	slogHandler := slog.NewTextHandler(os.Stdout, opts)
	logger := slog.New(slogHandler)
	slog.SetDefault(logger)

	var tasks []task.Task
	var maxTaskID int
	err := files.LoadData("todo.json", &tasks, &maxTaskID)
	if err != nil {
		slog.Error("Failed to load data", "error", err)
		os.Exit(1)
	}

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	defer saveTasks()

	go RunCLI()

	<-stop
	fmt.Println("Shutting down gracefully...")
}

func saveTasks() {
	tasks, _ := task.GetManagerTasks()
	if err := files.SaveData("todo.json", tasks); err != nil {
		slog.Error("Failed to save tasks to file", "error", err)
	} else {
		slog.Info("Tasks saved successfully")
	}
}
