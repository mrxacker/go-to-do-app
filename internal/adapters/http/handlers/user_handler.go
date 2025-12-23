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

type UserHandler struct {
	uc *usecase.UserUseCase
}

func NewUserHandler(uc *usecase.UserUseCase) *UserHandler {
	return &UserHandler{uc: uc}
}

func (h *UserHandler) RegisterRoutes(rg *gin.RouterGroup) {
	rg.POST("/register", h.CreateUser)
	rg.POST("/login", h.LoginUser)
}

func (h *UserHandler) CreateUser(c *gin.Context) {
	var req dto.RegisterUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	user := models.User{
		Username:     req.Username,
		Email:        req.Email,
		PasswordHash: req.Password,
	}

	_, err := h.uc.CreateUser(c.Request.Context(), user)
	if err != nil {
		if errors.Is(err, e.ErrUserAlreadyExists) {
			c.JSON(http.StatusConflict, gin.H{"error": "User already exists"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	c.Status(http.StatusCreated)
}

func (h *UserHandler) LoginUser(c *gin.Context) {
	var req dto.LoginUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	token, err := h.uc.LoginUser(c.Request.Context(), req.Email, req.Password)
	if err != nil {
		if errors.Is(err, e.ErrUserNotFound) || errors.Is(err, e.ErrInvalidIdentifier) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to login"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"access": token,
	})
}
