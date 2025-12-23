package usecase

import (
	"context"

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
	user, err := u.userRepo.GetUserByEmail(ctx, user.Email)
	if err != nil {
		return 0, err
	}

	if user != (models.User{}) {
		return 0, e.ErrUserAlreadyExists
	}
	user, err := u.userRepo.GetUserByEmail(ctx, user.Email)
	if err != nil {
		return 0, err
	}

	if user != (models.User{}) {
		return 0, e.ErrUserAlreadyExists
	}



	return u.userRepo.CreateUser(ctx, user)
}
