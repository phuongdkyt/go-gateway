package api

import "strings"

const (
	DevelopmentEndpointPrefix = "/api/gateway/dev"
)

const (
	SecretInitEndpoint      = "/api/gateway/init"
	LanguageRefreshEndpoint = "/api/gateway/language/refresh"
	LoginEndpoint           = "/api/users/login"
)

func IsRequiredEncryptionPath(path string) bool {
	if isDevelopmentEndpoint(path) ||
		isProbesEndpoint(path) ||
		strings.HasPrefix(path, SecretInitEndpoint) ||
		strings.HasPrefix(path, LanguageRefreshEndpoint) {
		return false
	}
	return true
}

func IsAuthorizedEndpoints(path string) bool {
	if isDevelopmentEndpoint(path) ||
		isProbesEndpoint(path) ||
		strings.HasPrefix(path, SecretInitEndpoint) ||
		strings.HasPrefix(path, LoginEndpoint) ||
		strings.HasPrefix(path, LanguageRefreshEndpoint) {
		return false
	}
	return true
}

func isProbesEndpoint(path string) bool {
	if strings.HasPrefix(path, "/api/liveness") ||
		strings.HasPrefix(path, "/api/readiness") {
		return true
	}
	return false
}

func isDevelopmentEndpoint(path string) bool {
	if strings.Contains(path, DevelopmentEndpointPrefix) {
		return true
	}
	return false
}
