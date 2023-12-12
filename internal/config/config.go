package config

const (
	AppEnv = "APP_ENV"
)

type Config struct {
	Server struct {
		GRPC struct {
			Port string
			Host string
		}
		HTTP struct {
			Port string
			Host string
		}
		Swagger struct {
			Port string
			Host string
		}
	}
	PG struct {
		DSN string
	}
	Redis struct {
		DSN string
	}
	Kafka struct {
		DSN string
	}
	AuthTokenSignKey string
}
