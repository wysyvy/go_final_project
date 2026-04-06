package api

import (
	"diploma/pkg/db"
	"encoding/json"
	"net/http"
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
	var t Task
	if err := json.NewDecoder(r.Body).Decode(&t); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if t.Title == "" {
		http.Error(w, "Title is required", http.StatusBadRequest)
		return
	}
	if t.Date == "" {
		http.Error(w, "Date is required", http.StatusBadRequest)
		return
	}

	result, err := db.DB.Exec(
		"INSERT INTO scheduler (title, comment, date, repeat) VALUES (?, ?, ?, ?)",
		t.Title, t.Comment, t.Date, t.Repeat,
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	id, _ := result.LastInsertId()
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{"id": id})
}

func RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /api/tasks", GetAllTasks)
	mux.HandleFunc("POST /api/task", AddTask)
}
