package postgres

import (
	"context"
	"database/sql"
	"errors"

	e "github.com/mrxacker/go-to-do-app/internal/errors"
	"github.com/mrxacker/go-to-do-app/internal/models"
)

type UserRepo struct {
	db *sql.DB
}

func NewUserRepo(db *sql.DB) *UserRepo {
	return &UserRepo{db: db}
}

func (r *UserRepo) CreateUser(ctx context.Context, user models.User) (models.UserID, error) {
	var id models.UserID
	err := r.db.QueryRowContext(ctx,
		"INSERT INTO users (username, email, password_hash) VALUES ($1, $2, $3) RETURNING id",
		user.Username, user.Email, user.PasswordHash).Scan(&id)
	return id, err
}

func (r *UserRepo) getUser(ctx context.Context, query string, arg any) (models.User, error) {

	var user models.User
	const baseUserSelect = `SELECT id, username, email, password_hash FROM users`
	err := r.db.QueryRowContext(ctx, baseUserSelect+" "+query, arg).
		Scan(&user.ID, &user.Username, &user.Email, &user.PasswordHash)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.User{}, e.ErrUserNotFound
		}
		return models.User{}, err
	}

	return user, nil
}

func (r *UserRepo) GetUserByID(ctx context.Context, id models.UserID) (models.User, error) {
	return r.getUser(ctx, "WHERE id = $1", id)
}

func (r *UserRepo) GetUserByEmail(ctx context.Context, email string) (models.User, error) {
	return r.getUser(ctx, "WHERE email = $1", email)
}

func (r *UserRepo) GetUserByUsername(ctx context.Context, username string) (models.User, error) {
	return r.getUser(ctx, "WHERE username = $1", username)
}
