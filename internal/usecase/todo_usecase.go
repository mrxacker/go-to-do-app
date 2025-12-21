package usecase

import "github.com/mrxacker/go-to-do-app/internal/ports/repository"

type TodoUsecase struct {
	repo repository.TodoRepository
}

func NewTodoUsecase(r repository.TodoRepository) *TodoUsecase {
	return &TodoUsecase{repo: r}
}

func (u *TodoUsecase) CreateTodo(title, description string) error {
	return u.repo.CreateTodo(title, description)
}
