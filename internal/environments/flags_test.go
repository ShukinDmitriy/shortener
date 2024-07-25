package environments_test

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/ShukinDmitriy/shortener/internal/environments"

	"github.com/stretchr/testify/assert"
)

func TestParseFlagsFromJSON(t *testing.T) {
	filename := "./configuration.json"

	// Установка переменных окружения для тестирования
	os.Setenv("CONFIG", filename)

	// Создание временного файла конфигурации
	conf := map[string]interface{}{
		"server_address":    "127.0.0.1:8080",
		"base_url":          "http://127.0.0.1:8080",
		"log_level":         "info",
		"file_storage_path": "/tmp/short-url-db.json",
		"database_dsn":      "postgres://user:password@host:port/dbname",
		"enable_https":      true,
		"trusted_subnet":    "127.0.0.1/24",
	}
	data, err := json.Marshal(conf)
	if err != nil {
		t.Error(err)
	}
	err = os.WriteFile(filename, data, 0o666)
	if err != nil {
		t.Error(err)
	}
	// Удаление временного файла
	defer os.Remove(filename)

	// Вызов функции ParseFlags
	configuration := environments.ParseFlags()

	// Проверка установленных переменных
	assert.Equal(t, "127.0.0.1:8080", configuration.RunAddr)
	assert.Equal(t, "http://127.0.0.1:8080", configuration.BaseAddr)
	assert.Equal(t, "info", configuration.LogLevel)
	assert.Equal(t, "/tmp/short-url-db.json", configuration.FileStoragePath)
	assert.Equal(t, "postgres://user:password@host:port/dbname", configuration.DatabaseDSN)
	assert.Equal(t, true, configuration.EnableHTTPS)
	assert.Equal(t, "127.0.0.1/24", configuration.TrustedSubnet)

	// Очистка переменных окружения
	os.Unsetenv("CONFIG")
}

func TestParseFlags(t *testing.T) {
	// Установка переменных окружения для тестирования
	os.Setenv("SERVER_ADDRESS", "127.0.0.1:8080")
	os.Setenv("BASE_URL", "http://127.0.0.1:8080")
	os.Setenv("LOG_LEVEL", "info")
	os.Setenv("FILE_STORAGE_PATH", "/tmp/short-url-db.json")
	os.Setenv("DATABASE_DSN", "postgres://user:password@host:port/dbname")
	os.Setenv("ENABLE_HTTPS", "true")
	os.Setenv("TRUSTED_SUBNET", "127.0.0.1/24")

	// Вызов функции ParseFlags
	configuration := environments.ParseFlags()

	// Проверка установленных переменных
	assert.Equal(t, "127.0.0.1:8080", configuration.RunAddr)
	assert.Equal(t, "http://127.0.0.1:8080", configuration.BaseAddr)
	assert.Equal(t, "info", configuration.LogLevel)
	assert.Equal(t, "/tmp/short-url-db.json", configuration.FileStoragePath)
	assert.Equal(t, "postgres://user:password@host:port/dbname", configuration.DatabaseDSN)
	assert.Equal(t, true, configuration.EnableHTTPS)
	assert.Equal(t, "127.0.0.1/24", configuration.TrustedSubnet)

	// Очистка переменных окружения
	os.Unsetenv("SERVER_ADDRESS")
	os.Unsetenv("BASE_URL")
	os.Unsetenv("LOG_LEVEL")
	os.Unsetenv("FILE_STORAGE_PATH")
	os.Unsetenv("DATABASE_DSN")
	os.Unsetenv("TRUSTED_SUBNET")
}
