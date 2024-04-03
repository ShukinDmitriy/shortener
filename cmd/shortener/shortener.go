package main

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/ShukinDmitriy/shortener/internal/auth"
	"github.com/ShukinDmitriy/shortener/internal/models"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
	"io"
	"net/http"
	"reflect"
	"time"
)

type PgxConnPinger interface {
	Ping(context.Context) error
}

type URLShortener struct {
	URLRepository models.URLRepository
	conn          PgxConnPinger

	// канал для отложенного удаления
	eDeletedEvent chan models.DeleteRequestBatch

	// канал для уведомления об окончании работы
	shutdownChan chan chan struct{}
}

func newURLShortener(
	urlRepository models.URLRepository,
	conn PgxConnPinger,
) *URLShortener {
	instance := &URLShortener{
		URLRepository: urlRepository,
		conn:          conn,
		eDeletedEvent: make(chan models.DeleteRequestBatch, 100),
		shutdownChan:  make(chan chan struct{}),
	}

	go instance.deleteEvents()

	return instance
}

func (us *URLShortener) HandleShorten(ctx echo.Context) error {
	originalURL, err := io.ReadAll(ctx.Request().Body)
	defer ctx.Request().Body.Close()

	if err != nil {
		ctx.Logger().Error(err)
		return echo.NewHTTPError(http.StatusBadRequest, "can't read body. internal error")
	}

	if string(originalURL) == "" {
		err := "empty url"
		ctx.Logger().Error(err)
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	// Generate a unique shortened key for the original URL
	shortKey := models.GenerateShortKey()

	status := http.StatusCreated
	events := []*models.Event{{
		ShortKey:    shortKey,
		OriginalURL: string(originalURL),
		UserID:      auth.GetUserID(ctx),
	}}
	err = us.URLRepository.Save(ctx.Request().Context(), events)

	if errors.Is(err, models.ErrURLExist) {
		ctx.Logger().Error(err)
		shortKey = events[0].ShortKey
		return ctx.String(http.StatusConflict, models.PrepareFullURL(ctx, shortKey))
	}

	if err != nil {
		ctx.Logger().Error(err)
		return echo.NewHTTPError(http.StatusBadRequest, "can't save url. internal error1"+err.Error())
	}

	result := models.PrepareFullURL(ctx, shortKey)

	ctx.Response().Header().Set("Content-Type", "text/plain; charset=utf-8")

	return ctx.String(status, result)
}

func (us *URLShortener) HandleCreateShorten(ctx echo.Context) error {
	// десериализуем запрос в структуру модели
	zap.L().Debug("decoding request")
	var req models.CreateRequest
	dec := json.NewDecoder(ctx.Request().Body)
	defer ctx.Request().Body.Close()

	if err := dec.Decode(&req); err != nil {
		zap.L().Debug("cannot decode request JSON body", zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, "invalid JSON")
	}

	// проверяем, что пришёл запрос понятного типа
	if string(req.URL) == "" {
		err := "empty url"
		ctx.Logger().Error(err)
		zap.L().Debug("unsupported request url", zap.String("url", req.URL))
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	// Generate a unique shortened key for the original URL
	shortKey := models.GenerateShortKey()

	status := http.StatusCreated
	events := []*models.Event{{
		ShortKey:    shortKey,
		OriginalURL: req.URL,
		UserID:      auth.GetUserID(ctx),
	}}

	err := us.URLRepository.Save(ctx.Request().Context(), events)
	if err != nil {
		ctx.Logger().Error(err)

		if errors.Is(err, models.ErrURLExist) {
			status = http.StatusConflict
			shortKey = events[0].ShortKey
		} else {
			return echo.NewHTTPError(http.StatusBadRequest, "can't save url. internal error")
		}
	}

	// заполняем модель ответа
	resp := models.CreateResponse{
		Result: models.PrepareFullURL(ctx, shortKey),
	}

	return ctx.JSON(status, resp)
}

func (us *URLShortener) HandleCreateShortenBatch(ctx echo.Context) error {
	// десериализуем запрос в структуру модели
	zap.L().Debug("decoding request")
	var req []models.CreateRequestBatch
	dec := json.NewDecoder(ctx.Request().Body)
	if err := dec.Decode(&req); err != nil {
		zap.L().Debug("cannot decode request JSON body", zap.Error(err))
		return echo.NewHTTPError(http.StatusInternalServerError, "invalid JSON")
	}

	// События для сохранения
	var events []*models.Event
	// заполняем модель ответа
	var resp []models.CreateResponseBatch

	userID := auth.GetUserID(ctx)

	for _, cr := range req {
		// проверяем, что пришёл запрос понятного типа
		if string(cr.OriginalURL) == "" || string(cr.CorrelationID) == "" {
			err := "empty original_url or correlation_id"
			ctx.Logger().Error(err)
			zap.L().Debug(
				"unsupported request url",
				zap.String("original_url", cr.OriginalURL),
				zap.String("correlation_id", cr.CorrelationID),
			)
			return echo.NewHTTPError(http.StatusBadRequest, err)
		}

		// Generate a unique shortened key for the original URL
		shortKey := models.GenerateShortKey()

		events = append(events, &models.Event{
			ShortKey:      shortKey,
			OriginalURL:   cr.OriginalURL,
			CorrelationID: cr.CorrelationID,
			UserID:        userID,
		})
	}

	status := http.StatusCreated
	err := us.URLRepository.Save(ctx.Request().Context(), events)
	if err != nil {
		ctx.Logger().Error(err)

		if errors.Is(err, models.ErrURLExist) {
			status = http.StatusConflict
		} else {
			return echo.NewHTTPError(http.StatusBadRequest, "can't save url. internal error")
		}
	}

	for _, event := range events {
		resp = append(resp, models.CreateResponseBatch{
			CorrelationID: event.CorrelationID,
			ShortURL:      models.PrepareFullURL(ctx, event.ShortKey),
		})
	}

	return ctx.JSON(status, resp)
}

func (us *URLShortener) HandleRedirect(ctx echo.Context) error {
	shortKey := ctx.Param("id")

	if shortKey == "" {
		ctx.Logger().Error("empty id")
		return echo.NewHTTPError(http.StatusBadRequest, "")
	}

	// Retrieve the original URL from the `urls` map using the shortened key
	event, found := us.URLRepository.Get(shortKey)
	if !found {
		err := "URL not found"
		ctx.Logger().Error(err)
		return echo.NewHTTPError(http.StatusNotFound, err)
	}

	if event.DeletedFlag {
		return ctx.String(http.StatusGone, "")
	}

	return ctx.Redirect(http.StatusTemporaryRedirect, event.OriginalURL)
}

func (us *URLShortener) HandlePing(ctx echo.Context) error {
	if reflect.ValueOf(us.conn).IsNil() {
		ctx.Logger().Error("No connect to db")
		return echo.NewHTTPError(http.StatusInternalServerError, "Internal Server Error")
	}
	err := us.conn.Ping(ctx.Request().Context())

	if err != nil {
		ctx.Logger().Error("Lost connect to db")
		return echo.NewHTTPError(http.StatusInternalServerError, "Internal Server Error")
	}

	return ctx.String(http.StatusOK, "OK")
}

func (us *URLShortener) HandleUserURLGet(ctx echo.Context) error {
	userID := auth.GetUserID(ctx)

	if userID == "" {
		return ctx.JSON(http.StatusNoContent, nil)
	}

	events := us.URLRepository.GetEventsByUserID(ctx.Request().Context(), userID)

	// заполняем модель ответа
	var resp []models.GetUserURLsResponse

	for _, event := range events {
		resp = append(resp, models.GetUserURLsResponse{
			ShortURL:    models.PrepareFullURL(ctx, event.ShortKey),
			OriginalURL: event.OriginalURL,
		})
	}

	return ctx.JSON(http.StatusOK, resp)
}

func (us *URLShortener) HandleUserURLDelete(ctx echo.Context) error {
	req := models.DeleteRequestBatch{
		UserID:    auth.GetUserID(ctx),
		ShortKeys: []string{},
	}
	dec := json.NewDecoder(ctx.Request().Body)
	defer ctx.Request().Body.Close()

	if err := dec.Decode(&req.ShortKeys); err != nil {
		zap.L().Debug("cannot decode request JSON body", zap.Error(err))
		return echo.NewHTTPError(http.StatusBadRequest, "invalid JSON")
	}

	zap.L().Info("delete", zap.Any("req", req))

	us.eDeletedEvent <- req

	return ctx.JSON(http.StatusAccepted, "Accepted")
}

func (us *URLShortener) Shutdown(ctx context.Context) chan struct{} {
	res := make(chan struct{})

	go func() {
		defer close(res)

		successShutdown := make(chan struct{})
		us.shutdownChan <- successShutdown

		<-successShutdown
		close(successShutdown)
		res <- struct{}{}
	}()

	return res
}

func (us *URLShortener) deleteEvents() {
	var events []models.DeleteRequestBatch
	ticker := time.NewTicker(2 * time.Second)

	for {
		select {
		case event := <-us.eDeletedEvent:
			events = append(events, event)
		case success := <-us.shutdownChan:
			if len(events) == 0 {
				success <- struct{}{}
				return
			}

			// Сброс на диск очереди на удаление
			fileProducer, err := models.NewProducer("/tmp/short-deleted-db.json")
			if err != nil {
				zap.L().Error("create file backup", zap.String("err", err.Error()))
				success <- struct{}{}
				return
			}

			err = fileProducer.WriteEvent(events)
			if err != nil {
				zap.L().Error("backup deleted events", zap.String("err", err.Error()))
			}

			events = nil
			success <- struct{}{}

			return
		case <-ticker.C:
			if len(events) == 0 {
				continue
			}

			copyEvents := make([]models.DeleteRequestBatch, len(events))

			copy(copyEvents, events)

			go func(events []models.DeleteRequestBatch) {
				err := us.URLRepository.Delete(context.TODO(), events)
				if err != nil {
					zap.L().Error("cannot delete events", zap.String("err", err.Error()))
				}
			}(copyEvents)

			events = nil
		}
	}
}
