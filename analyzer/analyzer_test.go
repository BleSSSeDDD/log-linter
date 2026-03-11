package analyzer

import (
	"testing"

	"golang.org/x/tools/go/analysis/analysistest"
)

func TestAnalyzer(t *testing.T) {
	testdata := analysistest.TestData()
	analysistest.Run(t, testdata, Analyzer, "testlintdata/test")
}

func TestIsValidLogMessage(t *testing.T) {
	tests := []struct {
		name string
		msg  string
		want bool
	}{
		{"валидное сообщение", "server started", true},
		{"с цифрами", "server v2 started", true},
		{"с дефисом", "user-login failed", true},
		{"с подчеркиванием", "api_key invalid", true},
		{"со слешем", "file /tmp/test not found", true},

		{"большая буква в начале", "Server started", false},
		{"русский текст", "сервер запущен", false},
		{"восклицательный знак", "server started!", false},
		{"вопрос", "server started?", false},
		{"точка", "server started.", false},
		{"эмодзи", "server started 🚀", false},
		{"спецсимволы", "error: failed!!!", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isValidLogMessage(tt.msg); got != tt.want {
				t.Errorf("isValidLogMessage(%q) = %v, want %v", tt.msg, got, tt.want)
			}
		})
	}
}
