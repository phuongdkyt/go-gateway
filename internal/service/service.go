package service

import (
	"github.com/google/wire"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"phuongnd/gateway/internal/cryto"
	"phuongnd/gateway/internal/session"
	api "phuongnd/gateway/proto"
)

var ProviderServiceSet = wire.NewSet(NewService)

var _ IService = (*Service)(nil)

type Service struct {
	log          *zap.Logger
	crytoService cryto.CrytoInterface
	sessionImpl  session.SessionInterface
	api.UnimplementedDevelopmentServiceServer
}

func NewService(log *zap.Logger, cryto cryto.CrytoInterface, sessionImpl session.SessionInterface) IService {
	svc := &Service{
		log:          log,
		crytoService: cryto,
		sessionImpl:  sessionImpl,
	}

	return svc
}

func (s *Service) RegisterServiceServer(grpcServer *grpc.Server) {
	api.RegisterDevelopmentServiceServer(grpcServer, s)
	api.RegisterGatewayServiceServer(grpcServer, s)
}
