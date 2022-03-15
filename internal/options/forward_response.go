package options

import (
	"context"
	"github.com/golang/protobuf/proto"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"net/http"
	"strconv"
)

const (
	GrpcMetadataHttpCode        = "x-http-code"
	GrpcMetadataHttpCodeForward = "Grpc-Metadata-X-Http-Code"
	GrpcMetadataContentType     = "Grpc-Metadata-Content-Type"
	GrpcMetadataAcceptEncoding  = "Grpc-Metadata-Grpc-Accept-Encoding"
)

// responseHttpHeaderModifier remove gRPC-related headers from HTTP response
func NewResponseHttpHeaderModifier(ctx context.Context, w http.ResponseWriter, p proto.Message) error {
	md, ok := runtime.ServerMetadataFromContext(ctx)
	if !ok {
		return nil
	}

	// set http status code
	if vals := md.HeaderMD.Get(GrpcMetadataHttpCode); len(vals) > 0 {
		code, err := strconv.Atoi(vals[0])
		if err != nil {
			return err
		}
		w.WriteHeader(code)
		// delete the headers to not expose any grpc-metadata in http response
		delete(md.HeaderMD, GrpcMetadataHttpCode)
		delete(w.Header(), GrpcMetadataHttpCodeForward)
		//delete(w.Header(), GrpcMetadataContentType)
		//delete(w.Header(), GrpcMetadataAcceptEncoding)
	}
	return nil
}
