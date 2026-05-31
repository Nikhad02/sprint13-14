package server

import (
	"log"
	"net/http"
	"os"
	"sprint13-14/pkg/api"
)

func Start(webDir string) {
	port := os.Getenv("TODO_PORT")
	if port == "" {
		port = "7540"
	}

	fileServer := http.FileServer(http.Dir(webDir))
	http.Handle("/", fileServer)

	api.Init()
	log.Printf("Сервер запущен на http://localhost:%s", port)

	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		log.Fatalf("Ошибка запуска сервера: %v", err)
	}
}
