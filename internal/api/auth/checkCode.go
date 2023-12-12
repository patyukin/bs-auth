package auth

import (
	"context"
	"github.com/golang-jwt/jwt/v5"
	desc "github.com/patyukin/bs-auth/pkg/auth_v1"
	"time"
)

const (
	TokenSignKey = "TOKEN_SIGN_KEY"
)

type UserClaims struct {
	Id    int64  `json:"id"`
	Token string `json:"token"`
}

type Claims struct {
	Id int64 `json:"id"`
	jwt.RegisteredClaims
}

func (i *Implementation) CheckCode(ctx context.Context, req *desc.CheckCodeRequest) (*desc.CheckCodeResponse, error) {
	// generate claims
	res, err := i.authService.CheckCode(ctx, req)
	if err != nil {
		return nil, err
	}

	mySigningKey := []byte(i.cfg.AuthTokenSignKey)
	claims := Claims{
		user.ID,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(1 * time.Hour)),
		},
	}

	// generate access token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	accessToken, err := token.SignedString(mySigningKey)
	if err != nil {
		return nil, err
	}

	// generate refresh token
	refreshToken, errTx := s.authRepository.Create(ctx, user.ID)
	if errTx != nil {
		return nil, errTx
	}
	return &desc.CheckCodeResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}
}
