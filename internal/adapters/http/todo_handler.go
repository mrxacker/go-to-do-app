package http

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mrxacker/go-to-do-app/internal/dto"
	e "github.com/mrxacker/go-to-do-app/internal/errors"
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
}

func (h *TodoHandler) CreateTodo(c *gin.Context) {
	var req dto.CreateTodoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
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
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	todo, err := h.uc.GetTodoByID(c.Request.Context(), req.ID)
	if err != nil {
		if errors.Is(err, e.ErrTodoNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "todo not found"})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
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
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	todos, err := h.uc.ListTodos(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
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
