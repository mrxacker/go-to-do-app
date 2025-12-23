package usecase

import (
	"context"
	"errors"

	e "github.com/mrxacker/go-to-do-app/internal/errors"
	"github.com/mrxacker/go-to-do-app/internal/models"
	"github.com/mrxacker/go-to-do-app/internal/ports/repository"
)

type UserUseCase struct {
	userRepo repository.UserRepository
}

func NewUserUseCase(r repository.UserRepository) *UserUseCase {
	return &UserUseCase{userRepo: r}
}

func (u *UserUseCase) CreateUser(ctx context.Context, user models.User) (models.UserID, error) {
	// Check email uniqueness
	_, err := u.userRepo.GetUserByEmail(ctx, user.Email)
	if err == nil {
		return 0, e.ErrUserAlreadyExists
	}
	if !errors.Is(err, e.ErrUserNotFound) {
		return 0, err
	}

	// Check username uniqueness
	_, err = u.userRepo.GetUserByUsername(ctx, user.Username)
	if err == nil {
		return 0, e.ErrUserAlreadyExists
	}
	if !errors.Is(err, e.ErrUserNotFound) {
		return 0, err
	}

	// Create user
	return u.userRepo.CreateUser(ctx, user)
}
