package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
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
}

func (h *TodoHandler) CreateTodo(c *gin.Context) {
	h.uc.CreateTodo(c.Request.Context(), "Sample Title", "Sample Description")
	c.Status(http.StatusCreated)
}
