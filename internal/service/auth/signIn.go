package auth

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/patyukin/bs-auth/internal/model"
	desc "github.com/patyukin/bs-auth/pkg/auth_v1"
	"time"
)

type HashRedis struct {
	RefreshToken string `json:"refresh_token"`
	Code         string `json:"code"`
}

func (s *serv) SignIn(ctx context.Context, user *model.User, fingerprint string) (*desc.SignInResponse, error) {
	// Генерация случайного UUID
	var key string
	fpCodeKey := fingerprint

	// fingerprint
	fpValue, err := s.cacher.Get(ctx, fpCodeKey)
	if err != nil {
		return nil, fmt.Errorf("can't get fingerprint from cache: %w", err)
	}

	err = s.cacher.Set(ctx, fpCodeKey, fpValue, 30*time.Minute)
	if err != nil {
		return nil, fmt.Errorf("can't set fingerprint in cache: %w", err)
	}

	// checking code
	uuid := uuid.New()
	// find unique key for caching fingerprint
	for {
		key, err = s.cacher.Get(ctx, uuid.String())
		if err != nil {
			return nil, err
		}

		if key != "" {
			continue
		}

		break
	}

	return &desc.SignInResponse{
		Base32:     "",
		OtpAuthUrl: "",
	}, nil
}
