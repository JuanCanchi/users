package repository

import (
	"context"
	"personal-jujuy/users/internal/domain"
)

type UserRepository interface {
	Save(ctx context.Context, user *domain.User) error
	FindByEmail(ctx context.Context, email string) (*domain.User, error)
}
