// @internal/transport/middleware/logger.go

package middleware

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
)

func Logger() runtime.Middleware {
	return func(next runtime.HandlerFunc) runtime.HandlerFunc {
		return func(writer http.ResponseWriter, request *http.Request, pathParams map[string]string) {
			begin := time.Now()
			ctx := request.Context()
			defer func() {
				slog.InfoContext(ctx, "gateway request", slog.String("path", request.URL.Path), slog.Duration("cost", time.Since(begin)))
			}()
			next(writer, request, pathParams)
		}
	}
}
