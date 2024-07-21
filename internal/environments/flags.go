// Package environments configuration
package environments

import (
	"flag"
	"os"

	"github.com/joho/godotenv"
	"github.com/labstack/gommon/log"
)

// FlagRunAddr содержит адрес и порт для запуска сервера
var FlagRunAddr string

// FlagBaseAddr содержит базовый адрес результирующего сокращённого URL
var FlagBaseAddr string

// FlagLogLevel содержит уровень логгирования
var FlagLogLevel string

// FlagFileStoragePath содержит путь до файла хранения
var FlagFileStoragePath string

// FlagDatabaseDSN содержит путь до бд
var FlagDatabaseDSN string

// EnableHTTPS включен ли HTTPS
var EnableHTTPS bool

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

	// регистрируем переменную EnableHTTPS
	// как аргумент -s с ложным значением по умолчанию
	flag.Bool("s", EnableHTTPS, "enable https")

	// парсим переданные серверу аргументы в зарегистрированные переменные
	flag.Parse()

	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}

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

	// для случаев, когда в переменной окружения ENABLE_HTTPS присутствует значение,
	// переопределим включение HTTPS,
	// даже если он был передан через аргумент командной строки
	if envEnableHTTPS, isExist := os.LookupEnv("ENABLE_HTTPS"); isExist {
		EnableHTTPS = envEnableHTTPS == "true" || envEnableHTTPS == "1"
	}
}
