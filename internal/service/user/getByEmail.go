package user

import (
	"context"
	"github.com/patyukin/bs-auth/internal/model"
)

func (s *serv) GetByEmail(ctx context.Context, email string) (*model.User, error) {
	user, err := s.userRepository.GetByEmail(ctx, email)
	if err != nil {
		return nil, err
	}

	return user, nil
}
