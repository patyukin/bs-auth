package auth

import (
	"context"
	"fmt"
	desc "github.com/patyukin/bs-auth/pkg/auth_v1"
)

// CheckCode request => dto
func (i *Implementation) CheckCode(ctx context.Context, req *desc.CheckCodeRequest) (*desc.CheckCodeResponse, error) {
	key := req.GetFingerprint()
	err := i.rl.Increment(key)
	if err != nil {
		return nil, fmt.Errorf("failed to increment counter: %w", err)
	}

	userId, err := i.authService.CheckCode(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to check code: %w", err)
	}

	// save refresh token into redis db
	tokens, err := i.authService.SaveSession(ctx, userId, i.cfg.AuthTokenSignKey)
	if err != nil {
		return nil, fmt.Errorf("failed to save session: %w", err)
	}

	return &desc.CheckCodeResponse{
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
	}, nil
}
