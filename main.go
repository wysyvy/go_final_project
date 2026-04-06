package main

import (
	"diploma/pkg/db"
	"diploma/pkg/server"
	"log"
)

func main() {
	// Инициализируем БД
	if err := db.Init("scheduler.db"); err != nil {
		log.Fatal("Database init failed:", err)
	}
	defer db.DB.Close()

	port := "7540"
	if err := server.Start(port); err != nil {
		log.Fatal("Server failed:", err)
	}
}
