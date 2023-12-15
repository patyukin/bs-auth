package service

import (
	"context"

	desc "github.com/patyukin/bs-auth/pkg/auth_v1"

	"github.com/patyukin/bs-auth/internal/model"
)

type UserService interface {
	Create(ctx context.Context, info *model.User) (int64, error)
	Get(ctx context.Context, id int64) (*model.User, error)
	GetByEmail(ctx context.Context, email string) (*model.User, error)
}

type AuthService interface {
	SignIn(ctx context.Context, params *model.User, fingerprint string) (*desc.SignInResponse, error)
	CheckCode(ctx context.Context, req *desc.CheckCodeRequest) (int64, error)
	SaveSession(ctx context.Context, userId int64, AuthTokenSignKey string) (*model.Tokens, error)
}
