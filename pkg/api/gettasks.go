package api

import (
	"diploma/pkg/db"
	"encoding/json"
	"log"
	"net/http"
)

// GetAllTasks возвращает список всех задач
func GetAllTasks(w http.ResponseWriter, r *http.Request) {
	tasks, err := db.GetAllTasks()
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(TasksResponse{Tasks: tasks}); err != nil {
		log.Printf("json encode error: %v", err)
	}
}
