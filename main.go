package main

import (
	"errors"
	"flag"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"todoapp/files"
	"todoapp/handlers"
	"todoapp/middleware"
	"todoapp/task"
	"todoapp/webserver"
)

func main() {
	requestChanSize := flag.Int("requestChanSize", 10, "Size of the request channel for TaskActor")
	flag.Parse()

	// TODO: add flag to set the logging level
	// TODO: add flag to set the json file name

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

	defer func() {
		if err := files.SaveData("todo.json", taskManager.Tasks); err != nil {
			slog.Error("Failed to save tasks to file", "error", err)
		} else {
			slog.Info("Tasks saved successfully")
		}
	}()

	taskActor := task.NewTaskActor(taskManager, *requestChanSize)
	handlers := &handlers.Handlers{
		TaskActor: taskActor,
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/create", handlers.CreateHandler)
	mux.HandleFunc("/get", handlers.GetHandler)
	mux.HandleFunc("/update", handlers.UpdateHandler)
	mux.HandleFunc("/delete/", handlers.DeleteHandler)

	// FIXME: Should this be moved somewhere else?
	webserver.ServeStaticPage(mux)
	webserver.ServeDynamicPage(mux, taskManager)

	wrappedMux := middleware.TraceIDMiddleware(mux)

	server := &http.Server{
		Addr:    ":8080",
		Handler: wrappedMux,
	}

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	go func() {
		slog.Info("Server starting on :8080")
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			slog.Error("Server failed to start", "error", err)
		}
	}()

	<-stop

	slog.Info("Shutting down server...")

	if err := server.Close(); err != nil {
		slog.Error("Server forced to shut down", "error", err)
	}
}
