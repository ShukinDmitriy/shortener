package models

import (
	"context"
	"database/sql"
	"github.com/ShukinDmitriy/shortener/internal/environments"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5"
	"go.uber.org/zap"
	"os"
	"path"
)

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

	m, _ := migrate.NewWithDatabaseInstance(
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

// Из-за того что тесты на 11 итерацию не проходят с новым полем
// correlation_id этот костыль с 2 одинаковыми методами
func (r *PGURLRepository) Save(events []Event) error {
	ctx := context.Background()
	tx, err := r.conn.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			tx.Rollback(ctx)
		} else {
			tx.Commit(ctx)
		}
	}()

	for _, event := range events {
		_, err = tx.Exec(
			ctx,
			`INSERT INTO public.url (short_key, original_url) VALUES ($1, $2)`,
			event.ShortKey, event.OriginalURL,
		)

		if err != nil {
			zap.L().Error(err.Error())
			return err
		}
	}

	return err
}

func (r *PGURLRepository) SaveBatch(events []Event) error {
	ctx := context.Background()
	tx, err := r.conn.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			tx.Rollback(ctx)
		} else {
			tx.Commit(ctx)
		}
	}()

	for _, event := range events {
		_, err = tx.Exec(
			ctx,
			`INSERT INTO public.url (short_key, original_url, correlation_id) VALUES ($1, $2, $3)`,
			event.ShortKey, event.OriginalURL, event.CorrelationID,
		)

		if err != nil {
			zap.L().Error(err.Error())
			return err
		}
	}

	return err
}
