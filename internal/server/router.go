package server

import (
	"context"
	"fmt"
	"github.com/spf13/viper"
	"net"
	"net/http"
	middleware "phuongnd/gateway/internal/middlewave"
	"time"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/kataras/muxie"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/encoding/protojson"
	api "phuongnd/gateway/proto"
)

// newHTTPServer create Gateway server
func newHTTPServer(ctx context.Context, log *zap.Logger,
	cryptoMiddleware *middleware.CryptoMiddleware,
	jwtMiddleware *middleware.JwtMiddleware,

) error {
	listenHTTP, err := net.Listen("tcp", ":"+viper.GetString("HTTP_PORT"))
	if err != nil {
		log.Error("net.Listen()", zap.Error(err))
		return err
	}

	// Create HTTP Server
	gwmux := runtime.NewServeMux(
		runtime.WithMarshalerOption(runtime.MIMEWildcard, &runtime.JSONPb{
			MarshalOptions: protojson.MarshalOptions{
				Multiline:       false,
				Indent:          "",
				AllowPartial:    false,
				UseProtoNames:   true,
				UseEnumNumbers:  false,
				EmitUnpopulated: true,
				Resolver:        nil,
			},
		}),
	)

	//conn, err := grpc.DialContext(
	//	ctx,
	//	grpcPort,
	//	grpc.WithBlock(),
	//	grpc.WithTransportCredentials(insecure.NewCredentials()),
	//)
	if err != nil {
		log.Error("Failed to dial Server:", zap.Error(err))
		return err
	}

	mux := muxie.NewMux()
	mux.Use(cryptoMiddleware.Wrap)
	mux.Handle("/*", gwmux)

	if err = registerServicesHandler(ctx, gwmux, ":"+viper.GetString("GRPC_PORT")); err != nil {
		log.Error("registerServicesHandler()", zap.Error(err))
		return err
	}

	return http.Serve(listenHTTP, mux)
}

// registerServicesHandler register services handler
func registerServicesHandler(ctx context.Context, mux *runtime.ServeMux, grpcPort string) error {
	var err error
	grpcServerUrl := "0.0.0.0" + grpcPort
	//TODO: support mTLS
	opts := []grpc.DialOption{grpc.WithInsecure(),
		grpc.WithChainUnaryInterceptor(),
		grpc.WithTimeout(2 * time.Minute),
	}
	// Register UserService Handler
	if err = api.RegisterGatewayServiceHandlerFromEndpoint(ctx, mux, grpcServerUrl, opts); err != nil {
		return fmt.Errorf("proto.RegisterGatewayServiceHandlerFromEndpoint(): %w", err)
	}
	if err = api.RegisterDevelopmentServiceHandlerFromEndpoint(ctx, mux, grpcServerUrl, opts); err != nil {
		return fmt.Errorf("proto.RegisterDevelopmentServiceHandlerFromEndpoint(): %w", err)
	}
	if err = api.RegisterUserServiceHandlerFromEndpoint(ctx, mux, viper.GetString("USER_SERVICE"), opts); err != nil {
		return fmt.Errorf("proto.RegisterUserServiceHandlerFromEndpoint(): %w", err)
	}
	return nil
}
