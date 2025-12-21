package models

import "time"

type ToDoID int64

type ToDo struct {
	ID          ToDoID    `db:"id"`
	Title       string    `db:"title"`
	Description string    `db:"description"`
	CreatedAt   time.Time `db:"created_at"`
	UpdatedAt   time.Time `db:"updated_at"`
}
