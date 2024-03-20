package models

import (
	"context"
	"database/sql"
	"errors"
	"github.com/ShukinDmitriy/shortener/internal/environments"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"go.uber.org/zap"
	"os"
	"path"
)

var ErrURLExist = errors.New("URL exist")

type PGURLRepository struct {
	conn *pgx.Conn
}

func (r *PGURLRepository) Initialize() error {
	cont := context.Background()
	var conn *pgx.Conn
	var err error

	conn, err = pgx.Connect(cont, environments.FlagDatabaseDSN)
	if err != nil {
		return err
	}
	r.conn = conn

	// Миграции TODO вынести в отдельный файл
	currentDir, _ := os.Getwd()
	zap.L().Info("current dir", zap.String("currentDir", currentDir))
	db, err := sql.Open("postgres", environments.FlagDatabaseDSN)
	if err != nil {
		zap.L().Error("can't connect to db", zap.String("err", err.Error()))
		return err
	}
	defer func() {
		db.Close()
	}()

	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		zap.L().Error("can't create driver", zap.String("err", err.Error()))
		return err
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file:///"+path.Join(currentDir, "db", "migrations"),
		"postgres", driver)
	if err != nil {
		zap.L().Error("can't create new migrate", zap.String("err", err.Error()))
		return err
	}

	err = m.Up()
	if err != nil {
		zap.L().Info("can't migrate up", zap.String("err", err.Error()))
	}

	zap.L().Info("migrate runned")

	return nil
}

func (r *PGURLRepository) Get(shortKey string) (string, bool) {
	var originalURL string

	row := r.conn.QueryRow(
		context.Background(),
		`SELECT original_url from public.url WHERE short_key = $1;`,
		shortKey,
	)

	err := row.Scan(&originalURL)
	if err != nil {
		zap.L().Error(err.Error())
	}

	return originalURL, err == nil && originalURL != ""
}

func (r *PGURLRepository) Save(events []*Event) error {
	ctx := context.Background()

	errs := []error{}

	for _, event := range events {
		_, err := r.conn.Exec(
			ctx,
			`INSERT INTO public.url (short_key, original_url, correlation_id)
VALUES ($1, $2, $3);`,
			event.ShortKey, event.OriginalURL, event.CorrelationID,
		)

		if err != nil {
			zap.L().Error(err.Error())

			// проверяем, что ошибка сигнализирует о наличие данных в БД
			var pgErr *pgconn.PgError
			if errors.As(err, &pgErr) && pgerrcode.IsIntegrityConstraintViolation(pgErr.Code) {
				errs = append(errs, ErrURLExist)
				shortKey, _ := r.GetShortKeyByOriginalURL(event.OriginalURL)
				event.ShortKey = shortKey
			} else {
				return err
			}
		}
	}

	return errors.Join(errs...)
}

func (r *PGURLRepository) GetShortKeyByOriginalURL(originalURL string) (string, bool) {
	var shortKey string

	row := r.conn.QueryRow(
		context.Background(),
		`SELECT short_key from public.url WHERE original_url = $1;`,
		originalURL,
	)

	err := row.Scan(&shortKey)
	if err != nil {
		zap.L().Error(err.Error())
	}

	return shortKey, err == nil && shortKey != ""
}
