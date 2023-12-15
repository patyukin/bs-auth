package app

import (
	"context"
	"github.com/patyukin/bs-auth/internal/cacher"
	"github.com/patyukin/bs-auth/internal/cacher/redis"
	"log"

	"github.com/patyukin/bs-auth/internal/api/auth"
	"github.com/patyukin/bs-auth/internal/api/user"
	"github.com/patyukin/bs-auth/internal/queue/kafka"
	authRepository "github.com/patyukin/bs-auth/internal/repository/auth"

	"github.com/patyukin/bs-auth/internal/client/db"
	"github.com/patyukin/bs-auth/internal/client/db/pg"
	"github.com/patyukin/bs-auth/internal/client/db/transaction"
	"github.com/patyukin/bs-auth/internal/closer"
	"github.com/patyukin/bs-auth/internal/config"
	"github.com/patyukin/bs-auth/internal/repository"
	userRepository "github.com/patyukin/bs-auth/internal/repository/user"
	"github.com/patyukin/bs-auth/internal/service"
	authService "github.com/patyukin/bs-auth/internal/service/auth"
	userService "github.com/patyukin/bs-auth/internal/service/user"
)

type serviceProvider struct {
	config *config.Config

	dbClient  db.Client
	txManager db.TxManager

	userRepository repository.UserRepository
	authRepository repository.AuthRepository

	userService service.UserService
	authService service.AuthService

	userImpl *user.Implementation
	authImpl *auth.Implementation
	producer *kafka.KafkaProducer
	cacher   *redis.RedisClient
}

func newServiceProvider(cfg *config.Config) *serviceProvider {
	return &serviceProvider{
		config: cfg,
	}
}

func (s *serviceProvider) DBClient(ctx context.Context) db.Client {
	if s.dbClient == nil {
		cl, err := pg.New(ctx, s.config.PG.DSN)
		if err != nil {
			log.Fatalf("failed to create db client: %v", err)
		}

		err = cl.DB().Ping(ctx)
		if err != nil {
			log.Fatalf("ping error: %s", err.Error())
		}
		closer.Add(cl.Close)

		s.dbClient = cl
	}

	return s.dbClient
}

func (s *serviceProvider) TxManager(ctx context.Context) db.TxManager {
	if s.txManager == nil {
		s.txManager = transaction.NewTransactionManager(s.DBClient(ctx).DB())
	}

	return s.txManager
}

func (s *serviceProvider) UserRepository(ctx context.Context) repository.UserRepository {
	if s.userRepository == nil {
		s.userRepository = userRepository.NewRepository(s.DBClient(ctx))
	}

	return s.userRepository
}

func (s *serviceProvider) UserService(ctx context.Context) service.UserService {
	if s.userService == nil {
		s.userService = userService.NewService(
			s.UserRepository(ctx),
			s.TxManager(ctx),
		)
	}

	return s.userService
}

func (s *serviceProvider) UserImpl(ctx context.Context) *user.Implementation {
	if s.userImpl == nil {
		s.userImpl = user.NewImplementation(s.UserService(ctx))
	}

	return s.userImpl
}

func (s *serviceProvider) AuthRepository(ctx context.Context) repository.AuthRepository {
	if s.authRepository == nil {
		s.authRepository = authRepository.NewRepository(s.DBClient(ctx))
	}

	return s.authRepository
}

func (s *serviceProvider) AuthService(ctx context.Context) service.AuthService {
	if s.authService == nil {
		s.authService = authService.NewService(
			s.AuthRepository(ctx),
			s.UserRepository(ctx),
			s.TxManager(ctx),
			s.Producer(ctx),
			s.Cacher(ctx),
		)
	}

	return s.authService
}

func (s *serviceProvider) Producer(_ context.Context) *kafka.KafkaProducer {
	var err error
	if s.producer == nil {
		s.producer, err = kafka.NewSyncProducer([]string{s.config.Kafka.DSN}, "my-topic")
		if err != nil {
			log.Fatalf("failed to create kafka producer: %v", err)
		}
	}

	return s.producer
}

func (s *serviceProvider) Cacher(_ context.Context) cacher.Cacher {
	address := s.config.Redis.DSN
	if s.cacher == nil {
		client := redis.NewRedis(address)

		closer.Add(client.Close)

		s.cacher = client
	}

	return s.cacher
}

func (s *serviceProvider) AuthImpl(ctx context.Context) *auth.Implementation {
	if s.authImpl == nil {
		s.authImpl = auth.NewImplementation(s.AuthService(ctx), s.UserService(ctx), s.config)
	}

	return s.authImpl
}
