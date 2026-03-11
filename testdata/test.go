package test

import (
	"log"
	"log/slog"
)

func main() {
	// не должен ругаться
	log.Println("starting server")
	slog.Info("server started")
	log.Printf("user %s logged in", "alice")

	// большие буквы
	log.Println("Starting server")  // ERROR: должно быть "starting"
	slog.Error("Failed to connect") // ERROR: должно быть "failed"

	// русский текст
	log.Println("запуск сервера")   // ERROR: русский текст
	slog.Warn("ошибка подключения") // ERROR: русский текст

	// спецсимволы
	log.Println("server started!")    // ERROR: !
	slog.Info("connection failed!!!") // ERROR: !!!
	log.Println("loading...")         // ERROR: ...
	slog.Debug("done 🚀")              // ERROR: эмодзи

	// чувствительные данные
	log.Println("user password: 12345") // ERROR: password
	slog.Info("api_key=abc123")         // ERROR: api_key
	log.Printf("token: %s", "secret")   // ERROR: token
}
