package api

import (
	"diploma/pkg/db"
	"diploma/pkg/nextdate"
	"encoding/json"
	"log"
	"net/http"
	"time"
)

// DoneTask отмечает задачу как выполненную
func DoneTask(w http.ResponseWriter, r *http.Request) {
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

	if task.Repeat == "" {
		if err := db.DeleteTaskByID(idStr); err != nil {
			writeError(w, http.StatusInternalServerError, err.Error())
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(map[string]interface{}{}); err != nil {
			log.Printf("json encode error: %v", err)
		}
		return
	}

	now := time.Now()
	nextDate, err := nextdate.NextDate(now, task.Date, task.Repeat)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid repeat rule: "+err.Error())
		return
	}

	if err := db.UpdateTaskDate(idStr, nextDate); err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(map[string]interface{}{}); err != nil {
		log.Printf("json encode error: %v", err)
	}
}
