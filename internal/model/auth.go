package model

import (
	desc "github.com/patyukin/bs-auth/pkg/auth_v1"
	"time"
)

type RefreshToken struct {
	ID        int64
	Token     string
	ExpiredAt time.Time
	CreatedAt time.Time
}

type AuthUser struct {
	Email    string
	Password string
}

type Tokens struct {
	AccessToken  string
	RefreshToken string
}

type SignInResponse = desc.SignInResponse

type CheckCodeRequest = desc.CheckCodeRequest
