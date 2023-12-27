package auth

import (
	"context"
	"encoding/json"
	"fmt"
	sq "github.com/Masterminds/squirrel"
	"github.com/brianvoe/gofakeit"
	"github.com/patyukin/bs-auth/internal/client/db"
	"github.com/patyukin/bs-auth/internal/service/auth"
	desc "github.com/patyukin/bs-auth/pkg/auth_v1"
	"github.com/pquerna/otp/totp"
	"golang.org/x/crypto/bcrypt"
	"time"
)

func (s *SignInGRPCTestSuite) TestAuthCheckCode() {
	ctx := context.Background()

	fingerprint := gofakeit.UUID()
	name := gofakeit.Name()
	email := gofakeit.Email()
	password := gofakeit.Password(true, false, false, false, false, 10)
	hashBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	passwordHash := string(hashBytes)

	builder := sq.Insert("users").
		PlaceholderFormat(sq.Dollar).
		Columns("name", "email", "password_hash").
		Values(name, email, passwordHash).
		Suffix("RETURNING id")

	query, args, err := builder.ToSql()
	s.Require().NoError(err)
	s.T().Helper()

	q := db.Query{Name: "user_repository.Create", QueryRaw: query}
	var userID int64
	err = s.dbClient.DB().QueryRowContext(ctx, q, args...).Scan(&userID)
	s.Require().NoError(err)

	// generate code
	totpKey, err := totp.Generate(totp.GenerateOpts{
		Issuer:      "example.com",
		AccountName: "admin@admin.com",
		SecretSize:  15,
	})

	otpHashRedis := auth.OtpHashRedis{
		UserId:     userID,
		OtpAuthUrl: totpKey.URL(),
		OtpSecret:  totpKey.Secret(),
	}

	// key for fingerprint
	// key = fingerprint:sdlkfjklsdf
	fpKey := fmt.Sprintf("fingerprint:%s", fingerprint)
	fpValue, err := json.Marshal(otpHashRedis)
	s.Require().NoError(err)

	err = s.cacher.Set(ctx, fpKey, string(fpValue), 1*time.Minute)
	s.Require().NoError(err)

	code, err := totp.GenerateCode(totpKey.Secret(), time.Now())
	s.Require().NoError(err)

	res, err := s.client.CheckCode(ctx, &desc.CheckCodeRequest{
		Code:        code,
		Fingerprint: fingerprint,
	})

	s.Require().NoError(err)
	s.Require().NotNil(res)
	s.Assert().NotNil(res.AccessToken)
	s.Assert().NotNil(res.RefreshToken)
}
