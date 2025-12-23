package repository

import (
	"context"

	"github.com/mrxacker/go-to-do-app/internal/models"
)

type TodoRepository interface {
	CreateTodo(ctx context.Context, todo models.ToDo) (models.ToDoID, error)
	GetTodoByID(ctx context.Context, id models.ToDoID) (models.ToDo, error)
	ListTodos(ctx context.Context, userID models.UserID, limit, offset int) ([]models.ToDo, error)
	DeleteTodoByID(ctx context.Context, id models.ToDoID) error
	UpdateTodo(ctx context.Context, todo models.ToDo) error
}
