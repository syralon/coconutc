// @file: internal/transport/interceptor/logger.go

package interceptor

import (
	"context"
	"log/slog"
	"time"

	"google.golang.org/grpc"
)

func LogServerInterceptor() func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
		start := time.Now()
		defer func() {
			slog.InfoContext(ctx, "grpc server", slog.String("method", info.FullMethod), slog.Duration("cost", time.Since(start)))
		}()
		return handler(ctx, req)
	}
}
