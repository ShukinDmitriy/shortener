package models

import (
	"context"
	"database/sql"
	"errors"
	"os"
	"path"
	"strings"

	"github.com/ShukinDmitriy/shortener/internal/environments"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

// ErrURLExist default error
var ErrURLExist = errors.New("URL exist")

// PGURLRepository repository for working with a database
type PGURLRepository struct {
	pool *pgxpool.Pool
}

// Initialize repository
func (r *PGURLRepository) Initialize(configuration environments.Configuration) error {
	cont := context.Background()
	var pool *pgxpool.Pool
	var err error

	pool, err = pgxpool.New(cont, configuration.DatabaseDSN)
	if err != nil {
		return err
	}
	r.pool = pool

	db, err := sql.Open("postgres", configuration.DatabaseDSN)
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

	currentDir, _ := os.Getwd()
	currentDir = strings.TrimSuffix(currentDir, "/internal/models")
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

// Get event by short key
func (r *PGURLRepository) Get(shortKey string) (Event, bool) {
	var originalURL string
	var isDeleted bool

	row := r.pool.QueryRow(
		context.Background(),
		`SELECT original_url, is_deleted from public.url WHERE short_key = $1;`,
		shortKey,
	)

	err := row.Scan(&originalURL, &isDeleted)
	if err != nil {
		zap.L().Error(err.Error())
	}

	return Event{
		ShortKey:    shortKey,
		OriginalURL: originalURL,
		DeletedFlag: isDeleted,
	}, err == nil && originalURL != ""
}

// Save batch save events
func (r *PGURLRepository) Save(ctx context.Context, events []*Event) error {
	var errs []error

	for _, event := range events {
		_, err := r.pool.Exec(
			ctx,
			`INSERT INTO public.url (short_key, original_url, correlation_id, user_id)
VALUES ($1, $2, $3, $4);`,
			event.ShortKey, event.OriginalURL, event.CorrelationID, event.UserID,
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

// Delete batch delete event
func (r *PGURLRepository) Delete(ctx context.Context, events []DeleteRequestBatch) error {
	var errs []error
	var shortKeys []string

	for _, deletedEvent := range events {
		shortKeys = append(shortKeys, deletedEvent.ShortKeys...)

		_, err := r.pool.Exec(
			ctx,
			`UPDATE public.url SET is_deleted = true WHERE user_id = $1 AND short_key = any($2);`,
			deletedEvent.UserID,
			shortKeys,
		)
		if err != nil {
			zap.L().Error(err.Error())
			errs = append(errs, err)
		}
	}

	return errors.Join(errs...)
}

// GetShortKeyByOriginalURL get short link from full link
func (r *PGURLRepository) GetShortKeyByOriginalURL(originalURL string) (string, bool) {
	var shortKey string

	row := r.pool.QueryRow(
		context.Background(),
		`SELECT short_key from public.url WHERE original_url = $1 and is_deleted is false;`,
		originalURL,
	)

	err := row.Scan(&shortKey)
	if err != nil {
		zap.L().Error(err.Error())
	}

	return shortKey, err == nil && shortKey != ""
}

// GetEventsByUserID get events by user ID
func (r *PGURLRepository) GetEventsByUserID(ctx context.Context, userID string) []*Event {
	var events []*Event

	rows, err := r.pool.Query(
		ctx,
		`SELECT short_key, original_url from public.url WHERE user_id = $1 and is_deleted is false;`,
		userID,
	)
	if err != nil {
		zap.L().Error(err.Error())
		return events
	}

	for rows.Next() {
		var shortKey string
		var originalURL string

		err := rows.Scan(&shortKey, &originalURL)
		if err != nil {
			zap.L().Error(err.Error())
			return events
		}

		events = append(events, &Event{
			ShortKey:    shortKey,
			OriginalURL: originalURL,
		})
	}

	if rows.Err() != nil {
		zap.L().Error(rows.Err().Error())
	}

	return events
}
