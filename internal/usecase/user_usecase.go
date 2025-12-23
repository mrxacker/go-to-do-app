package usecase

import (
	"context"
	"errors"

	e "github.com/mrxacker/go-to-do-app/internal/errors"
	"github.com/mrxacker/go-to-do-app/internal/infrastructure/auth"
	"github.com/mrxacker/go-to-do-app/internal/models"
	"github.com/mrxacker/go-to-do-app/internal/ports/repository"
)

type UserUseCase struct {
	userRepo   repository.UserRepository
	jwtService *auth.JWTService
}

func NewUserUseCase(r repository.UserRepository, jwtService *auth.JWTService) *UserUseCase {
	return &UserUseCase{userRepo: r, jwtService: jwtService}
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
	hashedPassword, err := auth.HashPassword(user.PasswordHash, auth.DefaultArgonParams)
	if err != nil {
		return 0, err
	}
	user.PasswordHash = hashedPassword
	return u.userRepo.CreateUser(ctx, user)
}

func (u *UserUseCase) LoginUser(ctx context.Context, email, password string) (string, error) {
	user, err := u.userRepo.GetUserByEmail(ctx, email)
	if err != nil {
		return "", err
	}

	// Here you would normally check the password hash
	ok, err := auth.VerifyPassword(password, user.PasswordHash)
	if err != nil {
		return "", err
	}
	if !ok {
		return "", e.ErrInvalidIdentifier
	}

	token, err := u.jwtService.GenerateToken(user)
	if err != nil {
		return "", err
	}

	return token, nil
}
