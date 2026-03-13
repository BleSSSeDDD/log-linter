package analyzer

import (
	"testing"

	"golang.org/x/tools/go/analysis/analysistest"
)

func TestAnalyzer(t *testing.T) {
	testdata := analysistest.TestData()
	analysistest.Run(t, testdata, LogAnalyzer, "testlintdata/test")
}

func TestIsValidLogMessage(t *testing.T) {
	tests := []struct {
		name string
		msg  string
		want bool
	}{
		{"валидный", "server started", true},
		{"большая буква", "Server started", false},
		{"русский", "сервер запущен", false},
		{"спецсимвол", "server started!", false},
		{"пустая строка", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isValidLogMessage(tt.msg); got != tt.want {
				t.Errorf("isValidLogMessage(%q) = %v, want %v", tt.msg, got, tt.want)
			}
		})
	}
}

func TestContainsSensitiveData(t *testing.T) {
	tests := []struct {
		name string
		text string
		want bool
	}{
		{"пароль", "user password", true},
		{"токен", "token: abc", true},
		{"безопасно", "server started", false},
		{"похожее слово", "keyboard", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := containsSensitiveData(tt.text); got != tt.want {
				t.Errorf("containsSensitiveData(%q) = %v, want %v", tt.text, got, tt.want)
			}
		})
	}
}
