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
	requestChanSize := flag.Int("requestChanSize", 10, "Size of the request channel")
	flag.Parse()

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
		return
	}

	task.SetTasks(tasks, maxTaskID)
	task.InitChannel(*requestChanSize)

	defer func() {
		tasks, maxTaskID = task.GetManagerTasks()
		if err := files.SaveData("todo.json", tasks); err != nil {
			slog.Error("Failed to save tasks to file", "error", err)
		} else {
			slog.Info("Tasks saved successfully")
		}
	}()

	mux := http.NewServeMux()
	mux.HandleFunc("/create", handlers.CreateHandler)
	mux.HandleFunc("/get", handlers.GetHandler)
	mux.HandleFunc("/update", handlers.UpdateHandler)
	mux.HandleFunc("/delete/", handlers.DeleteHandler)

	webserver.ServeStaticPage(mux)
	webserver.ServeDynamicPage(mux)

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
