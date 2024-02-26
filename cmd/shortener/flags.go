package main

import (
	"flag"
	"os"
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
	flag.StringVar(&flagBaseAddr, "b", "http://127.0.0.1:8080", "base server address and port")

	// парсим переданные серверу аргументы в зарегистрированные переменные
	flag.Parse()

	// для случаев, когда в переменной окружения SERVER_ADDRESS присутствует непустое значение,
	// переопределим адрес запуска сервера,
	// даже если он был передан через аргумент командной строки
	if envRunAddr := os.Getenv("SERVER_ADDRESS"); envRunAddr != "" {
		flagRunAddr = envRunAddr
	}

	// для случаев, когда в переменной окружения BASE_URL присутствует непустое значение,
	// переопределим адрес запуска сервера,
	// даже если он был передан через аргумент командной строки
	if envBaseAddr := os.Getenv("BASE_URL"); envBaseAddr != "" {
		flagBaseAddr = envBaseAddr
	}

}
