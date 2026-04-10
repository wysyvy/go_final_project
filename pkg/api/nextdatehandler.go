package api

import (
	"diploma/pkg/nextdate"
	"net/http"
	"time"
)

// NextDateHandler обрабатывает запросы на вычисление следующей даты
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
			http.Error(w, "invalid now date format", http.StatusBadRequest)
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
