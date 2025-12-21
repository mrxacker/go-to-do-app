package repository

type TodoRepository interface {
	CreateTodo(title, description string) error
}
