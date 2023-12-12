package auth

import (
	"context"
	"github.com/patyukin/bs-auth/internal/client/db"
	"github.com/patyukin/bs-auth/internal/client/db/pg"
	"golang.org/x/crypto/bcrypt"
	"log"
	"os"
	"testing"
	"time"

	"github.com/brianvoe/gofakeit"
	"github.com/patyukin/bs-auth/internal/api/auth"
	appInternal "github.com/patyukin/bs-auth/internal/app"
	desc "github.com/patyukin/bs-auth/pkg/auth_v1"
	"github.com/stretchr/testify/suite"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type GRPCTestSuite struct {
	suite.Suite

	server        *auth.Implementation
	clientConn    *grpc.ClientConn
	serviceClient desc.AuthV1Client
	impl          desc.UnimplementedAuthV1Server
	dbClient      db.Client
}

var (
	name         string
	email        string
	password     string
	passwordHash string
)

func (s *GRPCTestSuite) SetupSuite() {
	err := os.Setenv(appInternal.EnvFilePath, "./../../.env.test-dev")
	if err != nil {
		s.FailNow("Failed to init app: " + err.Error())
	}

	ctx := context.Background()
	app, err := appInternal.NewApp(ctx)
	if err != nil {
		s.FailNow("Failed to init app: " + err.Error())
	}

	go app.Run()

	conn, err := grpc.Dial("0.0.0.0:11110", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		s.FailNow("Failed to init app: " + err.Error())
	}

	s.clientConn = conn
	s.serviceClient = desc.NewAuthV1Client(conn)

	cl, err := pg.New(ctx, os.Getenv("PG_DSN"))
	if err != nil {
		log.Fatalf("failed to create db client: %v", err)
	}

	err = cl.DB().Ping(ctx)
	if err != nil {
		log.Fatalf("ping error: %s", err.Error())
	}

	s.dbClient = cl
}

func (s *GRPCTestSuite) TearDownSuite() {
	s.clientConn.Close()
	s.dbClient.Close()
}

func (s *GRPCTestSuite) TestSignIn() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*1000)
	defer cancel()

	name = gofakeit.Name()
	email = gofakeit.Email()
	password = gofakeit.Password(true, false, false, false, false, 10)
	hashBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	passwordHash = string(hashBytes)

	_, err = s.dbClient.DB().ExecContext(
		ctx,
		db.Query{QueryRaw: "INSERT INTO users (name, email, password_hash) VALUES ($1, $2, $3)"},
		name,
		email,
		passwordHash,
	)

	s.Require().NoError(err)

	res, err := s.serviceClient.SignIn(ctx, &desc.AuthRequest{
		Email:    email,
		Password: password,
	})

	s.Require().NoError(err)
	s.Require().NotNil(res)
}

func (s *GRPCTestSuite) TestSignInService() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*1000)
	defer cancel()

	res, err := s.serviceClient.SignIn(ctx, &desc.AuthRequest{
		Email:    email,
		Password: password,
	})

	s.Require().NoError(err)
	s.Require().NotNil(res)
}

func TestGRPCTestSuite(t *testing.T) {
	suite.Run(t, new(GRPCTestSuite))
}
