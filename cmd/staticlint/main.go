package main

import (
	"go/ast"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/multichecker"
	"golang.org/x/tools/go/analysis/passes/printf"
	"golang.org/x/tools/go/analysis/passes/shadow"
	"golang.org/x/tools/go/analysis/passes/structtag"
	"honnef.co/go/tools/staticcheck"
)

// OsExitAnalyzer - Анализатор использования прямого вызова os.Exit в функции main пакета main
var OsExitAnalyzer = &analysis.Analyzer{
	Name: "OsExitAnalyzer",
	Doc:  "проверка использования os.Exit в функции main пакета main",
	Run:  run,
}

func run(pass *analysis.Pass) (interface{}, error) {

	for _, file := range pass.Files {
		ast.Inspect(file, func(node ast.Node) bool {
			switch x := node.(type) {
			case *ast.File:
				if x.Name.Name == "main" {
					return true
				}
				return false
			case *ast.FuncDecl:
				if x.Name.Name == "main" {
					return true
				}
				return false
			case *ast.CallExpr:
				if s, ok := x.Fun.(*ast.SelectorExpr); ok {
					if p, ok := s.X.(*ast.Ident); ok {
						if p.Name == "os" && s.Sel.Name == "Exit" {
							pass.Reportf(x.Pos(), "Вызов os.Exit запрещен")
						}
					}
				}
			}
			return true
		})
	}
	return nil, nil
}

func main() {
	mychecks := []*analysis.Analyzer{OsExitAnalyzer, printf.Analyzer, shadow.Analyzer, structtag.Analyzer}
	for _, v := range staticcheck.Analyzers {
		mychecks = append(mychecks, v.Analyzer)
	}

	multichecker.Main(mychecks...)
}
