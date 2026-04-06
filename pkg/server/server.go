package server

import (
	"diploma/pkg/api"
	"fmt"
	"net/http"
)

func Start(port string) error {
	mux := http.NewServeMux()

	// Статика
	webDir := "./web"
	fs := http.FileServer(http.Dir(webDir))
	mux.Handle("/", fs)

	// API (регистрируем обработчики из пакета api)
	api.RegisterRoutes(mux)

	fmt.Printf("Server starting on :%s\n", port)
	return http.ListenAndServe(":"+port, mux)
}
