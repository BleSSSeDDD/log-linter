package analyzer

import (
	"go/ast"
	"go/token"
	"strings"
	"unicode"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

var LogAnalyzer = &analysis.Analyzer{
	Name:     "loglinter",
	Doc:      "проверка лог-сообщенгий на валидность",
	Run:      run,
	Requires: []*analysis.Analyzer{inspect.Analyzer},
}

func run(pass *analysis.Pass) (interface{}, error) {
	inspect := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)

	// нам нужны только вызовы функций
	nodeFilter := []ast.Node{
		(*ast.CallExpr)(nil),
	}

	inspect.Preorder(nodeFilter, func(node ast.Node) {
		call := node.(*ast.CallExpr)

		// проверяем, что это вызов методов логгера
		if !isLogMethodCall(pass, call) {
			return
		}

		if len(call.Args) == 0 {
			return
		}

		for i, arg := range call.Args {
			// рекурсивно получаем текст (для поиска чувствительных данных даже в конкатенации или именах переменных)
			fullContent := recursiveExtractText(arg)

			// проверка на чувствительные данные
			if containsSensitiveData(fullContent) {
				pass.Reportf(arg.Pos(), "аргумент лога содержит потенциально чувствительные данные: %q", fullContent)
			}
			//проверка строковых литералов
			if argLiteral, ok := arg.(*ast.BasicLit); ok && argLiteral.Kind == token.STRING {
				strArg := strings.Trim(argLiteral.Value, "`\"")

				if i == 0 {
					// для первого аргумента (основное сообщение)
					if !isValidLogMessage(strArg) {
						pass.Reportf(argLiteral.Pos(), "лог-сообщение должно начинаться со строчной буквы, быть на английском и не содержать спецсимволы: %q", strArg)
					}
				} else {
					// для остальных аргументов (ключи/значения) проверяем только английский и отсутствие спецсимволов
					if !isEnglishAndSafe(strArg, false) {
						pass.Reportf(argLiteral.Pos(), "дополнительный аргумент лога содержит недопустимые символы или не на английском: %q", strArg)
					}
				}
			}
		}
	})

	return nil, nil
}

// определяет, принадлежит ли вызов логгерам через анализ типов
func isLogMethodCall(pass *analysis.Pass, call *ast.CallExpr) bool {
	var sel *ast.SelectorExpr

	switch fun := call.Fun.(type) {
	case *ast.SelectorExpr:
		sel = fun
	case *ast.Ident:
		return false
	default:
		return false
	}

	// чтобы видеть методы через интерфейсы используем TypesInfo
	selection, ok := pass.TypesInfo.Selections[sel]
	if !ok {
		// eсли это вызов функции пакета, а не метода объекта (например log.Println)
		obj := pass.TypesInfo.Uses[sel.Sel]
		if obj == nil || obj.Pkg() == nil {
			return false
		}
		return isLogPackage(obj.Pkg().Path()) && isLogMethodName(sel.Sel.Name)
	}

	pkg := selection.Obj().Pkg()
	if pkg == nil {
		return false
	}

	return isLogPackage(pkg.Path()) && isLogMethodName(selection.Obj().Name())
}

func isLogPackage(path string) bool {
	return path == "log" || path == "log/slog" || path == "go.uber.org/zap"
}

func isLogMethodName(name string) bool {
	methods := map[string]bool{
		"Info": true, "Infof": true, "Infoln": true,
		"Error": true, "Errorf": true, "Errorln": true,
		"Warn": true, "Warnf": true, "Warnln": true,
		"Debug": true, "Debugf": true, "Debugln": true,
		"Fatal": true, "Fatalf": true, "Fatalln": true,
		"Print": true, "Printf": true, "Println": true,
		"Panic": true, "Panicf": true, "Panicln": true,
	}
	return methods[name]
}

// проверяет только алфавит и отсутствие спецсимволов
func isEnglishAndSafe(s string, isItFirstArg bool) bool {
	if isItFirstArg {
		for _, r := range s {
			isEnglish := (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z')
			isDigit := unicode.IsDigit(r)
			isSpace := r == ' '

			if !isEnglish && !isDigit && !isSpace {
				return false
			}
		}
	} else {
		// тут уже может быть подчеркивание
		for _, r := range s {
			isLetter := (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z')
			isDigit := r >= '0' && r <= '9'
			isSpace := r == ' '
			isUnderscore := r == '_'

			if !isLetter && !isDigit && !isSpace && !isUnderscore {
				return false
			}
		}
	}

	return true
}

// проверяет сообщение по полному набору правил
func isValidLogMessage(s string) bool {
	if s == "" {
		return false
	}

	runes := []rune(s)

	if !unicode.IsLower(runes[0]) || runes[0] < 'a' || runes[0] > 'z' {
		return false
	}

	return isEnglishAndSafe(s, true)
}

// достает текст из литералов, конкатенаций и имен переменных
func recursiveExtractText(n ast.Node) string {
	switch x := n.(type) {
	case *ast.BasicLit:
		if x.Kind == token.STRING {
			return strings.ToLower(strings.Trim(x.Value, `"`))
		}
	case *ast.BinaryExpr:
		if x.Op == token.ADD {
			return recursiveExtractText(x.X) + " " + recursiveExtractText(x.Y)
		}
	case *ast.Ident:
		// проверяем имя переменной
		return strings.ToLower(x.Name)
	case *ast.CallExpr:
		// если внутри вызова лога есть другой вызов по типу Sprintf, проверяем его аргументы
		var res []string
		for _, arg := range x.Args {
			res = append(res, recursiveExtractText(arg))
		}
		return strings.Join(res, " ")
	}
	return ""
}

func containsSensitiveData(s string) bool {
	lower := strings.ToLower(s)
	sensitive := []string{
		"password", "token", "secret", "key", "auth",
		"credential",
	}

	for _, word := range sensitive {
		words := strings.Fields(lower)
		for _, w := range words {
			if w == word {
				return true
			}
			// вхождение с пунктуацией
			if strings.Contains(w, word) && (strings.HasSuffix(w, ":") || strings.HasSuffix(w, "=")) {
				return true
			}
		}
	}
	return false
}
