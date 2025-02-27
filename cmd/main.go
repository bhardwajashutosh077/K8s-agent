package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/bhardwajashutosh077/k8s-agent/internal/agent"
)

func handleLogs(w http.ResponseWriter, r *http.Request) {
	// Check if the log file exists, if not, create it
	if _, err := os.Stat("scaling.log"); os.IsNotExist(err) {
		err := os.WriteFile("scaling.log", []byte("Log file created.\n"), 0644)
		if err != nil {
			http.Error(w, "Could not create log file", http.StatusInternalServerError)
			return
		}
	}

	// Read the log file
	logData, err := os.ReadFile("scaling.log")
	if err != nil {
		http.Error(w, "Could not read logs", http.StatusInternalServerError)
		return
	}
	fmt.Fprint(w, string(logData))
}

func main() {
	// Ensure the log file exists when the application starts
	if _, err := os.Stat("scaling.log"); os.IsNotExist(err) {
		os.WriteFile("scaling.log", []byte("Server started.\n"), 0644)
	}

	// Start the scaling agent in a separate goroutine
	go agent.StartAgent()

	http.HandleFunc("/logs", handleLogs)
	http.Handle("/", http.FileServer(http.Dir("./web/static")))

	fmt.Println("Server running on port 8080...")
	http.ListenAndServe(":8080", nil)
}
