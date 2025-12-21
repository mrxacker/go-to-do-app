package repository

import "context"

type TodoRepository interface {
	CreateTodo(ctx context.Context, title, description string) error
}
