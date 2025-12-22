package dto

import "github.com/mrxacker/go-to-do-app/internal/models"

type CreateTodoRequest struct {
	Title       string `json:"title" binding:"required"`
	Description string `json:"description" binding:"required"`
}

type UpdateTodoURI struct {
	ID models.ToDoID `uri:"id" binding:"required"`
}

type CreateTodoResponse struct {
	ID models.ToDoID `json:"id"`
}

type GetTodoByIDRequest struct {
	ID models.ToDoID `uri:"id" binding:"required"`
}

type GetListTodosRequest struct {
	Limit  int `form:"limit"`
	Offset int `form:"offset"`
}

type ListTodosResponse struct {
	Todos []TodoItem `json:"todos"`
}

type TodoItem struct {
	ID          models.ToDoID `json:"id"`
	Title       string        `json:"title"`
	Description string        `json:"description"`
}
