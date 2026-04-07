package main

import (
	"diploma/pkg/api"
	"diploma/pkg/db"
	"log"
	"net/http"
	"strings"
)

// loggingMiddleware логирует все входящие запросы
func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("REQUEST: %s %s", r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
	})
}

func main() {
	if err := db.Init("scheduler.db"); err != nil {
		log.Fatal("Database init failed:", err)
	}
	defer db.DB.Close()

	mux := http.NewServeMux()

	// api регистрируем
	api.RegisterRoutes(mux)

	fs := http.FileServer(http.Dir("./web"))
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if strings.HasPrefix(r.URL.Path, "/api") {
			log.Printf("STATIC HANDLER: rejecting /api request")
			http.NotFound(w, r)
			return
		}
		fs.ServeHTTP(w, r)
	})

	handler := loggingMiddleware(mux)

	port := "7540"
	log.Printf("Server starting on :%s", port)
	if err := http.ListenAndServe(":"+port, handler); err != nil {
		log.Fatal("Server failed:", err)
	}
}
