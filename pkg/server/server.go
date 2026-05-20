package server

import (
	"log"
	"net/http"
	"os"
	"sprint13-14/pkg/api"
)

// Start запускает наш веб-сервер
func Start(webDir string) {
	// Ищем порт в переменных окружения
	port := os.Getenv("TODO_PORT")
	if port == "" {
		port = "7540" // Порт по умолчанию
	}

	// Настраиваем файловый сервер
	fileServer := http.FileServer(http.Dir(webDir))
	http.Handle("/", fileServer)

	api.Init()
	log.Printf("Сервер запущен на http://localhost:%s", port)

	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		log.Fatalf("Ошибка запуска сервера: %v", err)
	}
}
