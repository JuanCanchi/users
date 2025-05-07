package repository

import (
	"context"
	"github.com/juancanchi/users/internal/domain"
)

type UserRepository interface {
	Save(ctx context.Context, user *domain.User) error
	FindByEmail(ctx context.Context, email string) (*domain.User, error)
}
