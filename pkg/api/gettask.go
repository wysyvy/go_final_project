package api

import (
	"diploma/pkg/db"
	"encoding/json"
	"log"
	"net/http"
)

// GetTask возвращает задачу по ID
func GetTask(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		writeError(w, http.StatusBadRequest, "id is required")
		return
	}

	task, err := db.GetTaskByID(idStr)
	if err != nil {
		writeError(w, http.StatusNotFound, "task not found")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(task); err != nil {
		log.Printf("json encode error: %v", err)
	}
}
