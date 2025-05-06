package v1

import (
	"net/http"
	"routinist/internal/dto"
	"routinist/internal/middleware"
	"routinist/internal/usecase"
	"routinist/pkg/logger"

	"github.com/gin-gonic/gin"
)

type HabitHandler struct {
	usecase usecase.HabitUsecase
	logger  logger.Interface
}

func NewHabitRoutes(handler *gin.RouterGroup, t usecase.HabitUsecase, l logger.Interface) {
	r := &HabitHandler{t, l}

	h1 := handler.Group("/habit")
	{
		h1.GET("/random", r.getRandomHabits)
	}

	auth := handler.Group("/protected/habit", middleware.JWTAuthMiddleware())
	{
		auth.GET("/today", r.getTodayHabits)
	}
}

func (r *HabitHandler) getRandomHabits(c *gin.Context) {
	response := dto.Response{}

	habits, err := r.usecase.GetRandomHabits()
	if err != nil {
		r.logger.Error(err)
		response.SetMessage("Failed to get random habits")
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	response.Data = habits
	c.JSON(http.StatusOK, response)
}

func (r *HabitHandler) getTodayHabits(c *gin.Context) {
	response := dto.Response{}

	userIDVal, exists := c.Get("user_id")
	if !exists {
		response.SetMessage("Unauthorized")
		c.JSON(http.StatusUnauthorized, response)
		return
	}

	userId := userIDVal.(uint)

	habits, err := r.usecase.GetTodayHabits(userId)
	if err != nil {
		r.logger.Error(err)
		response.SetMessage("Failed to get today habits")
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	response.Data = habits
	c.JSON(http.StatusOK, response)
}
