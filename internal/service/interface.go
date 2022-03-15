package service

import "google.golang.org/grpc"
import api "phuongnd/gateway/proto"

type IService interface {
	RegisterServiceServer(grpcServer *grpc.Server)
	api.DevelopmentServiceServer
}
