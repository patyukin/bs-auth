package config

import (
	"github.com/joho/godotenv"
	"os"
)

func LoadEnvConfig() (*Config, error) {
	appEnv := os.Getenv(AppEnv)

	var configFilePath string
	// проверка возможных вариантов для appEnv
	switch appEnv {
	case "local-dev":
		configFilePath = ".env.local-dev"
	case "local":
		configFilePath = ".env.local"
	case "test":
		configFilePath = ".env.test"
	case "test-dev":
		configFilePath = ".env.test-dev"
	case "prod":
		configFilePath = ".env.prod"
	default:
		panic("config not found")
	}

	cfg := Config{}

	// Open config file
	err := godotenv.Load(configFilePath)
	if err != nil {
		return nil, err
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
