package analyzer

import (
	"go/ast"
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
		sel, ok := call.Fun.(*ast.SelectorExpr)
		if !ok {
			return
		}

		if sel.Sel.Name != "Println" {
			return
		}

		ident, ok := sel.X.(*ast.Ident)
		if !ok || ident.Name != "log" {
			return
		}

		pass.Reportf(call.Pos(), "нашел log функцию")
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
