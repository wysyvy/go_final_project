package api

import (
	"diploma/pkg/db"
	"diploma/pkg/nextdate"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// AddTask добавляет новую задачу
func AddTask(w http.ResponseWriter, r *http.Request) {
	var task db.Task
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		writeError(w, http.StatusBadRequest, "invalid json")
		return
	}

	if task.Title == "" {
		writeError(w, http.StatusBadRequest, "title is required")
		return
	}

	now := time.Now()

	if task.Date == "" {
		task.Date = now.Format("20060102")
	}

	parsedDate, err := time.Parse("20060102", task.Date)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid date format")
		return
	}

	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	dateForCompare := time.Date(parsedDate.Year(), parsedDate.Month(), parsedDate.Day(), 0, 0, 0, 0, parsedDate.Location())

	if dateForCompare.Before(today) {
		if task.Repeat == "" {
			task.Date = today.Format("20060102")
		} else {
			nextDate, err := nextdate.NextDate(now, task.Date, task.Repeat)
			if err != nil {
				writeError(w, http.StatusBadRequest, "invalid repeat rule")
				return
			}
			task.Date = nextDate
		}
	}

	if task.Repeat != "" {
		parts := strings.Fields(task.Repeat)
		if len(parts) == 0 || (parts[0] != "y" && parts[0] != "d") {
			writeError(w, http.StatusBadRequest, "unsupported repeat rule")
			return
		}
		if parts[0] == "d" {
			if len(parts) != 2 {
				writeError(w, http.StatusBadRequest, "invalid days format")
				return
			}
			days, err := strconv.Atoi(parts[1])
			if err != nil || days <= 0 || days > 400 {
				writeError(w, http.StatusBadRequest, "invalid days count")
				return
			}
		}
	}

	id, err := db.AddTask(&task)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(map[string]interface{}{"id": fmt.Sprintf("%d", id)}); err != nil {
		log.Printf("json encode error: %v", err)
	}
}
