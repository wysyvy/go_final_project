package server

import (
	"diploma/pkg/api"
	"fmt"
	"net/http"
)

func Start(port string) error {
	mux := http.NewServeMux()

	webDir := "./web"
	fs := http.FileServer(http.Dir(webDir))
	mux.Handle("/", fs)

	// api
	api.RegisterRoutes(mux)

	fmt.Printf("Server starting on :%s\n", port)
	return http.ListenAndServe(":"+port, mux)
}
