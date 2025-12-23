package repository

import (
	"context"

	"github.com/mrxacker/go-to-do-app/internal/models"
)

type UserRepository interface {
	CreateUser(ctx context.Context, user models.User) (models.UserID, error)
	GetUserByEmail(ctx context.Context, email string) (models.User, error)
}
