package middlewares

import (
	"context"
	"time"

	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
)

func NewLoggingInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		log.Info().Str("FullMethod", info.FullMethod).Msg("received request")
		start := time.Now()
		defer func() {
			log.Info().Dur("Duration", time.Since(start)).Msg("request completed")
		}()
		return handler(ctx, req)
	}
}
