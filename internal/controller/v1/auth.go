package v1

import (
	"errors"
	"net/http"
	domainErr "routinist/internal/domain/errors"
	"routinist/internal/domain/model"
	"routinist/internal/dto"
	"routinist/internal/usecase"
	"routinist/pkg/logger"
	"strings"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	t usecase.AuthUseCase
	l logger.Interface
}

func NewAuthRoutes(handler *gin.RouterGroup, t usecase.AuthUseCase, l logger.Interface) {
	r := &AuthHandler{t, l}

	h1 := handler.Group("/auth")
	{
		h1.POST("/register", r.register)
		h1.POST("/login", r.login)
	}
}

func (r *AuthHandler) register(c *gin.Context) {
	response := dto.Response{}

	var req model.RegisterRequestDTO
	if err := c.Bind(&req); err != nil {
		response.SetMessage("Invalid request")
		c.JSON(http.StatusBadRequest, response)
		return
	}

	// Validate email and password
	if req.Email == "" {
		response.SetMessage("Email is required")
		c.JSON(http.StatusBadRequest, response)
		return
	}

	if req.Password == "" {
		response.SetMessage("Password is required")
		c.JSON(http.StatusBadRequest, response)
		return
	}

	// Check if email is valid format
	if !strings.Contains(req.Email, "@") || !strings.Contains(req.Email, ".") {
		response.SetMessage("Invalid email format")
		c.JSON(http.StatusBadRequest, response)
		return
	}

	// Check if password meets minimum requirements (e.g. at least 6 characters)
	if len(req.Password) < 6 {
		response.SetMessage("Password must be at least 6 characters")
		c.JSON(http.StatusBadRequest, response)
		return
	}

	token, err := r.t.Register(&req)

	if err != nil {
		if errors.Is(err, domainErr.ErrEmailAlreadyExists) {
			response.SetMessage("User with this email already exists")
			c.JSON(http.StatusBadRequest, response)
		} else {
			response.SetMessage("Something went wrong")
			c.JSON(http.StatusInternalServerError, response)
		}
		return
	}

	response.Data = token
	c.JSON(http.StatusOK, response)
}

func (r *AuthHandler) login(c *gin.Context) {
	response := dto.Response{}

	var req model.LoginRequestDTO
	if err := c.Bind(&req); err != nil {
		response.SetMessage("Invalid request")
		c.JSON(http.StatusBadRequest, response)
		return
	}

	token, err := r.t.Login(&req)

	if err != nil {
		if errors.Is(err, domainErr.ErrInvalidCredentials) {
			response.SetMessage("Invalid username or password")
			c.JSON(http.StatusBadRequest, response)
		} else {
			response.SetMessage("Something went wrong")
			c.JSON(http.StatusInternalServerError, response)
		}
		return
	}

	response.Data = token
	c.JSON(http.StatusOK, response)
}
