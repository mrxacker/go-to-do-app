package usecase

import (
	"context"
	"errors"
	"strings"

	"github.com/mrxacker/go-to-do-app/internal/dto"
	e "github.com/mrxacker/go-to-do-app/internal/errors"
	"github.com/mrxacker/go-to-do-app/internal/models"
	"github.com/mrxacker/go-to-do-app/internal/ports/repository"
)

type TodoUsecase struct {
	repo repository.TodoRepository
}

func NewTodoUsecase(r repository.TodoRepository) *TodoUsecase {
	return &TodoUsecase{repo: r}
}

func (u *TodoUsecase) CreateTodo(ctx context.Context, req dto.CreateTodoRequest) (models.ToDoID, error) {
	if strings.TrimSpace(req.Title) == "" {
		return 0, errors.New("title is required")
	}

	if len(req.Title) > 200 {
		return 0, errors.New("title is too long")
	}

	todo := models.ToDo{
		UserID:      req.UserID,
		Title:       req.Title,
		Description: req.Description,
	}

	return u.repo.CreateTodo(ctx, todo)
}

func (u *TodoUsecase) GetTodoByID(ctx context.Context, id models.ToDoID) (models.ToDo, error) {
	todo, err := u.repo.GetTodoByID(ctx, id)
	if err != nil {
		return models.ToDo{}, err
	}

	return todo, nil
}

func (u *TodoUsecase) ListTodos(ctx context.Context, req dto.GetListTodosRequest) ([]models.ToDo, error) {
	return u.repo.ListTodos(ctx, req.UserID, req.Limit, req.Offset)
}

func (u *TodoUsecase) DeleteTodoByID(ctx context.Context, id models.ToDoID) error {
	todo, err := u.repo.GetTodoByID(ctx, id)
	if err != nil {
		return err
	}

	if todo.ID == 0 {
		return e.ErrTodoNotFound
	}

	return u.repo.DeleteTodoByID(ctx, id)
}

func (u *TodoUsecase) UpdateTodo(ctx context.Context, req models.ToDo) error {
	todo, err := u.repo.GetTodoByID(ctx, req.ID)
	if err != nil {
		return err
	}

	if todo.ID == 0 {
		return e.ErrTodoNotFound
	}

	return u.repo.UpdateTodo(ctx, models.ToDo{
		ID:          req.ID,
		Title:       req.Title,
		Description: req.Description,
	})
}
