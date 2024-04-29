package user

import (
	"context"

	userController "github.com/Go_CleanArch/interface_adapter/controller"
	userRepository "github.com/Go_CleanArch/interface_adapter/gateway/repository"
	userService "github.com/Go_CleanArch/usecase/service/user"
)

type UserContainer struct {
	UserController *userController.UserController
}

func NewContainer(ctx context.Context) (*UserContainer, error) {
	// DI注入
	userRepository, err := userRepository.NewUserRepository(ctx)
	if err != nil {
		return nil, err
	}
	userSvc := userService.NewUserService(userRepository)
	userCtrl := userController.NewUserController(*userSvc)

	return &UserContainer{
		UserController: userCtrl,
	}, nil
}
