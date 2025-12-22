package repository

import (
	"context"

	"github.com/mrxacker/go-to-do-app/internal/dto"
	"github.com/mrxacker/go-to-do-app/internal/models"
)

type TodoRepository interface {
	CreateTodo(ctx context.Context, todo dto.CreateTodoRequest) (models.ToDoID, error)
	GetTodoByID(ctx context.Context, id models.ToDoID) (models.ToDo, error)
	ListTodos(ctx context.Context, limit, offset int) ([]models.ToDo, error)
	DeleteTodoByID(ctx context.Context, id models.ToDoID) error
	UpdateTodo(ctx context.Context, todo models.ToDo) error
}
