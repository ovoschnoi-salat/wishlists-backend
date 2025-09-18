package middlewares

import (
	"backend/internal/store"
	"context"
	"strings"
	"time"

	"github.com/rs/zerolog/log"
	initdata "github.com/telegram-mini-apps/init-data-golang"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

var (
	errMissingMetadata = status.Errorf(codes.InvalidArgument, "missing metadata")
	errInvalidToken    = status.Errorf(codes.Unauthenticated, "invalid token")
	errMissingData     = status.Errorf(codes.Unauthenticated, "missing init data")
)

const initDataCtxKey = "init_data"

func NewTgLoginInterceptor(secretToken string, db store.Queries) grpc.UnaryServerInterceptor {
	// Define how long since init data generation date init data is valid.
	expIn := 10 * time.Minute

	return func(ctx context.Context, req any, _ *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		// authentication (token verification)
		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return nil, errMissingMetadata
		}

		initData := md["authorization"]
		log.Print("initData ", initData)
		if len(initData) != 1 {
			return nil, errInvalidToken
		}
		if strings.HasPrefix(initData[0], "tma") {
			return nil, errInvalidToken
		}
		err := initdata.Validate(strings.TrimPrefix(initData[0], "tma "), secretToken, expIn)
		if err != nil {
			return nil, err
		}
		data, err := initdata.Parse(strings.TrimPrefix(initData[0], "tma "))
		if err != nil {
			return nil, err
		}

		return handler(context.WithValue(ctx, initDataCtxKey, data), req)
	}
}

func GetInitDataFromContext(ctx context.Context) (initdata.InitData, error) {
	if value, ok := ctx.Value(initDataCtxKey).(initdata.InitData); ok {
		return value, nil
	}
	return initdata.InitData{}, errMissingData
}
