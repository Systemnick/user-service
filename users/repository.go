package users

import (
	"context"

	"github.com/Systemnick/user-service/domain"
)

type Repository interface {
	Store(ctx context.Context, user *domain.User) error
	Fetch(ctx context.Context, login, password string) (*domain.User, error)
}
