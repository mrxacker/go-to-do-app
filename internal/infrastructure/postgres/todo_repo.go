package postgres

import "database/sql"

type TodoRepo struct {
	db *sql.DB
}

func NewTodoRepo(db *sql.DB) *TodoRepo {
	return &TodoRepo{db: db}
}

func (r *TodoRepo) CreateTodo(title, description string) error {
	_, err := r.db.Exec("INSERT INTO todos (title, description) VALUES ($1, $2)", title, description)
	return err
}
