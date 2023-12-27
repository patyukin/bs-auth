package auth

import (
	"context"
	"github.com/brianvoe/gofakeit"
	"github.com/patyukin/bs-auth/internal/app"
	"github.com/patyukin/bs-auth/internal/cacher"
	"github.com/patyukin/bs-auth/internal/cacher/redis"
	"github.com/patyukin/bs-auth/internal/client/db"
	"github.com/patyukin/bs-auth/internal/client/db/pg"
	"github.com/patyukin/bs-auth/internal/config"
	desc "github.com/patyukin/bs-auth/pkg/auth_v1"
	"github.com/stretchr/testify/suite"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"testing"
	"time"
)

type SignInGRPCTestSuite struct {
	suite.Suite
	serverAddr string
	server     desc.AuthV1Server
	client     desc.AuthV1Client
	dbClient   db.Client
	cacher     cacher.Cacher
}

func (s *SignInGRPCTestSuite) SetupTest() {
	ctx := context.Background()

	cfg, err := config.LoadEnvConfig()
	if err != nil {
		s.FailNow("Failed to init app: " + err.Error())
	}

	a, err := app.NewApp(ctx, cfg)
	if err != nil {
		log.Fatalf("failed to init app: %s", err.Error())
	}

	go func() {
		err = a.RunGRPCServer()
		s.Require().NoError(err)
	}()

	// проверка запустится ли сервер grpc
	var conn *grpc.ClientConn
	s.Require().Eventually(func() bool {
		conn, err = grpc.Dial(cfg.Server.GRPC.Host+":"+cfg.Server.GRPC.Port, grpc.WithTransportCredentials(insecure.NewCredentials()))
		return err == nil
	}, 1*time.Second, 100*time.Millisecond)

	s.client = desc.NewAuthV1Client(conn)

	cl, err := pg.New(ctx, cfg.PG.DSN)
	if err != nil {
		log.Fatalf("failed to create db client: %v", err)
	}

	err = cl.DB().Ping(ctx)
	if err != nil {
		log.Fatalf("ping error: %s", err.Error())
	}

	s.dbClient = cl
	s.cacher = redis.NewRedis(cfg.Redis.DSN)
}

func (s *SignInGRPCTestSuite) TearDownTest() {
	// Остановите сервер и выполните другие задачи по очистке после каждого теста.
	// Например:
}

func TestSignInGRPCTestSuite(t *testing.T) {
	suite.Run(t, new(SignInGRPCTestSuite))
}

func (s *SignInGRPCTestSuite) TestAuthSignIn() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*1000)
	defer cancel()

	name := gofakeit.Name()
	email := gofakeit.Email()
	password := gofakeit.Password(true, false, false, false, false, 10)
	hashBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	passwordHash := string(hashBytes)

	_, err = s.dbClient.DB().ExecContext(
		ctx,
		db.Query{QueryRaw: "INSERT INTO users (name, email, password_hash) VALUES ($1, $2, $3)"},
		name,
		email,
		passwordHash,
	)

	s.Require().NoError(err)

	res, err := s.client.SignIn(ctx, &desc.SignInRequest{
		Email:       email,
		Password:    password,
		Fingerprint: gofakeit.UUID(),
	})

	s.Require().NoError(err)
	s.Require().NotNil(res)
	// TUIZJTWMTMMYHWDAS424UXT2
	// otpauth://totp/example.com:admin@admin.com?algorithm=SHA1&digits=6&issuer=example.com&period=30&secret=TUIZJTWMTMMYHWDAS424UXT2
}
