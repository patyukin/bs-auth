package user

import (
	"github.com/patyukin/bs-auth/internal/service"
	desc "github.com/patyukin/bs-auth/pkg/user_v1"
)

type Implementation struct {
	desc.UnimplementedUserV1Server
	userService service.UserService
}

func NewImplementation(userService service.UserService) *Implementation {
	return &Implementation{
		userService: userService,
	}
}
