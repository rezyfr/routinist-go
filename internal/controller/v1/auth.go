package v1

import (
	"errors"
	"net/http"
	domainErr "routinist/internal/domain/errors"
	"routinist/internal/dto/request"
	"routinist/internal/dto/response"
	"routinist/internal/middleware"
	"routinist/internal/usecase"
	"routinist/pkg/logger"
	"strings"
	"time"

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

	h2 := handler.Group("/auth/protected", middleware.JWTAuthMiddleware())
	{
		h2.GET("/check", r.CheckToken)
	}
}

func (h *AuthHandler) register(c *gin.Context) {
	r := response.Response{}

	var req request.RegisterRequestDTO
	if err := c.Bind(&req); err != nil {
		r.SetMessage("Invalid request")
		c.JSON(http.StatusBadRequest, r)
		return
	}

	// Validate email and password
	if req.Email == "" {
		r.SetMessage("Email is required")
		c.JSON(http.StatusBadRequest, r)
		return
	}

	if req.Password == "" {
		r.SetMessage("Password is required")
		c.JSON(http.StatusBadRequest, r)
		return
	}

	// Check if email is valid format
	if !strings.Contains(req.Email, "@") || !strings.Contains(req.Email, ".") {
		r.SetMessage("Invalid email format")
		c.JSON(http.StatusBadRequest, r)
		return
	}

	// Check if password meets minimum requirements (e.g. at least 6 characters)
	if len(req.Password) < 6 {
		r.SetMessage("Password must be at least 6 characters")
		c.JSON(http.StatusBadRequest, r)
		return
	}

	token, err := h.t.Register(&req)

	if err != nil {
		if errors.Is(err, domainErr.ErrEmailAlreadyExists) {
			r.SetMessage("User with this email already exists")
			c.JSON(http.StatusBadRequest, r)
		} else {
			r.SetMessage("Something went wrong")
			c.JSON(http.StatusInternalServerError, r)
		}
		return
	}

	r.Data = token
	c.JSON(http.StatusOK, r)
}

func (h *AuthHandler) login(c *gin.Context) {
	r := response.Response{}

	var req request.LoginRequestDTO
	if err := c.Bind(&req); err != nil {
		r.SetMessage("Invalid request")
		c.JSON(http.StatusBadRequest, r)
		return
	}

	token, err := h.t.Login(&req)

	if err != nil {
		if errors.Is(err, domainErr.ErrInvalidCredentials) {
			r.SetMessage("Invalid username or password")
			c.JSON(http.StatusBadRequest, r)
		} else {
			r.SetMessage("Something went wrong")
			c.JSON(http.StatusInternalServerError, r)
		}
		return
	}

	r.Data = token
	c.JSON(http.StatusOK, r)
}

func (h *AuthHandler) CheckToken(c *gin.Context) {
	expiredAt, exist := c.Get("expired_at")

	if !exist {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid or expired token"})
		return
	}

	// Check expiration
	if expiredAt != nil && expiredAt.(time.Time).Before(time.Now()) {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid or expired token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Token valid"})
}
