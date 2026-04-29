package main

import (
	"sprint13-14/pkg/server" // Импортируем наш созданный пакет
)

func main() {
	// Указываем, где лежит фронтенд
	webDir := "./web"

	// Запускаем сервер через наш пакет
	server.Start(webDir)
}
