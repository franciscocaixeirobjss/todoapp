package main

import (
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"

	"todoapp/files"
	"todoapp/handlers"
	"todoapp/logging"
	"todoapp/middleware"
	"todoapp/task"
	"todoapp/webserver"
)

func main() {
	requestChanSize := flag.Int("requestChanSize", 10, "Size of the request channel")
	port := flag.String("port", "8081", "Port to run the backend server on")
	flag.Parse()

	filename := filepath.Join("..", "files", "server_"+*port+".json")

	logging.InitLogging(*port)

	var tasks map[int][]task.Task
	var maxTaskIDs map[int]int
	err := files.LoadData(filename, &tasks, &maxTaskIDs)
	if err != nil {
		log.Printf("Failed to load data: %v", err)
		return
	}

	task.SetTasks(tasks, maxTaskIDs)
	task.InitChannel(*requestChanSize)

	defer func() {
		tasks, maxTaskIDs = task.GetManagerTasks()
		if err := files.SaveData(filename, tasks, maxTaskIDs); err != nil {
			log.Printf("Failed to save tasks to file: %v", err)
		} else {
			log.Println("Tasks saved successfully")
		}
	}()

	mux := http.NewServeMux()

	mux.HandleFunc("/create", func(w http.ResponseWriter, r *http.Request) {
		middleware.ChainMiddleware(
			http.HandlerFunc(handlers.CreateHandler),
			middleware.TraceIDMiddleware,
			middleware.UserIDMiddleware,
		).ServeHTTP(w, r)
	})
	mux.HandleFunc("/get", func(w http.ResponseWriter, r *http.Request) {
		middleware.ChainMiddleware(
			http.HandlerFunc(handlers.GetHandler),
			middleware.TraceIDMiddleware,
			middleware.UserIDMiddleware,
		).ServeHTTP(w, r)
	})

	mux.HandleFunc("/update", func(w http.ResponseWriter, r *http.Request) {
		middleware.ChainMiddleware(
			http.HandlerFunc(handlers.UpdateHandler),
			middleware.TraceIDMiddleware,
			middleware.UserIDMiddleware,
		).ServeHTTP(w, r)
	})
	mux.HandleFunc("/delete/", func(w http.ResponseWriter, r *http.Request) {
		middleware.ChainMiddleware(
			http.HandlerFunc(handlers.DeleteHandler),
			middleware.TraceIDMiddleware,
			middleware.UserIDMiddleware,
		).ServeHTTP(w, r)
	})

	webserver.ServeStaticPage(mux)
	webserver.ServeDynamicPage(mux)

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	server := &http.Server{
		Addr:    ":" + *port,
		Handler: mux,
	}

	go func() {
		log.Printf("Starting backend server on port %s...", *port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed: %v", err)
		}
	}()

	<-stop
	log.Println("Shutting down server...")
	if err := server.Close(); err != nil {
		log.Fatalf("Server shutdown failed: %v", err)
	}
	log.Println("Server gracefully stopped.")
}
