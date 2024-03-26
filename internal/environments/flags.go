package environments

import (
	"flag"
	"os"
)

// неэкспортированная переменная FlagRunAddr содержит адрес и порт для запуска сервера
var FlagRunAddr string

// неэкспортированная переменная FlagBaseAddr содержит базовый адрес результирующего сокращённого URL
var FlagBaseAddr string

// неэкспортированная переменная FlagLogLevel содержит уровень логгирования
var FlagLogLevel string

// FlagFileStoragePath содержит путь до файла хранения
var FlagFileStoragePath string

// неэкспортированная переменная FlagDatabaseDSN содержит путь до бд
var FlagDatabaseDSN string

// ParseFlags обрабатывает аргументы командной строки
// и сохраняет их значения в соответствующих переменных
func ParseFlags() {
	// регистрируем переменную FlagRunAddr
	// как аргумент -a со значением :8080 по умолчанию
	flag.StringVar(&FlagRunAddr, "a", "127.0.0.1:8080", "address and port to run server")

	// регистрируем переменную FlagBaseAddr
	// как аргумент -b со значением :8000 по умолчанию
	flag.StringVar(&FlagBaseAddr, "b", "http://127.0.0.1:8080", "base server address and port")

	// регистрируем переменную FlagLogLevel
	// как аргумент -l со значением info по умолчанию
	flag.StringVar(&FlagLogLevel, "l", "info", "log level")

	// регистрируем переменную FlagFileStoragePath
	// как аргумент -f с пустым значением по умолчанию
	flag.StringVar(&FlagFileStoragePath, "f", "/tmp/short-url-db.json", "db file path")

	// регистрируем переменную FlagDatabaseDSN
	// как аргумент -d с пустым значением по умолчанию
	flag.StringVar(&FlagDatabaseDSN, "d", "", "database DSN")

	// парсим переданные серверу аргументы в зарегистрированные переменные
	flag.Parse()

	// для случаев, когда в переменной окружения SERVER_ADDRESS присутствует непустое значение,
	// переопределим адрес запуска сервера,
	// даже если он был передан через аргумент командной строки
	if envRunAddr := os.Getenv("SERVER_ADDRESS"); envRunAddr != "" {
		FlagRunAddr = envRunAddr
	}

	// для случаев, когда в переменной окружения BASE_URL присутствует непустое значение,
	// переопределим адрес запуска сервера,
	// даже если он был передан через аргумент командной строки
	if envBaseAddr := os.Getenv("BASE_URL"); envBaseAddr != "" {
		FlagBaseAddr = envBaseAddr
	}

	// для случаев, когда в переменной окружения LOG_LEVEL присутствует непустое значение,
	// переопределим уровень логирования,
	// даже если он был передан через аргумент командной строки
	if envLogLevel := os.Getenv("LOG_LEVEL"); envLogLevel != "" {
		FlagLogLevel = envLogLevel
	}

	// для случаев, когда в переменной окружения FILE_STORAGE_PATH присутствует непустое значение,
	// переопределим уровень логирования,
	// даже если он был передан через аргумент командной строки
	if envFileStoragePath, isExist := os.LookupEnv("FILE_STORAGE_PATH"); isExist {
		FlagFileStoragePath = envFileStoragePath
	}

	// для случаев, когда в переменной окружения DATABASE_DSN присутствует непустое значение,
	// переопределим подключение для бд,
	// даже если он был передан через аргумент командной строки
	if envDatabaseDSN, isExist := os.LookupEnv("DATABASE_DSN"); isExist {
		FlagDatabaseDSN = envDatabaseDSN
	}
}
