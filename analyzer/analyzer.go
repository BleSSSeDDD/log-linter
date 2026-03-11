package analyzer

import (
	"go/ast"
	"go/token"
	"strings"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

var Analyzer = &analysis.Analyzer{
	Name:     "loglinter",
	Doc:      "вот бы оно завелось хотя бы",
	Run:      run,
	Requires: []*analysis.Analyzer{inspect.Analyzer},
}

func run(pass *analysis.Pass) (interface{}, error) {
	inspect := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)

	nodeFilter := []ast.Node{
		(*ast.CallExpr)(nil),
	}

	inspect.Preorder(nodeFilter, func(n ast.Node) {
		call := n.(*ast.CallExpr)
		_, ok := call.Fun.(*ast.SelectorExpr)
		if !ok {
			return
		}

		// === ВАЖНО: временно убираем жесткую фильтрацию, чтобы видеть все логи ===
		// if sel.Sel.Name != "Println" {
		// 	return
		// }

		// ident, ok := sel.X.(*ast.Ident)
		// if !ok || ident.Name != "log" {
		// 	return
		// }

		// pass.Reportf(call.Pos(), "нашел log функцию")

		// === НОВАЯ ЛОГИКА: проверяем, что это вообще вызов метода (X.Y) ===
		// Это выражение вида log.Info, logger.Println, slog.Error и т.д.
		// Мы пока не проверяем пакет, чтобы увидеть всё, что может быть логом.
		// Позже здесь нужно будет добавить вызов isLogFunction(pass, call)

		// Пытаемся извлечь сообщение (первый аргумент-строку)
		if len(call.Args) == 0 {
			return
		}
		firstArg := call.Args[0]
		lit, ok := firstArg.(*ast.BasicLit)
		if !ok || lit.Kind != token.STRING {
			// Если первый аргумент не строка (например, переменная или конкатенация),
			// мы пока не можем его проверить. Это задача на будущее.
			return
		}

		// Извлекаем текст сообщения, убирая кавычки
		msg := strings.Trim(lit.Value, `"`)

		// Применяем проверки
		if !isValidLogMessage(msg) {
			pass.Reportf(call.Pos(), "лог-сообщение содержит недопустимые символы (только англ буквы и пробелы, начинаться со строчной): %q", msg)
		}

		if containsSensitiveData(msg) {
			pass.Reportf(call.Pos(), "лог-сообщение содержит потенциально чувствительные данные: %q", msg)
		}
	})

	return nil, nil
}

func isValidLogMessage(msg string) bool {
	if msg == "" {
		return false
	}

	// в начале должна быть строчная английская буква
	first := msg[0]
	if first < 'a' || first > 'z' {
		return false
	}

	// а дальше в любом регистре или пробел получается, хотя если в логе например путь какой-то, то ещё слэши бы сюда
	// ну сказано без спецсимволов значит будет без спецсимволов
	for i := 1; i < len(msg); i++ {
		c := msg[i]
		if (c >= 'a' && c <= 'z') ||
			(c >= 'A' && c <= 'Z') ||
			c == ' ' {
			continue
		}
		return false
	}

	return true
}

// пока максимально тупая проверка потом можно думать
func containsSensitiveData(s string) bool {
	lower := strings.ToLower(s)
	sensitive := []string{
		"password", "token", "secret", "key", "auth",
	}
	for _, word := range sensitive {
		if strings.Contains(lower, word) {
			return true
		}
	}
	return false
}
