package v1

import (
	"net/http"
	"routinist/internal/dto"
	"routinist/internal/entity"
	"routinist/internal/usecase"
	"routinist/pkg/logger"
	"strings"

	"github.com/gin-gonic/gin"
)

type authRoutes struct {
	t usecase.Auth
	l logger.Interface
}

func NewAuthRoutes(handler *gin.RouterGroup, t usecase.Auth, l logger.Interface) {
	r := &authRoutes{t, l}

	h1 := handler.Group("/auth")
	{
		h1.POST("/register", r.register)
		h1.POST("/login", r.login)
	}
}

func (r *authRoutes) register(c *gin.Context) {
	response := dto.Response{}

	var req entity.RegisterRequestDTO
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

	token, err := r.t.Register(
		c.Request.Context(),
		req,
	)

	if err != nil {
		response.SetMessage(err.Error())
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	response.Data = token
	c.JSON(http.StatusOK, response)
}

func (r *authRoutes) login(c *gin.Context) {
	response := dto.Response{}

	var req entity.LoginRequestDTO
	if err := c.Bind(&req); err != nil {
		response.SetMessage("Invalid request")
		c.JSON(http.StatusBadRequest, response)
		return
	}

	token, err := r.t.Login(
		c.Request.Context(),
		req,
	)

	if err != nil {
		response.SetMessage(err.Error())
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	response.Data = token
	c.JSON(http.StatusOK, response)
}
