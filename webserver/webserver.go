package webserver

import (
	"html/template"
	"log/slog"
	"net/http"
	"strconv"
	"strings"
	"todoapp/task"
)

type PageData struct {
	UserID int
	Tasks  []task.Task
}

// ServeStaticPage serves a static "about" page
func ServeStaticPage(mux *http.ServeMux) {
	mux.HandleFunc("/about", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "../webserver/static/about.html")
	})
}

// ServeDynamicPage serves a dynamic "list" page with all tasks for a specific user
func ServeDynamicPage(mux *http.ServeMux) {
	mux.HandleFunc("/user/", func(w http.ResponseWriter, r *http.Request) {
		slog.Info("Received request for user tasks", "method", r.Method, "path", r.URL.Path)
		if r.Method != http.MethodGet {
			slog.Error("Invalid HTTP method", "method", r.Method)
			http.Error(w, "Invalid HTTP method.", http.StatusMethodNotAllowed)
			return
		}

		pathParts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
		if len(pathParts) != 3 || pathParts[0] != "user" || pathParts[2] != "list" {
			http.Error(w, "Invalid URL pattern. Expected /user/{id}/list", http.StatusBadRequest)
			return
		}

		userID, err := strconv.Atoi(pathParts[1])
		if err != nil {
			http.Error(w, "Invalid UserID", http.StatusBadRequest)
			return
		}

		w.Header().Set("Content-Type", "text/html")

		tmpl, err := template.ParseFiles("../webserver/templates/list.html")
		if err != nil {
			slog.Error("Failed to parse template", "error", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		pageData := PageData{
			UserID: userID,
			Tasks:  task.GetTasks(userID),
		}

		err = tmpl.Execute(w, pageData)
		if err != nil {
			slog.Error("Failed to execute template", "error", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
	})
}
