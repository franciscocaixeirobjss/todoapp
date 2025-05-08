package webserver

import (
	"html/template"
	"log/slog"
	"net/http"
	"todoapp/task"
)

// ServeStaticPage serves a static "about" page
func ServeStaticPage(mux *http.ServeMux) {
	mux.HandleFunc("/about", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./webserver/static/about.html")
	})
}

// ServeDynamicPage serves a dynamic "list" page with all tasks
func ServeDynamicPage(mux *http.ServeMux, taskManager *task.TaskManager) {
	mux.HandleFunc("/list", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			slog.Error("Invalid HTTP method", "method", r.Method)
			http.Error(w, "Invalid HTTP method.", http.StatusMethodNotAllowed)
			return
		}

		w.Header().Set("Content-Type", "text/html")

		tmpl, err := template.ParseFiles("./webserver/templates/list.html")
		if err != nil {
			slog.Error("Failed to parse template", "error", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		tasks := taskManager.GetTasks()
		err = tmpl.Execute(w, tasks)
		if err != nil {
			slog.Error("Failed to execute template", "error", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
	})
}
