//go:build wireinject
// +build wireinject

package main

import (
	"context"
	"github.com/google/wire"
	"go.uber.org/zap"
	"phuongnd/gateway/internal/cryto"
	middleware "phuongnd/gateway/internal/middlewave"
	"phuongnd/gateway/internal/server"
	"phuongnd/gateway/internal/service"
	"phuongnd/gateway/internal/session"
	"phuongnd/gateway/pkg/redis"
)

func InitializeServer(ctx context.Context, log *zap.Logger) server.IServer {
	wire.Build(
		redis.ProviderRedisSet,
		session.ProviderSessionSet,
		cryto.ProviderCrytoSet,
		middleware.ProvideMiddlewareSet,
		service.ProviderServiceSet,
		server.ProviderServerSet,
	)
	return &server.Server{}
}
