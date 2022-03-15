package middleware

import (
	"crypto/rsa"
	"encoding/base64"
	"fmt"
	"github.com/kataras/jwt"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"net/http"
	"phuongnd/gateway/api"
	"strings"
)

const (
	UserID = "X-User-ID"
)

func init() {
	jwt.Unmarshal = jwt.UnmarshalWithRequired
}

type UserClaim struct {
	Expiry  int64  `json:"exp"`
	IssueAt int64  `json:"iat"`
	UserID  string `json:"userId,required"`
}

type JwtMiddleware struct {
	//app    *iris.Application
	logger *zap.Logger
	pubKey *rsa.PublicKey
}

func NewJwtMiddleware(logger *zap.Logger) *JwtMiddleware {
	key, err := base64.StdEncoding.DecodeString(viper.GetString("AUTH_PUBLIC_KEY"))
	if err != nil {
		logger.Fatal(fmt.Sprintf("base64 decode failed: %v", err))
	}
	rsaPubKey, err := jwt.ParsePublicKeyRSA(key)
	if err != nil {
		logger.Fatal(fmt.Sprintf("RSA authentication public key parse failed: %v", err))
	}
	return &JwtMiddleware{logger: logger, pubKey: rsaPubKey}
}

func (jw *JwtMiddleware) Wrap(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if api.IsAuthorizedEndpoints(r.URL.Path) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				jw.logger.Error("missing authentication token")
				writeErrorResponse(w, http.StatusUnauthorized, "Unauthorized")
				return
			}
			// TODO: Make this a bit more robust, parsing-wise
			authHeaderParts := strings.Split(authHeader, " ")
			if len(authHeaderParts) != 2 || strings.ToLower(authHeaderParts[0]) != "bearer" {
				jw.logger.Error("invalid JWT token format")
				writeErrorResponse(w, http.StatusUnauthorized, "Unauthorized")
				return
			}
			verifiedToken, err := jwt.Verify(jwt.RS256, jw.pubKey, []byte(authHeaderParts[1]))
			if err != nil {
				jw.logger.Error(fmt.Sprintf("invalid token: %v", err))
				writeErrorResponse(w, http.StatusUnauthorized, "Unauthorized")
				return
			}
			var claims UserClaim
			if err = verifiedToken.Claims(&claims); err != nil {
				jw.logger.Error(fmt.Sprintf("claim decode failed: %v", err))
				writeErrorResponse(w, http.StatusUnauthorized, "Unauthorized")
				return
			}
			r.Header.Add(UserID, claims.UserID)
		}
		next.ServeHTTP(w, r)
	})
}
