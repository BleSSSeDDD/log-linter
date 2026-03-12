package main

import (
	analyzer "github.com/BleSSSeDDD/log-linter/log-analyzer"
	"golang.org/x/tools/go/analysis"
)

func New(conf any) ([]*analysis.Analyzer, error) {
	// TODO: конфиг потом тут будет

	return []*analysis.Analyzer{
		analyzer.LogAnalyzer,
	}, nil
}
