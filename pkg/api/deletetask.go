package api

import (
	"diploma/pkg/db"
	"encoding/json"
	"log"
	"net/http"
)

// DeleteTaskHandler удаляет задачу
func DeleteTaskHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		writeError(w, http.StatusBadRequest, "id is required")
		return
	}

	if err := db.DeleteTaskByID(idStr); err != nil {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(map[string]interface{}{}); err != nil {
		log.Printf("json encode error: %v", err)
	}
}
