package api

import "net/http"

// RegisterRoutes регистрирует все обработчики api
func RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /api/tasks", GetAllTasks)
	mux.HandleFunc("POST /api/task", AddTask)
	mux.HandleFunc("GET /api/nextdate", NextDateHandler)
	mux.HandleFunc("GET /api/task", GetTask)
	mux.HandleFunc("PUT /api/task", UpdateTask)
	mux.HandleFunc("POST /api/task/done", DoneTask)
	mux.HandleFunc("DELETE /api/task", DeleteTaskHandler)
}
