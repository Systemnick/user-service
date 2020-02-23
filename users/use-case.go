package users

import (
	"context"

	"github.com/Systemnick/user-service/domain"
)

type UseCase interface {
	Login(ctx context.Context, cred *UserCredentials) (*domain.User, error)
	Register(ctx context.Context, params *UserParams) error
}

type UserParams struct {
	Login    string `validate:"required,max=255"`
	Email    string `validate:"required,email,max=255"`
	Password string `validate:"required,max=255"`
	Phone    string `validate:"required,max=255"`
}

type UserCredentials struct {
	Login    string `validate:"required,max=255"`
	Password string `validate:"required,max=255"`
}

type Response struct {
	StatusCode    int      `json:"status_code"`
	StatusMessage string   `json:"status_message"`
	Errors        []string `json:"errors,omitempty"`
}
