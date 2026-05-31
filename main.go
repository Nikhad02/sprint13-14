package main

import (
	"log"
	"os"
	"sprint13-14/pkg/db"
	"sprint13-14/pkg/server"
)

func main() {
	dbFile := os.Getenv("TODO_DBFILE")
	if dbFile == "" {
		dbFile = "scheduler.db"
	}

	err := db.Init(dbFile)
	if err != nil {
		log.Fatalf("Ошибка инициализации БД: %v", err)
	}
	defer db.DB.Close()

	webDir := "./web"
	server.Start(webDir)
}
