package environments_test

import (
	"os"
	"testing"

	"github.com/ShukinDmitriy/shortener/internal/environments"

	"github.com/stretchr/testify/assert"
)

func TestParseFlags(t *testing.T) {
	// Установка переменных окружения для тестирования
	os.Setenv("SERVER_ADDRESS", ":8081")
	os.Setenv("BASE_URL", "http://localhost:8081/")
	os.Setenv("LOG_LEVEL", "debug")
	os.Setenv("FILE_STORAGE_PATH", "/tmp/test-file.json")
	os.Setenv("DATABASE_DSN", "postgres://user:password@host:port/dbname")

	// Вызов функции ParseFlags
	environments.ParseFlags()

	// Проверка установленных переменных
	assert.Equal(t, ":8081", environments.FlagRunAddr)
	assert.Equal(t, "http://localhost:8081/", environments.FlagBaseAddr)
	assert.Equal(t, "debug", environments.FlagLogLevel)
	assert.Equal(t, "/tmp/test-file.json", environments.FlagFileStoragePath)
	assert.Equal(t, "postgres://user:password@host:port/dbname", environments.FlagDatabaseDSN)

	// Очистка переменных окружения
	os.Unsetenv("SERVER_ADDRESS")
	os.Unsetenv("BASE_URL")
	os.Unsetenv("LOG_LEVEL")
	os.Unsetenv("FILE_STORAGE_PATH")
	os.Unsetenv("DATABASE_DSN")
}
