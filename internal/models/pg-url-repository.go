package models

import (
	"context"
	"github.com/ShukinDmitriy/shortener/internal/environments"
	"github.com/jackc/pgx/v5"
	"go.uber.org/zap"
)

type PGURLRepository struct {
	conn *pgx.Conn
}

func (r *PGURLRepository) Initialize() error {
	// urlExample := "postgres://username:password@localhost:5432/database_name"
	cont := context.Background()
	var conn *pgx.Conn
	var err error

	conn, err = pgx.Connect(cont, environments.FlagDatabaseDSN)
	if err != nil {
		return err
	}

	r.conn = conn

	// Создаем БД и таблицу если их нет (TODO: по идее это делается в миграциях, но таковы требования)
	_, err = r.conn.Exec(
		cont,
		`create table if not exists public.url
(
    short_key    varchar not null
        constraint url_pk
            primary key
        constraint url_pk_2
            unique,
    original_url varchar
);`,
	)
	if err != nil {
		return err
	}

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
