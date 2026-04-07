package api

import (
	"diploma/pkg/db"
	"diploma/pkg/nextdate"
	"encoding/json"
	"log"
	"net/http"
	"time"
)

type Task struct {
	ID      int    `json:"id"`
	Title   string `json:"title"`
	Comment string `json:"comment"`
	Date    string `json:"date"`
	Repeat  string `json:"repeat"`
}

type TasksResponse struct {
	Tasks []Task `json:"tasks"`
}

func GetAllTasks(w http.ResponseWriter, r *http.Request) {
	rows, err := db.DB.Query("SELECT id, title, comment, date, repeat FROM scheduler ORDER BY date LIMIT 50")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var tasks []Task
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
	log.Println("AddTask called, method:", r.Method, "path:", r.URL.Path)

	var bodyBytes []byte
	r.Body.Read(bodyBytes)
	log.Println("Request body:", string(bodyBytes))

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

	// если дата не указана, то ставим сегодня
	if task.Date == "" {
		task.Date = now.Format("20060102")
	}

	dateTime, err := time.Parse("20060102", task.Date)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid date format"})
		return
	}

	if dateTime.Before(now) {
		if task.Repeat == "" {
			task.Date = now.Format("20060102")
		} else if task.Repeat == "d 1" {
			task.Date = now.Format("20060102")
		} else {
			nextDate, err := nextdate.NextDate(now, task.Date, task.Repeat)
			if err != nil {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(map[string]string{"error": "Invalid repeat rule: " + err.Error()})
				return
			}
			task.Date = nextDate
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
	json.NewEncoder(w).Encode(map[string]interface{}{"id": id})
}

func RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /api/tasks", GetAllTasks)
	mux.HandleFunc("POST /api/task", AddTask)
	mux.HandleFunc("GET /api/nextdate", NextDateHandler)
}

// NextDateHandler обрабатывает запрос
func NextDateHandler(w http.ResponseWriter, r *http.Request) {
	nowStr := r.URL.Query().Get("now")
	date := r.URL.Query().Get("date")
	repeat := r.URL.Query().Get("repeat")

	if date == "" || repeat == "" {
		http.Error(w, "Missing parameters: date and repeat required", http.StatusBadRequest)
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
