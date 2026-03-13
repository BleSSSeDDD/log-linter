package test

import (
	"log"
	"log/slog"
	// zap убран из тестов, но поддерживается в коде
)

func main() {
	log.Println("server started")
	slog.Info("user logged in", "id", 123)

	log.Println("Server started")   // want "лог-сообщение должно начинаться со строчной буквы"
	slog.Error("Failed to connect") // want "лог-сообщение должно начинаться со строчной буквы"
	log.Println("сервер запущен")   // want "лог-сообщение должно начинаться со строчной буквы"
	log.Println("server started!")  // want "лог-сообщение должно начинаться со строчной буквы"

	log.Println("user password is 12345") // want "аргумент лога содержит потенциально чувствительные данные"
	slog.Info("token abc123")             // want "аргумент лога содержит потенциально чувствительные данные"
}
