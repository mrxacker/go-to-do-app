package repository

import (
	"context"

	"github.com/mrxacker/go-to-do-app/internal/models"
)

type UserRepository interface {
	CreateUser(ctx context.Context, user models.User) (models.UserID, error)
	GetUserByEmail(ctx context.Context, email string) (models.User, error)
	GetUserByUsername(ctx context.Context, username string) (models.User, error)
	GetUserByID(ctx context.Context, id models.UserID) (models.User, error)
}
