package api

import (
	"diploma/pkg/db"
	"diploma/pkg/nextdate"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type Task struct {
	ID      string `json:"id"`
	Title   string `json:"title"`
	Comment string `json:"comment"`
	Date    string `json:"date"`
	Repeat  string `json:"repeat"`
}

type TasksResponse struct {
	Tasks []Task `json:"tasks"`
}

func RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /api/tasks", GetAllTasks)
	mux.HandleFunc("POST /api/task", AddTask)
	mux.HandleFunc("GET /api/nextdate", NextDateHandler)
	mux.HandleFunc("GET /api/task", GetTask)
	mux.HandleFunc("PUT /api/task", UpdateTask)
	mux.HandleFunc("POST /api/task/done", DoneTask)
	mux.HandleFunc("DELETE /api/task", DeleteTaskHandler)
}

func convertDate(dateStr string) string {
	parts := strings.Split(dateStr, ".")
	if len(parts) != 3 {
		return dateStr
	}
	return parts[2] + parts[1] + parts[0]
}

func GetAllTasks(w http.ResponseWriter, r *http.Request) {
	// fmt.Println("SELECT query executed")
	rows, err := db.DB.Query("SELECT id, title, comment, date, repeat FROM scheduler ORDER BY date LIMIT 50")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	tasks := []Task{}

	for rows.Next() {
		var t Task
		err := rows.Scan(&t.ID, &t.Title, &t.Comment, &t.Date, &t.Repeat)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		tasks = append(tasks, t)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(TasksResponse{Tasks: tasks})
}

func AddTask(w http.ResponseWriter, r *http.Request) {
	var task db.Task
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid JSON"})
		return
	}

	if task.Title == "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Title is required"})
		return
	}

	now := time.Now()

	if task.Date == "" {
		task.Date = now.Format("20060102")
	}

	_, err := time.Parse("20060102", task.Date)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid date format"})
		return
	}

	if task.Repeat != "" {
		parts := strings.Fields(task.Repeat)
		if len(parts) == 0 || (parts[0] != "y" && parts[0] != "d") {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"error": "Unsupported repeat rule"})
			return
		}
		if parts[0] == "d" {
			if len(parts) != 2 {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(map[string]string{"error": "Invalid days format"})
				return
			}
			days, err := strconv.Atoi(parts[1])
			if err != nil || days <= 0 || days > 400 {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(map[string]string{"error": "Invalid days count"})
				return
			}
		}
	}

	id, err := db.AddTask(&task)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{"id": fmt.Sprintf("%d", id)})
}

// NextDateHandler обрабатывает запрос
func NextDateHandler(w http.ResponseWriter, r *http.Request) {
	nowStr := r.URL.Query().Get("now")
	date := r.URL.Query().Get("date")
	repeat := r.URL.Query().Get("repeat")

	if date == "" || repeat == "" {
		http.Error(w, "missing parameters date and repeat required", http.StatusBadRequest)
		return
	}

	var now time.Time
	if nowStr == "" {
		now = time.Now()
	} else {
		var err error
		now, err = time.Parse("20060102", nowStr)
		if err != nil {
			http.Error(w, "Invalid now date format", http.StatusBadRequest)
			return
		}
	}

	next, err := nextdate.NextDate(now, date, repeat)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte(next))
}

func GetTask(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "id is required"})
		return
	}

	var task db.Task
	row := db.DB.QueryRow("SELECT id, date, title, comment, repeat FROM scheduler WHERE id = ?", idStr)
	err := row.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"error": "Task not found"})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(task)
}

func UpdateTask(w http.ResponseWriter, r *http.Request) {
	var task db.Task
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid JSON"})
		return
	}

	if task.Title == "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Title is required"})
		return
	}

	if len(task.Date) != 8 {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid date format"})
		return
	}

	for _, ch := range task.Date {
		if ch < '0' || ch > '9' {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"error": "Invalid date format"})
			return
		}
	}

	_, err := time.Parse("20060102", task.Date)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid date format"})
		return
	}

	if task.Repeat != "" {
		now := time.Now()
		_, err := nextdate.NextDate(now, task.Date, task.Repeat)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"error": "Invalid repeat rule: " + err.Error()})
			return
		}
	}

	result, err := db.DB.Exec(
		"UPDATE scheduler SET date = ?, title = ?, comment = ?, repeat = ? WHERE id = ?",
		task.Date, task.Title, task.Comment, task.Repeat, task.ID,
	)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"error": "Task not found"})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{})
}

func DoneTask(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "id is required"})
		return
	}

	task, err := db.GetTask(idStr)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"error": "Task not found"})
		return
	}

	if task.Repeat == "" {
		if err := db.DeleteTask(idStr); err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{})
		return
	}

	now := time.Now()
	nextDate, err := nextdate.NextDate(now, task.Date, task.Repeat)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid repeat rule: " + err.Error()})
		return
	}

	if err := db.UpdateTaskDate(idStr, nextDate); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{})
}

func DeleteTaskHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "id is required"})
		return
	}

	if err := db.DeleteTask(idStr); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{})
}
