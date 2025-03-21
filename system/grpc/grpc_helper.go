package grpc

import (
	"context"
	"log/slog"

	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	"google.golang.org/grpc"
)

const (
	maxMessageSize = 1024 * 1024 * 24
)

// InterceptorLogger adapts slog logger to interceptor logger.
// This code is simple enough to be copied and not imported.
func InterceptorLogger(l *slog.Logger) logging.Logger {
	return logging.LoggerFunc(func(ctx context.Context, lvl logging.Level, msg string, fields ...any) {
		l.Log(ctx, slog.Level(lvl), msg, fields...)
	})
}

func CreateGrpcConnection(ctx context.Context, address string) (*grpc.ClientConn, error) {
	opts := []logging.Option{
		logging.WithLevels(logging.DefaultClientCodeToLevel),
		logging.WithLogOnEvents(logging.FinishCall),
	}

	dialOpts := []grpc.DialOption{
		grpc.WithChainUnaryInterceptor(
			logging.UnaryClientInterceptor(InterceptorLogger(slog.Default()), opts...),
		),
		grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(maxMessageSize)),
	}

	conn, err := grpc.NewClient(address, dialOpts...)
	if err != nil {
		slog.ErrorContext(ctx, "grcp client connection failed", slog.Any("error", err), slog.String("address", address))
		return nil, err
	}
	return conn, nil
}
