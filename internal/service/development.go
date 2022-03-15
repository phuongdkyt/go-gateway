package service

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/google/tink/go/aead"
	"github.com/google/tink/go/insecurecleartextkeyset"
	"github.com/google/tink/go/keyset"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/types/known/emptypb"
	api "phuongnd/gateway/proto"
	"time"
)

func (d *Service) GenerateSecretKey(context.Context, *emptypb.Empty) (*api.GenerateSecretKeyResponse, error) {
	kh, err := keyset.NewHandle(aead.AES256GCMKeyTemplate())
	if err != nil {
		return nil, err
	}
	exportedPriv := &keyset.MemReaderWriter{}
	if err := insecurecleartextkeyset.Write(kh, exportedPriv); err != nil {
		return nil, err
	}
	ksPriv, err := proto.Marshal(exportedPriv.Keyset)
	if err != nil {
		return nil, err
	}
	rng := rand.Reader
	ciphertext, err := rsa.EncryptOAEP(sha256.New(), rng, d.crytoService.GetRsaPublicKey(), ksPriv, nil)
	return &api.GenerateSecretKeyResponse{
		Data:      base64.StdEncoding.EncodeToString(ciphertext),
		SecretKey: base64.StdEncoding.EncodeToString(ksPriv),
	}, nil
}

func (d *Service) DecryptSecretSessionId(_ context.Context, request *api.DecryptSecretSessionIdRequest) (*api.DecryptSecretSessionIdResponse, error) {
	cipherText, err := d.crytoService.AeadDecrypt(request.Data, "", request.SecretKey)
	if err != nil {
		return nil, err
	}
	return &api.DecryptSecretSessionIdResponse{
		Data: cipherText,
	}, nil
}
func (d *Service) EncryptData(_ context.Context, request *api.DevEncryptDataRequest) (*api.DevEncryptDataResponse, error) {
	cipherText, err := d.crytoService.AeadEncrypt(request.Data, "", request.SecretKey)
	if err != nil {
		return nil, err
	}
	return &api.DevEncryptDataResponse{Data: cipherText}, nil
}
func (d *Service) DecryptData(_ context.Context, request *api.DevDecryptDataRequest) (*api.DevDecryptDataResponse, error) {
	clearText, err := d.crytoService.AeadDecrypt(request.Data, "", request.SecretKey)
	if err != nil {
		return nil, err
	}
	return &api.DevDecryptDataResponse{Data: clearText}, nil
}
func (d *Service) SetSecretSessionTimeout(ctx context.Context, request *api.SetSecretSessionTimeoutRequest) (*emptypb.Empty, error) {
	headers, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, errors.New("read metadata failed")
	}
	sessionId, ok := headers["x-secret-session-id"]
	if !ok {
		return nil, errors.New("missing secret session id")
	}
	timeout, err := time.ParseDuration(request.Timeout)
	if err != nil {
		return nil, err
	}
	if err := d.sessionImpl.SetTimeout(sessionId[0], timeout); err != nil {
		return nil, err
	}
	return &empty.Empty{}, nil
}
