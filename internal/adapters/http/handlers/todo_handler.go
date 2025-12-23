package http

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mrxacker/go-to-do-app/internal/dto"
	e "github.com/mrxacker/go-to-do-app/internal/errors"
	"github.com/mrxacker/go-to-do-app/internal/models"
	"github.com/mrxacker/go-to-do-app/internal/usecase"
)

type TodoHandler struct {
	uc *usecase.TodoUsecase
}

func NewTodoHandler(uc *usecase.TodoUsecase) *TodoHandler {
	return &TodoHandler{uc: uc}
}

func (h *TodoHandler) RegisterRoutes(rg *gin.RouterGroup) {
	rg.POST("/", h.CreateTodo)
	rg.GET("/", h.ListTodos)
	rg.GET("/:id", h.GetTodoByID)
	rg.DELETE("/:id", h.DeleteTodoByID)
	rg.PUT("/:id", h.UpdateTodo)
}

func (h *TodoHandler) CreateTodo(c *gin.Context) {
	var req dto.CreateTodoRequest
	var userID models.UserID
	userIDValue, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	userID = userIDValue.(models.UserID)

	req.UserID = userID

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	id, err := h.uc.CreateTodo(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var res dto.CreateTodoResponse
	res.ID = id

	c.JSON(http.StatusCreated, res)
}

func (h *TodoHandler) GetTodoByID(c *gin.Context) {
	var req dto.GetTodoByIDRequest
	if err := c.ShouldBindUri(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid todo ID"})
		return
	}

	todo, err := h.uc.GetTodoByID(c.Request.Context(), req.ID)
	if err != nil {
		if errors.Is(err, e.ErrTodoNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "todo not found"})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get todo"})
		return
	}

	res := dto.TodoItem{
		ID:          todo.ID,
		Title:       todo.Title,
		Description: todo.Description,
	}

	c.JSON(http.StatusOK, res)
}

func (h *TodoHandler) ListTodos(c *gin.Context) {
	var req dto.GetListTodosRequest
	var userID models.UserID
	userIDValue, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	userID = userIDValue.(models.UserID)

	req.UserID = userID

	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid query parameters"})
		return
	}

	todos, err := h.uc.ListTodos(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to list todos"})
		return
	}

	res := make([]dto.TodoItem, len(todos))
	for i, todo := range todos {
		res[i] = dto.TodoItem{
			ID:          todo.ID,
			Title:       todo.Title,
			Description: todo.Description,
		}
	}

	c.JSON(http.StatusOK, res)
}

func (h *TodoHandler) DeleteTodoByID(c *gin.Context) {
	var req dto.GetTodoByIDRequest
	if err := c.ShouldBindUri(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid todo ID"})
		return
	}

	err := h.uc.DeleteTodoByID(c.Request.Context(), req.ID)
	if err != nil {
		if errors.Is(err, e.ErrTodoNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "todo not found"})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete todo"})
		return
	}

	c.Status(http.StatusNoContent)
}

func (h *TodoHandler) UpdateTodo(c *gin.Context) {
	var uri dto.UpdateTodoURI
	if err := c.ShouldBindUri(&uri); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid todo ID"})
		return
	}

	var req dto.CreateTodoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	err := h.uc.UpdateTodo(c.Request.Context(), models.ToDo{
		ID:          uri.ID,
		Title:       req.Title,
		Description: req.Description,
	})
	if err != nil {
		if errors.Is(err, e.ErrTodoNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "todo not found"})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update todo"})
		return
	}

	c.Status(http.StatusOK)
}
