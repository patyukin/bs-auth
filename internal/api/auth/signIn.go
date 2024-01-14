package auth

import (
	"context"
	"github.com/patyukin/bs-auth/internal/model"
	desc "github.com/patyukin/bs-auth/pkg/auth_v1"
	"golang.org/x/crypto/bcrypt"
)

func (i *Implementation) SignIn(ctx context.Context, req *desc.SignInRequest) (*model.SignInResponse, error) {
	user, err := i.userService.GetByEmail(ctx, req.GetEmail())
	if err != nil {
		return nil, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.GetPassword()))
	if err != nil {
		return nil, err
	}

	// TODO response
	return i.authService.SignIn(ctx, user, req.GetFingerprint())
}
