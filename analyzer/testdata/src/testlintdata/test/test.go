package test

import (
	"log"
	"log/slog"
)

func main() {
	// ===== ДОЛЖНЫ ПРОХОДИТЬ (валидные логи) =====
	log.Println("server started")
	log.Println("user logged in")
	log.Println("connection established")
	slog.Info("request processed")
	slog.Debug("cache hit")

	// ===== ДОЛЖНЫ ПАДАТЬ ПО ФОРМАТУ =====
	// большая буква в начале
	log.Println("Server started")
	slog.Info("Failed to connect")

	// русский язык
	log.Println("сервер запущен")
	slog.Warn("ошибка подключения")

	// цифры (по текущим правилам - нельзя)
	log.Println("server v2 started")
	slog.Info("port 8080 listening")

	// спецсимволы
	log.Println("server started!")
	slog.Info("connection failed!!!")
	log.Println("loading...")
	slog.Debug("done 🚀")

	// дефисы, подчеркивания, слеши
	log.Println("user-login failed")
	slog.Info("api_key invalid")
	log.Println("file /tmp/test not found")

	// ===== ДОЛЖНЫ ПАДАТЬ ПО ЧУВСТВИТЕЛЬНЫМ ДАННЫМ =====
	// пароли
	log.Println("user password: 12345")
	log.Println("passwd: secret")
	log.Println("pwd: qwerty")

	// токены
	slog.Info("token: abc123")
	log.Println("jwt: eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9")

	// ключи и секреты
	slog.Debug("api_key=abc123")
	log.Println("secret key: 12345")
	log.Println("auth: basic")

	// ===== СЛОЖНЫЕ СЛУЧАИ =====
	// чувствительные данные в валидном формате
	log.Println("invalid token") // формат ок, но есть token
	slog.Info("password reset")  // формат ок, но есть password

	// конкатенация (пока не ловим, но помечаем)
	log.Println("token: " + "secret")

	// форматирование (тоже сложный случай)
	log.Printf("user %s logged in with password %s", "alice", "12345")
}
