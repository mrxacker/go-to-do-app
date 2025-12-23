package postgres

import (
	"context"
	"database/sql"
	"errors"

	e "github.com/mrxacker/go-to-do-app/internal/errors"
	"github.com/mrxacker/go-to-do-app/internal/models"
)

type TodoRepo struct {
	db *sql.DB
}

func NewTodoRepo(db *sql.DB) *TodoRepo {
	return &TodoRepo{db: db}
}

func (r *TodoRepo) CreateTodo(ctx context.Context, todo models.ToDo) (models.ToDoID, error) {
	var id models.ToDoID
	err := r.db.QueryRowContext(ctx,
		"INSERT INTO to_do (user_id, title, description) VALUES ($1, $2, $3) RETURNING id",
		todo.UserID, todo.Title, todo.Description).Scan(&id)
	return id, err
}

func (r *TodoRepo) GetTodoByID(ctx context.Context, id models.ToDoID) (models.ToDo, error) {
	var todo models.ToDo
	err := r.db.QueryRowContext(ctx,
		"SELECT id, title, description, created_at, updated_at FROM to_do WHERE id = $1",
		id).Scan(&todo.ID, &todo.Title, &todo.Description, &todo.CreatedAt, &todo.UpdatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.ToDo{}, e.ErrTodoNotFound
		}
		return models.ToDo{}, err
	}
	return todo, nil
}

func (r *TodoRepo) ListTodos(ctx context.Context, userID models.UserID, limit, offset int) ([]models.ToDo, error) {
	if limit <= 0 {
		limit = 20
	}
	if offset < 0 {
		offset = 0
	}

	rows, err := r.db.QueryContext(ctx, "SELECT id, title, description FROM to_do WHERE user_id = $1 ORDER BY id LIMIT $2 OFFSET $3", userID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	todos := make([]models.ToDo, 0)
	for rows.Next() {
		var todo models.ToDo
		if err := rows.Scan(&todo.ID, &todo.Title, &todo.Description); err != nil {
			return nil, err
		}
		todos = append(todos, todo)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return todos, nil
}

func (r *TodoRepo) DeleteTodoByID(ctx context.Context, id models.ToDoID) error {
	result, err := r.db.ExecContext(ctx, "DELETE FROM to_do WHERE id = $1", id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return e.ErrTodoNotFound
	}

	return nil
}

func (r *TodoRepo) UpdateTodo(ctx context.Context, todo models.ToDo) error {
	result, err := r.db.ExecContext(ctx,
		"UPDATE to_do SET title = $1, description = $2, updated_at = NOW() WHERE id = $3",
		todo.Title, todo.Description, todo.ID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return e.ErrTodoNotFound
	}

	return nil
}
