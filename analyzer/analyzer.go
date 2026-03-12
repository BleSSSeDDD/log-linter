package analyzer

import (
	"go/ast"
	"go/token"
	"strings"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

var LogAnalyzer = &analysis.Analyzer{
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

	inspect.Preorder(nodeFilter, func(node ast.Node) {
		call := node.(*ast.CallExpr)

		// проверяем, что это вызов метода, а не просто какой-то функции
		selector, isMethod := call.Fun.(*ast.SelectorExpr)
		if !isMethod {
			return
		}

		// проверяем, что метод вызывается у обычной переменной
		ident, ok := selector.X.(*ast.Ident)
		if !ok {
			return
		}

		// берем тип переменной
		obj := pass.TypesInfo.ObjectOf(ident)
		if obj == nil {
			return
		}

		// проверяем, что переменная импортирована из пакета log или slog
		if pkg := obj.Pkg(); pkg == nil || pkg.Path() != "log" || pkg.Path() != "log/slog" {
			return
		}

		if len(call.Args) == 0 {
			return
		}

		for _, arg := range call.Args {
			// проверяем что аргумент это именно строка, а не переменная, число, вызов функции и тп
			argLiteral, isLiteral := arg.(*ast.BasicLit)
			if !isLiteral || argLiteral.Kind != token.STRING {
				continue
			}

			str_arg := strings.Trim(argLiteral.Value, `"`)
			if !isValidLogMessage(str_arg) {
				pass.Reportf(call.Pos(), "лог-сообщение содержит недопустимые символы %q", str_arg)
			}

			if containsSensitiveData(str_arg) {
				pass.Reportf(call.Pos(), "лог-сообщение содержит потенциально чувствительные данные: %q", str_arg)
			}
		}
	})

	return nil, nil
}

func isValidLogMessage(str_arg string) bool {
	if str_arg == "" {
		return false
	}

	for i := range []rune(str_arg) {
		// в начале должна быть строчная английская буква
		if i == 0 {
			if str_arg[i] < 'a' || str_arg[i] > 'z' {
				return false
			}
			continue
		}
		// а дальше в любом регистре или пробел получается, хотя если в логе например путь какой-то, то ещё слэши бы сюда
		// ну сказано без спецсимволов значит будет без спецсимволов
		if (str_arg[i] >= 'a' && str_arg[i] <= 'z') ||
			(str_arg[i] >= 'A' && str_arg[i] <= 'Z') ||
			(str_arg[i] >= '0' && str_arg[i] <= '9') ||
			str_arg[i] == ' ' {
			continue
		}

		return false
	}

	return true
}

// пока максимально тупая проверка потом можно думать
func containsSensitiveData(str_arg string) bool {
	lower := strings.ToLower(str_arg)
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
