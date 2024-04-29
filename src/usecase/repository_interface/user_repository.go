package repository

import (
	"context"

	entity "github.com/Go_CleanArch/interface_adapter/gateway/entity"
)

type UserRepository interface {
	FindUserByEmail(ctx context.Context, email string) (*entity.User, error)
	CreateUser(ctx context.Context, userJson []byte) (*entity.User, error)
}
