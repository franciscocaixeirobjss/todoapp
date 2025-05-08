package main

import (
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"todoapp/api"
	"todoapp/files"
	"todoapp/logging"
	"todoapp/task"
	"todoapp/webserver"
)

func main() {
	opts := &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}

	slogHandler := slog.NewTextHandler(os.Stdout, opts)
	logger := slog.New(slogHandler)
	slog.SetDefault(logger)

	taskManager := &task.TaskManager{}
	err := files.LoadData("todo.json", &taskManager.Tasks, &taskManager.MaxTaskID)
	if err != nil {
		slog.Error("Failed to load data", "error", err)
		return
	}

	taskActor := task.NewTaskActor(taskManager, 10)

	handlers := &api.Handlers{
		TaskActor: taskActor,
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/create", handlers.CreateHandler)
	mux.HandleFunc("/get", handlers.GetHandler)
	mux.HandleFunc("/update", handlers.UpdateHandler)
	mux.HandleFunc("/delete/", handlers.DeleteHandler)

	webserver.ServeStaticPage(mux)
	webserver.ServeDynamicPage(mux, taskManager)

	wrappedMux := logging.TraceIDMiddleware(mux)

	server := &http.Server{
		Addr:    ":8080",
		Handler: wrappedMux,
	}

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	go func() {
		slog.Info("Server starting on :8080")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("Server failed to start", "error", err)
		}
	}()

	<-stop
	slog.Info("Shutting down server...")

	if err := server.Close(); err != nil {
		slog.Error("Server forced to shut down", "error", err)
	}

	if err := files.SaveData("todo.json", taskManager.Tasks); err != nil {
		slog.Error("Failed to save tasks to file", "error", err)
	} else {
		slog.Info("Tasks saved successfully")
	}
}
