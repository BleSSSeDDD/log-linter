package main

import (
	"github.com/BleSSSeDDD/log-linter/analyzer"
	"golang.org/x/tools/go/analysis"
)

func New(conf any) ([]*analysis.Analyzer, error) {
	// TODO: конфиг потом тут будет

	return []*analysis.Analyzer{
		analyzer.LogAnalyzer,
	}, nil
}
