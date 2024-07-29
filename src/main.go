package main

import (
	"time"
	"strings"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"github.com/google/uuid"
)

type Task struct {
	Description string
	ReqID       string
	Logs []string
}

var tasks = make(map[string]Task)

func main() {
	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/submit", submitHandler)
	http.HandleFunc("/logs", logsHandler)

	fmt.Println("Server is running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("./src/templates/index.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, nil)
}

func submitHandler(response http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodPost {
		http.Error(response, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	description := request.FormValue("description")
	reqID := uuid.New().String()
	fmt.Println("Created a new task : desc -> ", description, ", req_id ->" , reqID)

	task := Task{
		Description: description,
		ReqID:       reqID,
	}

	tasks[reqID] = task

	// Loading the loading page

	tmpl, err := template.ParseFiles("./src/templates/loading.html")
	if err != nil {
		http.Error(response, err.Error(), http.StatusInternalServerError)
		return
	}
	
	response.Header().Set("Content-Type", "text/html")
	tmpl.Execute(response, reqID)
	// fmt.Fprintf(w, "<p>Task submitted! Req ID: %s</p>", reqID)
}

func logsHandler(w http.ResponseWriter, r *http.Request) {
	reqID := r.URL.Query().Get("reqId")
	task, ok := tasks[reqID]
	if !ok {
		http.Error(w, "Task not found", http.StatusNotFound)
		return
	}

	// Simulate log updates
	if len(task.Logs) < 5 {
		newLog := fmt.Sprintf("Processing step %d at %s", len(task.Logs)+1, time.Now().Format("15:04:05"))
		task.Logs = append(task.Logs, newLog)
		tasks[reqID] = task
	} else if len(task.Logs) == 5 {
		task.Logs = append(task.Logs, fmt.Sprintf("Task completed at %s", time.Now().Format("15:04:05")))
		tasks[reqID] = task
	}

	// Convert logs to HTML paragraphs
	var htmlLogs []string
	for _, logEntry := range task.Logs {
		htmlLog := fmt.Sprintf("<p>%s</p>", logEntry)
		htmlLogs = append(htmlLogs, htmlLog)
	}

	// Join all HTML logs into a single string
	logsHTML := strings.Join(htmlLogs, "")

	// Set content type to HTML and write the response
	w.Header().Set("Content-Type", "text/html")
	fmt.Fprintf(w, logsHTML)
}
