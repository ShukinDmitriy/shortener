# shortener

## Запуск тестов

Для запуска unit тестов необходимо запустить docker контейнер с бд

```docker compose -f test.docker-compose.yml -p shortener start shortener-test-postgres```

## Запуск подсчета покрытия кода тестами

Для запуска подсчета покрытия кода тестами можно выполнить

```shell
make test-cover-cli
```

или напрямую
```shell
./test-cover.sh
```

Из расчета покрытия исключены (файл .covignore)
- test_spec.go - пустой файл, необходим для запуска mockery
- файлы моков
- файлы cmd/staticlint - инструмент статического анализа кода
- файлы cmd/shortener/main.go - входная точка проекта, настройка роутинга и middleware
- файлы cmd/shortener/prof.go - профилирование