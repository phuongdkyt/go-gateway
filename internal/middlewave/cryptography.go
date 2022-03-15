package middleware

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/google/wire"
	jsoniter "github.com/json-iterator/go"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"io"
	"io/ioutil"
	"net/http"
	"phuongnd/gateway/api"
	"phuongnd/gateway/internal/cryto"
	"phuongnd/gateway/internal/session"
)

var ProvideMiddlewareSet = wire.NewSet(NewJwtMiddleware, NewCryptoMiddleware)

const (
	SecretSessionIDHeaderValue = "X-Secret-Session-Id"
)

// ResponseWriter wrap http response to CryptoService response body
type ResponseWriter struct {
	http.ResponseWriter
	Buf *bytes.Buffer
}

func (erw *ResponseWriter) Write(p []byte) (int, error) {
	return erw.Buf.Write(p)
}

type CryptoMiddleware struct {
	SessionService session.SessionInterface
	CryptoService  cryto.CrytoInterface
	Logger         *zap.Logger
}

func NewCryptoMiddleware(
	sessionService session.SessionInterface,
	cryptoService cryto.CrytoInterface,
	logger *zap.Logger) *CryptoMiddleware {
	return &CryptoMiddleware{
		SessionService: sessionService,
		CryptoService:  cryptoService,
		Logger:         logger,
	}
}

func (cm *CryptoMiddleware) Wrap(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		if !api.IsRequiredEncryptionPath(path) {
			next.ServeHTTP(w, r)
			return
		}

		sessionId := r.Header.Get(SecretSessionIDHeaderValue)
		if sessionId == "" {
			writeErrorResponse(w, http.StatusBadRequest, "invalid secret session id")
			return
		}

		// Get corresponding secret key from SecretSessionService
		secretKey, err := cm.SessionService.Get(sessionId)
		if err != nil {
			if errors.Is(err, session.ErrSessionKeyNotFound) {
				writeErrorResponse(w, http.StatusPreconditionFailed, "invalid secret session")
				return
			} else {
				cm.Logger.Error(fmt.Sprintf("failed to get secret key: %v", err))
				writeErrorResponse(w, http.StatusInternalServerError, "fail to get secret key")
				return
			}
		}

		// Refresh secret SecretSessionService
		go func() {
			if err := cm.SessionService.Refresh(sessionId); err != nil {
				// deliberately failed silently
				cm.Logger.Error(fmt.Sprintf("secret session refresh failed: %v", err))
			}
		}()

		// Decrypt request body
		var payload api.EncryptedRequestPayload
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			cm.Logger.Error(fmt.Sprintf("failed to marshall request payload: %v", err))
			writeErrorResponse(w, http.StatusBadRequest, "fail to marshal request payload")
			return
		}
		decryptedBody, err := cm.CryptoService.AeadDecrypt(payload.Data, "", secretKey)
		if err != nil {
			cm.Logger.Error(fmt.Sprintf("failed to decrypt request payload: %v", err))
			writeErrorResponse(w, http.StatusBadRequest, "fail to decrypt request payload")
			return
		}
		r.Body = ioutil.NopCloser(bytes.NewBuffer([]byte(decryptedBody)))

		// Wrap response write by ResponseWriter
		responseWriter := &ResponseWriter{
			ResponseWriter: w,
			Buf:            &bytes.Buffer{},
		}

		// Serve request through grpc-gateway
		next.ServeHTTP(responseWriter, r)

		// Copy responseWriter.Buf to w(default response writer)
		cipherText, err := cm.CryptoService.AeadEncrypt(string(responseWriter.Buf.Bytes()), "", secretKey)
		if err != nil {
			cm.Logger.Error(fmt.Sprintf("failed to encrypt response payload: %v", err))
			writeErrorResponse(w, http.StatusInternalServerError, "fail to encrypt response body")
			return
		}
		encryptedResponse := api.EncryptedResponsePayload{
			Data: cipherText,
		}
		responseBody, _ := jsoniter.Marshal(encryptedResponse)
		if _, err := io.Copy(w, bytes.NewBuffer(responseBody)); err != nil {
			cm.Logger.Error(fmt.Sprintf("failed to write response body: %v", err))
			writeErrorResponse(w, http.StatusInternalServerError, "fail to write response body")
			return
		}
	})
}

func writeErrorResponse(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = jsoniter.NewEncoder(w).Encode(api.ErrorResponse{
		Message: message,
	})
}
