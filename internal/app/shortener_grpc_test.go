package app_test

import (
	"context"
	"errors"
	"log"
	"net"
	"testing"

	"github.com/ShukinDmitriy/shortener/internal/app"
	"github.com/ShukinDmitriy/shortener/internal/environments"
	"github.com/ShukinDmitriy/shortener/internal/models"
	models2 "github.com/ShukinDmitriy/shortener/mocks/internal_/models"
	pb "github.com/ShukinDmitriy/shortener/proto"
	"github.com/pashagolub/pgxmock/v3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/test/bufconn"
)

func dialer(shortenerGRPC *app.URLShortenerGRPC) func(context.Context, string) (net.Conn, error) {
	listener := bufconn.Listen(1024 * 1024)

	server := grpc.NewServer()

	pb.RegisterURLServer(server, shortenerGRPC)

	go func() {
		if err := server.Serve(listener); err != nil {
			log.Fatal(err)
		}
	}()

	return func(context.Context, string) (net.Conn, error) {
		return listener.Dial()
	}
}

func TestURLShortenerGRPC_Create(t *testing.T) {
	type want struct {
		status   string
		response string
		err      *error
	}
	testError := errors.New("test error")
	tests := []struct {
		name string
		want want
		body string
	}{
		{
			name: "positive test #1",
			want: want{
				status:   "created",
				response: "http://example.com/",
			},
			body: "http://example.com/",
		},
		{
			name: "negative test #1",
			want: want{
				status:   "bad request",
				response: "",
			},
			body: "",
		},
		{
			name: "negative test #2",
			want: want{
				status:   "conflict",
				response: "",
				err:      &models.ErrURLExist,
			},
			body: "http://example.com/",
		},
		{
			name: "negative test #3",
			want: want{
				status:   "bad request",
				response: "",
				err:      &testError,
			},
			body: "http://example.com/",
		},
	}

	environments.BaseAddr = "http://example.com"
	repository := new(models2.URLRepository)
	mockConn, _ := pgxmock.NewConn()
	defer mockConn.Close(context.Background())
	_, subnet, err := net.ParseCIDR("127.0.0.1/24")
	if err != nil {
		t.Error(err)
	}
	shortenerGRPC := app.NewURLShortenerGRPC(repository, mockConn, subnet)

	ctx := context.Background()
	conn, err := grpc.NewClient(
		"passthrough://bufnet",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithContextDialer(dialer(shortenerGRPC)),
	)
	if err != nil {
		t.Fatalf("Failed to dial bufnet: %v", err)
	}
	defer conn.Close()

	client := pb.NewURLClient(conn)

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			repository.ExpectedCalls = nil
			if test.want.err != nil {
				repository.EXPECT().Save(
					mock.Anything,
					mock.Anything,
				).Return(*test.want.err)
			} else {
				repository.EXPECT().Save(
					mock.Anything,
					mock.Anything,
				).Return(nil)
			}

			resp, err := client.Create(ctx, &pb.CreateRequest{
				UserId:      "testUser",
				OriginalUrl: test.body,
			})
			if err != nil {
				t.Fatalf("gRPC Create failed: %v", err)
			}

			assert.Equal(t, test.want.status, resp.Status)
			assert.Contains(t, resp.ResponseUrl, test.want.response)
		})
	}
}

func TestURLShortenerGRPC_CreateBatch(t *testing.T) {
	type want struct {
		status string
		err    *error
	}
	testError := errors.New("test error")
	tests := []struct {
		name string
		want want
		body []*pb.CreateBatchRequest_URL
	}{
		{
			name: "positive test #1",
			want: want{
				status: "created",
			},
			body: []*pb.CreateBatchRequest_URL{
				{
					CorrelationId: "847b5414-7f41-4363-be2a-e316fbfc2b33",
					OriginalUrl:   "https://practicum.yandex.ru",
				},
				{
					CorrelationId: "022d3f81-2fb5-4fda-bb19-e89bad595b09",
					OriginalUrl:   "https://yandex.ru",
				},
				{
					CorrelationId: "847b5414-7f41-4363-be2a-e316fbfc2b33",
					OriginalUrl:   "https://music.yandex.ru",
				},
			},
		},
		{
			name: "negative test #1",
			want: want{
				status: "bad request",
			},
			body: []*pb.CreateBatchRequest_URL{
				{
					CorrelationId: "",
					OriginalUrl:   "https://practicum.yandex.ru",
				},
			},
		},
		{
			name: "negative test #2",
			want: want{
				status: "conflict",
				err:    &models.ErrURLExist,
			},
			body: []*pb.CreateBatchRequest_URL{
				{
					CorrelationId: "847b5414-7f41-4363-be2a-e316fbfc2b33",
					OriginalUrl:   "https://practicum.yandex.ru",
				},
				{
					CorrelationId: "022d3f81-2fb5-4fda-bb19-e89bad595b09",
					OriginalUrl:   "https://yandex.ru",
				},
				{
					CorrelationId: "847b5414-7f41-4363-be2a-e316fbfc2b33",
					OriginalUrl:   "https://music.yandex.ru",
				},
			},
		},
		{
			name: "negative test #3",
			want: want{
				status: "bad request",
				err:    &testError,
			},
			body: []*pb.CreateBatchRequest_URL{
				{
					CorrelationId: "847b5414-7f41-4363-be2a-e316fbfc2b33",
					OriginalUrl:   "https://practicum.yandex.ru",
				},
				{
					CorrelationId: "022d3f81-2fb5-4fda-bb19-e89bad595b09",
					OriginalUrl:   "https://yandex.ru",
				},
				{
					CorrelationId: "847b5414-7f41-4363-be2a-e316fbfc2b33",
					OriginalUrl:   "https://music.yandex.ru",
				},
			},
		},
	}

	environments.BaseAddr = "http://example.com"
	repository := new(models2.URLRepository)
	mockConn, _ := pgxmock.NewConn()
	defer mockConn.Close(context.Background())
	_, subnet, err := net.ParseCIDR("127.0.0.1/24")
	if err != nil {
		t.Error(err)
	}
	shortenerGRPC := app.NewURLShortenerGRPC(repository, mockConn, subnet)

	ctx := context.Background()
	conn, err := grpc.NewClient(
		"passthrough://bufnet",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithContextDialer(dialer(shortenerGRPC)),
	)
	if err != nil {
		t.Fatalf("Failed to dial bufnet: %v", err)
	}
	defer conn.Close()

	client := pb.NewURLClient(conn)

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			repository.ExpectedCalls = nil
			if test.want.err != nil {
				repository.EXPECT().Save(
					mock.Anything,
					mock.Anything,
				).Return(*test.want.err)
			} else {
				repository.EXPECT().Save(
					mock.Anything,
					mock.Anything,
				).Return(nil)
			}

			resp, err := client.CreateBatch(ctx, &pb.CreateBatchRequest{
				UserId: "testUser",
				Urls:   test.body,
			})
			if err != nil {
				t.Fatalf("gRPC CreateBatch failed: %v", err)
			}

			assert.Equal(t, test.want.status, resp.Status)
			if test.want.status == "created" {
				assert.Equal(t, len(test.body), len(resp.Urls))
			}
		})
	}
}

func TestURLShortenerGRPC_Redirect(t *testing.T) {
	type want struct {
		status string
		event  models.Event
	}
	tests := []struct {
		name string
		want want
		body string
	}{
		{
			name: "positive test #1",
			want: want{
				status: "ok",
				event: models.Event{
					DeletedFlag: false,
					OriginalURL: "http://example.com/",
				},
			},
			body: "http://example.com/",
		},
		{
			name: "negative test #1",
			want: want{
				status: "bad request",
			},
			body: "",
		},
		{
			name: "negative test #2",
			want: want{
				status: "not found",
			},
			body: "http://example.com/",
		},
		{
			name: "negative test #3",
			want: want{
				status: "gone",
				event: models.Event{
					DeletedFlag: true,
					OriginalURL: "http://example.com/",
				},
			},
			body: "http://example.com/",
		},
	}

	environments.BaseAddr = "http://example.com"
	repository := new(models2.URLRepository)
	mockConn, _ := pgxmock.NewConn()
	defer mockConn.Close(context.Background())
	_, subnet, err := net.ParseCIDR("127.0.0.1/24")
	if err != nil {
		t.Error(err)
	}
	shortenerGRPC := app.NewURLShortenerGRPC(repository, mockConn, subnet)

	ctx := context.Background()
	conn, err := grpc.NewClient(
		"passthrough://bufnet",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithContextDialer(dialer(shortenerGRPC)),
	)
	if err != nil {
		t.Fatalf("Failed to dial bufnet: %v", err)
	}
	defer conn.Close()

	client := pb.NewURLClient(conn)

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			repository.ExpectedCalls = nil
			repository.EXPECT().Get(
				mock.Anything,
			).Return(test.want.event, test.want.event.OriginalURL != "")

			resp, err := client.Redirect(ctx, &pb.RedirectRequest{
				ShortUrl: test.body,
			})
			if err != nil {
				t.Fatalf("gRPC Redirect failed: %v", err)
			}

			assert.Equal(t, test.want.status, resp.Status)
			if resp.Status == "ok" {
				assert.Equal(t, test.want.event.OriginalURL, resp.RedirectUrl)
			}
		})
	}
}

func TestURLShortenerGRPC_GetUserURLs(t *testing.T) {
	type want struct {
		status string
		events []*models.Event
	}
	tests := []struct {
		name string
		want want
		body string
	}{
		{
			name: "positive test #1",
			want: want{
				status: "ok",
				events: []*models.Event{
					{
						DeletedFlag: false,
						OriginalURL: "http://example.com/",
					},
				},
			},
			body: "testUserID",
		},
		{
			name: "negative test #1",
			want: want{
				status: "no content",
				events: []*models.Event{},
			},
			body: "",
		},
	}

	environments.BaseAddr = "http://example.com"
	repository := new(models2.URLRepository)
	mockConn, _ := pgxmock.NewConn()
	defer mockConn.Close(context.Background())
	_, subnet, err := net.ParseCIDR("127.0.0.1/24")
	if err != nil {
		t.Error(err)
	}
	shortenerGRPC := app.NewURLShortenerGRPC(repository, mockConn, subnet)

	ctx := context.Background()
	conn, err := grpc.NewClient(
		"passthrough://bufnet",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithContextDialer(dialer(shortenerGRPC)),
	)
	if err != nil {
		t.Fatalf("Failed to dial bufnet: %v", err)
	}
	defer conn.Close()

	client := pb.NewURLClient(conn)

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			repository.ExpectedCalls = nil
			repository.EXPECT().GetEventsByUserID(
				mock.Anything,
				mock.Anything,
			).Return(test.want.events)

			resp, err := client.GetUserURLs(ctx, &pb.GetUserURLsRequest{
				UserId: test.body,
			})
			if err != nil {
				t.Fatalf("gRPC GetUserURLs failed: %v", err)
			}

			assert.Equal(t, test.want.status, resp.Status)
			if resp.Status == "ok" {
				assert.Equal(t, len(test.want.events), len(resp.Urls))
			}
		})
	}
}

func TestURLShortenerGRPC_DeleteBatch(t *testing.T) {
	type want struct {
		status string
		events []string
	}
	tests := []struct {
		name string
		want want
		body string
	}{
		{
			name: "positive test #1",
			want: want{
				status: "accepted",
				events: []string{"SYqDJ3", "4SwGPJ", "z3e7av"},
			},
		},
	}

	environments.BaseAddr = "http://example.com"
	repository := new(models2.URLRepository)
	mockConn, _ := pgxmock.NewConn()
	defer mockConn.Close(context.Background())
	_, subnet, err := net.ParseCIDR("127.0.0.1/24")
	if err != nil {
		t.Error(err)
	}
	shortenerGRPC := app.NewURLShortenerGRPC(repository, mockConn, subnet)

	ctx := context.Background()
	conn, err := grpc.NewClient(
		"passthrough://bufnet",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithContextDialer(dialer(shortenerGRPC)),
	)
	if err != nil {
		t.Fatalf("Failed to dial bufnet: %v", err)
	}
	defer conn.Close()

	client := pb.NewURLClient(conn)

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			repository.ExpectedCalls = nil
			repository.EXPECT().Delete(
				mock.Anything,
				mock.Anything,
			).Return(nil)

			resp, err := client.DeleteBatch(ctx, &pb.DeleteBatchRequest{
				UserId: "testUserID",
				Urls:   test.want.events,
			})
			if err != nil {
				t.Fatalf("gRPC DeleteBatch failed: %v", err)
			}

			assert.Equal(t, test.want.status, resp.Status)
		})
	}
}

func TestURLShortenerGRPC_GetStats(t *testing.T) {
	type want struct {
		status    string
		countUser int
		countURL  int
		err       *error
	}
	testError := errors.New("test error")
	tests := []struct {
		name      string
		want      want
		ipAddress string
	}{
		{
			name: "positive test #1",
			want: want{
				status:    "ok",
				countUser: 10,
				countURL:  158,
			},
			ipAddress: "127.0.0.1",
		},
		{
			name: "negative test #1",
			want: want{
				status: "forbidden",
			},
			ipAddress: "192.168.0.1",
		},
		{
			name: "negative test #2",
			want: want{
				status:    "internal server error",
				countUser: 10,
				countURL:  158,
				err:       &testError,
			},
			ipAddress: "127.0.0.1",
		},
	}

	environments.BaseAddr = "http://example.com"
	repository := new(models2.URLRepository)
	mockConn, _ := pgxmock.NewConn()
	defer mockConn.Close(context.Background())
	_, subnet, err := net.ParseCIDR("127.0.0.1/24")
	if err != nil {
		t.Error(err)
	}
	shortenerGRPC := app.NewURLShortenerGRPC(repository, mockConn, subnet)

	ctx := context.Background()
	conn, err := grpc.NewClient(
		"passthrough://bufnet",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithContextDialer(dialer(shortenerGRPC)),
	)
	if err != nil {
		t.Fatalf("Failed to dial bufnet: %v", err)
	}
	defer conn.Close()

	client := pb.NewURLClient(conn)

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			repository.ExpectedCalls = nil
			if test.want.err != nil {
				repository.EXPECT().GetStats(
					mock.Anything,
				).Return(test.want.countUser, test.want.countURL, *test.want.err)
			} else {
				repository.EXPECT().GetStats(
					mock.Anything,
				).Return(test.want.countUser, test.want.countURL, nil)
			}

			resp, err := client.GetStats(ctx, &pb.GetStatsRequest{
				IpAddress: test.ipAddress,
			})
			if err != nil {
				t.Fatalf("gRPC GetStats failed: %v", err)
			}

			assert.Equal(t, test.want.status, resp.Status)
			if resp.Status == "ok" {
				assert.Equal(t, int32(test.want.countURL), resp.Urls)
				assert.Equal(t, int32(test.want.countUser), resp.Users)
			}
		})
	}
}
