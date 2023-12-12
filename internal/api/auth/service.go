package auth

import (
	"github.com/patyukin/bs-auth/internal/config"
	"github.com/patyukin/bs-auth/internal/service"
	desc "github.com/patyukin/bs-auth/pkg/auth_v1"
)

type Implementation struct {
	desc.UnimplementedAuthV1Server
	authService service.AuthService
	userService service.UserService
	cfg         *config.Config
}

func NewImplementation(authService service.AuthService, userService service.UserService, cfg *config.Config) *Implementation {
	return &Implementation{
		authService: authService,
		userService: userService,
		cfg:         cfg,
	}
}
