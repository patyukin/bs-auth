package auth

import (
	"context"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/patyukin/bs-auth/internal/model"
	"strconv"
	"time"
)

type UserClaims struct {
	Id    int64  `json:"id"`
	Token string `json:"token"`
}

type Claims struct {
	Id int64 `json:"id"`
	jwt.RegisteredClaims
}

func (s *serv) SaveSession(ctx context.Context, userId int64, AuthTokenSignKey string) (*model.Tokens, error) {
	mySigningKey := []byte(AuthTokenSignKey)
	claims := Claims{
		userId,
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

	// check unique key in cacher
	var u string
	for {
		u = uuid.New().String()
		uKey := fmt.Sprintf("token:refresh:%s", u)
		isExist, err := s.cacher.Exist(ctx, uKey)
		if err != nil {
			return nil, fmt.Errorf("failed to check unique key in cacher: %w", err)
		}

		if !isExist {
			err = s.cacher.Set(ctx, uKey, strconv.FormatInt(userId, 10), 1*time.Hour)
			if err != nil {
				return nil, fmt.Errorf("failed to set unique key in cacher: %w", err)
			}

			break
		}
	}

	return &model.Tokens{
		AccessToken:  accessToken,
		RefreshToken: u,
	}, nil

}
