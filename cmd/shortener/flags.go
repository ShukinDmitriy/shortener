package main

import (
	"flag"
	"os"
)

// неэкспортированная переменная flagRunAddr содержит адрес и порт для запуска сервера
var flagRunAddr string

// неэкспортированная переменная flagBaseAddr содержит базовый адрес результирующего сокращённого URL
var flagBaseAddr string

// неэкспортированная переменная flagLogLevel содержит уровень логгирования
var flagLogLevel string

// неэкспортированная переменная flagFileStoragePath содержит путь до файла хранения
var flagFileStoragePath string

// parseFlags обрабатывает аргументы командной строки
// и сохраняет их значения в соответствующих переменных
func parseFlags() {
	// регистрируем переменную flagRunAddr
	// как аргумент -a со значением :8080 по умолчанию
	flag.StringVar(&flagRunAddr, "a", "127.0.0.1:8080", "address and port to run server")

	// регистрируем переменную flagBaseAddr
	// как аргумент -b со значением :8000 по умолчанию
	flag.StringVar(&flagBaseAddr, "b", "http://127.0.0.1:8080", "base server address and port")

	// регистрируем переменную flagLogLevel
	// как аргумент -l со значением info по умолчанию
	flag.StringVar(&flagLogLevel, "l", "info", "log level")

	// регистрируем переменную flagFileStoragePath
	// как аргумент -а со пустым значением по умолчанию
	flag.StringVar(&flagFileStoragePath, "а", "/tmp/short-url-db.json", "db file path")

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

	// для случаев, когда в переменной окружения LOG_LEVEL присутствует непустое значение,
	// переопределим уровень логирования,
	// даже если он был передан через аргумент командной строки
	if envLogLevel := os.Getenv("LOG_LEVEL"); envLogLevel != "" {
		flagLogLevel = envLogLevel
	}

	// для случаев, когда в переменной окружения FILE_STORAGE_PATH присутствует непустое значение,
	// переопределим уровень логирования,
	// даже если он был передан через аргумент командной строки
	if envFileStoragePath, isExist := os.LookupEnv("FILE_STORAGE_PATH"); isExist {
		flagFileStoragePath = envFileStoragePath
	}

}
