package app

import (
	"context"
	"errors"
	"github.com/ShukinDmitriy/shortener/internal/models"
	pb "github.com/ShukinDmitriy/shortener/proto"
	"go.uber.org/zap"
	"net"
)

// URLShortenerGRPC the application
type URLShortenerGRPC struct {
	pb.UnimplementedURLServer

	URLRepository models.URLRepository
	conn          PgxConnPinger

	// разрешенная подсеть
	subnet *net.IPNet
}

// NewURLShortenerGRPC application's constructor
func NewURLShortenerGRPC(
	urlRepository models.URLRepository,
	conn PgxConnPinger,
	subnet *net.IPNet,
) *URLShortenerGRPC {
	instance := &URLShortenerGRPC{
		URLRepository: urlRepository,
		conn:          conn,
		subnet:        subnet,
	}

	return instance
}

// Create handler for create short link
func (us *URLShortenerGRPC) Create(ctx context.Context, req *pb.CreateRequest) (*pb.CreateResponse, error) {
	// проверяем, что пришёл запрос понятного типа
	if string(req.OriginalUrl) == "" {
		zap.L().Debug("unsupported request url", zap.String("url", req.OriginalUrl))
		return &pb.CreateResponse{
			Status: "bad request",
		}, nil
	}
	// Generate a unique shortened key for the original URL
	shortKey := models.GenerateShortKey()

	status := "created"
	events := []*models.Event{{
		ShortKey:    shortKey,
		OriginalURL: req.OriginalUrl,
		UserID:      req.UserId,
	}}

	err := us.URLRepository.Save(ctx, events)
	if err != nil {
		zap.L().Error(err.Error())

		if errors.Is(err, models.ErrURLExist) {
			status = "conflict"
			shortKey = events[0].ShortKey
		} else {
			zap.L().Debug("can't save url. internal error", zap.String("url", req.OriginalUrl))
			return &pb.CreateResponse{
				Status: "bad request",
			}, nil
		}
	}

	// заполняем модель ответа
	return &pb.CreateResponse{
		Status:      status,
		ResponseUrl: models.PrepareFullURL(shortKey, ""),
	}, nil

}

func (us *URLShortenerGRPC) CreateBatch(ctx context.Context, req *pb.CreateBatchRequest) (*pb.CreateBatchResponse, error) {
	// События для сохранения
	events := make([]*models.Event, len(req.Urls))

	for i, cr := range req.Urls {
		// проверяем, что пришёл запрос понятного типа
		if cr.OriginalUrl == "" || cr.CorrelationId == "" {
			err := "empty original_url or correlation_id"
			zap.L().Error(err)
			zap.L().Debug(
				"unsupported request url",
				zap.String("original_url", cr.OriginalUrl),
				zap.String("correlation_id", cr.CorrelationId),
			)
			return &pb.CreateBatchResponse{
				Status: "bad request",
			}, nil
		}

		// Generate a unique shortened key for the original URL
		shortKey := models.GenerateShortKey()

		events[i] = &models.Event{
			ShortKey:      shortKey,
			OriginalURL:   cr.OriginalUrl,
			CorrelationID: cr.CorrelationId,
			UserID:        req.UserId,
		}
	}

	status := "created"
	err := us.URLRepository.Save(ctx, events)
	if err != nil {
		zap.L().Error(err.Error())

		if errors.Is(err, models.ErrURLExist) {
			status = "conflict"
		} else {
			return &pb.CreateBatchResponse{
				Status: "bad request",
			}, nil
		}
	}

	// заполняем модель ответа
	resp := make([]*pb.CreateBatchResponse_URL, len(events))
	for i, event := range events {
		resp[i] = &pb.CreateBatchResponse_URL{
			CorrelationId: event.CorrelationID,
			ShortUrl:      models.PrepareFullURL(event.ShortKey, ""),
		}
	}

	// заполняем модель ответа
	return &pb.CreateBatchResponse{
		Status: status,
		Urls:   resp,
	}, nil
}

func (us *URLShortenerGRPC) Redirect(_ context.Context, req *pb.RedirectRequest) (*pb.RedirectResponse, error) {
	if req.ShortUrl == "" {
		zap.L().Error("empty id")
		return &pb.RedirectResponse{
			Status: "bad request",
		}, nil
	}

	// Retrieve the original URL from the `urls` map using the shortened key
	event, found := us.URLRepository.Get(req.ShortUrl)
	if !found {
		err := "URL not found"
		zap.L().Error(err)
		return &pb.RedirectResponse{
			Status: "not found",
		}, nil
	}

	if event.DeletedFlag {
		err := "URL deleted"
		zap.L().Error(err)
		return &pb.RedirectResponse{
			Status: "gone",
		}, nil
	}

	return &pb.RedirectResponse{
		Status:      "ok",
		RedirectUrl: event.OriginalURL,
	}, nil
}

func (us *URLShortenerGRPC) GetUserURLs(ctx context.Context, req *pb.GetUserURLsRequest) (*pb.GetUserURLsResponse, error) {
	if req.UserId == "" {
		zap.L().Error("empty user_id")
		return &pb.GetUserURLsResponse{
			Status: "no content",
		}, nil
	}

	events := us.URLRepository.GetEventsByUserID(ctx, req.UserId)

	// заполняем модель ответа
	resp := make([]*pb.GetUserURLsResponse_URL, len(events))

	for i, event := range events {
		resp[i] = &pb.GetUserURLsResponse_URL{
			ShortUrl:    models.PrepareFullURL(event.ShortKey, ""),
			OriginalUrl: event.OriginalURL,
		}
	}

	return &pb.GetUserURLsResponse{
		Status: "ok",
		Urls:   resp,
	}, nil
}

func (us *URLShortenerGRPC) DeleteBatch(_ context.Context, req *pb.DeleteBatchRequest) (*pb.DeleteBatchResponse, error) {
	zap.L().Info("delete", zap.Any("req", req))

	go func() {
		err := us.URLRepository.Delete(context.TODO(), []models.DeleteRequestBatch{
			{
				UserID:    req.UserId,
				ShortKeys: req.Urls,
			},
		})
		if err != nil {
			zap.L().Error("cannot delete events", zap.String("err", err.Error()))
		}
	}()

	return &pb.DeleteBatchResponse{
		Status: "accepted",
	}, nil
}

func (us *URLShortenerGRPC) GetStats(ctx context.Context, req *pb.GetStatsRequest) (*pb.GetStatsResponse, error) {
	// Проверяем доступ
	ip := net.ParseIP(req.IpAddress)
	if us.subnet == nil || !us.subnet.Contains(ip) {
		return &pb.GetStatsResponse{
			Status: "forbidden",
		}, nil
	}

	countUser, countURL, err := us.URLRepository.GetStats(ctx)
	if err != nil {
		zap.L().Error(err.Error())
		return &pb.GetStatsResponse{
			Status: "internal server error",
		}, nil
	}

	return &pb.GetStatsResponse{
		Status: "ok",
		Users:  int32(countUser),
		Urls:   int32(countURL),
	}, nil
}
