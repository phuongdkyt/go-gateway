package service

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"phuongnd/gateway/internal/options"
	api "phuongnd/gateway/proto"
)

func (b *Service) InitSecretKey(ctx context.Context, request *api.EncryptedGatewayRequest) (*api.EncryptedGatewayResponse, error) {
	key, err := b.crytoService.RsaDecrypt(request.Data)
	if err != nil {
		b.log.Error(fmt.Sprintf("secret-session/rsa: fail to decrypt request payload: %v", err))
		errorStatus := status.New(codes.InvalidArgument, "fail to decrypt request payload")
		return nil, errorStatus.Err()
	}

	session, err := b.sessionImpl.Create(key)
	if err != nil {
		b.log.Error(fmt.Sprintf("secret-session: fail to create: %v", err))
		errorStatus := status.New(codes.Internal, "error creating secret session")
		return nil, errorStatus.Err()
	}
	encryptId, err := b.crytoService.AeadEncrypt(session.Id, "", session.Key)
	if err != nil {
		b.log.Error(fmt.Sprintf("secret-session/aead: fail to encrypt: %v", err))
		errorStatus := status.New(codes.Internal, "error encrypting response")
		return nil, errorStatus.Err()
	}
	if err := grpc.SetHeader(ctx, metadata.Pairs(options.GrpcMetadataHttpCode, "201")); err != nil {
		b.log.Error(fmt.Sprintf("secret-session: fail to set gRPC header: %v", err))
	}
	return &api.EncryptedGatewayResponse{
		Data: encryptId,
	}, nil
}
