# log-linter

Линтер для Go, совместимый с golangci-lint. Проверяет лог-записи на соответствие правилам.

## Что проверяем
- Строчная буква в начале сообщения
- Только английский язык (латиница)
- Отсутствие спецсимволов и эмодзи
- Отсутствие чувствительных данных (password, token, secret, key, auth, credential)

## Поддерживаемые логгеры
- log
- log/slog
- go.uber.org/zap

## Что не обрабатывается
- Не обрабатываем zap.L().Info() и подобные цепочки вызовов
- Не отслеживаем переменные (только прямые строки и конкатенацию)
- Нет автофиксов
- Нет конфигурации через файл

## Установка
```bash
go build -buildmode=plugin -o loglinter.so plugin/plugin.go
```

В .golangci.yml:
```yaml
version: "2"

linters:
  default: none
  enable:
    - loglinter

  settings:
    custom:
      loglinter:
        path: ./loglinter.so
        description: "Линтер для проверки лог-сообщений"
        original-url: github.com/BleSSSeDDD/log-linter
```

## Запуск
```bash
golangci-lint run
```

## Тесты
```bash
go test ./analyzer -v
```
