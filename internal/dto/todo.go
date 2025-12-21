package dto

type CreateTodoRequest struct {
	Title       string `json:"title" binding:"required"`
	Description string `json:"description" binding:"required"`
}

type CreateTodoResponse struct {
	ID int64 `json:"id"`
}

type ListTodosResponse struct {
	Todos []TodoItem `json:"todos"`
}

type TodoItem struct {
	ID          int64  `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}
