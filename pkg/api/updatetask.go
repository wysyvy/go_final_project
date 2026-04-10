package api

import (
	"diploma/pkg/db"
	"diploma/pkg/nextdate"
	"encoding/json"
	"log"
	"net/http"
	"time"
)

// UpdateTask обновляет существующую задачу
func UpdateTask(w http.ResponseWriter, r *http.Request) {
	var task db.Task
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		writeError(w, http.StatusBadRequest, "invalid json")
		return
	}

	if task.Title == "" {
		writeError(w, http.StatusBadRequest, "title is required")
		return
	}

	if len(task.Date) != 8 {
		writeError(w, http.StatusBadRequest, "invalid date format")
		return
	}

	for _, ch := range task.Date {
		if ch < '0' || ch > '9' {
			writeError(w, http.StatusBadRequest, "invalid date format")
			return
		}
	}

	_, err := time.Parse("20060102", task.Date)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid date format")
		return
	}

	if task.Repeat != "" {
		now := time.Now()
		_, err := nextdate.NextDate(now, task.Date, task.Repeat)
		if err != nil {
			writeError(w, http.StatusBadRequest, "invalid repeat rule: "+err.Error())
			return
		}
	}

	if err := db.UpdateTask(&task); err != nil {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(map[string]interface{}{}); err != nil {
		log.Printf("json encode error: %v", err)
	}
}
