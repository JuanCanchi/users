package usecase

import (
	"context"
	"errors"
	"github.com/juancanchi/users/internal/delivery/http/middleware"
	"time"

	"github.com/google/uuid"
	"github.com/juancanchi/users/internal/domain"
	"github.com/juancanchi/users/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

type UserUsecase struct {
	Repo      repository.UserRepository
	JWTSecret string
}

func NewUserUsecase(repo repository.UserRepository, jwtSecret string) *UserUsecase {
	return &UserUsecase{
		Repo:      repo,
		JWTSecret: jwtSecret,
	}
}

// Registro
func (u *UserUsecase) Register(ctx context.Context, email, password string) (*domain.User, error) {
	// verificar si ya existe
	exists, _ := u.Repo.FindByEmail(ctx, email)
	if exists != nil {
		return nil, errors.New("email ya registrado")
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &domain.User{
		ID:        uuid.New().String(),
		Email:     email,
		Password:  string(hashed),
		CreatedAt: time.Now(),
	}

	err = u.Repo.Save(ctx, user)
	if err != nil {
		return nil, err
	}

	// No devolvemos el hash
	user.Password = ""
	return user, nil
}

// Login
func (u *UserUsecase) Login(ctx context.Context, email, password string) (string, error) {
	user, err := u.Repo.FindByEmail(ctx, email)
	if err != nil || user == nil {
		return "", errors.New("credenciales inválidas")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return "", errors.New("credenciales inválidas")
	}

	token, err := middleware.GenerateJWT(user, u.JWTSecret)
	if err != nil {
		return "", err
	}

	return token, nil
}

func (u *UserUsecase) UpdateUserRole(ctx context.Context, id string, role string) error {
	return u.Repo.UpdateRole(ctx, id, role)
}
