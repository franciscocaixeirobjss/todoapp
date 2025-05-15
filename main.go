package main

import (
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"todoapp/handlers"
	"todoapp/logging"
	"todoapp/middleware"
)

func main() {
	port := flag.String("port", "8080", "Port to run the backend server on")
	flag.Parse()

	logging.InitLogging(*port)

	mux := createMux()

	wrappedMux := middleware.ChainMiddleware(mux,
		middleware.LoadBalancerMiddleware,
		middleware.TraceIDMiddleware,
		middleware.UserIDMiddleware)

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	go func() {
		log.Printf("Starting server on %s", *port)
		if err := http.ListenAndServe(":"+*port, wrappedMux); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed: %v", err)
		}
	}()

	<-stop
	log.Println("Server shutting down...")
}

func createMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/create", handlers.CreateHandler)
	mux.HandleFunc("/get", handlers.GetHandler)
	mux.HandleFunc("/update", handlers.UpdateHandler)
	mux.HandleFunc("/delete/", handlers.DeleteHandler)

	return mux
}
