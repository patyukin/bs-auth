package auth

import (
	"context"
	"fmt"
	desc "github.com/patyukin/bs-auth/pkg/auth_v1"
)

func (s *serv) CheckCode(ctx context.Context, req *desc.CheckCodeRequest) (*desc.CheckCodeResponse, error) {
	key := req.GetFingerprint()
	value, err := s.cacher.Get(ctx, key)
	if err != nil {
		return nil, fmt.Errorf("can't get fingerprint from cache: %w", err)
	}

	if value != req.GetCode() {
		return nil, fmt.Errorf("wrong code")
	}

	return &desc.CheckCodeResponse{
		AccessToken:  "",
		RefreshToken: "",
	}, nil
}
