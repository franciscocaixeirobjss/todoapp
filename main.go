package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	"todoapp/api"
	"todoapp/files"
	"todoapp/logging"
)

func main() {
	opts := &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}

	handler := slog.NewTextHandler(os.Stdout, opts)
	logger := slog.New(handler)
	slog.SetDefault(logger)

	var listOfTasks []api.Task
	var maxTaskID int

	err := files.LoadData("todo.json", &listOfTasks, &maxTaskID)
	if err != nil {
		slog.Error("Failed to load data", "error", err)
		return
	}

	// Initialize handlers with shared state
	handlers := &api.Handlers{
		Tasks:     &listOfTasks,
		MaxTaskID: &maxTaskID,
	}

	mux := http.NewServeMux()

	mux.HandleFunc("/create", handlers.CreateHandler)
	mux.HandleFunc("/get", handlers.GetHandler)
	mux.HandleFunc("/update", handlers.UpdateHandler)
	mux.HandleFunc("/delete/", handlers.DeleteHandler)

	// Wrap the mux with the TraceIDMiddleware
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

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		slog.Error("Server forced to shut down", "error", err)
	}

	if err := files.SaveData("todo.json", listOfTasks); err != nil {
		slog.Error("Failed to save tasks to file", "error", err)
	} else {
		slog.Info("Tasks saved successfully")
	}
}
