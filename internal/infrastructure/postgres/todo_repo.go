package postgres

import (
	"context"
	"database/sql"
)

type TodoRepo struct {
	db *sql.DB
}

func NewTodoRepo(db *sql.DB) *TodoRepo {
	return &TodoRepo{db: db}
}

func (r *TodoRepo) CreateTodo(ctx context.Context, title, description string) error {
	_, err := r.db.ExecContext(ctx, "INSERT INTO todos (title, description) VALUES ($1, $2)", title, description)
	return err
}

func (r *TodoRepo) ListTodos() error {