package config

import (
	"fmt"
	"github.com/joho/godotenv"
	"os"
)

func LoadEnvConfig() (*Config, error) {
	configFilePath := os.Getenv(EnvFilePath)

	cfg := Config{}

	// Open config file
	err := godotenv.Load(configFilePath)
	if err != nil {
		return nil, fmt.Errorf("failed to load %s file: %w", configFilePath, err)
	}

	cfg.Server.HTTP.Host = os.Getenv("HTTP_HOST")
	cfg.Server.HTTP.Port = os.Getenv("HTTP_PORT")
	cfg.Server.GRPC.Host = os.Getenv("GRPC_HOST")
	cfg.Server.GRPC.Port = os.Getenv("GRPC_PORT")
	cfg.Server.Swagger.Host = os.Getenv("SWAGGER_HOST")
	cfg.Server.Swagger.Port = os.Getenv("SWAGGER_PORT")

	cfg.PG.DSN = os.Getenv("PG_DSN")
	cfg.Redis.DSN = os.Getenv("REDIS_DSN")
	cfg.Kafka.DSN = os.Getenv("KAFKA_DSN")

	cfg.AuthTokenSignKey = os.Getenv("TOKEN_SIGN_KEY")
	return &cfg, nil
}
