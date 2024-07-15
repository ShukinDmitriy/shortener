// Package osexitcheck analyzer for check use os.Exit
package osexitcheck

import (
	"golang.org/x/tools/go/analysis"
)

// OsExitCheckAnalyzer export analyzer
var OsExitCheckAnalyzer = &analysis.Analyzer{
	Name: "OsExitCheckAnalyzer",
	Doc:  "check for os.exit in main package§	",
	Run:  run,
}

func run(pass *analysis.Pass) (interface{}, error) {
	if pass.Pkg.Name() == "main" {
		for _, file := range pass.Files {
			if file.Name.Name == "main" {
				for i, v := range pass.TypesInfo.Uses {
					if v.String() == "func os.Exit(code int)" {
						pass.Reportf(i.Pos(), "os.Exit in main file")
					}
				}
			}
		}
	}
	return nil, nil
}
