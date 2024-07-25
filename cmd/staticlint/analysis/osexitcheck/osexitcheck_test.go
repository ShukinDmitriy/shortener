package osexitcheck_test

import (
	"testing"

	"github.com/ShukinDmitriy/shortener/cmd/staticlint/analysis/osexitcheck"
	"golang.org/x/tools/go/analysis/analysistest"
)

func TestOsExitCheck(t *testing.T) {
	// функция analysistest.Run применяет тестируемый анализатор OsExitCheckAnalyzer
	// к пакетам из папки testdata и проверяет ожидания
	// ./... — проверка всех поддиректорий в testdata
	// можно указать ./pkg1 для проверки только pkg1
	testData := analysistest.TestData()
	analyzer := osexitcheck.OsExitCheckAnalyzer
	analysistest.Run(t, testData, analyzer, "./...")
}
