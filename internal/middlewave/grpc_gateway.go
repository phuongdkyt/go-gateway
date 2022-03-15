package middleware

import (
	"net/http"
)

type GrpcGatewayRouterMiddleware interface {
	Wrap(next http.Handler) http.Handler
}
