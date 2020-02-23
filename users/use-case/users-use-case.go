package use_case

import (
	"context"
	"time"

	"github.com/Systemnick/user-service/domain"
	"github.com/Systemnick/user-service/users"
	"github.com/twinj/uuid"
)

type usersUseCase struct {
	userRepo       users.Repository
	contextTimeout time.Duration
}

func New(cRepo users.Repository, timeout time.Duration) users.UseCase {

	return &usersUseCase{
		userRepo:       cRepo,
		contextTimeout: timeout,
	}
}

func (u usersUseCase) Register(ctx context.Context, params *users.UserParams) error {
	user := &domain.User{
		ID:           domain.UserID(uuid.NewV4().String()),
		Login:        params.Login,
		Password:     params.Password,
		Email:        params.Email,
		Phone:        params.Phone,
		CreationTime: time.Now(),
	}

	err := u.userRepo.Store(ctx, user)

	return err
}

func (u usersUseCase) Login(ctx context.Context, credentials *users.UserCredentials) (*domain.User, error) {
	user, err := u.userRepo.Fetch(ctx, credentials.Login, credentials.Password)

	return user, err
}
