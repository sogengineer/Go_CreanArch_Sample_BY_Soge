package container

import (
	"context"

	user "github.com/Go_CleanArch/infrastructure/container/user"
)

type Container struct {
	UserContainer *user.UserContainer
}

func NewContainer(ctx context.Context) (*Container, error) {
	userContainer, err := user.NewContainer(ctx)
	if err != nil {
		return nil, err
	}

	return &Container{
		UserContainer: userContainer,
	}, nil
}
