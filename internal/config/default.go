package config

import (
	"github.com/spf13/viper"
)

// defaultConfig is the default configuration for the application.
func defaultConfig() {
	// APP
	viper.SetDefault("APP_NAME", "Golang gRPC Base Project")
	viper.SetDefault("APP_VERSION", "0.0.0")
	viper.SetDefault("APP_PORT", "8088")
	viper.SetDefault("GRPC_PORT", "3100")
	viper.SetDefault("HTTP_PORT", "3200")

	viper.SetDefault("LOG_PAYLOAD", true)
	viper.SetDefault("APP_KEY", "./config/key.pem")
	viper.SetDefault("APP_CERT", "./config/cert.pem")

	// DATABASE
	viper.SetDefault("MONGODB_URI", "mongodb://root:123456aA%40@localhost:27017")
	viper.SetDefault("MONGODB_DBNAME", "base")

	// REDIS
	viper.SetDefault("REDIS_URL", "localhost:16379")
	viper.SetDefault("REDIS_PASSWORD", "")
	viper.SetDefault("REDIS_DB", 0)

	viper.SetDefault("RSA_PRIVATE_KEY", "")
	viper.SetDefault("AUTH_PUBLIC_KEY", "")
	viper.SetDefault("SECRET_SESSION_TIMEOUT", "1h")
}
