// Package environments configuration
package environments

import (
	"encoding/json"
	"flag"
	"io"
	"os"

	"github.com/joho/godotenv"
	"github.com/labstack/gommon/log"
)

// configFile структура с конфигурацией в файле
type configFile struct {
	ServerAddress   string `json:"server_address"`
	BaseURL         string `json:"base_url"`
	LogLevel        string `json:"log_level"`
	FileStoragePath string `json:"file_storage_path"`
	DatabaseDSN     string `json:"database_dsn"`
	EnableHTTPS     bool   `json:"enable_https"`
	TrustedSubnet   string `json:"trusted_subnet"`
}

// getConfigFromFile Чтение конфигурации из файла
func getConfigFromFile(fileName string) configFile {
	f, err := os.Open(fileName)
	if err != nil {
		log.Fatal(err)
	}

	data, err := io.ReadAll(f)
	if err != nil {
		log.Fatal(err)
	}

	cfg := configFile{}
	err = json.Unmarshal(data, &cfg)
	if err != nil {
		log.Fatal(err)
	}

	return cfg
}

// Configuration структура с конфигурацией приложения
type Configuration struct {
	RunAddr         string
	BaseAddr        string
	LogLevel        string
	FileStoragePath string
	DatabaseDSN     string
	EnableHTTPS     bool
	TrustedSubnet   string
}

// flagConfig содержит путь к файлу конфигурации в формате JSON
var flagConfig string

// flagRunAddr содержит адрес и порт для запуска сервера
var flagRunAddr string

// BaseAddr содержит базовый адрес результирующего сокращённого URL
var BaseAddr string

// flagBaseAddr содержит базовый адрес результирующего сокращённого URL
var flagBaseAddr string

// flagLogLevel содержит уровень логгирования
var flagLogLevel string

// flagFileStoragePath содержит путь до файла хранения
var flagFileStoragePath string

// flagDatabaseDSN содержит путь до бд
var flagDatabaseDSN string

// flagEnableHTTPS включен ли HTTPS
var flagEnableHTTPS bool

// flagTrustedSubnet строковое представление бесклассовой адресации
var flagTrustedSubnet string

// ParseFlags обрабатывает аргументы командной строки
// и сохраняет их значения в соответствующих переменных
func ParseFlags() Configuration {
	// регистрируем переменную flagConfig
	// как аргумент -с с пустым значением по умолчанию
	if flag.Lookup("c") == nil {
		flag.StringVar(&flagConfig, "c", flagConfig, "path to config JSON file")
	}

	// регистрируем переменную flagConfig
	// как аргумент -config с пустым значением по умолчанию
	if flag.Lookup("config") == nil {
		flag.StringVar(&flagConfig, "config", flagConfig, "path to config JSON file")
	}

	// регистрируем переменную flagRunAddr
	// как аргумент -a со значением :8080 по умолчанию
	if flag.Lookup("a") == nil {
		flag.StringVar(&flagRunAddr, "a", "127.0.0.1:8080", "address and port to run server")
	}

	// регистрируем переменную flagBaseAddr
	// как аргумент -b со значением :8000 по умолчанию
	if flag.Lookup("b") == nil {
		flag.StringVar(&flagBaseAddr, "b", "http://127.0.0.1:8080", "base server address and port")
	}

	// регистрируем переменную flagLogLevel
	// как аргумент -l со значением info по умолчанию
	if flag.Lookup("l") == nil {
		flag.StringVar(&flagLogLevel, "l", "info", "log level")
	}

	// регистрируем переменную flagFileStoragePath
	// как аргумент -f с пустым значением по умолчанию
	if flag.Lookup("f") == nil {
		flag.StringVar(&flagFileStoragePath, "f", "/tmp/short-url-db.json", "db file path")
	}

	// регистрируем переменную flagDatabaseDSN
	// как аргумент -d с пустым значением по умолчанию
	if flag.Lookup("d") == nil {
		flag.StringVar(&flagDatabaseDSN, "d", "", "database DSN")
	}

	// регистрируем переменную flagEnableHTTPS
	// как аргумент -s с ложным значением по умолчанию
	if flag.Lookup("s") == nil {
		flag.BoolVar(&flagEnableHTTPS, "s", false, "enable https")
	}

	// регистрируем переменную flagEnableHTTPS
	// как аргумент -t с пустым значением по умолчанию
	if flag.Lookup("t") == nil {
		flag.StringVar(&flagTrustedSubnet, "t", "", "CIDR")
	}

	// парсим переданные серверу аргументы в зарегистрированные переменные
	flag.Parse()

	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}

	// для случаев, когда в переменной окружения CONFIG присутствует непустое значение,
	// переопределим путь к файлу конфигурации в формате JSON,
	// даже если он был передан через аргумент командной строки
	if envConfig := os.Getenv("CONFIG"); envConfig != "" {
		flagConfig = envConfig
	}

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

	// для случаев, когда в переменной окружения DATABASE_DSN присутствует непустое значение,
	// переопределим подключение для бд,
	// даже если он был передан через аргумент командной строки
	if envDatabaseDSN, isExist := os.LookupEnv("DATABASE_DSN"); isExist {
		flagDatabaseDSN = envDatabaseDSN
	}

	// для случаев, когда в переменной окружения ENABLE_HTTPS присутствует значение,
	// переопределим включение HTTPS,
	// даже если он был передан через аргумент командной строки
	if envEnableHTTPS, isExist := os.LookupEnv("ENABLE_HTTPS"); isExist {
		flagEnableHTTPS = envEnableHTTPS == "true" || envEnableHTTPS == "1"
	}

	// для случаев, когда в переменной окружения TRUSTED_SUBNET присутствует значение,
	// переопределим строковое представление бесклассовой адресации,
	// даже если он был передан через аргумент командной строки
	if envTrustedSubnet, isExist := os.LookupEnv("TRUSTED_SUBNET"); isExist {
		flagTrustedSubnet = envTrustedSubnet
	}

	fileConfig := configFile{}
	if flagConfig != "" {
		fileConfig = getConfigFromFile(flagConfig)
	}
	flagConfig = ""

	configuration := Configuration{}
	if configuration.RunAddr = flagRunAddr; configuration.RunAddr == "" {
		configuration.RunAddr = fileConfig.ServerAddress
	}
	if configuration.BaseAddr = flagBaseAddr; configuration.BaseAddr == "" {
		configuration.BaseAddr = fileConfig.BaseURL
	}
	BaseAddr = configuration.BaseAddr
	if configuration.LogLevel = flagLogLevel; configuration.LogLevel == "" {
		configuration.LogLevel = fileConfig.LogLevel
	}
	if configuration.FileStoragePath = flagFileStoragePath; configuration.FileStoragePath == "" {
		configuration.FileStoragePath = fileConfig.FileStoragePath
	}
	if configuration.DatabaseDSN = flagDatabaseDSN; configuration.DatabaseDSN == "" {
		configuration.DatabaseDSN = fileConfig.DatabaseDSN
	}
	if configuration.EnableHTTPS = flagEnableHTTPS; !configuration.EnableHTTPS {
		configuration.EnableHTTPS = fileConfig.EnableHTTPS
	}
	if configuration.TrustedSubnet = flagTrustedSubnet; configuration.TrustedSubnet == "" {
		configuration.TrustedSubnet = fileConfig.TrustedSubnet
	}

	return configuration
}
