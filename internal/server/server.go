package server

import (
	"context"
	"fmt"
	"github.com/google/wire"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"net"
	middleware "phuongnd/gateway/internal/middlewave"
	"phuongnd/gateway/internal/service"
)

// ProviderServerSet is server providers.
var ProviderServerSet = wire.NewSet(NewServer, NewGRPCServer, newHTTPServer)
var _ IServer = (*Server)(nil)

type Server struct {
	ctx     context.Context
	log     *zap.Logger
	grpc    *grpc.Server
	service service.IService
}

func NewServer(
	ctx context.Context, log *zap.Logger, grpc *grpc.Server, service service.IService,
	cryptoMiddleware *middleware.CryptoMiddleware,
	jwtMiddleware *middleware.JwtMiddleware) IServer {
	// register service to grpc server
	service.RegisterServiceServer(grpc)

	grpcPort := fmt.Sprintf(":%d", viper.GetInt("GRPC_PORT"))
	log.Info(fmt.Sprintf("Serving gRPC on http://localhost%s", grpcPort))

	listenGRPC, err := net.Listen("tcp", grpcPort)
	if err != nil {
		log.Panic("net.Listen()", zap.Error(err))
	}

	go func() {
		if err = grpc.Serve(listenGRPC); err != nil {
			log.Panic("grpc.Serve()", zap.Error(err))
		}
	}()

	httpPort := fmt.Sprintf(":%d", viper.GetInt("HTTP_PORT"))
	log.Info(fmt.Sprintf("Serving gRPC-Gateway on http://localhost%s", httpPort))

	go func() {
		if err = newHTTPServer(ctx, log, cryptoMiddleware, jwtMiddleware); err != nil {
			log.Panic("newHTTPServer()", zap.Error(err))
		}
	}()

	return &Server{
		ctx:     ctx,
		log:     log,
		grpc:    grpc,
		service: service,
	}
}

func (s *Server) Close() error {
	return nil
}
