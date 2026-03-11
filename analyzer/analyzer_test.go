package analyzer

import (
	"testing"

	"golang.org/x/tools/go/analysis/analysistest"
)

func TestAnalyzer(t *testing.T) {
	testdata := analysistest.TestData()
	analysistest.Run(t, testdata, Analyzer, "testlintdata/test")
}
func TestNewLogRules(t *testing.T) {
	tests := []struct {
		name      string
		msg       string
		wantValid bool
		wantSens  bool
	}{
		{"валидный", "server started", true, false},
		{"с большой буквы", "Server started", false, false},
		{"с цифрой", "server v2 started", false, false},
		{"с дефисом", "user-login failed", false, false},
		{"русский", "сервер запущен", false, false},
		{"спецсимвол", "server started!", false, false},
		{"чувствительный пароль", "user password 123", false, true},
		{"чувствительный токен", "token is valid", false, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isValidLogMessage(tt.msg); got != tt.wantValid {
				t.Errorf("isValidLogMessage() = %v, want %v for %q", got, tt.wantValid, tt.msg)
			}
			if got := containsSensitiveData(tt.msg); got != tt.wantSens {
				t.Errorf("containsSensitiveData() = %v, want %v for %q", got, tt.wantSens, tt.msg)
			}
		})
	}
}
