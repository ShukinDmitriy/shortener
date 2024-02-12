package main

import (
	"flag"
)

// неэкспортированная переменная flagRunAddr содержит адрес и порт для запуска сервера
var flagRunAddr string

// неэкспортированная переменная flagBaseAddr содержит базовый адрес результирующего сокращённого URL
var flagBaseAddr string

// parseFlags обрабатывает аргументы командной строки
// и сохраняет их значения в соответствующих переменных
func parseFlags() {
	// регистрируем переменную flagRunAddr
	// как аргумент -a со значением :8080 по умолчанию
	flag.StringVar(&flagRunAddr, "a", "127.0.0.1:8080", "address and port to run server")

	// регистрируем переменную flagBaseAddr
	// как аргумент -b со значением :8000 по умолчанию
	flag.StringVar(&flagBaseAddr, "b", "127.0.0.1:8000", "base server address and port")

	// парсим переданные серверу аргументы в зарегистрированные переменные
	flag.Parse()
}
