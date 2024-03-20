package main

import (
	"context"
	"encoding/json"
	"github.com/ShukinDmitriy/shortener/internal/models"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
	"io"
	"net/http"
	"reflect"
)

type IPgxConn interface {
	Ping(context.Context) error
}

type URLShortener struct {
	URLRepository models.URLRepository
	conn          IPgxConn
}

func newURLShortener(
	urlRepository models.URLRepository,
	conn IPgxConn,
) *URLShortener {
	return &URLShortener{
		URLRepository: urlRepository,
		conn:          conn,
	}
}

func (us *URLShortener) HandleShorten(ctx echo.Context) error {
	originalURL, err := io.ReadAll(ctx.Request().Body)

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
	shortKey := generateShortKey()

	err = us.URLRepository.Save([]models.Event{{
		ShortKey:    shortKey,
		OriginalURL: string(originalURL),
	}})
	if err != nil {
		ctx.Logger().Error(err)
		return echo.NewHTTPError(http.StatusBadRequest, "can't save url. internal error")
	}

	result := prepareFullURL(shortKey, ctx)

	ctx.Response().Header().Set("Content-Type", "text/plain; charset=utf-8")

	return ctx.String(http.StatusCreated, result)
}

func (us *URLShortener) HandleCreateShorten(ctx echo.Context) error {
	// десериализуем запрос в структуру модели
	zap.L().Debug("decoding request")
	var req models.CreateRequest
	dec := json.NewDecoder(ctx.Request().Body)
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
	shortKey := generateShortKey()

	err := us.URLRepository.Save([]models.Event{{
		ShortKey:    shortKey,
		OriginalURL: req.URL,
	}})
	if err != nil {
		ctx.Logger().Error(err)
		return echo.NewHTTPError(http.StatusBadRequest, "can't save url. internal error")
	}

	// заполняем модель ответа
	resp := models.CreateResponse{
		Result: prepareFullURL(shortKey, ctx),
	}

	return ctx.JSON(http.StatusCreated, resp)
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
	var events []models.Event
	// заполняем модель ответа
	var resp []models.CreateResponseBatch

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
		shortKey := generateShortKey()

		events = append(events, models.Event{
			ShortKey:      shortKey,
			OriginalURL:   cr.OriginalURL,
			CorrelationID: cr.CorrelationID,
		})

		resp = append(resp, models.CreateResponseBatch{
			CorrelationID: cr.CorrelationID,
			ShortURL:      shortKey,
		})
	}

	us.URLRepository.Save(events)

	return ctx.JSON(http.StatusCreated, resp)
}

func (us *URLShortener) HandleRedirect(ctx echo.Context) error {
	shortKey := ctx.Param("id")

	if shortKey == "" {
		ctx.Logger().Error("empty id")
		return echo.NewHTTPError(http.StatusBadRequest, "")
	}

	// Retrieve the original URL from the `urls` map using the shortened key
	originalURL, found := us.URLRepository.Get(shortKey)
	if !found {
		err := "URL not found"
		ctx.Logger().Error(err)
		return echo.NewHTTPError(http.StatusNotFound, err)
	}

	return ctx.Redirect(http.StatusTemporaryRedirect, originalURL)
}

func (us *URLShortener) HandlePing(ctx echo.Context) error {
	if reflect.ValueOf(us.conn).IsNil() {
		ctx.Logger().Error("No connect to db")
		return echo.NewHTTPError(http.StatusInternalServerError, "Internal Server Error")
	}
	err := us.conn.Ping(context.Background())

	if err != nil {
		ctx.Logger().Error("Lost connect to db")
		return echo.NewHTTPError(http.StatusInternalServerError, "Internal Server Error")
	}

	return ctx.String(http.StatusOK, "OK")
}
