package api

import (
	"diploma/pkg/db"
	"encoding/json"
	"log"
	"net/http"
)

// TasksResponse структура ответа для списка задач
type TasksResponse struct {
	Tasks []db.Task `json:"tasks"`
}

// writeError отправляет ошибку в json формате
func writeError(w http.ResponseWriter, statusCode int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	if err := json.NewEncoder(w).Encode(map[string]string{"error": message}); err != nil {
		log.Printf("json encode error: %v", err)
	}
}
