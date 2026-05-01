package main

import (
	"log"
	"os"
	"sprint13-14/pkg/db"
	"sprint13-14/pkg/server"
)

func main() {
	// Определяем файл базы данных
	dbFile := os.Getenv("TODO_DBFILE")
	if dbFile == "" {
		dbFile = "scheduler.db"
	}

	// Инициализируем БД
	err := db.Init(dbFile)
	if err != nil {
		log.Fatalf("Ошибка инициализации БД: %v", err)
	}
	defer db.DB.Close()

	// Запускаем сервер
	webDir := "./web"
	server.Start(webDir)
}
