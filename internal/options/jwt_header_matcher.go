package options

//
//import (
//	"github.com/grpc-ecosystem/grpc-gateway/runtime"
//	"multi-acquire-api-gateway.nextpay.global/middleware"
//)
//
//// jwtTokenMatcher injects origin mPOS claim information into gRPC metadata
//func NewJwtTokenMatcher(key string) (string, bool) {
//	switch key {
//	case middleware.ClaimMposId, middleware.ClaimMposUserId, middleware.ClaimMposMerchantId, middleware.ClaimMposTid, middleware.SecretSessionIDHeaderValue:
//		return key, true
//	default:
//		return runtime.DefaultHeaderMatcher(key)
//	}
//}
