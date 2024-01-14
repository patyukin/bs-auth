package auth

import (
	"github.com/patyukin/bs-auth/internal/cacher"
	"github.com/patyukin/bs-auth/internal/client/db"
	"github.com/patyukin/bs-auth/internal/queue/kafka"
	"github.com/patyukin/bs-auth/internal/repository"
	"github.com/patyukin/bs-auth/internal/service"
)

type serv struct {
	userRepository repository.UserRepository
	txManager      db.TxManager
	authRepository repository.AuthRepository
	producer       *kafka.KafkaProducer
	cacher         cacher.Cacher
}

func NewService(authRepository repository.AuthRepository, userRepository repository.UserRepository, txManager db.TxManager, producer *kafka.KafkaProducer, cacher cacher.Cacher) service.AuthService {
	return &serv{
		authRepository: authRepository,
		userRepository: userRepository,
		txManager:      txManager,
		producer:       producer,
		cacher:         cacher,
	}
}
